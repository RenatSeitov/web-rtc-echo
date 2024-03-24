package controllers

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"errors"
	"web-rtc-echo/interfaces"
	"web-rtc-echo/utils"
	"net/http"
	"sync"
    "fmt"
	"github.com/labstack/echo/v4"
	"github.com/gorilla/websocket"
)

// ConnectSession - Given a host and a password returns the session object.
func ConnectSession(c echo.Context) error {
    db := c.Get("db").(*sql.DB)

    url := c.Param("url")
    password := c.QueryParam("password")
    fmt.Println(password)
    row := db.QueryRow("SELECT sessionID FROM sockets WHERE hashedURL = ?", url)
    var sessionID string
    if err := row.Scan(&sessionID); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Socket connection not found."})
    }

    row = db.QueryRow("SELECT title, host, password FROM sessions WHERE id = ?", sessionID)
    var title, host, dbPassword string
    if err := row.Scan(&title, &host, &dbPassword); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Session not found."})
    }
    fmt.Println(dbPassword) 
    if !utils.ComparePasswords(dbPassword, []byte(password)) {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid password."})
    }

    return c.JSON(http.StatusOK, map[string]string{"title": title, "socket": host})
}


// GetSession - Checks if session exists.
func GetSession(c echo.Context) error {
	db := c.Get("db").(*sql.DB)

	id := c.QueryParam("url")

	row := db.QueryRow("SELECT id FROM sockets WHERE hashedURL = ?", id)
	var sessionID string
	if err := row.Scan(&sessionID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Socket connection not found."})
	}

	return c.NoContent(http.StatusOK)
}

// CreateSocket - Creates socket connection with given session
func CreateSocket(session interfaces.Session, c echo.Context, id string) (string, error) {
	db := c.Get("db").(*sql.DB)

	hashURL := hashSession(session.Host + session.Title)
	socketURL := hashSession(session.Host + session.Password)

	stmt, err := db.Prepare("INSERT INTO sockets(sessionID, hashedURL, socketURL) VALUES(?, ?, ?)")
	if err != nil {
		return "", err // обработка ошибок подготовки запроса
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, hashURL, socketURL)
	if err != nil {
		return "", err // обработка ошибок вставки
	}

	return hashURL, nil
}

func hashSession(str string) string {
	hash := sha1.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

// ConnectionManager - Manages websocket connections
type ConnectionManager struct {
	connections map[string]*interfaces.Connection
	mu          sync.Mutex
}

// NewConnectionManager - Creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*interfaces.Connection),
	}
}

// AddConnection - Adds a new websocket connection
func (cm *ConnectionManager) AddConnection(sessionID string, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.connections[sessionID] = &interfaces.Connection{Socket: conn}
}

// RemoveConnection - Removes a websocket connection
func (cm *ConnectionManager) RemoveConnection(sessionID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	delete(cm.connections, sessionID)
}

// SendMessage - Sends a message to a websocket connection
func (cm *ConnectionManager) SendMessage(sessionID string, message interfaces.Message) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	conn, ok := cm.connections[sessionID]
	if !ok {
		return errors.New("connection not found")
	}

	return conn.Send(message)
}

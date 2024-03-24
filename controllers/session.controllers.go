package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"web-rtc-echo/interfaces"
	"web-rtc-echo/utils"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

// CreateSession - Creates user session
func CreateSession(c echo.Context) error {
	db := c.Get("db").(*sql.DB)

	var session interfaces.Session
	if err := c.Bind(&session); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	session.Password = utils.HashPassword(session.Password)

	stmt, err := db.Prepare("INSERT INTO sessions(host, title, password) VALUES(?, ?, ?)")
	if err != nil {
		return err // обработка ошибок подготовки запроса
	}
	defer stmt.Close()

	result, err := stmt.Exec(session.Host, session.Title, session.Password)
	if err != nil {
		return err // обработка ошибок вставки
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return err // обработка ошибок получения ID вставленной записи
	}

	// Преобразуем insertedID в строку
	insertedIDString := fmt.Sprintf("%d", insertedID)

	// Передаем строку в функцию CreateSocket
	url, err := CreateSocket(session, c, insertedIDString)
	if err != nil {
		return err // обработка ошибок создания сокета
	}

	fmt.Printf(url)
	return c.JSON(http.StatusOK, map[string]string{"socket": url})
}

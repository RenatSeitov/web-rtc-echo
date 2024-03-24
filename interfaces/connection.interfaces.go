package interfaces

import (
    "github.com/gorilla/websocket"
    "sync"
)

// Connection - WebSocket соединение
type Connection struct {
    Socket *websocket.Conn
    mu     sync.Mutex
}

// Send - отправка сообщения с обработкой конкурентных операций
func (c *Connection) Send(message Message) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.Socket.WriteJSON(message)
}

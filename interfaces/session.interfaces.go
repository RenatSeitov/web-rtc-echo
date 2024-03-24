package interfaces

// Session - структура сессии
type Session struct {
    Host     string `json:"host"`
    Title    string `json:"title"`
    Password string `json:"password"`
}

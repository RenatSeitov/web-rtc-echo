package interfaces

// Socket - структура сокета
type Socket struct {
    SessionID string `json:"sessionID"`
    HashedURL string `json:"hashedURL"`
    SocketURL string `json:"socketURL"`
}

// Message - структура сообщения
type Message struct {
    Type        string `json:"type"`
    UserID      string `json:"userID"`
    Description string `json:"description"`
    Candidate   string `json:"candidate"`
    To          string `json:"to"`
}
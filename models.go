package main

import "time"

// Message представляет сообщение в чате
type Message struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

// CreateMessageRequest представляет запрос на создание сообщения
type CreateMessageRequest struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

package main

import "time"

type Message struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

type CreateMessageRequest struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

package main

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Storage struct {
	mu        sync.RWMutex
	messages  []Message
	nextID    int
	clients   map[*websocket.Conn]bool
	broadcast chan Message
}

func NewStorage() *Storage {
	storage := &Storage{
		messages:  make([]Message, 0),
		nextID:    1,
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan Message, 256),
	}
	go storage.handleBroadcast()
	return storage
}

func (s *Storage) GetAllMessages() []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Message, len(s.messages))
	copy(result, s.messages)
	return result
}

func (s *Storage) GetMessageByID(id int) *Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.messages {
		if s.messages[i].ID == id {
			return &s.messages[i]
		}
	}
	return nil
}

func (s *Storage) AddMessage(username, text string) Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	message := Message{
		ID:        s.nextID,
		Username:  username,
		Text:      text,
		Timestamp: time.Now(),
	}

	s.messages = append(s.messages, message)
	s.nextID++

	select {
	case s.broadcast <- message:
	default:
	}

	return message
}

func (s *Storage) DeleteMessage(id int) bool {
	s.mu.Lock()
	var deletedMessage *Message
	for i := range s.messages {
		if s.messages[i].ID == id {
			deletedMessage = &s.messages[i]
			s.messages = append(s.messages[:i], s.messages[i+1:]...)
			break
		}
	}
	s.mu.Unlock()

	if deletedMessage != nil {
		deleteNotification := Message{
			ID:        -id, // Отрицательный ID означает удаление
			Username:  deletedMessage.Username,
			Text:      deletedMessage.Text,
			Timestamp: deletedMessage.Timestamp,
		}
		select {
		case s.broadcast <- deleteNotification:
		default:
		}
		return true
	}
	return false
}

func (s *Storage) RegisterClient(conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[conn] = true
}

func (s *Storage) UnregisterClient(conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, conn)
	conn.Close()
}

func (s *Storage) handleBroadcast() {
	for message := range s.broadcast {
		messageJSON, err := json.Marshal(message)
		if err != nil {
			continue
		}

		s.mu.RLock()
		clients := make([]*websocket.Conn, 0, len(s.clients))
		for conn := range s.clients {
			clients = append(clients, conn)
		}
		s.mu.RUnlock()

		for _, conn := range clients {
			if err := conn.WriteMessage(websocket.TextMessage, messageJSON); err != nil {
				s.UnregisterClient(conn)
			}
		}
	}
}

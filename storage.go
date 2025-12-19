package main

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Storage хранит сообщения в памяти
type Storage struct {
	mu        sync.RWMutex
	messages  []Message
	nextID    int
	clients   map[*websocket.Conn]bool
	broadcast chan Message
}

// NewStorage создает новое хранилище
func NewStorage() *Storage {
	storage := &Storage{
		messages:  make([]Message, 0),
		nextID:    1,
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan Message, 256),
	}
	// Запускаем горутину для рассылки сообщений
	go storage.handleBroadcast()
	return storage
}

// GetAllMessages возвращает все сообщения
func (s *Storage) GetAllMessages() []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Возвращаем копию слайса
	result := make([]Message, len(s.messages))
	copy(result, s.messages)
	return result
}

// GetMessageByID возвращает сообщение по ID
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

// AddMessage добавляет новое сообщение
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
	
	// Отправляем сообщение в канал для рассылки через WebSocket
	select {
	case s.broadcast <- message:
	default:
		// Если канал переполнен, просто пропускаем
	}
	
	return message
}

// DeleteMessage удаляет сообщение по ID
func (s *Storage) DeleteMessage(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	for i := range s.messages {
		if s.messages[i].ID == id {
			s.messages = append(s.messages[:i], s.messages[i+1:]...)
			return true
		}
	}
	return false
}

// RegisterClient регистрирует нового WebSocket клиента
func (s *Storage) RegisterClient(conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[conn] = true
}

// UnregisterClient удаляет WebSocket клиента
func (s *Storage) UnregisterClient(conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, conn)
	conn.Close()
}

// handleBroadcast рассылает сообщения всем подключенным клиентам
func (s *Storage) handleBroadcast() {
	for message := range s.broadcast {
		s.mu.RLock()
		clients := make([]*websocket.Conn, 0, len(s.clients))
		for conn := range s.clients {
			clients = append(clients, conn)
		}
		s.mu.RUnlock()
		
		// Сериализуем сообщение в JSON
		messageJSON, err := json.Marshal(message)
		if err != nil {
			continue
		}
		
		// Рассылаем всем клиентам
		for _, conn := range clients {
			s.mu.RLock()
			_, exists := s.clients[conn]
			s.mu.RUnlock()
			
			if !exists {
				continue
			}
			
			if err := conn.WriteMessage(websocket.TextMessage, messageJSON); err != nil {
				s.UnregisterClient(conn)
			}
		}
	}
}

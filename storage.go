package main

import (
	"sync"
	"time"
)

// Storage хранит сообщения в памяти
type Storage struct {
	mu        sync.RWMutex
	messages  []Message
	nextID    int
}

// NewStorage создает новое хранилище
func NewStorage() *Storage {
	return &Storage{
		messages: make([]Message, 0),
		nextID:   1,
	}
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

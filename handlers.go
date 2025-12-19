package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Handlers содержит обработчики HTTP запросов
type Handlers struct {
	storage *Storage
}

// NewHandlers создает новый экземпляр обработчиков
func NewHandlers(storage *Storage) *Handlers {
	return &Handlers{
		storage: storage,
	}
}

// GetAllMessages обрабатывает GET /api/messages
func (h *Handlers) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	messages := h.storage.GetAllMessages()
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(messages)
}

// GetMessageByID обрабатывает GET /api/messages/:id
func (h *Handlers) GetMessageByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из пути
	path := strings.TrimPrefix(r.URL.Path, "/api/messages/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	message := h.storage.GetMessageByID(id)
	if message == nil {
		http.Error(w, "Сообщение не найдено", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(message)
}

// CreateMessage обрабатывает POST /api/messages
func (h *Handlers) CreateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var req CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Text == "" {
		http.Error(w, "Имя пользователя и текст сообщения обязательны", http.StatusBadRequest)
		return
	}

	message := h.storage.AddMessage(req.Username, req.Text)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

// DeleteMessage обрабатывает DELETE /api/messages/:id
func (h *Handlers) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из пути
	path := strings.TrimPrefix(r.URL.Path, "/api/messages/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if !h.storage.DeleteMessage(id) {
		http.Error(w, "Сообщение не найдено", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]string{"message": "Сообщение удалено"})
}

// HandleOptions обрабатывает OPTIONS запросы для CORS
func (h *Handlers) HandleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}

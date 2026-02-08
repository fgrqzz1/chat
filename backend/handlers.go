package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handlers struct {
	storage *Storage
}

func NewHandlers(storage *Storage) *Handlers {
	return &Handlers{
		storage: storage,
	}
}

func (h *Handlers) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	messages := h.storage.GetAllMessages()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}

func (h *Handlers) GetMessageByID(w http.ResponseWriter, r *http.Request) {
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
	if err := json.NewEncoder(w).Encode(message); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}

func (h *Handlers) CreateMessage(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var req CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Text == "" {
		http.Error(w, "Имя пользователя и текст сообщения обязательны", http.StatusBadRequest)
		return
	}

	if len(req.Username) > 50 {
		http.Error(w, "Имя пользователя слишком длинное (максимум 50 символов)", http.StatusBadRequest)
		return
	}

	if len(req.Text) > 1000 {
		http.Error(w, "Текст сообщения слишком длинный (максимум 1000 символов)", http.StatusBadRequest)
		return
	}

	message := h.storage.AddMessage(req.Username, req.Text)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}

func (h *Handlers) DeleteMessage(w http.ResponseWriter, r *http.Request) {
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
	response := map[string]string{"message": "Сообщение удалено"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}

func (h *Handlers) HandleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка обновления до WebSocket: %v", err)
		return
	}
	defer conn.Close()

	h.storage.RegisterClient(conn)
	defer h.storage.UnregisterClient(conn)

	messages := h.storage.GetAllMessages()
	for _, msg := range messages {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
			return
		}
	}

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Ошибка WebSocket: %v", err)
			}
			break
		}
	}
}

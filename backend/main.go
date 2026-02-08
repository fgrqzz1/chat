package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	storage := NewStorage()

	handlers := NewHandlers(storage)

	http.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/messages" {
			switch r.Method {
			case http.MethodGet:
				handlers.GetAllMessages(w, r)
			case http.MethodPost:
				handlers.CreateMessage(w, r)
			case http.MethodOptions:
				handlers.HandleOptions(w, r)
			default:
				http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			}
			return
		}

		if strings.HasPrefix(r.URL.Path, "/api/messages/") {
			switch r.Method {
			case http.MethodGet:
				handlers.GetMessageByID(w, r)
			case http.MethodDelete:
				handlers.DeleteMessage(w, r)
			case http.MethodOptions:
				handlers.HandleOptions(w, r)
			default:
				http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			}
			return
		}

		http.NotFound(w, r)
	})

	http.HandleFunc("/ws", handlers.HandleWebSocket)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	fmt.Printf("Сервер запущен на порту %s\n", port)
	fmt.Printf("API доступно по адресу http://localhost:%s/api/messages\n", port)
	fmt.Printf("WebSocket доступен по адресу ws://localhost:%s/ws\n", port)
	fmt.Printf("Health check: http://localhost:%s/health\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

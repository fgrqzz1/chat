package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Создаем хранилище сообщений
	storage := NewStorage()

	// Создаем обработчики
	handlers := NewHandlers(storage)

	// Настраиваем маршруты
	http.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, является ли путь точно /api/messages
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
		
		// Обработка /api/messages/:id
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

	// Получаем порт из переменной окружения или используем 8080 по умолчанию
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Сервер запущен на порту %s\n", port)
	fmt.Printf("API доступно по адресу http://localhost:%s/api/messages\n", port)
	
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

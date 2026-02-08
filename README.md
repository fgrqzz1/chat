# Chat

Полнофункциональное приложение чата с бэкендом на Go и фронтендом на React.

## Структура проекта

```
chat-1/
├── backend/          # Go-сервер (API, WebSocket)
│   ├── main.go       # Точка входа
│   ├── handlers.go   # HTTP-обработчики
│   ├── models.go     # Модели данных
│   ├── storage.go    # Хранилище сообщений
│   ├── go.mod
│   └── go.sum
├── frontend/         # React + Vite
│   ├── src/
│   │   ├── api/      # API-клиент
│   │   ├── components/
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── index.html
│   ├── package.json
│   └── vite.config.js
├── .github/workflows/ # CI (бэкенд и фронтенд)
├── package.json      # Корневой: скрипты для запуска всего
└── README.md
```

## Старт

### Требования

- **Go** 1.23+
- **Node.js** 16+
- **npm**

### Установка и запуск всего проекта

1. **Зависимости бэкенда:**
   ```bash
   cd backend && go mod download && cd ..
   ```
   Если Go не установлен: `brew install go`

2. **Зависимости фронтенда и корневые (для `dev:all`):**
   ```bash
   npm install
   cd frontend && npm install && cd ..
   ```

3. **Запуск бэкенда и фронтенда одним скриптом:**
   ```bash
   npm run dev:all
   ```
   - Бэкенд: `http://localhost:8080`
   - Фронтенд: `http://localhost:3000`

Откройте в браузере `http://localhost:3000`.

### Запуск по отдельности

- Только бэкенд: `npm run dev:backend` или `cd backend && go run .`
- Только фронтенд: `npm run dev:frontend` или `cd frontend && npm run dev`

### Сборка

- Фронтенд: `npm run build` или `cd frontend && npm run build`
- Бэкенд: `npm run build:backend` или `cd backend && go build -o chat-backend .`

---

## API

### GET /api/messages
Получить все сообщения.

### POST /api/messages
Отправить сообщение. Тело: `{"username": "...", "text": "..."}`.

### GET /api/messages/:id
Получить сообщение по ID.

### DELETE /api/messages/:id
Удалить сообщение по ID.

### WebSocket /ws
Подключение к `ws://localhost:8080/ws` для сообщений в реальном времени. При подключении клиент получает все текущие сообщения; новые рассылаются автоматически.

### Health check
`GET /health` — ответ `{"status":"ok"}`.

---

## Технологии

- **Бэкенд:** Go 1.23, net/http, gorilla/websocket
- **Фронтенд:** React 18, Vite 5, axios

## Примечания

- Сообщения хранятся в памяти и теряются при перезапуске сервера.
- В режиме разработки Vite проксирует `/api` и `/ws` на бэкенд (порт 8080).

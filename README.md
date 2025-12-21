# Chat

## Backend

Простой бэкенд для чата на Go.

### Установка

Убедитесь, что у вас установлен Go версии 1.19 или выше.

```bash
go mod download
```

### Запуск

```bash
go run .
```

Или скомпилируйте и запустите:

```bash
go build -o chat-backend
./chat-backend
```

Сервер запустится на порту 8080 (или на порту, указанном в переменной окружения PORT).

### API Endpoints

#### GET /api/messages
Получить все сообщения.

**Ответ:**
```json
[
  {
    "id": 1,
    "username": "user1",
    "text": "Привет!",
    "timestamp": "2024-01-01T12:00:00.000Z"
  }
]
```

#### POST /api/messages
Отправить новое сообщение.

**Тело запроса:**
```json
{
  "username": "user1",
  "text": "Привет!"
}
```

**Ответ:**
```json
{
  "id": 1,
  "username": "user1",
  "text": "Привет!",
  "timestamp": "2024-01-01T12:00:00.000Z"
}
```

#### GET /api/messages/:id
Получить сообщение по ID.

**Ответ:**
```json
{
  "id": 1,
  "username": "user1",
  "text": "Привет!",
  "timestamp": "2024-01-01T12:00:00.000Z"
}
```

#### DELETE /api/messages/:id
Удалить сообщение по ID.

**Ответ:**
```json
{
  "message": "Сообщение удалено"
}
```

#### WebSocket /ws
Подключиться к WebSocket для получения сообщений в реальном времени.

**Подключение:**
```
ws://localhost:8080/ws
```

**Поведение:**
- При подключении клиент получает все существующие сообщения
- При создании нового сообщения через POST /api/messages, все подключенные клиенты получают его автоматически
- Сообщения отправляются в формате JSON

**Пример сообщения:**
```json
{
  "id": 1,
  "username": "user1",
  "text": "Привет!",
  "timestamp": "2024-01-01T12:00:00.000Z"
}
```

**Пример использования в JavaScript:**
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onmessage = function(event) {
  const message = JSON.parse(event.data);
  console.log('Новое сообщение:', message);
};

ws.onerror = function(error) {
  console.error('WebSocket ошибка:', error);
};

ws.onclose = function() {
  console.log('WebSocket соединение закрыто');
};
```

### Технологии

- Go 1.23+
- Стандартная библиотека net/http
- gorilla/websocket для WebSocket поддержки

### Примечания

- Сообщения хранятся в памяти и будут потеряны при перезапуске сервера
- WebSocket соединения автоматически закрываются при отключении клиента
- Новые сообщения автоматически рассылаются всем подключенным клиентам через WebSocket

### Структура проекта

```
main.go       Точка входа приложения
models.go     Модели данных
storage.go    Хранилище сообщений
handlers.go   HTTP обработчики
go.mod        Зависимости Go
```

## Frontend

Простой фронтенд для чата на React.

### Установка

```bash
npm install
```

### Запуск

```bash
npm run dev
```

Приложение откроется на `http://localhost:3000`

### Сборка

```bash
npm run build
```

### Документация

- [BACKEND_INTEGRATION.md](./BACKEND_INTEGRATION.md) - описание работы с бекендом
- [FRONTEND.md](./FRONTEND.md) - подробный разбор фронтенд части

### Требования

- Node.js 16+
- npm или yarn
- Бекенд должен быть запущен на `http://localhost:8080`

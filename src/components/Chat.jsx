import React, { useState, useEffect } from 'react'
import MessageList from './MessageList'
import MessageInput from './MessageInput'
import { getMessages, sendMessage } from '../api/chatApi'
import './Chat.css'

function Chat() {
  const [messages, setMessages] = useState([])
  const [username, setUsername] = useState('')
  const [loading, setLoading] = useState(true)

  // Загружаем сообщения при монтировании компонента
  useEffect(() => {
    loadMessages()
    
    // Обновляем сообщения каждые 2000 милисекунды
    const interval = setInterval(loadMessages, 2000)
    
    return () => clearInterval(interval)
  }, [])

  // Загружаем сообщения с сервера
  const loadMessages = async () => {
    try {
      const data = await getMessages()
      setMessages(data)
      setLoading(false)
    } catch (error) {
      console.error('Ошибка загрузки сообщений:', error)
      setLoading(false)
    }
  }

  // Отправляем новое сообщение
  const handleSendMessage = async (text) => {
    if (!text.trim() || !username.trim()) return

    try {
      await sendMessage(username, text)
      // Обновляем сообщения после отправки
      await loadMessages()
    } catch (error) {
      console.error('Ошибка отправки сообщения:', error)
    }
  }

  // Устанавливаем имя пользователя
  const handleSetUsername = (name) => {
    setUsername(name)
  }

  if (loading) {
    return <div className="loading">Загрузка...</div>
  }

  return (
    <div className="chat">
      <div className="chat-header">
        <h1>Чат</h1>
        {!username && (
          <div className="username-setup">
            <input
              type="text"
              placeholder="Введите ваше имя"
              onKeyPress={(e) => {
                if (e.key === 'Enter' && e.target.value.trim()) {
                  handleSetUsername(e.target.value.trim())
                }
              }}
              className="username-input"
            />
          </div>
        )}
        {username && (
          <div className="username-display">
            Вы: <strong>{username}</strong>
          </div>
        )}
      </div>
      <MessageList messages={messages} currentUsername={username} />
      {username && <MessageInput onSend={handleSendMessage} />}
    </div>
  )
}

export default Chat


import React from 'react'
import './Message.css'

function Message({ message, isOwn }) {
  const formatTime = (timestamp) => {
    const date = new Date(timestamp)
    return date.toLocaleTimeString('ru-RU', { 
      hour: '2-digit', 
      minute: '2-digit' 
    })
  }

  return (
    <div className={`message ${isOwn ? 'message-own' : ''}`}>
      <div className="message-header">
        <span className="message-username">{message.username}</span>
        <span className="message-time">{formatTime(message.timestamp)}</span>
      </div>
      <div className="message-text">{message.text}</div>
    </div>
  )
}

export default Message


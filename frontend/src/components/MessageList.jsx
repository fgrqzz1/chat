import React from 'react'
import Message from './Message'
import './MessageList.css'

function MessageList({ messages, currentUsername }) {
  return (
    <div className="message-list">
      {messages.length === 0 ? (
        <div className="empty-messages">Нет сообщений</div>
      ) : (
        messages.map((message) => (
          <Message
            key={message.id}
            message={message}
            isOwn={message.username === currentUsername}
          />
        ))
      )}
    </div>
  )
}

export default MessageList

import React, { useState } from 'react'
import './MessageInput.css'

function MessageInput({ onSend }) {
  const [text, setText] = useState('')

  const handleSubmit = (e) => {
    e.preventDefault()
    if (text.trim()) {
      onSend(text)
      setText('')
    }
  }

  return (
    <form className="message-input" onSubmit={handleSubmit}>
      <input
        type="text"
        value={text}
        onChange={(e) => setText(e.target.value)}
        placeholder="Введите сообщение..."
        className="message-input-field"
      />
      <button type="submit" className="message-input-button">
        Отправить
      </button>
    </form>
  )
}

export default MessageInput

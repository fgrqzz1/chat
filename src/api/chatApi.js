import axios from 'axios'

const API_URL = 'http://localhost:8080/api'

// Получить все сообщения
export const getMessages = async () => {
  const response = await axios.get(`${API_URL}/messages`)
  return response.data
}

// Отправить новое сообщение
export const sendMessage = async (username, text) => {
  const response = await axios.post(`${API_URL}/messages`, {
    username,
    text
  })
  return response.data
}


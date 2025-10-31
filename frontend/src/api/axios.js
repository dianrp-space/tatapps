import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor - add token to requests
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor - handle errors
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    const status = error.response?.status

    if (status === 401) {
      localStorage.removeItem('token')
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    } else if (status === 403) {
      if (window.location.pathname !== '/error/forbidden') {
        window.location.href = '/error/forbidden'
      }
    } else if (status && status >= 500) {
      if (window.location.pathname !== '/error/server') {
        window.location.href = '/error/server'
      }
    }
    return Promise.reject(error)
  }
)

export default api

import { createApp } from 'vue'
import axios from 'axios'
import 'bootstrap/dist/css/bootstrap.min.css'
import 'bootstrap/dist/js/bootstrap.bundle.min.js'
import App from './App.vue'
import router from './router'

// 配置全局 axios
axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL || ''
axios.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response && error.response.status === 401) {
      const errorMsg = error.response.data.error
      if (errorMsg === 'token_expired' || errorMsg === 'Invalid or expired token') {
        alert('登录已过期，请重新登录')
        localStorage.clear()
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

createApp(App).use(router).mount('#app')

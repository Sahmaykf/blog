<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()
const username = ref('')
const userId = ref('')
const avatar = ref('')
const isLoggedIn = ref(false)

const fetchUserProfile = async () => {
  const token = localStorage.getItem('token')
  const storedUserId = localStorage.getItem('user_id')
  
  if (token && storedUserId) {
    isLoggedIn.value = true
    username.value = localStorage.getItem('username')
    userId.value = storedUserId
    
    try {
      const res = await axios.get(`/api/v1/users/${storedUserId}`, {
        headers: { Authorization: `Bearer ${token}` }
      })
      if (res.data.code === 200) {
        avatar.value = res.data.data.avatar
        username.value = res.data.data.username
        localStorage.setItem('username', username.value)
      }
    } catch (e) {
      console.error('Failed to fetch user profile in App.vue', e)
    }
  } else {
    isLoggedIn.value = false
    avatar.value = ''
  }
}

onMounted(() => {
  fetchUserProfile()
  // Listen for profile updates from UserSettings.vue
  window.addEventListener('profile-updated', fetchUserProfile)
  // Listen for login events if they happen without reload
  window.addEventListener('login-success', fetchUserProfile)
})

const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('username')
  localStorage.removeItem('user_id')
  localStorage.removeItem('user_role')
  isLoggedIn.value = false
  username.value = ''
  userId.value = ''
  router.push('/login')
}
</script>

<template>
  <div class="app-container">
    <!-- 导航栏 -->
    <nav class="navbar navbar-expand-lg navbar-light bg-white shadow-sm sticky-top">
      <div class="container">
        <router-link class="navbar-brand fw-bold text-primary" to="/">
          <i class="bi bi-journal-text me-2"></i>Simple Blog
        </router-link>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarNav">
          <ul class="navbar-nav me-auto mb-2 mb-lg-0">
            <li class="nav-item">
              <router-link class="nav-link" to="/">首页</router-link>
            </li>
            <li class="nav-item" v-if="isLoggedIn">
              <router-link class="nav-link" :to="`/user/${userId}`">我的主页</router-link>
            </li>
            <li class="nav-item" v-if="isLoggedIn">
              <router-link class="nav-link" to="/admin">后台管理</router-link>
            </li>
          </ul>
          <div class="d-flex align-items-center">
            <template v-if="isLoggedIn">
              <router-link to="/settings" class="navbar-text me-3 text-decoration-none hover-primary d-flex align-items-center py-0">
                <div v-if="avatar" class="avatar-circle-sm me-2 border shadow-sm">
                  <img :src="avatar" class="w-100 h-100 object-fit-cover">
                </div>
                <i v-else class="bi bi-person-circle me-2 fs-4 text-primary"></i>
                <span class="d-none d-sm-inline">你好, <span class="fw-bold text-dark">{{ username }}</span></span>
              </router-link>
              <button class="btn btn-outline-danger btn-sm rounded-pill px-3" @click="logout">退出</button>
            </template>
            <template v-else>
              <router-link to="/login" class="btn btn-primary btn-sm rounded-pill px-4">登录 / 注册</router-link>
            </template>
          </div>
        </div>
      </div>
    </nav>

    <main class="container py-5">
      <router-view></router-view>
    </main>
    
    <footer class="text-center text-muted py-4 mt-auto border-top bg-white">
      <small>&copy; 2025 Simple Blog. All rights reserved.</small>
    </footer>
  </div>
</template>

<style>
body {
  background-color: #f8f9fa;
  font-family: 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  color: #333;
}

.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.navbar-brand {
  font-size: 1.5rem;
  letter-spacing: -0.5px;
}

.nav-link {
  font-weight: 500;
  color: #555 !important;
}

.nav-link.router-link-active {
  color: #0d6efd !important;
}

.avatar-circle-sm {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f8f9fa;
}

.hover-primary:hover {
  color: #0d6efd !important;
}
</style>

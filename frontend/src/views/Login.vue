<script setup>
import { ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()
const authMode = ref('login') // 'login' or 'register'
const authForm = ref({
  username: '',
  password: '',
  email: ''
})
const errorMsg = ref('')
const successMsg = ref('')

const handleAuth = async () => {
  errorMsg.value = ''
  successMsg.value = ''
  const url = authMode.value === 'login' ? '/api/v1/login' : '/api/v1/register'
  
  try {
    const response = await axios.post(url, authForm.value)
    const res = response.data
    
    if (res.code !== 200) {
      errorMsg.value = res.msg || '操作失败'
      return
    }

    if (authMode.value === 'login') {
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('username', authForm.value.username)
      localStorage.setItem('user_id', res.data.user_id)
      localStorage.setItem('user_role', res.data.role)
      
      // 登录成功跳转到首页或管理页
      // 刷新页面以更新 App.vue 的状态（或者使用 Vuex/Pinia，但这里简单处理）
      window.location.href = '/' 
    } else {
      successMsg.value = '注册成功，请登录'
      authMode.value = 'login'
      authForm.value.password = ''
    }
  } catch (error) {
    errorMsg.value = error.response?.data?.msg || '操作失败'
  }
}
</script>

<template>
  <div class="auth-container">
    <div class="card shadow-sm border-0 rounded-4">
      <div class="card-header bg-white border-0 pt-4 pb-0 text-center">
        <h4 class="fw-bold text-dark mb-4" style="letter-spacing: -0.5px;">Simple Blog</h4>
        <div class="d-flex justify-content-center mb-2">
          <div class="bg-light p-1 rounded-pill d-inline-flex">
            <button class="btn btn-sm rounded-pill px-4 fw-medium transition-all" 
              :class="authMode === 'login' ? 'btn-white shadow-sm text-dark' : 'text-muted border-0'" 
              @click="authMode = 'login'">登录</button>
            <button class="btn btn-sm rounded-pill px-4 fw-medium transition-all" 
              :class="authMode === 'register' ? 'btn-white shadow-sm text-dark' : 'text-muted border-0'" 
              @click="authMode = 'register'">注册</button>
          </div>
        </div>
      </div>
      <div class="card-body p-4 pt-2">
        <form @submit.prevent="handleAuth" class="mt-3">
          <div class="mb-3">
            <label class="form-label text-secondary small fw-medium mb-1">用户名</label>
            <input type="text" class="form-control bg-light border-0 py-2 px-3" v-model="authForm.username" required placeholder="请输入用户名">
          </div>
          <div class="mb-3" v-if="authMode === 'register'">
            <label class="form-label text-secondary small fw-medium mb-1">邮箱</label>
            <input type="email" class="form-control bg-light border-0 py-2 px-3" v-model="authForm.email" required placeholder="name@example.com">
          </div>
          <div class="mb-4">
            <label class="form-label text-secondary small fw-medium mb-1">密码</label>
            <input type="password" class="form-control bg-light border-0 py-2 px-3" v-model="authForm.password" required placeholder="请输入密码">
          </div>
          <div class="d-grid pt-2">
            <button type="submit" class="btn btn-dark py-2 rounded-3 fw-medium shadow-sm hover-lift">{{ authMode === 'login' ? '立即登录' : '立即注册' }}</button>
          </div>
        </form>
        <div v-if="errorMsg" class="alert alert-danger mt-3 border-0 bg-danger bg-opacity-10 text-danger small py-2">{{ errorMsg }}</div>
        <div v-if="successMsg" class="alert alert-success mt-3 border-0 bg-success bg-opacity-10 text-success small py-2">{{ successMsg }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.auth-container { 
  max-width: 380px; 
  margin: 60px auto; 
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
}
.btn-white {
  background-color: #fff;
}
.transition-all {
  transition: all 0.2s ease;
}
.form-control:focus {
  box-shadow: 0 0 0 2px rgba(33, 37, 41, 0.1);
  background-color: #fff;
}
.hover-lift {
  transition: transform 0.2s;
}
.hover-lift:hover {
  transform: translateY(-1px);
}
</style>

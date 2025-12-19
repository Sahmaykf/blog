<template>
  <div class="container py-5" style="max-width: 600px;">
    <div class="card shadow-sm border-0">
      <div class="card-body p-4">
        <h3 class="fw-bold mb-4">个人设置</h3>
        
        <!-- Avatar Section -->
        <div class="text-center mb-5">
          <div class="position-relative d-inline-block">
            <div v-if="user?.avatar" class="avatar-preview rounded-circle shadow-sm" 
              :style="{ backgroundImage: `url(${user.avatar})` }">
            </div>
            <div v-else class="avatar-placeholder rounded-circle bg-primary text-white d-flex align-items-center justify-content-center fs-1 fw-bold shadow-sm mx-auto"
              style="width: 120px; height: 120px;">
              {{ user?.username?.charAt(0).toUpperCase() }}
            </div>
            <label for="avatarInput" class="btn btn-sm btn-light rounded-circle position-absolute bottom-0 end-0 shadow-sm border d-flex align-items-center justify-content-center" style="width: 36px; height: 36px;">
              <span v-if="uploading" class="spinner-border spinner-border-sm"></span>
              <i v-else class="bi bi-camera-fill fs-5"></i>
              <input type="file" id="avatarInput" class="d-none" @change="handleAvatarUpload" accept="image/*" :disabled="uploading">
            </label>
          </div>
          <p class="text-muted small mt-2">点击相机图标更换头像</p>
        </div>

        <!-- Profile Form -->
        <form @submit.prevent="updateProfile">
          <div class="mb-3">
            <label class="form-label text-muted small fw-bold">用户名</label>
            <input type="text" class="form-control bg-light" :value="user?.username" disabled>
            <div class="form-text">用户名暂不支持修改</div>
          </div>
          
          <div class="mb-4">
            <label class="form-label text-muted small fw-bold">电子邮箱</label>
            <input type="email" v-model="email" class="form-control" required placeholder="请输入邮箱">
          </div>

          <div class="d-grid">
            <button type="submit" class="btn btn-primary rounded-pill py-2 fw-bold" :disabled="saving">
              <span v-if="saving" class="spinner-border spinner-border-sm me-2"></span>
              保存修改
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const user = ref(null)
const email = ref('')
const tempAvatar = ref('')
const saving = ref(false)
const uploading = ref(false)

const fetchUser = async () => {
  const userId = localStorage.getItem('user_id')
  const token = localStorage.getItem('token')
  try {
    const res = await axios.get(`/api/v1/users/${userId}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (res.data.code === 200) {
      user.value = res.data.data
      email.value = res.data.data.email
      tempAvatar.value = res.data.data.avatar
    }
  } catch (err) {
    console.error(err)
  }
}

const updateProfile = async () => {
  saving.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await axios.put('/api/v1/user/profile', {
      email: email.value,
      avatar: tempAvatar.value
    }, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (res.data.code === 200) {
      alert('修改成功')
      // Emit event to notify App.vue to refresh avatar
      window.dispatchEvent(new CustomEvent('profile-updated'))
      fetchUser()
    }
  } catch (err) {
    alert(err.response?.data?.msg || '修改失败')
  } finally {
    saving.value = false
  }
}

const handleAvatarUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  // Basic validation
  if (file.size > 2 * 1024 * 1024) {
    alert('图片大小不能超过 2MB')
    return
  }

  const formData = new FormData()
  formData.append('avatar', file)

  uploading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await axios.post('/api/v1/user/avatar', formData, {
      headers: { 
        'Content-Type': 'multipart/form-data',
        Authorization: `Bearer ${token}` 
      }
    })
    if (res.data.code === 200) {
      // Only update the preview and temp path, don't save to DB yet
      tempAvatar.value = res.data.data.avatar_url
      user.value.avatar = res.data.data.avatar_url
      alert('头像已上传，请点击下方保存按钮以生效')
    }
  } catch (err) {
    alert('头像上传失败: ' + (err.response?.data?.msg || err.message))
  } finally {
    uploading.value = false
  }
}

onMounted(fetchUser)
</script>

<style scoped>
.avatar-preview {
  width: 120px;
  height: 120px;
  background-size: cover;
  background-position: center;
  margin: 0 auto;
}
</style>

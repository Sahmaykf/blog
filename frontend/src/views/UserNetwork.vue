<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const users = ref([])
const loading = ref(true)
const title = ref('')

const fetchUsers = async () => {
  loading.value = true
  const userId = route.params.id
  const type = route.params.type // 'followers' or 'following'
  
  title.value = type === 'followers' ? '粉丝列表' : '关注列表'
  
  try {
    const endpoint = type === 'followers' 
      ? `/api/v1/users/${userId}/followers` 
      : `/api/v1/users/${userId}/following`
      
    const res = await axios.get(endpoint)
    if (res.data.code === 200) {
      users.value = res.data.data || []
    }
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(fetchUsers)
</script>

<template>
  <div class="container py-5" style="max-width: 800px;">
    <div class="d-flex align-items-center mb-4">
      <button @click="router.go(-1)" class="btn btn-outline-secondary btn-sm me-3">
        <i class="bi bi-arrow-left"></i> 返回
      </button>
      <h3 class="fw-bold mb-0">{{ title }}</h3>
    </div>

    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary" role="status"></div>
    </div>

    <div v-else-if="users.length > 0" class="list-group shadow-sm">
      <router-link 
        v-for="user in users" 
        :key="user.ID" 
        :to="`/user/${user.ID}`"
        class="list-group-item list-group-item-action d-flex align-items-center p-3"
      >
        <div class="bg-primary text-white rounded-circle d-flex align-items-center justify-content-center me-3" style="width: 48px; height: 48px;">
          <span class="fs-5 fw-bold">{{ user.username ? user.username.charAt(0).toUpperCase() : '?' }}</span>
        </div>
        <div>
          <h5 class="mb-0 text-dark">{{ user.username }}</h5>
        </div>
        <i class="bi bi-chevron-right ms-auto text-muted"></i>
      </router-link>
    </div>

    <div v-else class="text-center py-5 bg-white rounded shadow-sm">
      <p class="text-muted mb-0">暂无数据</p>
    </div>
  </div>
</template>

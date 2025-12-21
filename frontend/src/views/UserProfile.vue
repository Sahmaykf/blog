<template>
  <div class="container mt-4">
    <div v-if="user" class="row justify-content-center">
      <div class="col-md-8">
        <!-- Profile Header -->
        <div class="card shadow-sm mb-4 border-0">
          <div class="card-body text-center py-5">
            <div class="mb-3">
              <div v-if="user.avatar" class="mx-auto shadow-sm rounded-circle overflow-hidden" style="width: 100px; height: 100px;">
                <img :src="user.avatar" class="w-100 h-100 object-fit-cover" :alt="user.username">
              </div>
              <div v-else class="avatar-placeholder rounded-circle bg-primary text-white d-flex align-items-center justify-content-center mx-auto fs-1 fw-bold shadow-sm"
                style="width: 100px; height: 100px;">
                {{ user.username.charAt(0).toUpperCase() }}
              </div>
            </div>
            <h2 class="card-title fw-bold mb-1">{{ user.username }}</h2>
            <p v-if="user.bio" class="text-secondary mb-3 px-4 mx-auto" style="max-width: 500px;">{{ user.bio }}</p>
            <p class="text-muted mb-4 small">加入时间: {{ formatDate(user.created_at) }}</p>
            
            <!-- Stats Row -->
            <div class="d-flex justify-content-center gap-5 mb-4">
              <div class="text-center cursor-pointer p-2 rounded hover-bg" @click="goToNetwork('following')">
                <div class="h4 fw-bold mb-0 text-primary">{{ followingCount }}</div>
                <div class="small text-muted">关注</div>
              </div>
              <div class="text-center cursor-pointer p-2 rounded hover-bg" @click="goToNetwork('followers')">
                <div class="h4 fw-bold mb-0 text-primary">{{ followersCount }}</div>
                <div class="small text-muted">粉丝</div>
              </div>
              <div class="text-center p-2">
                <div class="h4 fw-bold mb-0 text-dark">{{ posts.length }}</div>
                <div class="small text-muted">文章</div>
              </div>
            </div>

            <!-- Action Buttons -->
            <div v-if="isOwnProfile">
              <router-link to="/settings" class="btn btn-outline-secondary rounded-pill px-4">
                <i class="bi bi-pencil me-1"></i> 编辑资料
              </router-link>
            </div>
            <div v-else>
              <button 
                @click="toggleFollow" 
                class="btn rounded-pill px-4 transition-all"
                :class="isFollowing ? 'btn-secondary' : 'btn-primary'"
                :disabled="loading"
              >
                <span v-if="loading" class="spinner-border spinner-border-sm me-1"></span>
                <i class="bi" :class="isFollowing ? 'bi-person-check' : 'bi-person-plus'"></i>
                {{ isFollowing ? '已关注' : '关注' }}
              </button>
            </div>
          </div>
        </div>

        <!-- Content Tabs -->
        <div class="card shadow-sm border-0">
          <div class="card-header bg-white border-bottom-0 pt-3 px-0 mx-3">
            <ul class="nav nav-tabs card-header-tabs">
              <li class="nav-item">
                <a class="nav-link" :class="{ active: activeTab === 'posts' }" @click.prevent="activeTab = 'posts'" href="#">
                  <i class="bi bi-file-text me-1"></i> 我的文章
                </a>
              </li>
              <li class="nav-item">
                <a class="nav-link" :class="{ active: activeTab === 'likes' }" @click.prevent="activeTab = 'likes'" href="#">
                  <i class="bi bi-heart me-1"></i> 赞过的
                </a>
              </li>
              <li class="nav-item">
                <a class="nav-link" :class="{ active: activeTab === 'favorites' }" @click.prevent="activeTab = 'favorites'" href="#">
                  <i class="bi bi-star me-1"></i> 收藏的
                </a>
              </li>
            </ul>
          </div>
          
          <div class="card-body p-0">
            <!-- My Posts List -->
            <div v-if="activeTab === 'posts'">
              <div v-if="posts.length === 0" class="text-center py-5 text-muted">
                <i class="bi bi-journal-x fs-1 d-block mb-2"></i>
                暂无文章
              </div>
              <div v-else class="list-group list-group-flush">
                <router-link 
                  v-for="post in posts" 
                  :key="post.ID"
                  :to="'/post/' + post.ID"
                  class="list-group-item list-group-item-action p-4 border-bottom position-relative"
                >
                  <!-- 置顶标识 -->
                  <div v-if="post.is_top" class="position-absolute top-0 start-0 bg-primary text-white px-3 py-1 rounded-bottom-end small shadow-sm" style="z-index: 1; border-top-left-radius: 0;">
                    <i class="bi bi-pin-angle-fill me-1"></i>置顶
                  </div>
                  
                  <div class="d-flex justify-content-between align-items-start mb-2" :class="{'mt-3': post.is_top}">
                    <h5 class="mb-1 fw-bold text-break">{{ post.title }}</h5>
                    <small class="text-muted text-nowrap ms-2">{{ formatDate(post.CreatedAt) }}</small>
                  </div>
                  <p class="mb-1 text-muted">{{ truncateContent(post.content) }}</p>
                  <div class="mt-2 d-flex gap-3 text-muted small">
                    <span><i class="bi bi-chat-dots me-1"></i>{{ post.comment_count || 0 }} 评论</span>
                    <span><i class="bi bi-heart me-1"></i>{{ post.like_count || 0 }} 赞</span>
                  </div>
                </router-link>
              </div>
            </div>

            <!-- Liked Posts List -->
            <div v-if="activeTab === 'likes'">
              <div v-if="likedPosts.length === 0" class="text-center py-5 text-muted">
                <i class="bi bi-heart-break fs-1 d-block mb-2"></i>
                暂无赞过的文章
              </div>
              <div v-else class="list-group list-group-flush">
                <router-link 
                  v-for="post in likedPosts" 
                  :key="post.ID"
                  :to="'/post/' + post.ID"
                  class="list-group-item list-group-item-action p-4 border-bottom"
                >
                  <div class="d-flex justify-content-between align-items-start mb-2">
                    <h5 class="mb-1 fw-bold text-break">{{ post.title }}</h5>
                    <small class="text-muted text-nowrap ms-2">{{ formatDate(post.CreatedAt) }}</small>
                  </div>
                  <div class="d-flex align-items-center mt-2">
                    <div class="d-flex align-items-center me-3">
                      <div v-if="post.user?.avatar" class="rounded-circle overflow-hidden me-2" style="width: 24px; height: 24px;">
                        <img :src="post.user.avatar" class="w-100 h-100 object-fit-cover">
                      </div>
                      <div v-else class="rounded-circle bg-secondary text-white d-flex align-items-center justify-content-center me-2" style="width: 24px; height: 24px; font-size: 12px;">
                        {{ post.user?.username?.charAt(0).toUpperCase() || '?' }}
                      </div>
                      <small class="text-muted">{{ post.user?.username || 'Unknown' }}</small>
                    </div>
                  </div>
                </router-link>
              </div>
            </div>

            <!-- Favorite Posts List -->
            <div v-if="activeTab === 'favorites'">
              <div v-if="favoritePosts.length === 0" class="text-center py-5 text-muted">
                <i class="bi bi-star fs-1 d-block mb-2"></i>
                暂无收藏的文章
              </div>
              <div v-else class="list-group list-group-flush">
                <router-link 
                  v-for="post in favoritePosts" 
                  :key="post.ID"
                  :to="'/post/' + post.ID"
                  class="list-group-item list-group-item-action p-4 border-bottom"
                >
                  <div class="d-flex justify-content-between align-items-start mb-2">
                    <h5 class="mb-1 fw-bold text-break">{{ post.title }}</h5>
                    <small class="text-muted text-nowrap ms-2">{{ formatDate(post.CreatedAt) }}</small>
                  </div>
                  <div class="d-flex align-items-center mt-2">
                    <div class="d-flex align-items-center me-3">
                      <div v-if="post.user?.avatar" class="rounded-circle overflow-hidden me-2" style="width: 24px; height: 24px;">
                        <img :src="post.user.avatar" class="w-100 h-100 object-fit-cover">
                      </div>
                      <div v-else class="rounded-circle bg-secondary text-white d-flex align-items-center justify-content-center me-2" style="width: 24px; height: 24px; font-size: 12px;">
                        {{ post.user?.username?.charAt(0).toUpperCase() || '?' }}
                      </div>
                      <small class="text-muted">{{ post.user?.username || 'Unknown' }}</small>
                    </div>
                  </div>
                </router-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-else class="text-center mt-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const user = ref(null)
const posts = ref([])
const likedPosts = ref([])
const favoritePosts = ref([])
const isFollowing = ref(false)
const followersCount = ref(0)
const followingCount = ref(0)
const loading = ref(false)
const activeTab = ref('posts')

const currentUserId = computed(() => {
  const token = localStorage.getItem('token')
  if (!token) return null
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.user_id
  } catch {
    return null
  }
})

const isOwnProfile = computed(() => {
  return user.value && currentUserId.value === user.value.id
})

const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const truncateContent = (content) => {
  if (!content) return ''
  // 移除 Markdown 图片: ![alt](url)
  let text = content.replace(/!\[.*?\]\(.*?\)/g, '')
  // 移除 Markdown 链接: [text](url) -> text
  text = text.replace(/\[(.*?)\]\(.*?\)/g, '$1')
  // 移除标题符号 #
  text = text.replace(/#+\s/g, '')
  // 移除粗体/斜体
  text = text.replace(/[*_]{1,3}/g, '')
  
  return text.length > 100 ? text.substring(0, 100) + '...' : text
}

const goToNetwork = (type) => {
  router.push({
    name: 'UserNetwork',
    params: { id: user.value.id, type }
  })
}

const fetchUserData = async () => {
  try {
    const userId = route.params.id
    const token = localStorage.getItem('token')
    const headers = token ? { Authorization: `Bearer ${token}` } : {}

    // Fetch basic user info
    const userRes = await axios.get(`/api/v1/users/${userId}`, { headers })
    if (userRes.data.code === 200) {
      user.value = userRes.data.data
      isFollowing.value = userRes.data.data.is_following
      followersCount.value = userRes.data.data.follower_count
      followingCount.value = userRes.data.data.following_count
    }

    // Fetch user's posts
    const postsRes = await axios.get(`/api/v1/users/${userId}/posts`)
    if (postsRes.data.code === 200) {
      posts.value = postsRes.data.data
    }

    // Fetch liked posts
    const likedRes = await axios.get(`/api/v1/users/${userId}/liked-posts`)
    if (likedRes.data.code === 200) {
      likedPosts.value = likedRes.data.data
    }

    // Fetch favorite posts
    const favoriteRes = await axios.get(`/api/v1/users/${userId}/favorite-posts`)
    if (favoriteRes.data.code === 200) {
      favoritePosts.value = favoriteRes.data.data
    }

  } catch (error) {
    console.error('Failed to fetch user data:', error)
    if (error.response && error.response.status === 404) {
      // Check if we are trying to view our own profile
      if (currentUserId.value && String(route.params.id) === String(currentUserId.value)) {
        alert('用户信息已失效，请重新登录')
        localStorage.clear()
        router.push('/login')
        // Force reload to update App.vue state
        setTimeout(() => window.location.reload(), 100)
        return
      }
      alert('用户不存在')
      router.push('/')
    }
  }
}

const toggleFollow = async () => {
  if (!currentUserId.value) {
    alert('请先登录')
    router.push('/login')
    return
  }

  loading.value = true
  try {
    const action = isFollowing.value ? 'unfollow' : 'follow'
    const res = await axios.post(`/api/v1/users/${user.value.id}/${action}`, {}, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    })
    
    if (res.data.code === 200) {
      isFollowing.value = !isFollowing.value
      followersCount.value += isFollowing.value ? 1 : -1
    }
  } catch (error) {
    alert('操作失败: ' + (error.response?.data?.msg || error.message))
  } finally {
    loading.value = false
  }
}

watch(() => route.params.id, () => {
  if (route.name === 'UserProfile') {
    fetchUserData()
    activeTab.value = 'posts'
  }
})

onMounted(fetchUserData)
</script>

<style scoped>
.avatar-placeholder {
  font-size: 2.5rem;
  transition: transform 0.3s ease;
}
.avatar-placeholder:hover {
  transform: scale(1.05);
}
.hover-bg:hover {
  background-color: #f8f9fa;
  cursor: pointer;
}
.cursor-pointer {
  cursor: pointer;
}
.nav-link {
  color: #6c757d;
  cursor: pointer;
}
.nav-link.active {
  color: #0d6efd;
  font-weight: bold;
}
.transition-all {
  transition: all 0.3s ease;
}
</style>
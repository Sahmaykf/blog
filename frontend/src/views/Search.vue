<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const keyword = ref(route.query.q || '')
const activeTab = ref('posts')
const loading = ref(false)
const results = ref({
  posts: [],
  users: [],
  tags: []
})

const fetchResults = async () => {
  if (!keyword.value) return
  loading.value = true
  try {
    const response = await axios.get(`/api/v1/search?keyword=${encodeURIComponent(keyword.value)}`)
    if (response.data.code === 200) {
      results.value = response.data.data
    }
  } catch (error) {
    console.error('搜索失败', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  if (!keyword.value.trim()) return
  router.push({ path: '/search', query: { q: keyword.value } })
}

const viewPost = (id) => {
  router.push(`/post/${id}`)
}

const viewUser = (id) => {
  router.push(`/user/${id}`)
}

const filterByTag = (tagName) => {
  router.push({ path: '/', query: { tag: tagName } })
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString()
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
  
  return text.length > 150 ? text.substring(0, 150) + '...' : text
}

watch(() => route.query.q, (newQ) => {
  keyword.value = newQ || ''
  fetchResults()
})

onMounted(() => {
  fetchResults()
})
</script>

<template>
  <div class="search-container py-4">
    <div class="row justify-content-center mb-5">
      <div class="col-md-8">
        <div class="input-group input-group-lg shadow-sm rounded-pill overflow-hidden bg-white">
          <span class="input-group-text bg-white border-0 ps-4">
            <i class="bi bi-search text-muted"></i>
          </span>
          <input type="text" 
            class="form-control border-0 py-3 shadow-none" 
            v-model="keyword" 
            placeholder="搜索文章、用户或标签..."
            @keyup.enter="handleSearch">
          <button class="btn btn-primary px-4 rounded-pill m-1" @click="handleSearch">搜索</button>
        </div>
      </div>
    </div>

    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>

    <template v-else-if="keyword">
      <div class="d-flex justify-content-center mb-5">
        <div class="bg-light p-1 rounded-pill d-inline-flex shadow-sm">
          <button class="btn btn-sm rounded-pill px-4 fw-medium transition-all" 
            :class="activeTab === 'posts' ? 'btn-white shadow-sm text-dark' : 'text-muted border-0'" 
            @click="activeTab = 'posts'">
            文章 ({{ results.posts.length }})
          </button>
          <button class="btn btn-sm rounded-pill px-4 fw-medium transition-all" 
            :class="activeTab === 'users' ? 'btn-white shadow-sm text-dark' : 'text-muted border-0'" 
            @click="activeTab = 'users'">
            用户 ({{ results.users.length }})
          </button>
          <button class="btn btn-sm rounded-pill px-4 fw-medium transition-all" 
            :class="activeTab === 'tags' ? 'btn-white shadow-sm text-dark' : 'text-muted border-0'" 
            @click="activeTab = 'tags'">
            标签 ({{ results.tags.length }})
          </button>
        </div>
      </div>

      <div class="tab-content">
        <!-- 文章结果 -->
        <div v-if="activeTab === 'posts'" class="row g-4">
          <div v-if="results.posts.length === 0" class="text-center py-5 text-muted">
            未找到相关文章
          </div>
          <div v-for="post in results.posts" :key="post.ID" class="col-md-8 mx-auto">
            <div class="card border-0 shadow-sm hover-shadow transition-all" @click="viewPost(post.ID)" style="cursor: pointer;">
              <div class="card-body p-4">
                <h5 class="card-title fw-bold mb-2">{{ post.title }}</h5>
                <div class="d-flex align-items-center gap-3 text-muted small mb-3">
                  <span><i class="bi bi-person me-1"></i>{{ post.user?.username }}</span>
                  <span><i class="bi bi-calendar3 me-1"></i>{{ formatDate(post.CreatedAt) }}</span>
                  <span class="ms-auto">
                    <i class="bi bi-eye me-1"></i>{{ post.view_count }}
                    <i class="bi bi-heart ms-3 me-1"></i>{{ post.like_count }}
                  </span>
                </div>
                <p class="card-text text-secondary text-truncate-2">{{ truncateContent(post.content) }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 用户结果 -->
        <div v-if="activeTab === 'users'" class="row g-4">
          <div v-if="results.users.length === 0" class="text-center py-5 text-muted">
            未找到相关用户
          </div>
          <div v-for="user in results.users" :key="user.ID" class="col-md-6 mx-auto">
            <div class="card border-0 shadow-sm hover-shadow transition-all" @click="viewUser(user.ID)" style="cursor: pointer;">
              <div class="card-body d-flex align-items-center p-3">
                <div class="avatar-container me-3">
                  <img v-if="user.avatar" :src="user.avatar" class="rounded-circle shadow-sm" style="width: 60px; height: 60px; object-fit: cover;">
                  <div v-else class="rounded-circle bg-primary d-flex align-items-center justify-content-center text-white fw-bold shadow-sm" style="width: 60px; height: 60px; font-size: 1.5rem;">
                    {{ user.username.charAt(0).toUpperCase() }}
                  </div>
                </div>
                <div>
                  <h6 class="fw-bold mb-1">{{ user.username }}</h6>
                  <p class="text-muted small mb-0">{{ user.bio || '这个用户很懒，什么都没写' }}</p>
                </div>
                <button class="btn btn-outline-primary btn-sm rounded-pill ms-auto px-3">查看主页</button>
              </div>
            </div>
          </div>
        </div>

        <!-- 标签结果 -->
        <div v-if="activeTab === 'tags'" class="row g-4">
          <div v-if="results.tags.length === 0" class="text-center py-5 text-muted">
            未找到相关标签
          </div>
          <div class="col-md-8 mx-auto d-flex flex-wrap gap-3 justify-content-center">
            <button v-for="tag in results.tags" :key="tag.ID" 
              class="btn btn-light border rounded-pill px-4 py-2 shadow-sm hover-primary"
              @click="filterByTag(tag.name)">
              # {{ tag.name }}
            </button>
          </div>
        </div>
      </div>
    </template>
    
    <div v-else class="text-center py-5 text-muted">
      请输入关键词开始搜索
    </div>
  </div>
</template>

<style scoped>
.hover-shadow:hover {
  transform: translateY(-3px);
  box-shadow: 0 .5rem 1rem rgba(0,0,0,.1)!important;
}
.transition-all {
  transition: all 0.3s ease;
}
.btn-white {
  background-color: #fff;
}
.text-truncate-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.hover-primary:hover {
  background-color: #0d6efd;
  color: white;
  border-color: #0d6efd;
}
</style>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

const route = useRoute()
const router = useRouter()
const post = ref(null)
const comments = ref([])
const loading = ref(true)
const error = ref('')
const newComment = ref('')
const submitting = ref(false)
const isLiked = ref(false)
const likeCount = ref(0)
const isFavorited = ref(false)
const favoriteCount = ref(0)
const replyTo = ref(null) // { id: number, username: string, parentId: number }

// User info for permissions
const currentUserId = ref(localStorage.getItem('user_id'))
const currentUserRole = ref(localStorage.getItem('user_role'))
const isLoggedIn = computed(() => !!localStorage.getItem('token'))

const fetchPost = async () => {
  loading.value = true
  try {
    const response = await axios.get(`/api/v1/posts/${route.params.id}`)
    const res = response.data
    if (res.code === 200) {
      post.value = res.data
    } else {
      error.value = res.msg || '获取文章失败'
    }
  } catch (err) {
    error.value = '获取文章失败'
  } finally {
    loading.value = false
  }
}

const fetchComments = async () => {
  try {
    const response = await axios.get(`/api/v1/posts/${route.params.id}/comments`)
    if (response.data.code === 200) {
      comments.value = response.data.data || []
    }
  } catch (err) {
    console.error('获取评论失败', err)
  }
}

const submitComment = async () => {
  if (!newComment.value.trim()) return
  
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const payload = {
      post_id: parseInt(route.params.id),
      content: newComment.value
    }
    if (replyTo.value) {
      payload.parent_id = replyTo.value.parentId
      payload.reply_to_id = replyTo.value.id
    }

    const response = await axios.post('/api/v1/comments', payload, {
      headers: { Authorization: `Bearer ${token}` }
    })
    
    if (response.data.code === 200) {
      newComment.value = ''
      replyTo.value = null
      fetchComments()
      alert('评论成功')
    } else {
      alert(response.data.msg || '评论失败')
    }
  } catch (err) {
    if (err.response && err.response.status === 401) {
      alert('请先登录')
      router.push('/login')
    } else {
      alert('评论失败')
    }
  } finally {
    submitting.value = false
  }
}

const deleteComment = async (commentId) => {
  if (!confirm('确定要删除这条评论吗？')) return
  
  try {
    const token = localStorage.getItem('token')
    const response = await axios.delete(`/api/v1/comments/${commentId}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    
    if (response.data.code === 200) {
      fetchComments()
    } else {
      alert(response.data.msg || '删除失败')
    }
  } catch (err) {
    alert('删除失败')
  }
}

const canDelete = (comment) => {
  if (!isLoggedIn.value) return false
  if (currentUserRole.value === 'admin') return true
  // Allow if current user is the comment author
  if (String(comment.user_id) === String(currentUserId.value)) return true
  // Allow if current user is the post author
  if (post.value && String(post.value.user_id) === String(currentUserId.value)) return true
  return false
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString()
}

const backLink = computed(() => {
  if (route.query.from === 'user_profile' && post.value && post.value.user) {
    return `/user/${post.value.user.ID}`
  }
  return '/'
})

const backText = computed(() => {
  if (route.query.from === 'user_profile') {
    return '返回个人主页'
  }
  return '返回列表'
})

onMounted(() => {
  fetchPost()
  fetchComments()
  fetchLikeStatus()
})

const fetchLikeStatus = async () => {
  try {
    const token = localStorage.getItem('token')
    const headers = token ? { Authorization: `Bearer ${token}` } : {}
    const res = await axios.get(`/api/v1/posts/${route.params.id}/like`, { headers })
    if (res.data.code === 200) {
      likeCount.value = res.data.data.count
      isLiked.value = res.data.data.is_liked
      favoriteCount.value = res.data.data.favorite_count || 0
      isFavorited.value = res.data.data.is_favorited || false
    }
  } catch (e) { console.error(e) }
}

const toggleLike = async () => {
  if (!isLoggedIn.value) {
    alert('请先登录')
    router.push('/login')
    return
  }
  try {
    const res = await axios.post(`/api/v1/posts/${route.params.id}/like`, {}, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    })
    if (res.data.code === 200) {
      isLiked.value = !isLiked.value
      likeCount.value += isLiked.value ? 1 : -1
    }
  } catch (e) { alert('操作失败') }
}

const toggleFavorite = async () => {
  if (!isLoggedIn.value) {
    alert('请先登录')
    router.push('/login')
    return
  }
  try {
    const res = await axios.post(`/api/v1/posts/${route.params.id}/favorite`, {}, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    })
    if (res.data.code === 200) {
      isFavorited.value = !isFavorited.value
      favoriteCount.value += isFavorited.value ? 1 : -1
    }
  } catch (e) { alert('操作失败') }
}

const setReply = (comment, parentId = null) => {
  replyTo.value = { 
    id: comment.ID, 
    username: comment.user.username,
    parentId: parentId || comment.ID 
  }
  // 滚动到评论框
  document.querySelector('.comment-form').scrollIntoView({ behavior: 'smooth' })
}

const cancelReply = () => {
  replyTo.value = null
}
</script>

<template>
  <div class="container" style="max-width: 800px;">
    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>

    <div v-else-if="error" class="alert alert-danger shadow-sm border-0">
      {{ error }}
      <div class="mt-2">
        <router-link to="/" class="btn btn-outline-secondary btn-sm">返回首页</router-link>
      </div>
    </div>

    <div v-else-if="post" class="bg-white rounded shadow-sm p-5">
      <!-- Post Header -->
      <div class="mb-4 text-center">
        <div class="mb-2">
          <span v-if="post.is_system_top" class="badge bg-danger me-2 rounded-pill px-3">
            <i class="bi bi-pin-fill me-1"></i>全站置顶
          </span>
          <span v-if="post.is_top" class="badge bg-primary me-2 rounded-pill px-3">
            <i class="bi bi-pin-angle-fill me-1"></i>个人置顶
          </span>
        </div>
        <h1 class="fw-bold mb-3 text-dark">{{ post.title }}</h1>
        <div class="text-muted small mb-3 d-flex align-items-center justify-content-center gap-3">
          <span class="d-flex align-items-center">
            <div v-if="post.user && post.user.avatar" class="rounded-circle overflow-hidden me-2" style="width: 24px; height: 24px;">
              <img :src="post.user.avatar" class="w-100 h-100 object-fit-cover">
            </div>
            <i v-else class="bi bi-person me-1"></i>
            <router-link v-if="post.user" :to="`/user/${post.user.ID}`" class="text-decoration-none text-muted fw-medium">
              {{ post.user.username }}
            </router-link>
            <span v-else>未知</span>
          </span>
          <span><i class="bi bi-calendar3 me-1"></i>{{ formatDate(post.CreatedAt) }}</span>
        </div>
        <div v-if="post.tags && post.tags.length > 0">
          <span v-for="tag in post.tags" :key="tag.ID" class="badge bg-light text-secondary me-1 border fw-normal">
            #{{ tag.name }}
          </span>
        </div>
      </div>
      
      <!-- Post Content -->
      <div class="post-content py-4 border-top border-bottom">
        <MdPreview :modelValue="post.content" />
      </div>

      <!-- Like & Favorite Buttons -->
      <div class="d-flex justify-content-center gap-3 py-4">
        <button @click="toggleLike" class="btn rounded-pill px-4 py-2 d-flex align-items-center gap-2 transition-all"
          :class="isLiked ? 'btn-danger' : 'btn-outline-danger'">
          <i class="bi" :class="isLiked ? 'bi-heart-fill' : 'bi-heart'"></i>
          <span>{{ isLiked ? '已赞' : '点赞' }}</span>
          <span class="border-start ps-2 ms-1 border-danger-subtle">{{ likeCount }}</span>
        </button>

        <button @click="toggleFavorite" class="btn rounded-pill px-4 py-2 d-flex align-items-center gap-2 transition-all"
          :class="isFavorited ? 'btn-warning text-white' : 'btn-outline-warning'">
          <i class="bi" :class="isFavorited ? 'bi-star-fill' : 'bi-star'"></i>
          <span>{{ isFavorited ? '已收藏' : '收藏' }}</span>
          <span class="border-start ps-2 ms-1 border-warning-subtle">{{ favoriteCount }}</span>
        </button>
      </div>

      <!-- Comments Section -->
      <div class="comments-section mt-5">
        <h4 class="fw-bold mb-4">评论 ({{ comments.length }})</h4>
        
        <!-- Comment List -->
        <div v-if="comments.length > 0" class="mb-5">
          <div v-for="comment in comments" :key="comment.ID" class="mb-4 border-bottom pb-3">
            <div class="d-flex">
              <div class="flex-shrink-0">
                <div v-if="comment.user && comment.user.avatar" class="rounded-circle overflow-hidden" style="width: 40px; height: 40px;">
                  <img :src="comment.user.avatar" class="w-100 h-100 object-fit-cover">
                </div>
                <div v-else class="avatar bg-primary text-white rounded-circle d-flex align-items-center justify-content-center" style="width: 40px; height: 40px;">
                  <span class="fw-bold">{{ comment.user && comment.user.username ? comment.user.username.charAt(0).toUpperCase() : '?' }}</span>
                </div>
              </div>
              <div class="flex-grow-1 ms-3">
                <div class="d-flex justify-content-between align-items-center mb-1">
                  <div>
                    <span class="fw-bold me-2">{{ comment.user ? comment.user.username : '匿名用户' }}</span>
                    <small class="text-muted">{{ formatDate(comment.CreatedAt) }}</small>
                  </div>
                  <div class="d-flex gap-2">
                    <button v-if="isLoggedIn" @click="setReply(comment)" class="btn btn-link text-primary p-0 text-decoration-none" style="font-size: 0.875rem;">
                      <i class="bi bi-reply me-1"></i>回复
                    </button>
                    <button v-if="canDelete(comment)" @click="deleteComment(comment.ID)" class="btn btn-link text-danger p-0 text-decoration-none" style="font-size: 0.875rem;">
                      <i class="bi bi-trash me-1"></i>删除
                    </button>
                  </div>
                </div>
                <p class="mb-0 text-secondary">{{ comment.content }}</p>
                
                <!-- Replies -->
                <div v-if="comment.replies && comment.replies.length > 0" class="mt-3 ms-2 ps-3 border-start">
                  <div v-for="reply in comment.replies" :key="reply.ID" class="mb-3 small">
                    <div class="d-flex justify-content-between align-items-center">
                      <div>
                        <span class="fw-bold me-1">{{ reply.user ? reply.user.username : '匿名用户' }}</span>
                        <span v-if="reply.reply_to_user" class="text-muted me-1">回复 <span class="text-primary">@{{ reply.reply_to_user.username }}</span></span>
                        <small class="text-muted">{{ formatDate(reply.CreatedAt) }}</small>
                      </div>
                      <div class="d-flex gap-2">
                        <button v-if="isLoggedIn" @click="setReply(reply, comment.ID)" class="btn btn-link text-primary p-0 text-decoration-none" style="font-size: 0.75rem;">
                          回复
                        </button>
                        <button v-if="canDelete(reply)" @click="deleteComment(reply.ID)" class="btn btn-link text-danger p-0 text-decoration-none" style="font-size: 0.75rem;">
                          删除
                        </button>
                      </div>
                    </div>
                    <p class="mb-0 text-secondary">{{ reply.content }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="text-center text-muted py-4 mb-4 bg-light rounded">
          暂无评论，快来抢沙发吧！
        </div>

        <!-- Comment Form -->
        <div class="comment-form">
          <h5 class="fw-bold mb-3">发表评论</h5>
          <div v-if="isLoggedIn">
            <div v-if="replyTo" class="alert alert-info d-flex justify-content-between align-items-center py-2 px-3 mb-3 rounded-pill border-0" style="font-size: 0.875rem;">
              <span>正在回复 <strong>{{ replyTo.username }}</strong></span>
              <button @click="cancelReply" class="btn-close" style="font-size: 0.75rem;"></button>
            </div>
            <div class="mb-3">
              <textarea v-model="newComment" class="form-control" rows="4" :placeholder="replyTo ? `回复 ${replyTo.username}...` : '写下你的评论...'"></textarea>
            </div>
            <div class="d-flex justify-content-end">
              <button @click="submitComment" :disabled="submitting || !newComment.trim()" class="btn btn-primary px-4 rounded-pill">
                {{ submitting ? '提交中...' : '提交评论' }}
              </button>
            </div>
          </div>
          <div v-else class="text-center py-4 bg-light rounded border">
            <p class="mb-3 text-secondary">登录后参与评论</p>
            <router-link to="/login" class="btn btn-outline-primary rounded-pill px-4">去登录</router-link>
          </div>
        </div>
      </div>

      <div class="mt-5 pt-4 border-top d-flex justify-content-between align-items-center">
        <router-link :to="backLink" class="btn btn-outline-secondary rounded-pill px-4">
          <i class="bi bi-arrow-left me-1"></i>{{ backText }}
        </router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>
.post-content {
  font-size: 1.15rem;
  line-height: 1.8;
  color: #2c3e50;
  font-family: 'Georgia', serif;
}
.avatar {
  background-color: #f8f9fa;
}
</style>

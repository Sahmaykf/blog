<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'

const activeTab = ref('posts') // 'posts', 'comments', 'notifications'
const posts = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const totalPosts = ref(0)
const totalPages = computed(() => Math.ceil(totalPosts.value / pageSize.value) || 1)

const commentsReceived = ref([])
const commentsSent = ref([])
const notifications = ref([])
const unreadCount = ref(0)

const userId = ref(null)
const userRole = ref('user')
const isAdmin = computed(() => userRole.value === 'admin')

// 编辑/新建状态
const isEditing = ref(false)
const editForm = ref({
  id: null,
  title: '',
  content: '',
  tags: '',
  status: 'published'
})

onMounted(() => {
  userId.value = Number(localStorage.getItem('user_id'))
  userRole.value = localStorage.getItem('user_role') || 'user'
  fetchPosts()
  fetchUnreadCount()
})

const fetchPosts = async () => {
  try {
    const token = localStorage.getItem('token')
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      status: 'all' // 后端逻辑：如果是管理员或查询自己的，status=all 会返回所有
    }
    
    // 如果不是管理员，只查自己的
    if (!isAdmin.value) {
      params.user_id = userId.value
    }

    const response = await axios.get('/api/v1/posts', {
      params,
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const res = response.data
    if (res.code === 200) {
      posts.value = res.data.list
      if (res.data.meta) {
        totalPosts.value = res.data.meta.total
      }
    }
  } catch (error) {
    console.error('获取文章失败', error)
  }
}

const fetchComments = async () => {
  try {
    const token = localStorage.getItem('token')
    const headers = { 'Authorization': `Bearer ${token}` }
    
    const [resReceived, resSent] = await Promise.all([
      axios.get('/api/v1/my/post-comments', { headers }),
      axios.get('/api/v1/my/comments', { headers })
    ])
    
    if (resReceived.data.code === 200) commentsReceived.value = resReceived.data.data
    if (resSent.data.code === 200) commentsSent.value = resSent.data.data
  } catch (error) {
    console.error('获取评论失败', error)
  }
}

const fetchNotifications = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await axios.get('/api/v1/notifications', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (response.data.code === 200) {
      notifications.value = response.data.data
      fetchUnreadCount()
    }
  } catch (error) {
    console.error('获取通知失败', error)
  }
}

const fetchUnreadCount = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await axios.get('/api/v1/notifications/unread-count', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (response.data.code === 200) {
      unreadCount.value = response.data.data
    }
  } catch (error) {
    console.error('获取未读数失败', error)
  }
}

const markAsRead = async (id) => {
  try {
    const token = localStorage.getItem('token')
    await axios.put(`/api/v1/notifications/${id}/read`, {}, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    fetchNotifications()
  } catch (error) {
    console.error('标记已读失败', error)
  }
}

const changePage = (page) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  fetchPosts()
}

const startCreate = () => {
  editForm.value = { id: null, title: '', content: '', tags: '', status: 'published' }
  isEditing.value = true
}

const startEdit = (post) => {
  editForm.value = {
    id: post.ID,
    title: post.title,
    content: post.content,
    tags: post.tags ? post.tags.map(t => t.name).join(',') : '',
    status: post.status || 'published'
  }
  isEditing.value = true
}

const togglePostStatus = async (post) => {
  const newStatus = post.status === 'published' ? 'hidden' : 'published'
  try {
    await axios.put(`/api/v1/posts/${post.ID}`, {
      title: post.title,
      content: post.content,
      status: newStatus
    })
    fetchPosts()
  } catch (error) {
    alert('操作失败')
  }
}

const cancelEdit = () => {
  isEditing.value = false
  editForm.value = { id: null, title: '', content: '', tags: '', status: 'published' }
}

const savePost = async () => {
  const token = localStorage.getItem('token')
  if (!token) return

  const tagNames = editForm.value.tags.split(/[,，]/).map(t => t.trim()).filter(t => t)

  try {
    let response
    if (editForm.value.id) {
      response = await axios.put(`/api/v1/posts/${editForm.value.id}`, {
        title: editForm.value.title,
        content: editForm.value.content,
        tag_names: tagNames,
        status: editForm.value.status
      }, {
        headers: { 'Authorization': `Bearer ${token}` }
      })
    } else {
      response = await axios.post('/api/v1/posts', {
        title: editForm.value.title,
        content: editForm.value.content,
        tag_names: tagNames,
        status: editForm.value.status
      }, {
        headers: { 'Authorization': `Bearer ${token}` }
      })
    }

    const res = response.data
    if (res.code === 200) {
      alert(editForm.value.id ? '更新成功' : '发布成功')
      isEditing.value = false
      fetchPosts()
    } else {
      alert('操作失败: ' + (res.msg || '未知错误'))
    }
  } catch (error) {
    alert('操作失败: ' + (error.response?.data?.msg || '未知错误'))
  }
}

const deletePost = async (id) => {
  if (!confirm('确定要删除这篇文章吗？')) return
  
  try {
    const response = await axios.delete(`/api/v1/posts/${id}`)
    const res = response.data
    if (res.code === 200) {
      fetchPosts()
    } else {
      alert('删除失败: ' + (res.msg || '未知错误'))
    }
  } catch (error) {
    alert('删除失败: ' + (error.response?.data?.msg || '未知错误'))
  }
}

const deleteComment = async (id) => {
  if (!confirm('确定要删除这条评论吗？')) return
  try {
    await axios.delete(`/api/v1/comments/${id}`)
    fetchComments()
  } catch (error) {
    alert('删除失败')
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString()
}

const canManage = (post) => {
  return isAdmin.value || post.user_id === userId.value
}

const switchTab = (tab) => {
  activeTab.value = tab
  if (tab === 'posts') fetchPosts()
  if (tab === 'comments') fetchComments()
  if (tab === 'notifications') fetchNotifications()
}
</script>

<template>
  <div class="admin-container">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h2 class="fw-bold text-dark mb-0">后台管理</h2>
      <div v-if="!isEditing">
        <button class="btn btn-dark rounded-pill px-4 shadow-sm" @click="startCreate">
          <i class="bi bi-plus-lg me-1"></i>发布新文章
        </button>
      </div>
      <button v-else class="btn btn-secondary rounded-pill px-4" @click="cancelEdit">
        <i class="bi bi-arrow-left me-1"></i>返回列表
      </button>
    </div>

    <!-- 标签页导航 -->
    <ul v-if="!isEditing" class="nav nav-pills mb-4 bg-light p-1 rounded-pill shadow-sm d-inline-flex">
      <li class="nav-item">
        <button class="nav-link rounded-pill px-4" :class="{ active: activeTab === 'posts' }" @click="switchTab('posts')">
          <i class="bi bi-file-text me-2"></i>文章管理
        </button>
      </li>
      <li class="nav-item">
        <button class="nav-link rounded-pill px-4" :class="{ active: activeTab === 'comments' }" @click="switchTab('comments')">
          <i class="bi bi-chat-dots me-2"></i>评论管理
        </button>
      </li>
      <li class="nav-item">
        <button class="nav-link rounded-pill px-4 position-relative" :class="{ active: activeTab === 'notifications' }" @click="switchTab('notifications')">
          <i class="bi bi-bell me-2"></i>通知中心
          <span v-if="unreadCount > 0" class="position-absolute top-0 start-100 translate-middle badge rounded-pill bg-danger">
            {{ unreadCount }}
          </span>
        </button>
      </li>
    </ul>

    <!-- 编辑/发布表单 -->
    <div v-if="isEditing" class="card shadow-sm border-0">
      <div class="card-header bg-white border-bottom py-3">
        <h5 class="mb-0 fw-bold text-primary">{{ editForm.id ? '编辑文章' : '发布新文章' }}</h5>
      </div>
      <div class="card-body p-4">
        <form @submit.prevent="savePost">
          <div class="mb-3">
            <label class="form-label text-muted small">标题</label>
            <input type="text" class="form-control form-control-lg bg-light border-0" v-model="editForm.title" required placeholder="请输入文章标题">
          </div>
          <div class="mb-3">
            <label class="form-label text-muted small">标签</label>
            <input type="text" class="form-control bg-light border-0" v-model="editForm.tags" placeholder="输入标签，用逗号分隔">
          </div>
          <div class="mb-3">
            <label class="form-label text-muted small">状态</label>
            <select class="form-select bg-light border-0" v-model="editForm.status">
              <option value="published">发布</option>
              <option value="hidden">隐藏 (草稿)</option>
            </select>
          </div>
          <div class="mb-4">
            <label class="form-label text-muted small">内容</label>
            <textarea class="form-control bg-light border-0" rows="12" v-model="editForm.content" required placeholder="请输入文章内容..."></textarea>
          </div>
          <div class="d-flex gap-2">
            <button type="submit" class="btn btn-dark rounded-pill px-4">保存</button>
            <button type="button" class="btn btn-outline-secondary rounded-pill px-4" @click="cancelEdit">取消</button>
          </div>
        </form>
      </div>
    </div>

    <!-- 内容区域 -->
    <div v-else>
      <!-- 文章管理 -->
      <div v-if="activeTab === 'posts'">
        <div class="card shadow-sm border-0 overflow-hidden">
          <div class="table-responsive">
            <table class="table table-hover align-middle mb-0">
              <thead class="table-light">
                <tr>
                  <th class="ps-4">标题</th>
                  <th>状态</th>
                  <th>作者</th>
                  <th>发布时间</th>
                  <th class="pe-4 text-end">操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="post in posts" :key="post.ID">
                  <td class="ps-4">
                    <div class="fw-bold text-dark">{{ post.title }}</div>
                    <div class="text-muted small">ID: #{{ post.ID }}</div>
                  </td>
                  <td>
                    <span class="badge rounded-pill px-3" :class="post.status === 'published' ? 'bg-success-subtle text-success' : 'bg-warning-subtle text-warning'">
                      {{ post.status === 'published' ? '已发布' : '已隐藏' }}
                    </span>
                  </td>
                  <td>
                    <span class="badge bg-light text-dark border rounded-pill px-3">{{ post.user ? post.user.username : '未知' }}</span>
                  </td>
                  <td class="text-muted small">{{ formatDate(post.CreatedAt) }}</td>
                  <td class="pe-4 text-end">
                    <div v-if="canManage(post)">
                      <button class="btn btn-sm btn-link text-decoration-none me-2 text-dark" @click="togglePostStatus(post)">
                        {{ post.status === 'published' ? '隐藏' : '发布' }}
                      </button>
                      <button class="btn btn-sm btn-outline-dark me-2 rounded-pill px-3" @click="startEdit(post)">编辑</button>
                      <button class="btn btn-sm btn-outline-danger rounded-pill px-3" @click="deletePost(post.ID)">删除</button>
                    </div>
                    <span v-else class="text-muted small fst-italic">无权限</span>
                  </td>
                </tr>
                <tr v-if="posts.length === 0">
                  <td colspan="5" class="text-center py-5 text-muted">暂无文章</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        <!-- 分页 -->
        <nav v-if="posts.length > 0" class="mt-4">
          <ul class="pagination justify-content-center">
            <li class="page-item" :class="{ disabled: currentPage === 1 }">
              <button class="page-link border-0 rounded-circle mx-1" @click="changePage(currentPage - 1)">&laquo;</button>
            </li>
            <li class="page-item disabled"><span class="page-link border-0 bg-transparent fw-bold">{{ currentPage }} / {{ totalPages }}</span></li>
            <li class="page-item" :class="{ disabled: currentPage === totalPages }">
              <button class="page-link border-0 rounded-circle mx-1" @click="changePage(currentPage + 1)">&raquo;</button>
            </li>
          </ul>
        </nav>
      </div>

      <!-- 评论管理 -->
      <div v-if="activeTab === 'comments'">
        <div class="row">
          <div class="col-md-6">
            <h5 class="fw-bold mb-3"><i class="bi bi-chat-left-text me-2"></i>收到的评论</h5>
            <div class="card shadow-sm border-0">
              <div class="list-group list-group-flush">
                <div v-for="c in commentsReceived" :key="c.ID" class="list-group-item p-3">
                  <div class="d-flex justify-content-between align-items-start mb-2">
                    <div class="d-flex align-items-center">
                      <div class="avatar-sm bg-primary text-white rounded-circle me-2 d-flex align-items-center justify-content-center" style="width: 24px; height: 24px; font-size: 12px;">
                        {{ c.user?.username?.charAt(0).toUpperCase() }}
                      </div>
                      <span class="fw-bold small">{{ c.user?.username }}</span>
                    </div>
                    <small class="text-muted">{{ formatDate(c.CreatedAt) }}</small>
                  </div>
                  <p class="mb-2 small text-dark">{{ c.content }}</p>
                  <div class="d-flex justify-content-between align-items-center">
                    <router-link :to="'/post/'+c.post_id" class="text-muted small text-decoration-none">
                      <i class="bi bi-link-45deg"></i> 来自文章: {{ c.post?.title }}
                    </router-link>
                    <button class="btn btn-sm btn-link text-danger p-0 text-decoration-none small" @click="deleteComment(c.ID)">删除</button>
                  </div>
                </div>
                <div v-if="commentsReceived.length === 0" class="p-4 text-center text-muted small">暂无收到的评论</div>
              </div>
            </div>
          </div>
          <div class="col-md-6">
            <h5 class="fw-bold mb-3"><i class="bi bi-chat-right-text me-2"></i>发出的评论</h5>
            <div class="card shadow-sm border-0">
              <div class="list-group list-group-flush">
                <div v-for="c in commentsSent" :key="c.ID" class="list-group-item p-3">
                  <div class="d-flex justify-content-between align-items-start mb-2">
                    <span class="text-muted small">{{ formatDate(c.CreatedAt) }}</span>
                    <button class="btn btn-sm btn-link text-danger p-0 text-decoration-none small" @click="deleteComment(c.ID)">删除</button>
                  </div>
                  <p class="mb-2 small text-dark">{{ c.content }}</p>
                  <router-link :to="'/post/'+c.post_id" class="text-muted small text-decoration-none">
                    <i class="bi bi-link-45deg"></i> 评论文章: {{ c.post?.title }}
                  </router-link>
                </div>
                <div v-if="commentsSent.length === 0" class="p-4 text-center text-muted small">暂无发出的评论</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 通知中心 -->
      <div v-if="activeTab === 'notifications'">
        <div class="card shadow-sm border-0">
          <div class="list-group list-group-flush">
            <div v-for="n in notifications" :key="n.ID" class="list-group-item p-3 d-flex align-items-start" :class="{ 'bg-light': !n.is_read }">
              <div class="flex-shrink-0 me-3">
                <div class="bg-primary-subtle text-primary rounded-circle p-2">
                  <i class="bi" :class="n.type === 'comment' ? 'bi-chat-dots' : 'bi-bell'"></i>
                </div>
              </div>
              <div class="flex-grow-1">
                <div class="d-flex justify-content-between align-items-center mb-1">
                  <span class="fw-bold small">
                    {{ n.from_user?.username }} 评论了你的文章
                  </span>
                  <small class="text-muted">{{ formatDate(n.CreatedAt) }}</small>
                </div>
                <p class="text-muted small mb-2">{{ n.content }}</p>
                <div class="d-flex gap-3">
                  <router-link :to="'/post/'+n.post_id" class="btn btn-sm btn-outline-dark rounded-pill px-3 py-0 small" @click="markAsRead(n.ID)">查看</router-link>
                  <button v-if="!n.is_read" class="btn btn-sm btn-link text-muted p-0 text-decoration-none small" @click="markAsRead(n.ID)">标记已读</button>
                </div>
              </div>
            </div>
            <div v-if="notifications.length === 0" class="p-5 text-center text-muted">暂无通知</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.nav-pills .nav-link {
  color: #6c757d;
  font-weight: 500;
  transition: all 0.2s ease;
  border: none;
}
.nav-pills .nav-link.active {
  background-color: #fff;
  color: #212529;
  box-shadow: 0 2px 6px rgba(0,0,0,0.08);
}
.table th {
  font-weight: 600;
  text-transform: uppercase;
  font-size: 0.75rem;
  letter-spacing: 0.5px;
}
.list-group-item {
  transition: background-color 0.2s;
}
.list-group-item:hover {
  background-color: #f8f9fa;
}
</style>

<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()
const posts = ref([])
const hotPosts = ref([])
const tags = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const totalPosts = ref(0)
const totalPages = computed(() => Math.ceil(totalPosts.value / pageSize.value) || 1)
const currentTag = ref('')
const orderBy = ref('created_at')
const searchKeyword = ref('')
const loading = ref(false)

const fetchPosts = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      order_by: orderBy.value
    }
    if (currentTag.value) {
      params.tag = currentTag.value
    }
    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }

    const response = await axios.get('/api/v1/posts', { params })
    const res = response.data
    if (res.code === 200) {
      posts.value = res.data.list
      if (res.data.meta) {
        totalPosts.value = res.data.meta.total
      }
    }
  } catch (error) {
    console.error('获取文章失败', error)
  } finally {
    loading.value = false
  }
}

const fetchHotPosts = async () => {
  try {
    const response = await axios.get('/api/v1/posts/hot?limit=5')
    if (response.data.code === 200) {
      hotPosts.value = response.data.data
    }
  } catch (error) {
    console.error('获取热门文章失败', error)
  }
}

const fetchTags = async () => {
  try {
    const response = await axios.get('/api/v1/tags')
    const res = response.data
    if (res.code === 200) {
      tags.value = res.data
    }
  } catch (error) {
    console.error('获取标签失败', error)
  }
}

const filterByTag = (tagName) => {
  currentTag.value = tagName
  currentPage.value = 1
  fetchPosts()
}

const clearFilter = () => {
  currentTag.value = ''
  currentPage.value = 1
  fetchPosts()
}

const clearAllFilters = () => {
  currentTag.value = ''
  searchKeyword.value = ''
  currentPage.value = 1
  fetchPosts()
}

const handleOrderChange = (order) => {
  orderBy.value = order
  currentPage.value = 1
  fetchPosts()
}

const handleSearch = () => {
  if (!searchKeyword.value.trim()) return
  router.push({ path: '/search', query: { q: searchKeyword.value } })
}

const clearSearch = () => {
  searchKeyword.value = ''
  currentPage.value = 1
  fetchPosts()
}

const changePage = (page) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  fetchPosts()
}

const handlePageSizeChange = () => {
  currentPage.value = 1
  fetchPosts()
}

const viewPost = (id) => {
  router.push(`/post/${id}`)
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString()
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

onMounted(() => {
  fetchPosts()
  fetchHotPosts()
  fetchTags()
})
</script>

<template>
  <div class="home-container">
    <div class="row">
      <!-- 左侧文章列表 -->
      <div class="col-md-9">
        <div class="d-flex justify-content-between align-items-center mb-4">
          <h4 class="fw-bold text-dark mb-0">
            {{ searchKeyword ? `搜索: ${searchKeyword}` : (currentTag ? `标签: ${currentTag}` : '最新文章') }}
            <button v-if="currentTag || searchKeyword" class="btn btn-sm btn-outline-secondary ms-2 rounded-pill" @click="clearAllFilters">
              <i class="bi bi-x-lg me-1"></i>清除筛选
            </button>
          </h4>
          
          <!-- 排序切换 -->
          <div class="btn-group btn-group-sm shadow-sm rounded-pill overflow-hidden bg-white p-1">
            <button class="btn border-0 rounded-pill px-3" 
              :class="orderBy === 'created_at' ? 'btn-primary' : 'btn-white text-muted'"
              @click="handleOrderChange('created_at')">最新</button>
            <button class="btn border-0 rounded-pill px-3" 
              :class="orderBy === 'views' ? 'btn-primary' : 'btn-white text-muted'"
              @click="handleOrderChange('views')">热度</button>
            <button class="btn border-0 rounded-pill px-3" 
              :class="orderBy === 'likes' ? 'btn-primary' : 'btn-white text-muted'"
              @click="handleOrderChange('likes')">点赞</button>
          </div>
        </div>
        
        <div v-if="loading" class="text-center py-5">
          <div class="spinner-border text-primary" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
        </div>

        <template v-else>
          <div v-if="posts.length === 0" class="text-center py-5 text-muted bg-white rounded shadow-sm">
            <p class="mb-0">暂无文章...</p>
          </div>

          <div class="row g-4">
            <div v-for="post in posts" :key="post.ID" class="col-12">
              <div class="card h-100 border-0 shadow-sm hover-shadow transition-all position-relative" @click="viewPost(post.ID)" style="cursor: pointer;">
                <!-- 置顶标识 -->
                <div v-if="post.is_system_top" class="position-absolute top-0 start-0 bg-danger text-white px-3 py-1 rounded-bottom-end small shadow-sm" style="z-index: 1; border-top-left-radius: 0.375rem;">
                  <i class="bi bi-pin-fill me-1"></i>全站置顶
                </div>
                
                <div class="card-body p-4" :class="{'pt-5': post.is_system_top}">
                  <h5 class="card-title fw-bold text-dark mb-2 text-truncate-2">{{ post.title }}</h5>
                  <div class="mb-3 text-muted small d-flex align-items-center gap-3">
                    <span class="d-flex align-items-center">
                      <i class="bi bi-person me-1"></i>
                      <span v-if="post.user" @click.stop="router.push(`/user/${post.user.ID}`)" class="text-primary" style="cursor: pointer;">
                        {{ post.user.username }}
                      </span>
                      <span v-else>未知</span>
                    </span>
                    <span><i class="bi bi-calendar3 me-1"></i>{{ formatDate(post.CreatedAt) }}</span>
                    <span class="ms-auto d-flex gap-3">
                      <span><i class="bi bi-eye me-1"></i>{{ post.view_count || 0 }}</span>
                      <span><i class="bi bi-heart me-1"></i>{{ post.like_count || 0 }}</span>
                    </span>
                  </div>
                  
                  <p class="card-text text-secondary text-truncate-3 mb-3">{{ truncateContent(post.content) }}</p>
                  
                  <div class="d-flex justify-content-between align-items-center">
                    <div class="tags">
                      <span v-for="tag in post.tags" :key="tag.ID" 
                        class="badge bg-light text-secondary me-1 border fw-normal" 
                        @click.stop="filterByTag(tag.name)">
                        #{{ tag.name }}
                      </span>
                    </div>
                    <button class="btn btn-outline-primary btn-sm rounded-pill px-3">阅读全文</button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </template>

        <!-- 分页控件 -->
        <div class="d-flex justify-content-center mt-5 align-items-center" v-if="posts.length > 0">
          <nav aria-label="Page navigation">
            <ul class="pagination shadow-sm mb-0 rounded-pill overflow-hidden">
              <li class="page-item" :class="{ disabled: currentPage === 1 }">
                <button class="page-link border-0 px-3" @click="changePage(currentPage - 1)">
                  <span aria-hidden="true">&laquo;</span>
                </button>
              </li>
              <li class="page-item disabled">
                <span class="page-link border-0 bg-white text-dark fw-bold px-3">{{ currentPage }} / {{ totalPages }}</span>
              </li>
              <li class="page-item" :class="{ disabled: currentPage === totalPages }">
                <button class="page-link border-0 px-3" @click="changePage(currentPage + 1)">
                  <span aria-hidden="true">&raquo;</span>
                </button>
              </li>
            </ul>
          </nav>

          <div class="ms-3">
            <select class="form-select form-select-sm border-0 bg-light text-secondary rounded-pill ps-3 pe-4 py-1 shadow-sm" 
              style="width: auto; cursor: pointer; font-size: 0.85rem;" 
              v-model="pageSize" 
              @change="handlePageSizeChange">
              <option :value="5">5 条/页</option>
              <option :value="10">10 条/页</option>
              <option :value="20">20 条/页</option>
            </select>
          </div>
        </div>
      </div>

      <!-- 右侧侧边栏 -->
      <div class="col-md-3">
        <!-- 搜索框 -->
        <div class="card border-0 shadow-sm mb-4 overflow-hidden">
          <div class="card-body p-0">
            <div class="input-group">
              <span class="input-group-text bg-white border-0 ps-3">
                <i class="bi bi-search text-muted"></i>
              </span>
              <input type="text" 
                class="form-control border-0 py-3 shadow-none" 
                v-model="searchKeyword" 
                placeholder="搜索文章、用户或标签..."
                @keyup.enter="handleSearch">
              <button v-if="searchKeyword" 
                class="btn bg-white border-0 text-muted" 
                @click="clearSearch">
                <i class="bi bi-x-lg"></i>
              </button>
            </div>
          </div>
        </div>

        <!-- 热门文章 -->
        <div class="card border-0 shadow-sm mb-4">
          <div class="card-header bg-white border-0 pt-4 pb-2">
            <h5 class="fw-bold mb-0"><i class="bi bi-fire text-danger me-2"></i>热门文章</h5>
          </div>
          <div class="card-body pt-0">
            <div class="list-group list-group-flush">
              <div v-for="(post, index) in hotPosts" :key="post.ID" 
                class="list-group-item px-0 py-3 border-0 border-bottom-dashed cursor-pointer"
                @click="viewPost(post.ID)">
                <div class="d-flex align-items-start">
                  <span class="badge rounded-circle me-2 mt-1 d-flex align-items-center justify-content-center" 
                    :class="index < 3 ? 'bg-blue-soft' : 'bg-secondary'"
                    style="width: 20px; height: 20px; font-size: 0.7rem;">
                    {{ index + 1 }}
                  </span>
                  <div class="flex-grow-1">
                    <div class="fw-bold text-dark small text-truncate-2 mb-1 hover-text-primary">{{ post.title }}</div>
                    <div class="text-muted" style="font-size: 0.75rem;">
                      <i class="bi bi-person me-1"></i>{{ post.user?.username }}
                    </div>
                  </div>
                </div>
              </div>
              <div v-if="hotPosts.length === 0" class="text-center py-3 text-muted small">暂无热门文章</div>
            </div>
          </div>
        </div>

        <!-- 热门标签 -->
        <div class="card border-0 shadow-sm mb-4">
          <div class="card-header bg-white border-0 pt-4 pb-2">
            <h5 class="fw-bold mb-0"><i class="bi bi-tags me-2 text-primary"></i>分类导航</h5>
          </div>
          <div class="card-body">
            <div class="d-flex flex-wrap gap-2">
              <button v-for="tag in tags" :key="tag.ID" 
                class="btn btn-sm rounded-pill px-3"
                :class="currentTag === tag.name ? 'btn-primary' : 'btn-outline-secondary'"
                @click="filterByTag(tag.name)">
                {{ tag.name }}
              </button>
              <div v-if="tags.length === 0" class="text-muted small">暂无分类</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.hover-shadow:hover {
  transform: translateY(-5px);
  box-shadow: 0 .5rem 1rem rgba(0,0,0,.15)!important;
}
.transition-all {
  transition: all 0.3s ease;
}
.text-truncate-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.text-truncate-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.border-bottom-dashed {
  border-bottom: 1px dashed #dee2e6 !important;
}
.border-bottom-dashed:last-child {
  border-bottom: none !important;
}
.cursor-pointer {
  cursor: pointer;
}
.hover-text-primary:hover {
  color: #0d6efd !important;
}
.bg-blue-soft {
  background-color: #5c92d1 !important;
  color: white;
}
.rounded-bottom-end {
  border-bottom-right-radius: 1rem !important;
}
</style>

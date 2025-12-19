<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()
const posts = ref([])
const tags = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const totalPosts = ref(0)
const totalPages = computed(() => Math.ceil(totalPosts.value / pageSize.value) || 1)
const currentTag = ref('')
const loading = ref(false)

const fetchPosts = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (currentTag.value) {
      params.tag = currentTag.value
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
  return content.length > 150 ? content.substring(0, 150) + '...' : content
}

onMounted(() => {
  fetchPosts()
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
            {{ currentTag ? `标签: ${currentTag}` : '最新文章' }}
            <button v-if="currentTag" class="btn btn-sm btn-outline-secondary ms-2 rounded-pill" @click="clearFilter">
              <i class="bi bi-x-lg me-1"></i>清除筛选
            </button>
          </h4>
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
              <div class="card h-100 border-0 shadow-sm hover-shadow transition-all" @click="viewPost(post.ID)" style="cursor: pointer;">
                <div class="card-body p-4">
                  <h5 class="card-title fw-bold text-dark mb-2 text-truncate-2">{{ post.title }}</h5>
                  <div class="mb-3 text-muted small">
                    <span class="me-3">
                      <i class="bi bi-person me-1"></i>
                      <span v-if="post.user" @click.stop="router.push(`/user/${post.user.ID}`)" class="text-primary" style="cursor: pointer;">
                        {{ post.user.username }}
                      </span>
                      <span v-else>未知</span>
                    </span>
                    <span><i class="bi bi-calendar3 me-1"></i>{{ formatDate(post.CreatedAt) }}</span>
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

      <!-- 右侧标签栏 -->
      <div class="col-md-3">
        <div class="card border-0 shadow-sm mb-4">
          <div class="card-header bg-white border-0 pt-4 pb-2">
            <h5 class="fw-bold mb-0">热门标签</h5>
          </div>
          <div class="card-body">
            <div class="d-flex flex-wrap gap-2">
              <button v-for="tag in tags" :key="tag.ID" 
                class="btn btn-sm rounded-pill"
                :class="currentTag === tag.name ? 'btn-primary' : 'btn-outline-secondary'"
                @click="filterByTag(tag.name)">
                {{ tag.name }}
              </button>
              <div v-if="tags.length === 0" class="text-muted small">暂无标签</div>
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
</style>

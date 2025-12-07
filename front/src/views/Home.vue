<template>
  <Layout>
    <!-- 搜索栏 -->
    <div class="search-bar">
      <input
        v-model="searchKeyword"
        type="text"
        placeholder="搜索文章标题/内容..."
        class="search-input"
      />
      <button class="search-btn" @click="getArticleList">搜索</button>
    </div>

    <!-- 文章列表 -->
    <div class="article-list" v-if="!loading">
      <div class="article-card" v-for="article in articleList" :key="article.id">
        <router-link :to="`/article/${article.id}`" class="card-title">
          {{ article.title }}
        </router-link>
        <div class="card-meta">
          <span class="author">作者：{{ article.user?.username || '未知' }}</span>
          <span class="category">分类：{{ article.category?.name || '未分类' }}</span>
          <span class="view-count">阅读：{{ article.view_count }}</span>
          <span class="create-time">{{ formatTime(article.created_at) }}</span>
        </div>
        <div class="card-excerpt">
          {{ article.content.length > 150 ? article.content.slice(0, 150) + '...' : article.content }}
        </div>
        <router-link :to="`/article/${article.id}`" class="card-more">查看全文</router-link>
      </div>
    </div>

    <!-- 加载中 -->
    <div class="loading-tip" v-if="loading">加载中...</div>

    <!-- 空数据 -->
    <div class="empty-tip" v-if="!loading && articleList.length === 0">
      暂无文章，<router-link to="/publish" v-if="hasToken">快去发布第一篇文章吧</router-link>
      <span v-else>登录后发布你的第一篇文章吧</span>
    </div>

    <!-- 分页 -->
    <div class="pagination" v-if="!loading && total > 0">
      <button
        class="page-btn"
        :disabled="currentPage === 1"
        @click="handlePageChange(currentPage - 1)"
      >
        上一页
      </button>
      <span class="page-info">
        第 {{ currentPage }} 页 / 共 {{ totalPage }} 页（总计 {{ total }} 篇）
      </span>
      <button
        class="page-btn"
        :disabled="currentPage === totalPage"
        @click="handlePageChange(currentPage + 1)"
      >
        下一页
      </button>
    </div>
  </Layout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import Layout from '../components/Layout.vue'

const router = useRouter()
const loading = ref(false)
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const articleList = ref([])
const total = ref(0)
const hasToken = computed(() => !!localStorage.getItem('token'))

// 总页数
const totalPage = computed(() => {
  return Math.ceil(total.value / pageSize.value)
})

// 格式化时间
const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')}`
}

// 获取文章列表
const getArticleList = async () => {
  try {
    loading.value = true
    const res = await axios.get('/api/v1/article/list', {
      params: {
        page: currentPage.value,
        size: pageSize.value,
        keyword: searchKeyword.value
      }
    })
    if (res.data.code === 200) {
      articleList.value = res.data.data.list || []
      total.value = res.data.data.total || 0
    }
  } catch (err) {
    console.error('获取文章列表失败：', err)
  } finally {
    loading.value = false
  }
}

// 分页切换
const handlePageChange = (page) => {
  if (page < 1 || page > totalPage.value) return
  currentPage.value = page
  getArticleList()
}

onMounted(() => {
  getArticleList()
})
</script>

<style scoped>
.search-bar {
  display: flex;
  max-width: 600px;
  margin: 0 auto 30px;
  gap: 8px;
}
.search-input {
  flex: 1;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}
.search-input:focus {
  outline: none;
  border-color: #409eff;
}
.search-btn {
  padding: 0 16px;
  background-color: #409eff;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
}
.search-btn:hover {
  background-color: #337ecc;
}
.article-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}
.article-card {
  padding: 20px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: transform 0.2s, box-shadow 0.2s;
}
.article-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
}
.card-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  text-decoration: none;
  margin-bottom: 12px;
  display: block;
}
.card-title:hover {
  color: #409eff;
}
.card-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  font-size: 12px;
  color: #999;
  margin-bottom: 10px;
}
.card-excerpt {
  font-size: 14px;
  color: #666;
  line-height: 1.6;
  margin-bottom: 12px;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.card-more {
  font-size: 13px;
  color: #409eff;
  text-decoration: none;
}
.card-more:hover {
  text-decoration: underline;
}
.loading-tip, .empty-tip {
  text-align: center;
  padding: 40px 0;
  color: #666;
  font-size: 14px;
}
.empty-tip a {
  color: #409eff;
  text-decoration: none;
}
.empty-tip a:hover {
  text-decoration: underline;
}
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-top: 20px;
}
.page-btn {
  padding: 6px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background-color: #fff;
  color: #333;
  cursor: pointer;
  transition: all 0.2s;
}
.page-btn:disabled {
  color: #999;
  cursor: not-allowed;
  background-color: #f5f5f5;
}
.page-btn:hover:not(:disabled) {
  border-color: #409eff;
  color: #409eff;
}
.page-info {
  font-size: 14px;
  color: #666;
}
</style>
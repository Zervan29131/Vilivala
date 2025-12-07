<template>
  <Layout>
    <div class="detail-container" v-if="!loading">
      <!-- 文章标题 -->
      <h1 class="detail-title">{{ article.title }}</h1>
      
      <!-- 文章元信息 -->
      <div class="detail-meta">
        <span class="author">作者：{{ article.user?.username || '未知' }}</span>
        <span class="category">分类：{{ article.category?.name || '未分类' }}</span>
        <span class="view-count">阅读：{{ article.view_count }}</span>
        <span class="create-time">发布时间：{{ formatTime(article.created_at) }}</span>
      </div>

      <!-- 文章内容 -->
      <div class="detail-content">
        {{ article.content }}
      </div>

      <!-- 返回按钮 -->
      <button class="back-btn" @click="router.back()">返回列表</button>
    </div>

    <!-- 加载中 -->
    <div class="loading-tip" v-if="loading">加载中...</div>

    <!-- 文章不存在 -->
    <div class="empty-tip" v-if="!loading && !article.id">
      该文章不存在或已被删除
      <button class="back-btn" @click="router.push('/')">返回首页</button>
    </div>
  </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import axios from 'axios'
import Layout from '../components/Layout.vue'

const router = useRouter()
const route = useRoute()
const loading = ref(false)
const article = ref({})

// 格式化时间
const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')} ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}

// 获取文章详情
const getArticleDetail = async () => {
  try {
    loading.value = true
    const articleId = route.params.id
    const res = await axios.get(`/api/v1/article/${articleId}`)
    if (res.data.code === 200) {
      article.value = res.data.data || {}
    }
  } catch (err) {
    console.error('获取文章详情失败：', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  getArticleDetail()
})
</script>

<style scoped>
.detail-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px 0;
}
.detail-title {
  font-size: 28px;
  font-weight: bold;
  color: #333;
  text-align: center;
  margin-bottom: 20px;
  line-height: 1.4;
}
.detail-meta {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 16px;
  font-size: 14px;
  color: #999;
  margin-bottom: 30px;
  padding-bottom: 16px;
  border-bottom: 1px solid #eee;
}
.detail-content {
  font-size: 16px;
  color: #333;
  line-height: 1.8;
  white-space: pre-wrap; /* 保留换行 */
  margin-bottom: 40px;
}
.back-btn {
  padding: 8px 16px;
  background-color: #f5f7fa;
  border: 1px solid #ddd;
  border-radius: 4px;
  color: #333;
  cursor: pointer;
  transition: all 0.2s;
}
.back-btn:hover {
  background-color: #e5e9f0;
  border-color: #409eff;
  color: #409eff;
}
.loading-tip, .empty-tip {
  text-align: center;
  padding: 40px 0;
  color: #666;
  font-size: 14px;
}
.empty-tip .back-btn {
  margin-top: 16px;
}
</style>
<template>
  <div class="detail-container">
    <el-card>
      <!-- 文章标题 -->
      <div class="article-title">{{ article.title }}</div>
      <!-- 文章元信息 -->
      <div class="article-meta">
        <span>作者：{{ article.user.username }}</span>
        <span>分类：{{ article.category.name }}</span>
        <span>阅读量：{{ article.view_count }}</span>
        <span>{{ article.created_at }}</span>
      </div>
      <!-- 文章封面 -->
      <el-image :src="article.cover_img" fit="cover" style="width: 100%; height: 400px; margin: 20px 0"></el-image>
      <!-- 文章内容 -->
      <div class="article-content" v-html="article.content"></div>
    </el-card>

    <!-- 评论区（后续可扩展） -->
    <el-card style="margin-top: 20px">
      <template #header>
        <span>评论区</span>
      </template>
      <el-divider>暂无评论，快来抢沙发~</el-divider>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getArticleDetail } from '@/api/article'

const route = useRoute()
const article = ref({})

// 获取文章详情
const loadArticleDetail = () => {
  const id = route.params.id
  getArticleDetail(id).then(res => {
    article.value = res.data
  })
}

onMounted(() => {
  loadArticleDetail()
})
</script>

<style scoped>
.detail-container {
  padding: 20px 0;
}

.article-title {
  font-size: 24px;
  font-weight: bold;
  text-align: center;
  margin-bottom: 20px;
}

.article-meta {
  display: flex;
  justify-content: center;
  gap: 20px;
  color: #999;
  font-size: 14px;
  margin-bottom: 20px;
}

.article-content {
  line-height: 2;
  font-size: 16px;
  color: #333;
}
</style>
<template>
  <div class="home-container">
    <el-row :gutter="20">
      <!-- 左侧文章列表 -->
      <el-col :span="16">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>最新文章</span>
              <el-input v-model="searchKeyword" placeholder="搜索文章标题/内容" style="width: 200px; margin-left: 20px"></el-input>
              <el-button type="primary" icon="Search" @click="getArticleList()">搜索</el-button>
            </div>
          </template>

          <!-- 文章列表 -->
          <div class="article-list">
            <el-divider v-if="articleList.length === 0">暂无文章</el-divider>
            <el-card v-for="article in articleList" :key="article.id" class="article-item" @click="goToDetail(article.id)">
              <div class="article-title">{{ article.title }}</div>
              <div class="article-meta">
                <span>作者：{{ article.user.username }}</span>
                <span>分类：{{ article.category.name }}</span>
                <span>阅读量：{{ article.view_count }}</span>
                <span>{{ article.created_at }}</span>
              </div>
              <div class="article-content">{{ article.content.slice(0, 200) }}...</div>
            </el-card>
          </div>

          <!-- 分页 -->
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="size"
            :total="total"
            layout="prev, pager, next, jumper, ->, total"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
            style="margin-top: 20px; text-align: center"
          >
          </el-pagination>
        </el-card>
      </el-col>

      <!-- 右侧分类/标签 -->
      <el-col :span="8">
        <el-card>
          <template #header>
            <span>文章分类</span>
          </template>
          <el-menu>
            <el-menu-item index="0">全部文章</el-menu-item>
            <el-menu-item v-for="category in categoryList" :key="category.id" :index="category.id.toString()">
              {{ category.name }}
            </el-menu-item>
          </el-menu>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getArticleList } from '@/api/article'

const router = useRouter()

// 分页参数
const page = ref(1)
const size = ref(10)
const total = ref(0)
// 搜索参数
const searchKeyword = ref('')
// 文章列表/分类列表
const articleList = ref([])
const categoryList = ref([])

// 获取文章列表
const loadArticleList = () => {
  getArticleList({
    page: page.value,
    size: size.value,
    keyword: searchKeyword.value
  }).then(res => {
    articleList.value = res.data.list
    total.value = res.data.total
  })
}

// 分页大小改变
const handleSizeChange = (val) => {
  size.value = val
  loadArticleList()
}

// 当前页改变
const handleCurrentChange = (val) => {
  page.value = val
  loadArticleList()
}

// 跳转到文章详情
const goToDetail = (id) => {
  router.push(`/article/${id}`)
}

// 页面加载时获取数据
onMounted(() => {
  loadArticleList()
  // 后续可补充获取分类列表的接口调用
})
</script>

<style scoped>
.home-container {
  padding: 20px 0;
}

.article-list {
  margin-top: 10px;
}

.article-item {
  margin-bottom: 10px;
  cursor: pointer;
  transition: all 0.3s;
}

.article-item:hover {
  transform: translateY(-5px);
}

.article-title {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 10px;
  color: #409eff;
}

.article-meta {
  display: flex;
  justify-content: space-between;
  color: #999;
  font-size: 12px;
  margin-bottom: 10px;
}

.article-content {
  color: #666;
  line-height: 1.6;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
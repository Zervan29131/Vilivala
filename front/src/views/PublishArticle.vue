<template>
  <Layout>
    <div class="publish-container">
      <h2 class="publish-title">发布文章</h2>
      <form @submit.prevent="handlePublish">
        <!-- 标题 -->
        <div class="form-item">
          <label class="form-label">文章标题</label>
          <input
            v-model="form.title"
            type="text"
            class="form-input"
            placeholder="请输入文章标题（1-100位）"
            required
          />
        </div>

        <!-- 分类 -->
        <div class="form-item">
          <label class="form-label">文章分类</label>
          <select v-model="form.category_id" class="form-select" required>
            <option value="">请选择分类</option>
            <option v-for="category in categoryList" :key="category.id" :value="category.id">
              {{ category.name }}
            </option>
          </select>
        </div>

        <!-- 内容 -->
        <div class="form-item">
          <label class="form-label">文章内容</label>
          <textarea
            v-model="form.content"
            class="form-textarea"
            rows="10"
            placeholder="请输入文章内容"
            required
          ></textarea>
        </div>

        <!-- 是否发布 -->
        <div class="form-item form-radio">
          <label class="form-label">发布状态</label>
          <label class="radio-item">
            <input type="radio" v-model="form.is_publish" value="true" checked /> 发布
          </label>
          <label class="radio-item">
            <input type="radio" v-model="form.is_publish" value="false" /> 草稿
          </label>
        </div>

        <!-- 错误提示 -->
        <div class="form-error" v-if="errorMsg">{{ errorMsg }}</div>

        <!-- 提交按钮 -->
        <button type="submit" class="btn-submit" :disabled="loading">
          <span v-if="loading">发布中...</span>
          <span v-else>发布文章</span>
        </button>
      </form>
    </div>
  </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import Layout from '../components/Layout.vue'

const router = useRouter()
const loading = ref(false)
const errorMsg = ref('')
const categoryList = ref([])
const form = ref({
  title: '',
  category_id: '',
  content: '',
  is_publish: 'true'
})

// 获取分类列表
const getCategoryList = async () => {
  try {
    const res = await axios.get('/api/v1/category/list')
    if (res.data.code === 200) {
      categoryList.value = res.data.data.list || []
    }
  } catch (err) {
    console.error('获取分类列表失败：', err)
  }
}

// 发布文章
const handlePublish = async () => {
  try {
    loading.value = true
    errorMsg.value = ''
    // 表单校验
    if (!form.value.title || form.value.title.length > 100) {
      errorMsg.value = '文章标题不能为空且长度不超过100位'
      return
    }
    if (!form.value.category_id) {
      errorMsg.value = '请选择文章分类'
      return
    }
    if (!form.value.content) {
      errorMsg.value = '文章内容不能为空'
      return
    }
    // 请求发布接口
    const res = await axios.post('/api/v1/article', {
      title: form.value.title,
      category_id: form.value.category_id,
      content: form.value.content,
      is_publish: form.value.is_publish === 'true'
    })
    if (res.data.code === 200) {
      alert('发布成功！')
      router.push('/')
    } else {
      errorMsg.value = res.data.msg || '发布失败'
    }
  } catch (err) {
    errorMsg.value = '服务器错误，请稍后重试'
    console.error('发布文章失败：', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  getCategoryList()
})
</script>

<style scoped>
.publish-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px 0;
}
.publish-title {
  font-size: 24px;
  color: #333;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
}
.form-item {
  margin-bottom: 20px;
}
.form-label {
  display: block;
  margin-bottom: 8px;
  color: #333;
  font-size: 14px;
  font-weight: 500;
}
.form-input, .form-select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
}
.form-input:focus, .form-select:focus {
  outline: none;
  border-color: #409eff;
}
.form-textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
  resize: vertical;
}
.form-textarea:focus {
  outline: none;
  border-color: #409eff;
}
.form-radio {
  display: flex;
  align-items: center;
  gap: 16px;
}
.radio-item {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #666;
  font-size: 14px;
}
.form-error {
  color: #f56c6c;
  font-size: 13px;
  margin-bottom: 10px;
}
.btn-submit {
  padding: 10px 20px;
  background-color: #409eff;
  color: #fff;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.2s;
}
.btn-submit:disabled {
  background-color: #a0cfff;
  cursor: not-allowed;
}
.btn-submit:hover:not(:disabled) {
  background-color: #337ecc;
}
</style>
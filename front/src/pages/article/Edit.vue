<template>
  <div class="edit-container">
    <el-card>
      <el-form :model="articleForm" :rules="articleRules" ref="articleFormRef" label-width="80px">
        <!-- 其他表单项不变 -->
        <el-form-item label="文章内容" prop="content">
          <!-- 替换为wangEditor -->
          <div style="border: 1px solid #ccc; margin-top: 10px">
            <Editor
              v-model="articleForm.content"
              placeholder="请输入文章内容"
              :toolbarConfig="toolbarConfig"
              :editorConfig="editorConfig"
              style="height: 500px"
            />
          </div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSubmit">保存文章</el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
// 导入wangEditor
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import '@wangeditor/editor/dist/css/style.css' // 引入样式

import { createArticle, updateArticle, getArticleDetail } from '@/api/article'

const router = useRouter()
const route = useRoute()
const articleFormRef = ref(null)
const isEdit = ref(!!route.params.id)

// 文章表单（不变）
const articleForm = ref({
  title: '',
  content: '',
  cover_img: '',
  category_id: '',
  is_publish: true,
  tag_ids: []
})

// 表单校验规则（不变）
const articleRules = ref({
  title: [{ required: true, message: '请输入文章标题', trigger: 'blur' }],
  category_id: [{ required: true, message: '请选择文章分类', trigger: 'change' }],
  content: [{ required: true, message: '请输入文章内容', trigger: 'blur' }]
})

// 分类列表（模拟数据，后续替换为接口）
const categoryList = ref([
  { id: 1, name: '技术博客' },
  { id: 2, name: '生活随笔' },
  { id: 3, name: '读书笔记' }
])

// wangEditor配置
const toolbarConfig = ref({})
const editorConfig = ref({
  placeholder: '请输入文章内容',
  MENU_CONF: {
    uploadImage: {
      // 后续可扩展图片上传，暂时注释
      // server: '/api/v1/upload/image', // 后端上传接口
      // fieldName: 'file',
    }
  }
})

// 加载文章详情（不变）
const loadArticleDetail = () => {
  const id = route.params.id
  getArticleDetail(id).then(res => {
    articleForm.value = {
      title: res.data.title,
      content: res.data.content,
      cover_img: res.data.cover_img,
      category_id: res.data.category_id,
      is_publish: res.data.is_publish,
      tag_ids: res.data.tags?.map(tag => tag.id) || []
    }
  })
}

// 提交文章（不变）
const handleSubmit = () => {
  articleFormRef.value.validate((valid) => {
    if (!valid) return
    if (!isEdit.value) {
      createArticle(articleForm.value).then(() => {
        ElMessage.success('文章发布成功')
        router.push('/')
      })
    } else {
      updateArticle(route.params.id, articleForm.value).then(() => {
        ElMessage.success('文章更新成功')
        router.push('/')
      })
    }
  })
}

// 返回上一页（不变）
const goBack = () => {
  router.go(-1)
}

onMounted(() => {
  if (isEdit.value) {
    loadArticleDetail()
  }
})
</script>

<style scoped>
.edit-container {
  padding: 20px 0;
}
</style>
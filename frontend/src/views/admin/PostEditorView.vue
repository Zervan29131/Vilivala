<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-xl font-bold text-[#213547] dark:text-white">{{ isEdit ? '编辑文章' : '新建文章' }}</h1>
      <router-link to="/admin/posts" class="text-sm text-gray-500 dark:text-gray-400 hover:text-[#42b883]">← 返回列表</router-link>
    </div>

    <div class="card p-5 space-y-4 mb-4">
      <div class="grid sm:grid-cols-[1fr_auto] gap-4 items-end">
        <div>
          <label class="lbl">标题</label>
          <input v-model="form.title" class="input-field" placeholder="文章标题" />
        </div>
        <label class="inline-flex items-center gap-2 text-sm text-gray-600 dark:text-gray-300 pb-2 cursor-pointer select-none">
          <input v-model="form.is_top" type="checkbox" class="rounded accent-[#42b883] w-4 h-4" />
          置顶文章
        </label>
      </div>

      <div class="grid sm:grid-cols-2 gap-4">
        <div>
          <label class="lbl">Slug（留空自动生成）</label>
          <input v-model="form.slug" class="input-field font-mono" placeholder="my-first-post" />
        </div>
        <div>
          <label class="lbl">分类（支持三级目录）</label>
          <select v-model="form.category_id" class="input-field">
            <option :value="0">未分类</option>
            <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.path }}</option>
          </select>
        </div>
      </div>

      <div>
        <label class="lbl">摘要</label>
        <textarea v-model="form.summary" rows="2" class="input-field" placeholder="一两句话介绍文章内容" />
      </div>

      <div>
        <label class="lbl">标签</label>
        <div class="flex flex-wrap gap-2 items-center">
          <button
            v-for="t in tags"
            :key="t.id"
            type="button"
            class="badge !px-3 !py-1 border transition-colors"
            :class="
              selectedTags.includes(t.id)
                ? 'bg-[#42b883] text-white border-[#42b883]'
                : 'bg-white dark:bg-[#1e1e1e] text-gray-600 dark:text-gray-300 border-gray-200 dark:border-white/15 hover:border-[#42b883]'
            "
            @click="toggleTag(t.id)"
          >
            # {{ t.name }}
          </button>
          <input
            v-model="newTag"
            class="input-field !w-36 !py-1"
            placeholder="新建标签 ↵"
            @keydown.enter.prevent="createTag"
          />
        </div>
      </div>
    </div>

    <MdEditor
      v-model="form.content"
      language="zh-CN"
      :theme="theme.isDark ? 'dark' : 'light'"
      :preview="isWideEditor"
      preview-theme="default"
      style="height: 560px"
      :toolbars-exclude="['github', 'mermaid', 'katex', 'htmlPreview', 'catalog', 'pageFullscreen', 'sub', 'sup', 'prettier', 'save']"
      @on-upload-img="onUploadImg"
    />

    <div class="flex items-center gap-3 mt-4 sticky bottom-0 bg-[#f6f6f7]/90 dark:bg-[#1a1a1a]/90 backdrop-blur py-3">
      <button class="btn-secondary" :disabled="saving" @click="save('draft')">保存草稿</button>
      <button class="btn-primary" :disabled="saving" @click="save('published')">
        {{ isEdit && originalStatus === 'published' ? '更新并发布' : '发布' }}
      </button>
      <router-link
        v-if="isEdit && form.slug"
        :to="`/posts/${form.slug}`"
        target="_blank"
        class="btn-secondary ml-auto"
      >
        前台预览 ↗
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import {
  adminCreatePost,
  adminCreateTag,
  adminGetPost,
  adminUpdatePost,
  adminUploadImage,
  getCategories,
  getTags,
} from '@/api'
import type { Category, Tag } from '@/api/types'
import { useThemeStore } from '@/stores/theme'
import { makeSlug } from '@/utils/format'
import { toast } from '@/utils/toast'

const route = useRoute()
const router = useRouter()
const theme = useThemeStore()
const isWideEditor = ref(typeof window !== 'undefined' ? window.innerWidth >= 768 : true)
function onResize() {
  isWideEditor.value = window.innerWidth >= 768
}

const id = computed(() => (route.params.id ? Number(route.params.id) : 0))
const isEdit = computed(() => id.value > 0)

const categories = ref<Category[]>([])
const tags = ref<Tag[]>([])
const selectedTags = ref<number[]>([])
const newTag = ref('')
const saving = ref(false)
const originalStatus = ref<'draft' | 'published'>('draft')
const form = reactive({ title: '', slug: '', summary: '', content: '', category_id: 0, is_top: false })

function toggleTag(tid: number) {
  const i = selectedTags.value.indexOf(tid)
  if (i >= 0) selectedTags.value.splice(i, 1)
  else selectedTags.value.push(tid)
}

async function createTag() {
  const name = newTag.value.trim()
  if (!name) return
  const exist = tags.value.find((t) => t.name === name)
  if (exist) {
    if (!selectedTags.value.includes(exist.id)) selectedTags.value.push(exist.id)
  } else {
    const res = await adminCreateTag({ name, slug: makeSlug(name, 'tag') })
    tags.value.push(res.data)
    selectedTags.value.push(res.data.id)
  }
  newTag.value = ''
}

async function onUploadImg(files: File[], callback: (urls: string[]) => void) {
  try {
    const urls = await Promise.all(files.map(async (f) => (await adminUploadImage(f)).data.url))
    callback(urls)
    toast.success('图片上传成功')
  } catch {
    /* error toast shown by interceptor */
  }
}

async function save(status: 'draft' | 'published') {
  if (!form.title.trim()) return toast.error('请填写标题')
  if (!form.content.trim()) return toast.error('请填写正文')
  saving.value = true
  try {
    const payload = {
      title: form.title.trim(),
      slug: form.slug.trim() || makeSlug(form.title, 'post'),
      summary: form.summary,
      content: form.content,
      status,
      category_id: form.category_id || undefined,
      tag_ids: selectedTags.value,
      is_top: form.is_top,
    }
    if (isEdit.value) await adminUpdatePost(id.value, payload)
    else await adminCreatePost(payload)
    toast.success(status === 'published' ? '已发布' : '草稿已保存')
    router.push('/admin/posts')
  } catch {
    /* error toast shown by interceptor */
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  window.addEventListener('resize', onResize, { passive: true })
  try {
    const [c, t] = await Promise.all([getCategories(false), getTags()])
    categories.value = [...c.data].sort((a, b) => a.path.localeCompare(b.path, 'zh-CN'))
    tags.value = t.data
    if (isEdit.value) {
      const p = (await adminGetPost(id.value)).data
      form.title = p.title
      form.slug = p.slug
      form.summary = p.summary
      form.content = p.content || ''
      form.category_id = p.category_id || 0
      form.is_top = !!p.is_top
      selectedTags.value = (p.tags || []).map((x) => x.id)
      originalStatus.value = p.status
    }
  } catch {
    /* handled by interceptor */
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', onResize)
})
</script>

<style scoped>
.lbl {
  @apply block text-sm text-gray-600 dark:text-gray-300 mb-1;
}
</style>

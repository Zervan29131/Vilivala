<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-xl font-bold text-[#213547] dark:text-white">独立页面管理</h1>
      <button class="btn-primary" @click="startNew">+ 新建页面</button>
    </div>

    <div class="grid lg:grid-cols-[240px_1fr] gap-6 items-start">
      <!-- Page List -->
      <div class="card p-2 space-y-1">
        <button
          v-for="p in pages"
          :key="p.slug"
          type="button"
          class="w-full text-left px-3 py-2 rounded-lg text-sm transition-colors flex items-center justify-between"
          :class="
            activeSlug === p.slug
              ? 'bg-[#42b883] text-white'
              : 'text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-white/5'
          "
          @click="selectPage(p.slug)"
        >
          <span class="truncate">{{ p.title }}</span>
          <span class="text-xs opacity-70 font-mono">{{ p.slug }}</span>
        </button>
        <div v-if="!pages.length" class="p-4 text-xs text-gray-400 text-center">暂无页面</div>
      </div>

      <!-- Page Editor -->
      <div class="card p-5 space-y-4 min-w-0">
        <div class="grid sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm text-gray-600 dark:text-gray-300 mb-1">Slug（仅小写字母、数字、连字符）</label>
            <input
              v-model="form.slug"
              class="input-field font-mono"
              :disabled="!isNew"
              placeholder="about"
            />
          </div>
          <div>
            <label class="block text-sm text-gray-600 dark:text-gray-300 mb-1">页面标题</label>
            <input v-model="form.title" class="input-field" placeholder="关于" />
          </div>
        </div>

        <MdEditor
          v-model="form.content"
          language="zh-CN"
          :theme="theme.isDark ? 'dark' : 'light'"
          :preview="isWideEditor"
          preview-theme="default"
          style="height: 480px"
          :toolbars-exclude="['github', 'mermaid', 'katex', 'htmlPreview', 'catalog', 'pageFullscreen', 'sub', 'sup', 'prettier', 'save']"
          @on-upload-img="onUploadImg"
        />

        <div class="flex items-center justify-between">
          <div class="flex gap-2">
            <button class="btn-primary" :disabled="saving" @click="save">保存页面</button>
            <router-link
              v-if="!isNew && form.slug"
              :to="`/pages/${form.slug}`"
              target="_blank"
              class="btn-secondary"
            >
              预览
            </router-link>
          </div>
          <button v-if="!isNew && form.slug" class="text-sm text-red-500 hover:underline" @click="remove">
            删除页面
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { adminDeletePage, adminGetPages, adminPutPage, adminUploadImage, getPageBySlug } from '@/api'
import type { PageSummary } from '@/api/types'
import { useThemeStore } from '@/stores/theme'
import { toast } from '@/utils/toast'

const theme = useThemeStore()
const isWideEditor = ref(typeof window !== 'undefined' ? window.innerWidth >= 768 : true)
function onResize() {
  isWideEditor.value = window.innerWidth >= 768
}
const pages = ref<PageSummary[]>([])
const activeSlug = ref('')
const isNew = ref(false)
const saving = ref(false)
const form = reactive({ slug: '', title: '', content: '' })

async function loadList() {
  try {
    pages.value = (await adminGetPages()).data
  } catch {
    /* handled by interceptor */
  }
}

async function selectPage(slug: string) {
  isNew.value = false
  activeSlug.value = slug
  try {
    const d = (await getPageBySlug(slug)).data
    form.slug = d.slug
    form.title = d.title
    form.content = d.content
  } catch {
    /* handled by interceptor */
  }
}

function startNew() {
  isNew.value = true
  activeSlug.value = ''
  form.slug = ''
  form.title = ''
  form.content = ''
}

async function onUploadImg(files: File[], callback: (urls: string[]) => void) {
  try {
    const urls = await Promise.all(files.map(async (f) => (await adminUploadImage(f)).data.url))
    callback(urls)
    toast.success('图片上传成功')
  } catch {
    /* toast shown */
  }
}

async function save() {
  const slug = form.slug.trim()
  if (!slug || !form.title.trim()) return toast.error('请填写 slug 和标题')
  saving.value = true
  try {
    await adminPutPage(slug, { title: form.title.trim(), content: form.content })
    toast.success('页面已保存')
    await loadList()
    await selectPage(slug)
  } finally {
    saving.value = false
  }
}

async function remove() {
  if (!confirm(`确定删除页面「${form.title}」吗？`)) return
  await adminDeletePage(form.slug)
  toast.success('已删除')
  await loadList()
  if (pages.value.length) await selectPage(pages.value[0].slug)
  else startNew()
}

onMounted(async () => {
  window.addEventListener('resize', onResize, { passive: true })
  await loadList()
  if (pages.value.length) await selectPage(pages.value[0].slug)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', onResize)
})
</script>

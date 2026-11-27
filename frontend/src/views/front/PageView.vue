<template>
  <div class="max-w-4xl mx-auto px-4 py-8">
    <div v-if="loading" class="card h-64 animate-pulse bg-gray-200/60 dark:bg-white/5" />
    <div v-else-if="!page" class="text-center py-20 text-gray-400">
      页面不存在 · <router-link to="/" class="text-[#42b883]">返回首页</router-link>
    </div>
    <article v-else>
      <header class="mb-6">
        <h1 class="text-3xl font-bold text-[#213547] dark:text-white">{{ page.title }}</h1>
        <p class="mt-2 text-xs text-gray-400">最后更新于 {{ formatDate(page.updated_at) }}</p>
      </header>
      <div class="card p-6 sm:p-10">
        <div class="md-body" v-html="html" />
      </div>
    </article>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getPageBySlug } from '@/api'
import type { PageDetail } from '@/api/types'
import { siteConfig } from '@/site.config'
import { renderMarkdown } from '@/utils/markdown'
import { formatDate } from '@/utils/format'

const route = useRoute()
const page = ref<PageDetail | null>(null)
const html = ref('')
const loading = ref(true)

async function load(slug: string) {
  loading.value = true
  page.value = null
  try {
    const res = await getPageBySlug(slug)
    page.value = res.data
    html.value = renderMarkdown(res.data.content || '').html
    document.title = `${res.data.title} - ${siteConfig.title}`
  } catch {
    page.value = null
  } finally {
    loading.value = false
  }
}

watch(() => route.params.slug as string, (s) => s && load(s), { immediate: true })
</script>

<template>
  <div class="max-w-4xl mx-auto px-4 py-8">
    <LoadingState v-if="loading" :rows="1" height-class="h-64" />
    <ErrorState v-else-if="errorMsg" :message="errorMsg" show-home @retry="retryLoad" />
    <EmptyState v-else-if="!page" message="页面不存在">
      <router-link to="/" class="inline-block mt-2 text-sm text-[#42b883] hover:underline">返回首页</router-link>
    </EmptyState>
    <article v-else>
      <header class="mb-6">
        <h1 class="text-3xl font-bold text-[#213547] dark:text-white">{{ page.title }}</h1>
        <p class="mt-2 text-xs text-gray-400">最后更新于 {{ formatDate(page.updated_at) }}</p>
      </header>
      <div class="card p-6 sm:p-10">
        <div ref="bodyRef" class="md-body" v-html="html" @click="onBodyClick" />
      </div>
    </article>

    <ImageLightbox
      :open="lightboxOpen"
      :images="lightboxImages"
      :start-index="lightboxIndex"
      @close="lightboxOpen = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getPageBySlug } from '@/api'
import type { PageDetail } from '@/api/types'
import LoadingState from '@/components/LoadingState.vue'
import ErrorState from '@/components/ErrorState.vue'
import EmptyState from '@/components/EmptyState.vue'
import ImageLightbox, { type LightboxImage } from '@/components/ImageLightbox.vue'
import { renderMarkdown } from '@/utils/markdown'
import { formatDate } from '@/utils/format'
import { updatePageMeta } from '@/utils/seo'

const route = useRoute()
const bodyRef = ref<HTMLElement | null>(null)
const page = ref<PageDetail | null>(null)
const html = ref('')
const loading = ref(true)
const errorMsg = ref('')

const lightboxOpen = ref(false)
const lightboxImages = ref<LightboxImage[]>([])
const lightboxIndex = ref(0)

async function load(slug: string) {
  loading.value = true
  errorMsg.value = ''
  page.value = null
  try {
    const res = await getPageBySlug(slug)
    page.value = res.data
    html.value = renderMarkdown(res.data.content || '', { stripLeadingTitle: res.data.title }).html
    updatePageMeta({
      title: res.data.title,
      description: res.data.title,
      path: `/pages/${res.data.slug}`,
    })
  } catch (err) {
    const msg = err instanceof Error ? err.message : ''
    if (msg.includes('不存在')) page.value = null
    else errorMsg.value = msg || '加载页面失败'
  } finally {
    loading.value = false
  }
}

function retryLoad() {
  const s = route.params.slug as string
  if (s) load(s)
}

function onBodyClick(e: MouseEvent) {
  const img = (e.target as HTMLElement).closest('img') as HTMLImageElement | null
  if (img && bodyRef.value) {
    const imgs = Array.from(bodyRef.value.querySelectorAll('img'))
    const idx = imgs.indexOf(img)
    if (idx >= 0) {
      lightboxImages.value = imgs.map((el) => ({ src: el.src, alt: el.alt || '' }))
      lightboxIndex.value = idx
      lightboxOpen.value = true
    }
  }
}

watch(() => route.params.slug as string, (s) => s && load(s), { immediate: true })
</script>

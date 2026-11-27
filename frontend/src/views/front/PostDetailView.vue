<template>
  <div class="max-w-7xl mx-auto px-4 py-8">
    <div v-if="loading" class="max-w-3xl mx-auto card h-72 animate-pulse bg-gray-200/60 dark:bg-white/5" />
    <div v-else-if="!post" class="text-center py-20 text-gray-400">
      文章不存在 · <router-link to="/" class="text-[#42b883]">返回首页</router-link>
    </div>

    <!-- Vue-docs style: Left sticky TOC + Right main content -->
    <div v-else class="grid grid-cols-1 lg:grid-cols-[240px_1fr] gap-8 items-start">
      <!-- Left TOC -->
      <aside class="hidden lg:block sticky top-20">
        <div v-if="toc.length" class="card p-4 text-sm max-h-[calc(100vh-7rem)] overflow-y-auto">
          <h3 class="font-semibold text-xs uppercase tracking-wider text-gray-400 mb-3">文章目录</h3>
          <ul class="space-y-1 border-l border-gray-200 dark:border-white/10">
            <li
              v-for="item in toc"
              :key="item.id"
              :style="{ paddingLeft: `${(item.level - 1) * 10 + 10}px` }"
            >
              <a
                :href="`#${item.id}`"
                class="block -ml-px border-l-2 pl-2.5 py-1 text-[13px] leading-snug transition-colors"
                :class="
                  active === item.id
                    ? 'border-[#42b883] text-[#42b883] dark:text-[#4fc08d] font-medium'
                    : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-[#213547] dark:hover:text-white'
                "
                @click.prevent="jump(item.id)"
              >
                {{ item.text }}
              </a>
            </li>
          </ul>
        </div>
        <div v-else class="card p-4 text-xs text-gray-400">本文无章节目录</div>
      </aside>

      <!-- Main Article -->
      <article class="min-w-0 max-w-4xl">
        <!-- Breadcrumb -->
        <div class="flex items-center flex-wrap gap-1.5 text-xs text-gray-400 mb-3">
          <router-link to="/" class="hover:text-[#42b883]">首页</router-link>
          <template v-if="post.category?.path?.length">
            <template v-for="node in post.category.path" :key="node.id">
              <span>/</span>
              <router-link :to="`/categories/${node.slug}`" class="hover:text-[#42b883]">
                {{ node.name }}
              </router-link>
            </template>
          </template>
        </div>

        <header class="mb-6">
          <div class="flex items-center gap-2 mb-2" v-if="post.is_top">
            <span class="badge bg-[#42b883]/15 text-[#42b883]">置顶文章</span>
          </div>
          <h1 class="text-2xl sm:text-4xl font-bold text-[#213547] dark:text-white leading-tight">
            {{ post.title }}
          </h1>
          <div class="mt-3 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-gray-400">
            <span>发布于 {{ formatDate(post.created_at) }}</span>
            <span>·</span>
            <span>{{ post.word_count }} 字</span>
            <span>·</span>
            <span>预计阅读 {{ post.reading_time || 1 }} 分钟</span>
            <span>·</span>
            <span>{{ post.view_count }} 阅读</span>
          </div>
          <div v-if="post.tags?.length" class="mt-4 flex flex-wrap gap-1.5">
            <router-link
              v-for="t in post.tags"
              :key="t.id"
              :to="`/tags/${t.slug}`"
              class="badge bg-gray-100 dark:bg-white/10 text-gray-600 dark:text-gray-300 hover:bg-[#42b883] hover:text-white transition-colors"
            >
              # {{ t.name }}
            </router-link>
          </div>
        </header>

        <div class="card p-6 sm:p-10">
          <div class="md-body" v-html="html" @click="onBodyClick" />
        </div>

        <!-- Prev / Next Navigation -->
        <nav class="mt-6 grid grid-cols-1 sm:grid-cols-2 gap-4">
          <router-link
            v-if="post.prev"
            :to="`/posts/${post.prev.slug}`"
            class="card p-4 hover:border-[#42b883] transition-colors group"
          >
            <div class="text-xs text-gray-400 mb-1">← 上一篇</div>
            <div class="text-sm font-medium text-[#213547] dark:text-gray-100 group-hover:text-[#42b883] truncate">
              {{ post.prev.title }}
            </div>
          </router-link>
          <div v-else />

          <router-link
            v-if="post.next"
            :to="`/posts/${post.next.slug}`"
            class="card p-4 text-right hover:border-[#42b883] transition-colors group"
          >
            <div class="text-xs text-gray-400 mb-1">下一篇 →</div>
            <div class="text-sm font-medium text-[#213547] dark:text-gray-100 group-hover:text-[#42b883] truncate">
              {{ post.next.title }}
            </div>
          </router-link>
        </nav>
      </article>
    </div>
  </div>
</template>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getPostBySlug } from '@/api'
import type { Post } from '@/api/types'
import { siteConfig } from '@/site.config'
import { renderMarkdown, type TocItem } from '@/utils/markdown'
import { formatDate } from '@/utils/format'
import { toast } from '@/utils/toast'

const route = useRoute()
const post = ref<Post | null>(null)
const html = ref('')
const toc = ref<TocItem[]>([])
const loading = ref(true)
const active = ref('')

let observer: IntersectionObserver | null = null

function observeHeadings() {
  observer?.disconnect()
  observer = new IntersectionObserver(
    (entries) => {
      for (const e of entries) if (e.isIntersecting) active.value = e.target.id
    },
    { rootMargin: '-72px 0px -70% 0px' }
  )
  document.querySelectorAll('.md-body h1[id], .md-body h2[id], .md-body h3[id]').forEach((el) => observer!.observe(el))
}

async function load(slug: string) {
  loading.value = true
  post.value = null
  try {
    const res = await getPostBySlug(slug)
    post.value = res.data
    const r = renderMarkdown(res.data.content || '')
    html.value = r.html
    toc.value = r.toc
    active.value = r.toc[0]?.id ?? ''
    document.title = `${res.data.title} - ${siteConfig.title}`
  } catch {
    post.value = null
  } finally {
    loading.value = false
  }
  await nextTick()
  observeHeadings()
}

function jump(id: string) {
  document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  active.value = id
}

async function onBodyClick(e: MouseEvent) {
  const btn = (e.target as HTMLElement).closest('.copy-btn') as HTMLButtonElement | null
  if (!btn) return
  const code = btn.closest('.code-block')?.querySelector('code')?.textContent ?? ''
  try {
    await navigator.clipboard.writeText(code)
    btn.textContent = '已复制'
    setTimeout(() => (btn.textContent = '复制'), 1500)
  } catch {
    toast.error('复制失败')
  }
}

watch(() => route.params.slug as string, (s) => s && load(s), { immediate: true })
onBeforeUnmount(() => observer?.disconnect())
</script>

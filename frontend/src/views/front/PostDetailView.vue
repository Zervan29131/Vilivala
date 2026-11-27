<template>
  <div class="max-w-7xl mx-auto px-4 py-8">
    <LoadingState v-if="loading" :rows="1" height-class="h-72 max-w-3xl mx-auto" />
    <ErrorState v-else-if="errorMsg" :message="errorMsg" show-home @retry="retryLoad" />
    <EmptyState v-else-if="!post" message="文章不存在或已下架">
      <router-link to="/" class="inline-block mt-2 text-sm text-[#42b883] hover:underline">返回首页</router-link>
    </EmptyState>

    <!-- Vue-docs style: Left sticky TOC + Right main content -->
    <div v-else class="grid grid-cols-1 lg:grid-cols-[240px_1fr] gap-8 items-start">
      <!-- Desktop Left Sticky TOC -->
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
        <!-- Admin Draft Preview Banner -->
        <div
          v-if="post.is_preview"
          class="mb-4 px-4 py-2.5 rounded-xl border border-amber-500/40 bg-amber-500/10 text-amber-700 dark:text-amber-300 text-xs flex items-center justify-between gap-3 flex-wrap"
        >
          <span>📝 当前为管理员<strong>草稿预览</strong>模式，该文章尚未公开发布（不计入阅读量）。</span>
          <router-link :to="`/admin/posts/${post.id}/edit`" class="underline font-medium shrink-0">
            返回编辑器 →
          </router-link>
        </div>

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
          <div v-if="post.is_top" class="flex items-center gap-2 mb-2">
            <span class="badge bg-[#42b883]/15 text-[#42b883]">置顶文章</span>
          </div>
          <h1 class="text-2xl sm:text-4xl font-bold text-[#213547] dark:text-white leading-tight">
            {{ post.title }}
          </h1>
          <div class="mt-3 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-gray-400">
            <span>发布于 {{ formatDate(post.created_at) }}</span>
            <template v-if="showUpdatedDate">
              <span>·</span>
              <span>更新于 {{ formatDate(post.updated_at) }}</span>
            </template>
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

        <!-- Mobile / Tablet Collapsible TOC (< 1024px) -->
        <div v-if="toc.length" class="lg:hidden mb-4 card overflow-hidden">
          <button
            type="button"
            class="w-full px-4 py-3 flex items-center justify-between text-xs font-semibold text-[#213547] dark:text-gray-200"
            :aria-expanded="mobileTocOpen"
            @click="mobileTocOpen = !mobileTocOpen"
          >
            <span>📑 本页目录（共 {{ toc.length }} 节）</span>
            <span class="text-gray-400 transition-transform" :class="{ 'rotate-180': mobileTocOpen }">▾</span>
          </button>
          <ul
            v-if="mobileTocOpen"
            class="px-4 pb-3 pt-1 border-t border-gray-100 dark:border-white/10 max-h-64 overflow-y-auto space-y-1"
          >
            <li
              v-for="item in toc"
              :key="item.id"
              :style="{ paddingLeft: `${(item.level - 1) * 12}px` }"
            >
              <a
                :href="`#${item.id}`"
                class="block py-1 text-xs transition-colors"
                :class="
                  active === item.id
                    ? 'text-[#42b883] dark:text-[#4fc08d] font-medium'
                    : 'text-gray-600 dark:text-gray-400'
                "
                @click.prevent="jumpMobile(item.id)"
              >
                {{ item.text }}
              </a>
            </li>
          </ul>
        </div>

        <!-- Stale Technical Article Gentle Notice -->
        <div
          v-if="staleYears >= 2 && isTechPost"
          class="mb-4 px-4 py-2.5 rounded-xl border border-amber-500/25 bg-amber-500/[0.07] text-xs text-amber-800 dark:text-amber-300/90 leading-relaxed"
        >
          ⏳ 本文发布于 {{ staleYears }} 年前（{{ formatDate(post.created_at) }}），文中涉及的部分软件版本、环境或命令可能已随时间更新，阅读时请注意甄别。
        </div>

        <div class="card p-4 sm:p-8 lg:p-10">
          <div ref="bodyRef" class="md-body" v-html="html" @click="onBodyClick" />

          <!-- Reward & Share Actions -->
          <PostActions :title="post.title" :slug="post.slug" />

          <!-- Original Copyright Box -->
          <PostCopyright :slug="post.slug" :created-at="post.created_at" />
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
          <div v-else class="hidden sm:block" />

          <router-link
            v-if="post.next"
            :to="`/posts/${post.next.slug}`"
            class="card p-4 text-left sm:text-right hover:border-[#42b883] transition-colors group"
          >
            <div class="text-xs text-gray-400 mb-1">下一篇 →</div>
            <div class="text-sm font-medium text-[#213547] dark:text-gray-100 group-hover:text-[#42b883] truncate">
              {{ post.next.title }}
            </div>
          </router-link>
        </nav>

        <!-- Related Posts -->
        <section v-if="post.related?.length" class="mt-8">
          <h2 class="text-base font-bold text-[#213547] dark:text-white mb-3">相关推荐</h2>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <router-link
              v-for="r in post.related"
              :key="r.id"
              :to="`/posts/${r.slug}`"
              class="card p-4 hover:border-[#42b883] transition-colors group flex flex-col justify-between gap-2"
            >
              <div>
                <div class="text-sm font-semibold text-[#213547] dark:text-gray-100 group-hover:text-[#42b883] line-clamp-1">
                  {{ r.title }}
                </div>
                <p v-if="r.summary" class="mt-1 text-xs text-gray-500 dark:text-gray-400 line-clamp-2">
                  {{ r.summary }}
                </p>
              </div>
              <div class="text-[11px] text-gray-400 font-mono">
                {{ formatDate(r.created_at) }} · 约 {{ r.reading_time || 1 }} 分钟
              </div>
            </router-link>
          </div>
        </section>
      </article>
    </div>

    <ImageLightbox
      :open="lightboxOpen"
      :images="lightboxImages"
      :start-index="lightboxIndex"
      @close="lightboxOpen = false"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getPostBySlug } from '@/api'
import type { Post } from '@/api/types'
import LoadingState from '@/components/LoadingState.vue'
import ErrorState from '@/components/ErrorState.vue'
import EmptyState from '@/components/EmptyState.vue'
import ImageLightbox, { type LightboxImage } from '@/components/ImageLightbox.vue'
import PostActions from '@/components/PostActions.vue'
import PostCopyright from '@/components/PostCopyright.vue'
import { copyText } from '@/utils/clipboard'
import { attachMarkdownImageFallbacks, renderMarkdown, type TocItem } from '@/utils/markdown'
import { formatDate } from '@/utils/format'
import { updatePageMeta } from '@/utils/seo'
import { toast } from '@/utils/toast'

const route = useRoute()
const bodyRef = ref<HTMLElement | null>(null)
const post = ref<Post | null>(null)
const html = ref('')
const toc = ref<TocItem[]>([])
const loading = ref(true)
const errorMsg = ref('')
const active = ref('')
const mobileTocOpen = ref(false)

const lightboxOpen = ref(false)
const lightboxImages = ref<LightboxImage[]>([])
const lightboxIndex = ref(0)

const showUpdatedDate = computed(() => {
  if (!post.value?.updated_at || !post.value?.created_at) return false
  return formatDate(post.value.updated_at) !== formatDate(post.value.created_at)
})

const isTechPost = computed(() => {
  const topSlug = post.value?.category?.path?.[0]?.slug?.toLowerCase()
  return topSlug === 'code'
})

const staleYears = computed(() => {
  if (!post.value?.created_at) return 0
  const created = new Date(post.value.created_at).getTime()
  if (isNaN(created)) return 0
  return Math.floor((Date.now() - created) / (365.25 * 24 * 3600 * 1000))
})

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

function scrollToInitialHash() {
  const rawHash = window.location.hash.slice(1)
  if (!rawHash) return
  try {
    const id = decodeURIComponent(rawHash)
    const target = document.getElementById(id)
    if (target) {
      target.scrollIntoView({ behavior: 'auto', block: 'start' })
      active.value = id
    }
  } catch {
    /* ignore malformed hash */
  }
}

async function load(slug: string) {
  loading.value = true
  errorMsg.value = ''
  post.value = null
  mobileTocOpen.value = false
  try {
    const res = await getPostBySlug(slug)
    post.value = res.data
    const r = renderMarkdown(res.data.content || '', { stripLeadingTitle: res.data.title })
    html.value = r.html
    toc.value = r.toc
    active.value = r.toc[0]?.id ?? ''
    updatePageMeta({
      title: res.data.title,
      description: res.data.summary || res.data.title,
      path: `/posts/${res.data.slug}`,
      type: 'article',
      publishedAt: res.data.created_at,
      updatedAt: res.data.updated_at,
      tags: res.data.tags?.map((t) => t.name),
    })
  } catch (err) {
    const msg = err instanceof Error ? err.message : ''
    if (msg.includes('不存在')) {
      post.value = null
    } else {
      errorMsg.value = msg || '加载文章失败，请重试'
    }
  } finally {
    loading.value = false
  }
  await nextTick()
  attachMarkdownImageFallbacks(bodyRef.value)
  observeHeadings()
  scrollToInitialHash()
}

function retryLoad() {
  const s = route.params.slug as string
  if (s) load(s)
}

function jump(id: string) {
  const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  document.getElementById(id)?.scrollIntoView({ behavior: reduced ? 'auto' : 'smooth', block: 'start' })
  window.history.replaceState(null, '', `#${encodeURIComponent(id)}`)
  active.value = id
}

function jumpMobile(id: string) {
  mobileTocOpen.value = false
  jump(id)
}

async function onBodyClick(e: MouseEvent) {
  const target = e.target as HTMLElement

  // Heading `#` anchor copy & jump
  const anchor = target.closest('.header-anchor') as HTMLAnchorElement | null
  if (anchor) {
    e.preventDefault()
    const id = anchor.dataset.anchor || ''
    if (id) {
      jump(id)
      const full = `${window.location.origin}${window.location.pathname}#${encodeURIComponent(id)}`
      if (await copyText(full)) {
        toast.success('章节链接已复制')
      }
    }
    return
  }

  // Code copy button
  const btn = target.closest('.copy-btn') as HTMLButtonElement | null
  if (btn) {
    const code = btn.closest('.code-block')?.querySelector('code')?.textContent ?? ''
    if (await copyText(code)) {
      btn.textContent = '已复制'
      setTimeout(() => (btn.textContent = '复制'), 1500)
    } else {
      toast.error('复制失败')
    }
    return
  }

  // Image lightbox (exclude friend avatars)
  const img = target.closest('img[data-lightbox="true"]') as HTMLImageElement | null
  if (img && bodyRef.value) {
    const imgs = Array.from(bodyRef.value.querySelectorAll('img[data-lightbox="true"]')) as HTMLImageElement[]
    const idx = imgs.indexOf(img)
    if (idx >= 0) {
      lightboxImages.value = imgs.map((el) => ({ src: el.src, alt: el.alt || '' }))
      lightboxIndex.value = idx
      lightboxOpen.value = true
    }
  }
}

watch(() => route.params.slug as string, (s) => s && load(s), { immediate: true })
onBeforeUnmount(() => observer?.disconnect())
</script>

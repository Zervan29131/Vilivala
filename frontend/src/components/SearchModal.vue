<template>
  <teleport to="body">
    <transition name="modal-fade">
      <div
        v-if="open"
        class="fixed inset-0 z-[60] bg-black/50 backdrop-blur-sm flex items-start justify-center pt-[10vh] px-4"
        role="dialog"
        aria-modal="true"
        aria-label="全文搜索"
        @click.self="close"
      >
        <div class="w-full max-w-2xl card shadow-2xl overflow-hidden border border-gray-200 dark:border-white/15">
          <!-- Input -->
          <div class="flex items-center gap-3 px-4 py-3 border-b border-gray-200/80 dark:border-white/10">
            <svg class="w-5 h-5 text-[#42b883] shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-4.35-4.35m0 0A7.5 7.5 0 1010.5 18a7.5 7.5 0 006.15-3.35z" />
            </svg>
            <input
              ref="inputRef"
              v-model="query"
              type="text"
              class="flex-1 bg-transparent text-sm sm:text-base text-[#213547] dark:text-white placeholder:text-gray-400 focus:outline-none"
              placeholder="全文搜索文章标题、摘要或正文…（↑↓ 选择，Enter 打开，Esc 关闭）"
              @keydown.down.prevent="move(1)"
              @keydown.up.prevent="move(-1)"
              @keydown.enter.prevent="openSelected"
            />
            <button
              v-if="query"
              type="button"
              class="text-xs text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 px-1.5"
              @click="query = ''"
            >
              清空
            </button>
            <kbd class="hidden sm:inline-block px-1.5 py-0.5 text-[11px] font-mono rounded bg-gray-100 dark:bg-white/10 text-gray-500 dark:text-gray-400">
              ESC
            </kbd>
          </div>

          <!-- Results -->
          <div class="max-h-[60vh] overflow-y-auto p-2">
            <!-- Empty query: quick jump links -->
            <div v-if="!query.trim()" class="p-2">
              <div class="px-2 py-1 text-[11px] font-semibold uppercase tracking-wider text-gray-400">
                快速跳转
              </div>
              <button
                v-for="(item, idx) in quickLinks"
                :key="item.to"
                type="button"
                class="w-full text-left px-3 py-2.5 rounded-lg flex items-center justify-between text-sm transition-colors"
                :class="
                  activeIdx === idx
                    ? 'bg-[#42b883]/15 text-[#42b883] dark:text-[#4fc08d]'
                    : 'text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-white/5'
                "
                @mouseenter="activeIdx = idx"
                @click="go(item.to)"
              >
                <span class="flex items-center gap-2.5">
                  <span>{{ item.icon }}</span>
                  <span class="font-medium">{{ item.label }}</span>
                </span>
                <span class="text-xs text-gray-400 font-mono">{{ item.hint }}</span>
              </button>
            </div>

            <!-- Loading -->
            <div v-else-if="loading" class="p-6 text-center text-sm text-gray-400">正在检索…</div>

            <!-- Hits -->
            <div v-else-if="results.length" class="space-y-1">
              <div class="px-3 py-1 flex items-center justify-between text-xs text-gray-400">
                <span>找到 {{ total }} 篇匹配文章</span>
                <button
                  type="button"
                  class="text-[#42b883] hover:underline"
                  @click="goFullSearch"
                >
                  在首页查看全部 →
                </button>
              </div>
              <button
                v-for="(item, idx) in results"
                :key="item.id"
                type="button"
                class="w-full text-left px-3 py-2.5 rounded-lg transition-colors block"
                :class="
                  activeIdx === idx
                    ? 'bg-[#42b883]/15'
                    : 'hover:bg-gray-100/80 dark:hover:bg-white/5'
                "
                @mouseenter="activeIdx = idx"
                @click="go(`/posts/${item.slug}`)"
              >
                <div class="flex items-center justify-between gap-2">
                  <span class="font-medium text-sm text-[#213547] dark:text-white truncate">
                    <template v-for="(seg, si) in highlightParts(item.title, query)" :key="si">
                      <mark v-if="seg.hit" class="bg-[#42b883]/25 text-[#213547] dark:text-[#4fc08d] rounded px-0.5">{{ seg.text }}</mark>
                      <template v-else>{{ seg.text }}</template>
                    </template>
                  </span>
                  <span class="text-[11px] text-gray-400 shrink-0 font-mono">
                    {{ item.category_name ? `${item.category_name} · ` : '' }}{{ formatDate(item.created_at) }}
                  </span>
                </div>
                <p v-if="item.snippet" class="mt-1 text-xs text-gray-500 dark:text-gray-400 line-clamp-2 leading-relaxed">
                  <template v-for="(seg, si) in highlightParts(item.snippet, query)" :key="si">
                    <mark v-if="seg.hit" class="bg-[#42b883]/25 text-[#213547] dark:text-[#4fc08d] rounded px-0.5">{{ seg.text }}</mark>
                    <template v-else>{{ seg.text }}</template>
                  </template>
                </p>
              </button>
            </div>

            <!-- No result -->
            <div v-else class="p-8 text-center text-sm text-gray-400">
              未找到包含 “<span class="text-[#42b883]">{{ query }}</span>” 的文章
            </div>
          </div>

          <!-- Footer hints -->
          <div class="px-4 py-2 border-t border-gray-200/80 dark:border-white/10 bg-gray-50/70 dark:bg-white/[0.02] flex items-center justify-between text-[11px] text-gray-400">
            <span>支持中英文全文检索</span>
            <div class="hidden sm:flex items-center gap-3">
              <span>↑↓ 移动</span>
              <span>↵ 打开</span>
              <span>Esc 关闭</span>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { searchPosts } from '@/api'
import type { SearchResultItem } from '@/api/types'
import { formatDate } from '@/utils/format'

const props = defineProps<{ open: boolean }>()
const emit = defineEmits<{ close: [] }>()

const router = useRouter()
const inputRef = ref<HTMLInputElement | null>(null)
const query = ref('')
const results = ref<SearchResultItem[]>([])
const total = ref(0)
const loading = ref(false)
const activeIdx = ref(0)

const quickLinks = [
  { icon: '>_', label: '终端起始页（浏览器风格）', hint: '/terminal', to: '/terminal' },
  { icon: '🏠', label: '返回首页', hint: '/', to: '/' },
  { icon: '🗓️', label: '时间归档', hint: '/archives', to: '/archives' },
  { icon: '🗂️', label: '分类与标签', hint: '/categories', to: '/categories' },
  { icon: '👤', label: '关于本站', hint: '/pages/about', to: '/pages/about' },
]

let timer: number | undefined
let reqSeq = 0

function highlightParts(text: string, kw: string): { text: string; hit: boolean }[] {
  const q = kw.trim()
  if (!q) return [{ text, hit: false }]
  const escaped = q.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  const re = new RegExp(`(${escaped})`, 'gi')
  return text.split(re).filter(Boolean).map((part) => ({
    text: part,
    hit: part.toLowerCase() === q.toLowerCase(),
  }))
}

async function runSearch(q: string) {
  const trimmed = q.trim()
  if (!trimmed) {
    results.value = []
    total.value = 0
    activeIdx.value = 0
    return
  }
  const seq = ++reqSeq
  loading.value = true
  try {
    const res = await searchPosts(trimmed, 12)
    if (seq !== reqSeq) return
    results.value = res.data.list
    total.value = res.data.total
    activeIdx.value = 0
  } catch {
    if (seq !== reqSeq) return
    results.value = []
    total.value = 0
  } finally {
    if (seq === reqSeq) loading.value = false
  }
}

watch(query, (val) => {
  window.clearTimeout(timer)
  if (!val.trim()) {
    results.value = []
    total.value = 0
    activeIdx.value = 0
    return
  }
  timer = window.setTimeout(() => runSearch(val), 180)
})

watch(
  () => props.open,
  async (val) => {
    if (val) {
      activeIdx.value = 0
      await nextTick()
      inputRef.value?.focus()
      inputRef.value?.select()
    }
  }
)

function move(delta: number) {
  const len = query.value.trim() ? results.value.length : quickLinks.length
  if (!len) return
  activeIdx.value = (activeIdx.value + delta + len) % len
}

function go(path: string) {
  close()
  router.push(path)
}

function goFullSearch() {
  const q = query.value.trim()
  close()
  router.push({ path: '/', query: { keyword: q || undefined } })
}

function openSelected() {
  if (!query.value.trim()) {
    const item = quickLinks[activeIdx.value]
    if (item) go(item.to)
    return
  }
  const hit = results.value[activeIdx.value]
  if (hit) go(`/posts/${hit.slug}`)
  else goFullSearch()
}

function close() {
  emit('close')
}

function onKeydown(e: KeyboardEvent) {
  if (props.open && e.key === 'Escape') {
    e.preventDefault()
    close()
  }
}

onMounted(() => window.addEventListener('keydown', onKeydown))
onBeforeUnmount(() => {
  window.clearTimeout(timer)
  window.removeEventListener('keydown', onKeydown)
})
</script>

<style scoped>
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.16s ease;
}
.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}
</style>

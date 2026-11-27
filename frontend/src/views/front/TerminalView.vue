<template>
  <div class="min-h-screen flex flex-col text-gray-100 transition-colors duration-300" :class="bgClass">
    <!-- Browser-style Top Chrome -->
    <header class="sticky top-0 z-30 border-b border-white/10 bg-black/40 backdrop-blur-md">
      <div class="max-w-6xl mx-auto px-4 h-13 py-2 flex items-center justify-between gap-3">
        <!-- Left: traffic lights + brand -->
        <div class="flex items-center gap-3 min-w-0">
          <div class="hidden sm:flex items-center gap-1.5" aria-hidden="true">
            <span class="w-3 h-3 rounded-full bg-[#ff5f56]" />
            <span class="w-3 h-3 rounded-full bg-[#ffbd2e]" />
            <span class="w-3 h-3 rounded-full bg-[#27c93f]" />
          </div>
          <router-link to="/" class="flex items-center gap-2 hover:text-[#4fc08d] transition-colors shrink-0">
            <img :src="siteConfig.logo" :alt="siteConfig.title" class="w-6 h-6 rounded-full object-cover" />
            <span class="font-semibold text-sm text-white">{{ siteConfig.title }}</span>
          </router-link>
          <span class="hidden md:inline-flex items-center px-2 py-0.5 rounded bg-[#42b883]/20 text-[#4fc08d] text-[11px] font-mono">
            &gt;_ Start Page
          </span>
        </div>

        <!-- Center: Clock -->
        <div class="hidden sm:flex items-center gap-2 font-mono text-xs text-gray-300">
          <span>{{ clockTime }}</span>
          <span class="text-white/25">|</span>
          <span class="text-gray-400">{{ clockDate }}</span>
        </div>

        <!-- Right: Background switcher & blog links -->
        <div class="flex items-center gap-1.5 text-xs">
          <button
            v-for="preset in bgPresets"
            :key="preset.key"
            type="button"
            class="px-2 py-1 rounded transition-colors"
            :class="
              bgMode === preset.key
                ? 'bg-[#42b883] text-white'
                : 'text-gray-400 hover:text-white hover:bg-white/10'
            "
            @click="setBg(preset.key)"
          >
            {{ preset.label }}
          </button>
          <span class="text-white/15 mx-0.5">|</span>
          <router-link
            to="/"
            class="px-2.5 py-1 rounded bg-white/10 hover:bg-[#42b883] text-white transition-colors"
          >
            返回博客
          </router-link>
        </div>
      </div>
    </header>

    <!-- Main Container -->
    <div class="flex-1 max-w-6xl w-full mx-auto px-4 py-6 flex flex-col gap-5">
      <!-- Quick Launch Bookmarks Dock (Browser Start Page Feel) -->
      <section class="rounded-xl border border-white/10 bg-black/30 backdrop-blur-sm p-4">
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-2 text-xs text-gray-400">
            <span class="text-[#4fc08d] font-mono">★</span>
            <span>快速导航</span>
            <span class="text-white/25">·</span>
            <span class="hidden sm:inline font-mono text-[11px]">支持命令 <code class="text-[#4fc08d]">bm add 名称 网址</code> 自定义</span>
          </div>
          <div class="flex items-center gap-2 text-xs">
            <router-link to="/archives" class="text-gray-400 hover:text-[#4fc08d] transition-colors">归档</router-link>
            <router-link to="/categories" class="text-gray-400 hover:text-[#4fc08d] transition-colors">分类</router-link>
            <a href="/rss.xml" target="_blank" rel="noopener" class="text-gray-400 hover:text-[#4fc08d] transition-colors">RSS</a>
          </div>
        </div>
        <div class="flex flex-wrap gap-2">
          <template v-for="bm in bookmarks" :key="bm.name">
            <router-link
              v-if="bm.url.startsWith('/')"
              :to="bm.url"
              class="px-3 py-1.5 rounded-lg bg-white/[0.06] hover:bg-[#42b883]/20 border border-white/10 hover:border-[#42b883]/50 text-xs text-gray-200 hover:text-[#4fc08d] transition-all flex items-center gap-1.5"
            >
              <span>{{ bm.icon || '🔗' }}</span>
              <span>{{ bm.name }}</span>
            </router-link>
            <a
              v-else
              :href="bm.url"
              target="_blank"
              rel="noopener noreferrer"
              class="px-3 py-1.5 rounded-lg bg-white/[0.06] hover:bg-[#42b883]/20 border border-white/10 hover:border-[#42b883]/50 text-xs text-gray-200 hover:text-[#4fc08d] transition-all flex items-center gap-1.5"
            >
              <span>{{ bm.icon || '🌐' }}</span>
              <span>{{ bm.name }}</span>
            </a>
          </template>
        </div>
      </section>

      <!-- Interactive Terminal Window -->
      <section
        class="flex-1 min-h-[520px] rounded-xl border border-white/15 bg-[#121417]/90 shadow-2xl flex flex-col overflow-hidden"
        @click="focusInput"
      >
        <!-- Quick Command Chips -->
        <div class="px-4 py-2.5 border-b border-white/10 bg-white/[0.03] flex items-center justify-between gap-2 flex-wrap">
          <div class="flex items-center flex-wrap gap-1.5 text-xs font-mono">
            <span class="text-gray-500 mr-1">快捷指令:</span>
            <button
              v-for="chip in starterChips"
              :key="chip"
              type="button"
              class="px-2 py-0.5 rounded bg-white/[0.07] hover:bg-[#42b883]/25 text-gray-300 hover:text-[#4fc08d] transition-colors"
              @click.stop="runPreset(chip)"
            >
              {{ chip }}
            </button>
          </div>
          <span class="text-[11px] font-mono text-gray-500">Tab 补全 · ↑↓ 历史 · Ctrl+L 清屏</span>
        </div>

        <!-- Terminal Output Stream -->
        <div ref="outputRef" class="flex-1 overflow-y-auto p-4 sm:p-5 font-mono text-[13px] leading-relaxed space-y-3">
          <div v-for="entry in entries" :key="entry.id" class="space-y-1.5">
            <!-- Command line echo -->
            <div v-if="entry.command" class="flex items-center gap-2">
              <span class="text-[#4fc08d] shrink-0">zervan@vilivala</span>
              <span class="text-gray-500">:</span>
              <span class="text-sky-400">~</span>
              <span class="text-gray-400">$</span>
              <span class="text-white break-all">{{ entry.command }}</span>
            </div>

            <!-- Output lines -->
            <div class="pl-0 sm:pl-2 space-y-1">
              <div
                v-for="(line, li) in entry.lines"
                :key="li"
                class="break-words whitespace-pre-wrap"
                :class="lineToneClass(line.tone)"
              >
                <router-link
                  v-if="line.to"
                  :to="line.to"
                  class="text-[#4fc08d] hover:underline"
                  @click.stop
                >
                  {{ line.text }}
                </router-link>
                <a
                  v-else-if="line.href"
                  :href="line.href"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="text-sky-400 hover:underline"
                  @click.stop
                >
                  {{ line.text }}
                </a>
                <span v-else>{{ line.text }}</span>
                <span v-if="line.meta" class="ml-2 text-xs text-gray-500">{{ line.meta }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Prompt Input Bar -->
        <form
          class="px-4 sm:px-5 py-3 border-t border-white/10 bg-black/30 flex items-center gap-2 font-mono text-sm"
          @submit.prevent="onSubmit"
        >
          <span class="text-[#4fc08d] shrink-0">zervan@vilivala</span>
          <span class="text-gray-500">:</span>
          <span class="text-sky-400">~</span>
          <span class="text-gray-400">$</span>
          <div class="relative flex-1 min-w-0">
            <input
              ref="inputRef"
              v-model="input"
              type="text"
              autocomplete="off"
              autocapitalize="off"
              spellcheck="false"
              class="w-full bg-transparent text-white placeholder:text-gray-500 focus:outline-none"
              placeholder="输入 help 查看命令，或输入 s <关键词> 搜索博客 / g <关键词> 搜索全网…"
              @keydown.up.prevent="browseHistory(-1)"
              @keydown.down.prevent="browseHistory(1)"
              @keydown.tab.prevent="acceptCompletion"
              @keydown.ctrl.l.prevent="clearEntries"
            />
            <span
              v-if="ghostHint"
              class="pointer-events-none absolute left-0 top-0 text-gray-500/70 select-none whitespace-pre"
            >
              <span class="invisible">{{ input }}</span>{{ ghostHint }}
            </span>
          </div>
          <button
            type="submit"
            class="px-3 py-1 rounded bg-[#42b883] hover:bg-[#33a06f] text-white text-xs font-sans transition-colors shrink-0"
          >
            执行
          </button>
        </form>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getCategories, getPages, getPostBySlug, getPublishedPosts, searchPosts } from '@/api'
import { siteConfig } from '@/site.config'
import { useThemeStore, type ThemeMode } from '@/stores/theme'
import { formatDate } from '@/utils/format'
import { updatePageMeta } from '@/utils/seo'

type Tone = 'normal' | 'accent' | 'muted' | 'error' | 'warn'
type BgMode = 'grid' | 'emerald' | 'midnight' | 'slate'

interface OutputLine {
  text: string
  tone?: Tone
  to?: string
  href?: string
  meta?: string
}

interface TerminalEntry {
  id: number
  command?: string
  lines: OutputLine[]
}

interface Bookmark {
  name: string
  url: string
  icon?: string
}

interface TodoItem {
  id: number
  text: string
  done: boolean
}

const BOOKMARKS_KEY = 'vilivala_terminal_bookmarks'
const TODOS_KEY = 'vilivala_terminal_todos'
const BG_KEY = 'vilivala_terminal_bg'

const DEFAULT_BOOKMARKS: Bookmark[] = [
  { name: '博客首页', url: '/', icon: '🏠' },
  { name: 'GitHub', url: 'https://github.com', icon: '🐙' },
  { name: '知乎', url: 'https://www.zhihu.com', icon: '📘' },
  { name: 'Bilibili', url: 'https://www.bilibili.com', icon: '📺' },
  { name: 'Vue 文档', url: 'https://cn.vuejs.org', icon: '💚' },
  { name: 'MDN', url: 'https://developer.mozilla.org/zh-CN/', icon: '📚' },
  { name: 'Go 文档', url: 'https://go.dev/doc/', icon: '🐹' },
]

const COMMAND_NAMES = [
  'help',
  'ls',
  'search',
  's',
  'open',
  'cat',
  'google',
  'g',
  'baidu',
  'bd',
  'github',
  'gh',
  'bilibili',
  'bili',
  'goto',
  'bm',
  'bookmarks',
  'todo',
  'bg',
  'theme',
  'whoami',
  'about',
  'date',
  'history',
  'clear',
  'home',
]

const bgPresets: { key: BgMode; label: string }[] = [
  { key: 'grid', label: '网格' },
  { key: 'emerald', label: '翡翠' },
  { key: 'midnight', label: '午夜' },
  { key: 'slate', label: '极简' },
]

const starterChips = ['help', 'ls posts', 'ls categories', 's 博客', 'bm', 'todo', 'whoami']

const router = useRouter()
const theme = useThemeStore()

const inputRef = ref<HTMLInputElement | null>(null)
const outputRef = ref<HTMLElement | null>(null)
const input = ref('')
const entries = ref<TerminalEntry[]>([])
const history = ref<string[]>([])
const historyIdx = ref(-1)
const bgMode = ref<BgMode>((localStorage.getItem(BG_KEY) as BgMode) || 'grid')
const bookmarks = ref<Bookmark[]>(loadBookmarks())
const todos = ref<TodoItem[]>(loadTodos())
const clockTime = ref('')
const clockDate = ref('')
const siteStats = ref({ posts: 0, categories: 0 })

let seq = 0
let clockTimer: number | undefined

const bgClass = computed(() => {
  switch (bgMode.value) {
    case 'emerald':
      return 'bg-gradient-to-br from-[#0f241e] via-[#121a19] to-[#0d1117]'
    case 'midnight':
      return 'bg-gradient-to-br from-[#0d1326] via-[#131726] to-[#0a0c12]'
    case 'slate':
      return 'bg-[#18191c]'
    default:
      return 'hero-grid'
  }
})

const ghostHint = computed(() => {
  const val = input.value
  if (!val || val.includes(' ')) return ''
  const match = COMMAND_NAMES.find((c) => c.startsWith(val.toLowerCase()) && c !== val.toLowerCase())
  return match ? match.slice(val.length) : ''
})

function loadBookmarks(): Bookmark[] {
  try {
    const raw = localStorage.getItem(BOOKMARKS_KEY)
    if (raw) {
      const parsed = JSON.parse(raw)
      if (Array.isArray(parsed) && parsed.length) return parsed
    }
  } catch {
    /* ignore */
  }
  return [...DEFAULT_BOOKMARKS]
}

function saveBookmarks() {
  localStorage.setItem(BOOKMARKS_KEY, JSON.stringify(bookmarks.value))
}

function loadTodos(): TodoItem[] {
  try {
    const raw = localStorage.getItem(TODOS_KEY)
    if (raw) {
      const parsed = JSON.parse(raw)
      if (Array.isArray(parsed)) return parsed
    }
  } catch {
    /* ignore */
  }
  return []
}

function saveTodos() {
  localStorage.setItem(TODOS_KEY, JSON.stringify(todos.value))
}

function setBg(mode: BgMode) {
  bgMode.value = mode
  localStorage.setItem(BG_KEY, mode)
}

function lineToneClass(tone?: Tone) {
  switch (tone) {
    case 'accent':
      return 'text-[#4fc08d]'
    case 'muted':
      return 'text-gray-400'
    case 'error':
      return 'text-rose-400'
    case 'warn':
      return 'text-amber-300'
    default:
      return 'text-gray-200'
  }
}

function pushEntry(command: string | undefined, lines: OutputLine[]) {
  entries.value.push({ id: ++seq, command, lines })
  nextTick(() => {
    if (outputRef.value) {
      outputRef.value.scrollTop = outputRef.value.scrollHeight
    }
  })
}

function focusInput() {
  inputRef.value?.focus()
}

function acceptCompletion() {
  if (ghostHint.value) {
    input.value += ghostHint.value + ' '
  }
}

function browseHistory(delta: number) {
  if (!history.value.length) return
  if (historyIdx.value === -1) {
    historyIdx.value = delta < 0 ? history.value.length - 1 : -1
  } else {
    const next = historyIdx.value + delta
    if (next < 0) historyIdx.value = 0
    else if (next >= history.value.length) {
      historyIdx.value = -1
      input.value = ''
      return
    } else {
      historyIdx.value = next
    }
  }
  if (historyIdx.value >= 0) {
    input.value = history.value[historyIdx.value]
  }
}

function clearEntries() {
  entries.value = []
}

function normalizeUrl(raw: string): string {
  if (/^https?:\/\//i.test(raw) || raw.startsWith('/')) return raw
  return `https://${raw}`
}

function openWebSearch(engine: string, q: string, urlBuilder: (encoded: string) => string) {
  if (!q) {
    return [{ text: `用法：${engine} <关键词>`, tone: 'warn' as Tone }]
  }
  const url = urlBuilder(encodeURIComponent(q))
  window.open(url, '_blank', 'noopener,noreferrer')
  return [{ text: `已在新标签打开 ${engine} 搜索：${q}`, href: url }]
}

async function executeCommand(raw: string) {
  const trimmed = raw.trim()
  if (!trimmed) return

  history.value.push(trimmed)
  historyIdx.value = -1

  const parts = trimmed.split(/\s+/)
  const cmd = parts[0].toLowerCase()
  const argStr = trimmed.slice(parts[0].length).trim()
  const args = parts.slice(1)

  if (cmd === 'clear') {
    clearEntries()
    return
  }

  if (cmd === 'home' || cmd === 'exit') {
    pushEntry(trimmed, [{ text: '正在返回博客首页…', tone: 'accent' }])
    router.push('/')
    return
  }

  if (cmd === 'help') {
    pushEntry(trimmed, [
      { text: '=== 博客指令 ===', tone: 'accent' },
      { text: '  ls [posts|categories|pages]  列出最新文章 / 分类 / 独立页' },
      { text: '  s / search <关键词>           站内全文搜索文章' },
      { text: '  cat / open <slug>            预览或直接打开文章' },
      { text: '  whoami / about               查看博主名片与站点统计' },
      { text: '=== 浏览器起始页指令 ===', tone: 'accent' },
      { text: '  g / google <关键词>           Google 全网搜索' },
      { text: '  bd / baidu <关键词>           百度搜索' },
      { text: '  gh / github <关键词>          GitHub 搜索' },
      { text: '  bili / bilibili <关键词>      Bilibili 搜索' },
      { text: '  goto <书签名或网址>           直达网站' },
      { text: '  bm [add <名> <url>|rm <名>]  查看或管理书签' },
      { text: '  todo [add <事项>|done <n>|rm <n>]  本地待办清单' },
      { text: '=== 系统与外观 ===', tone: 'accent' },
      { text: '  bg [grid|emerald|midnight|slate]   切换终端背景' },
      { text: '  theme [system|light|dark]          切换全站主题' },
      { text: '  date / history / clear / home      实用工具' },
    ])
    return
  }

  if (cmd === 'date') {
    pushEntry(trimmed, [{ text: new Date().toLocaleString('zh-CN', { hour12: false }), tone: 'accent' }])
    return
  }

  if (cmd === 'history') {
    const lines: OutputLine[] = history.value.slice(-15).map((h, i) => ({
      text: `  ${String(i + 1).padStart(2, ' ')}  ${h}`,
      tone: 'muted',
    }))
    pushEntry(trimmed, lines.length ? lines : [{ text: '暂无历史记录', tone: 'muted' }])
    return
  }

  if (cmd === 'whoami' || cmd === 'about') {
    pushEntry(trimmed, [
      { text: `${siteConfig.blogger.name} — ${siteConfig.blogger.slogan}`, tone: 'accent' },
      { text: `站点：${siteConfig.title}（${siteConfig.hero.title}）` },
      { text: `统计：${siteStats.value.posts} 篇已发布文章 · ${siteStats.value.categories} 个分类`, tone: 'muted' },
      { text: '→ 打开完整「关于」页面', to: '/pages/about' },
    ])
    return
  }

  if (cmd === 'bg') {
    const target = args[0]?.toLowerCase() as BgMode | undefined
    if (!target || !['grid', 'emerald', 'midnight', 'slate'].includes(target)) {
      pushEntry(trimmed, [{ text: `当前背景：${bgMode.value}。可选：grid | emerald | midnight | slate`, tone: 'muted' }])
      return
    }
    setBg(target)
    pushEntry(trimmed, [{ text: `已切换终端背景为 ${target}`, tone: 'accent' }])
    return
  }

  if (cmd === 'theme') {
    const target = args[0]?.toLowerCase() as ThemeMode | undefined
    if (!target || !['system', 'light', 'dark'].includes(target)) {
      pushEntry(trimmed, [{ text: `当前全站主题：${theme.mode}。可选：system | light | dark`, tone: 'muted' }])
      return
    }
    theme.setMode(target)
    pushEntry(trimmed, [{ text: `已切换全站主题为 ${target}`, tone: 'accent' }])
    return
  }

  if (cmd === 'g' || cmd === 'google') {
    pushEntry(trimmed, openWebSearch('Google', argStr, (q) => `https://www.google.com/search?q=${q}`))
    return
  }
  if (cmd === 'bd' || cmd === 'baidu') {
    pushEntry(trimmed, openWebSearch('百度', argStr, (q) => `https://www.baidu.com/s?wd=${q}`))
    return
  }
  if (cmd === 'gh' || cmd === 'github') {
    pushEntry(trimmed, openWebSearch('GitHub', argStr, (q) => `https://github.com/search?q=${q}`))
    return
  }
  if (cmd === 'bili' || cmd === 'bilibili') {
    pushEntry(trimmed, openWebSearch('Bilibili', argStr, (q) => `https://search.bilibili.com/all?keyword=${q}`))
    return
  }

  if (cmd === 'goto') {
    if (!argStr) {
      pushEntry(trimmed, [{ text: '用法：goto <书签名或网址>，例如 goto github.com', tone: 'warn' }])
      return
    }
    const found = bookmarks.value.find((b) => b.name.toLowerCase() === argStr.toLowerCase())
    const targetUrl = found ? found.url : normalizeUrl(argStr)
    if (targetUrl.startsWith('/')) {
      pushEntry(trimmed, [{ text: `正在前往 ${targetUrl}`, to: targetUrl }])
      router.push(targetUrl)
    } else {
      window.open(targetUrl, '_blank', 'noopener,noreferrer')
      pushEntry(trimmed, [{ text: `已打开 ${targetUrl}`, href: targetUrl }])
    }
    return
  }

  if (cmd === 'bm' || cmd === 'bookmarks') {
    const sub = args[0]?.toLowerCase()
    if (sub === 'add' && args.length >= 3) {
      const name = args[1]
      const url = normalizeUrl(args[2])
      bookmarks.value = [...bookmarks.value.filter((b) => b.name !== name), { name, url, icon: '🔖' }]
      saveBookmarks()
      pushEntry(trimmed, [{ text: `已添加书签：${name} → ${url}`, tone: 'accent' }])
      return
    }
    if (sub === 'rm' && args[1]) {
      const before = bookmarks.value.length
      bookmarks.value = bookmarks.value.filter((b) => b.name.toLowerCase() !== args[1].toLowerCase())
      saveBookmarks()
      pushEntry(trimmed, [
        {
          text: bookmarks.value.length < before ? `已删除书签：${args[1]}` : `未找到书签：${args[1]}`,
          tone: bookmarks.value.length < before ? 'accent' : 'warn',
        },
      ])
      return
    }
    pushEntry(
      trimmed,
      bookmarks.value.map((b) => ({
        text: `${b.icon || '🔗'} ${b.name}`,
        meta: b.url,
        to: b.url.startsWith('/') ? b.url : undefined,
        href: b.url.startsWith('/') ? undefined : b.url,
      }))
    )
    return
  }

  if (cmd === 'todo') {
    const sub = args[0]?.toLowerCase()
    if (sub === 'add' && args.length >= 2) {
      const text = args.slice(1).join(' ')
      const nextId = todos.value.reduce((m, t) => Math.max(m, t.id), 0) + 1
      todos.value.push({ id: nextId, text, done: false })
      saveTodos()
      pushEntry(trimmed, [{ text: `已添加待办 #${nextId}：${text}`, tone: 'accent' }])
      return
    }
    if (sub === 'done' && args[1]) {
      const item = todos.value.find((t) => t.id === Number(args[1]))
      if (!item) {
        pushEntry(trimmed, [{ text: `未找到待办 #${args[1]}`, tone: 'warn' }])
        return
      }
      item.done = !item.done
      saveTodos()
      pushEntry(trimmed, [{ text: `待办 #${item.id} 已标记为 ${item.done ? '完成' : '未完成'}`, tone: 'accent' }])
      return
    }
    if (sub === 'rm' && args[1]) {
      const before = todos.value.length
      todos.value = todos.value.filter((t) => t.id !== Number(args[1]))
      saveTodos()
      pushEntry(trimmed, [
        {
          text: todos.value.length < before ? `已删除待办 #${args[1]}` : `未找到待办 #${args[1]}`,
          tone: todos.value.length < before ? 'accent' : 'warn',
        },
      ])
      return
    }
    if (!todos.value.length) {
      pushEntry(trimmed, [{ text: '暂无待办。使用 todo add <事项> 添加一条吧！', tone: 'muted' }])
      return
    }
    pushEntry(
      trimmed,
      todos.value.map((t) => ({
        text: `  [${t.done ? '✓' : ' '}] #${t.id} ${t.text}`,
        tone: t.done ? 'muted' : 'normal',
      }))
    )
    return
  }

  if (cmd === 'ls') {
    const target = (args[0] || 'posts').toLowerCase()
    try {
      if (target === 'categories' || target === 'cat') {
        const res = await getCategories(false)
        pushEntry(
          trimmed,
          res.data.map((c) => ({
            text: `  ${c.path}`,
            meta: `${c.total_count} 篇  (/categories/${c.slug})`,
            to: `/categories/${c.slug}`,
          }))
        )
        return
      }
      if (target === 'pages') {
        const res = await getPages()
        pushEntry(
          trimmed,
          res.data.map((p) => ({
            text: `  ${p.title}`,
            meta: `/pages/${p.slug}`,
            to: `/pages/${p.slug}`,
          }))
        )
        return
      }
      const res = await getPublishedPosts({ page: 1, page_size: 10 })
      pushEntry(trimmed, [
        { text: `最新 ${res.data.list.length} / ${res.data.total} 篇文章（点击标题或输入 open <slug> 打开）：`, tone: 'muted' },
        ...res.data.list.map((p) => ({
          text: `  ${p.is_top ? '[顶] ' : ''}${p.title}`,
          meta: `${p.slug} · ${formatDate(p.created_at)}`,
          to: `/posts/${p.slug}`,
        })),
      ])
    } catch {
      pushEntry(trimmed, [{ text: '获取数据失败，请检查后端服务状态', tone: 'error' }])
    }
    return
  }

  if (cmd === 's' || cmd === 'search') {
    if (!argStr) {
      pushEntry(trimmed, [{ text: '用法：s <关键词>', tone: 'warn' }])
      return
    }
    try {
      const res = await searchPosts(argStr, 10)
      if (!res.data.list.length) {
        pushEntry(trimmed, [{ text: `未找到包含「${argStr}」的文章`, tone: 'muted' }])
        return
      }
      const lines: OutputLine[] = [
        { text: `共找到 ${res.data.total} 篇匹配文章（展示前 ${res.data.list.length} 篇）：`, tone: 'accent' },
      ]
      for (const item of res.data.list) {
        lines.push({
          text: `  • ${item.title}`,
          meta: `${item.slug} · ${formatDate(item.created_at)}`,
          to: `/posts/${item.slug}`,
        })
        if (item.snippet) {
          lines.push({ text: `    ${item.snippet}`, tone: 'muted' })
        }
      }
      pushEntry(trimmed, lines)
    } catch {
      pushEntry(trimmed, [{ text: '检索失败，请稍后重试', tone: 'error' }])
    }
    return
  }

  if (cmd === 'open' || cmd === 'cat') {
    if (!argStr) {
      pushEntry(trimmed, [{ text: `用法：${cmd} <文章slug>`, tone: 'warn' }])
      return
    }
    if (cmd === 'open') {
      pushEntry(trimmed, [{ text: `正在打开 /posts/${argStr}`, to: `/posts/${argStr}` }])
      router.push(`/posts/${argStr}`)
      return
    }
    try {
      const res = await getPostBySlug(argStr)
      const p = res.data
      const preview = (p.summary || p.content || '').slice(0, 260)
      pushEntry(trimmed, [
        { text: `# ${p.title}`, tone: 'accent' },
        { text: `发布于 ${formatDate(p.created_at)} · ${p.word_count} 字 · 预计 ${p.reading_time || 1} 分钟`, tone: 'muted' },
        { text: preview + (preview.length >= 260 ? '…' : '') },
        { text: `→ 点击阅读完整文章 (/posts/${p.slug})`, to: `/posts/${p.slug}` },
      ])
    } catch {
      pushEntry(trimmed, [{ text: `未找到 slug 为「${argStr}」的文章`, tone: 'error' }])
    }
    return
  }

  pushEntry(trimmed, [
    { text: `未知命令：${cmd}。输入 help 查看可用指令，或使用 g ${trimmed} 搜索全网。`, tone: 'error' },
  ])
}

function runPreset(cmd: string) {
  executeCommand(cmd)
  focusInput()
}

function onSubmit() {
  const val = input.value
  input.value = ''
  executeCommand(val)
}

function tickClock() {
  const now = new Date()
  clockTime.value = now.toLocaleTimeString('zh-CN', { hour12: false })
  clockDate.value = now.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    weekday: 'short',
  })
}

onMounted(async () => {
  updatePageMeta({
    title: '终端起始页',
    description: `${siteConfig.title} 交互式终端浏览器起始页，支持站内全文检索、全网快捷搜索、书签与待办管理。`,
    path: '/terminal',
  })
  tickClock()
  clockTimer = window.setInterval(tickClock, 1000)

  pushEntry(undefined, [
    { text: `Welcome to ${siteConfig.title} Terminal Start Page (v1.1.1)`, tone: 'accent' },
    { text: `${siteConfig.hero.title} — ${siteConfig.hero.tagline}`, tone: 'muted' },
    { text: '输入 help 查看所有命令，或直接点击上方快捷指令开始探索。', tone: 'muted' },
  ])

  try {
    const [p, c] = await Promise.all([getPublishedPosts({ page: 1, page_size: 1 }), getCategories(false)])
    siteStats.value = { posts: p.data.total, categories: c.data.length }
  } catch {
    /* ignore */
  }
  focusInput()
})

onBeforeUnmount(() => {
  if (clockTimer) window.clearInterval(clockTimer)
})
</script>

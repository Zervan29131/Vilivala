<template>
  <div
    class="min-h-dvh-safe flex flex-col text-gray-100 transition-colors duration-300 relative"
    :class="bgImage ? 'bg-black' : bgClass"
    :style="bgWrapperStyle"
  >
    <!-- Dark overlay when custom wallpaper image is active -->
    <div v-if="bgImage" class="fixed inset-0 bg-black/60 backdrop-blur-[2px] pointer-events-none z-0" />

    <!-- Browser-style Top Chrome -->
    <header class="sticky top-0 z-30 border-b border-white/10 bg-black/45 backdrop-blur-md">
      <div class="max-w-6xl mx-auto px-4 h-13 py-2 flex items-center justify-between gap-3">
        <!-- Left: traffic lights + brand + current path -->
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
          <span class="hidden md:inline-flex items-center px-2 py-0.5 rounded bg-[#42b883]/20 text-[#4fc08d] text-[11px] font-mono truncate max-w-[220px]">
            &gt;_ {{ promptPath }}
          </span>
        </div>

        <!-- Center: Clock -->
        <div class="hidden sm:flex items-center gap-2 font-mono text-xs text-gray-300">
          <span>{{ clockTime }}</span>
          <span class="text-white/25">|</span>
          <span class="text-gray-400">{{ clockDate }}</span>
        </div>

        <!-- Right: Background switcher & blog links -->
        <div class="flex items-center gap-1.5 text-xs shrink-0">
          <!-- Compact select on small phones -->
          <select
            :value="bgMode"
            aria-label="切换终端背景"
            class="sm:hidden bg-white/10 text-gray-200 rounded px-2 py-1 text-xs border border-white/10 focus:outline-none"
            @change="setBgPreset(($event.target as HTMLSelectElement).value as BgMode)"
          >
            <option v-for="preset in bgPresets" :key="preset.key" :value="preset.key" class="bg-[#1a1a1a] text-white">
              {{ preset.label }}
            </option>
          </select>

          <!-- Buttons on sm+ -->
          <div class="hidden sm:flex items-center gap-1">
            <button
              v-for="preset in bgPresets"
              :key="preset.key"
              type="button"
              class="px-2 py-1 rounded transition-colors"
              :class="
                !bgImage && bgMode === preset.key
                  ? 'bg-[#42b883] text-white'
                  : 'text-gray-400 hover:text-white hover:bg-white/10'
              "
              @click="setBgPreset(preset.key)"
            >
              {{ preset.label }}
            </button>
          </div>
          <span class="hidden sm:inline text-white/15 mx-0.5">|</span>
          <router-link
            to="/"
            class="px-2.5 py-1 rounded bg-white/10 hover:bg-[#42b883] text-white transition-colors whitespace-nowrap"
          >
            返回博客
          </router-link>
        </div>
      </div>
    </header>

    <!-- Main Container -->
    <div class="relative z-10 flex-1 max-w-6xl w-full mx-auto px-4 py-6 flex flex-col gap-5">
      <!-- Quick Launch Bookmarks Dock -->
      <section class="rounded-xl border border-white/10 bg-black/35 backdrop-blur-sm p-4">
        <div class="flex items-center justify-between mb-3 gap-2 flex-wrap">
          <div class="flex items-center gap-2 text-xs text-gray-400">
            <span class="text-[#4fc08d] font-mono">★</span>
            <span>书签与目录空间</span>
            <span class="text-white/25">·</span>
            <span class="hidden sm:inline font-mono text-[11px]">
              试试 <code class="text-[#4fc08d]">cd code</code> 逛分类，或 <code class="text-[#4fc08d]">mkdir</code> / <code class="text-[#4fc08d]">touch</code> 管理书签
            </span>
          </div>
          <div class="flex items-center gap-2.5 text-xs">
            <button
              v-if="cwd !== '/'"
              type="button"
              class="text-[#4fc08d] hover:underline font-mono"
              @click="runPreset('cd ..')"
            >
              ↑ 返回上级 ({{ cwd }})
            </button>
            <router-link to="/archives" class="text-gray-400 hover:text-[#4fc08d] transition-colors">归档</router-link>
            <router-link to="/categories" class="text-gray-400 hover:text-[#4fc08d] transition-colors">分类</router-link>
            <a href="/rss.xml" target="_blank" rel="noopener" class="text-gray-400 hover:text-[#4fc08d] transition-colors">RSS</a>
          </div>
        </div>
        <div class="flex flex-wrap gap-2">
          <template v-for="bm in dockLinks" :key="bm.path">
            <router-link
              v-if="bm.url?.startsWith('/')"
              :to="bm.url"
              class="px-3 py-1.5 rounded-lg bg-white/[0.06] hover:bg-[#42b883]/20 border border-white/10 hover:border-[#42b883]/50 text-xs text-gray-200 hover:text-[#4fc08d] transition-all flex items-center gap-1.5"
            >
              <span>{{ bm.icon || '🔗' }}</span>
              <span>{{ bm.name }}</span>
            </router-link>
            <a
              v-else-if="bm.url"
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
        class="flex-1 min-h-[420px] sm:min-h-[540px] rounded-xl border border-white/15 bg-[#121417]/90 shadow-2xl flex flex-col overflow-hidden"
        @click="focusInput"
      >
        <!-- Quick Command Chips -->
        <div class="px-3.5 sm:px-4 py-2.5 border-b border-white/10 bg-white/[0.03] flex items-center justify-between gap-2">
          <div class="flex items-center flex-nowrap sm:flex-wrap gap-1.5 text-xs font-mono overflow-x-auto no-scrollbar min-w-0">
            <span class="text-gray-500 mr-0.5 shrink-0">快捷体验:</span>
            <button
              v-for="chip in starterChips"
              :key="chip"
              type="button"
              class="px-2 py-0.5 rounded bg-white/[0.07] hover:bg-[#42b883]/25 text-gray-300 hover:text-[#4fc08d] transition-colors shrink-0"
              @click.stop="runPreset(chip)"
            >
              {{ chip }}
            </button>
          </div>
          <div class="flex items-center gap-3 text-[11px] font-mono text-gray-500 shrink-0">
            <button
              v-if="entries.length > 1"
              type="button"
              class="hover:text-[#4fc08d] transition-colors"
              title="折叠或展开全部输出 (Ctrl+O)"
              @click.stop="toggleAllCollapse"
            >
              {{ anyExpanded ? '全部收起' : '全部展开' }}<span class="hidden md:inline"> (Ctrl+O)</span>
            </button>
            <span class="hidden lg:inline">Tab 补全 · Ctrl+L 清屏</span>
          </div>
        </div>

        <!-- Terminal Output Stream -->
        <div ref="outputRef" class="flex-1 overflow-y-auto p-4 sm:p-5 font-mono text-[13px] leading-relaxed space-y-3">
          <div v-for="entry in entries" :key="entry.id" class="space-y-1.5">
            <!-- Command line echo + collapse toggle -->
            <div v-if="entry.command" class="flex items-center justify-between gap-2 group">
              <div class="flex items-center gap-1.5 min-w-0 flex-wrap">
                <span class="text-[#4fc08d] shrink-0">zervan@vilivala</span>
                <span class="text-gray-500">:</span>
                <span class="text-sky-400">{{ entry.promptPath || '~' }}</span>
                <span class="text-gray-400">$</span>
                <span class="text-white break-all">{{ entry.command }}</span>
              </div>
              <button
                v-if="canCollapse(entry)"
                type="button"
                class="text-[11px] px-1.5 py-0.5 rounded text-gray-500 hover:text-[#4fc08d] hover:bg-white/5 shrink-0 transition-colors"
                @click.stop="entry.collapsed = !entry.collapsed"
              >
                {{ entry.collapsed ? `展开 (${entry.lines.length || 1}) ▸` : '收起 ▾' }}
              </button>
            </div>

            <!-- Collapsible Output Body -->
            <div v-show="!entry.collapsed" class="pl-0 sm:pl-2 space-y-1">
              <!-- Text / Link lines -->
              <div
                v-for="(line, li) in entry.lines"
                :key="li"
                class="break-words whitespace-pre-wrap"
                :class="lineToneClass(line.tone)"
              >
                <button
                  v-if="line.runCmd"
                  type="button"
                  class="text-sky-400 hover:text-[#4fc08d] hover:underline text-left"
                  @click.stop="runPreset(line.runCmd)"
                >
                  {{ line.text }}
                </button>
                <router-link
                  v-else-if="line.to"
                  :to="line.to"
                  class="text-[#4fc08d] hover:underline"
                  @click.stop
                >
                  {{ line.text }}
                </router-link>
                <a
                  v-else-if="line.href"
                  :href="line.href"
                  :target="line.selfTab ? '_self' : '_blank'"
                  rel="noopener noreferrer"
                  class="text-sky-400 hover:underline"
                  @click.stop
                >
                  {{ line.text }}
                </a>
                <span v-else>{{ line.text }}</span>
                <span v-if="line.meta" class="ml-2 text-xs text-gray-500">{{ line.meta }}</span>
              </div>

              <!-- Interactive Embedded Widget -->
              <TerminalHotBox v-if="entry.widget?.kind === 'hot'" />
              <TerminalTimingBox
                v-else-if="entry.widget?.kind === 'timing'"
                :seconds="entry.widget.seconds || 60"
                :label="entry.widget.label"
              />
              <TerminalVarbookBox
                v-else-if="entry.widget?.kind === 'varbook'"
                :query="entry.widget.query || ''"
              />
              <TerminalShortcutBox v-else-if="entry.widget?.kind === 'shortcut'" />
            </div>
          </div>
        </div>

        <!-- Prompt Input Bar + Live Syntax Hint -->
        <div class="border-t border-white/10 bg-black/35 px-3.5 sm:px-5 py-2.5 sm:py-3 space-y-1.5 pb-safe">
          <form class="flex items-center gap-1.5 sm:gap-2 font-mono text-sm" @submit.prevent="onSubmit">
            <span class="hidden sm:inline text-[#4fc08d] shrink-0">zervan@vilivala</span>
            <span class="hidden sm:inline text-gray-500">:</span>
            <span class="text-sky-400 shrink-0 max-w-[110px] sm:max-w-[220px] truncate">{{ promptPath }}</span>
            <span class="text-gray-400 shrink-0">$</span>
            <div class="relative flex-1 min-w-0">
              <input
                ref="inputRef"
                v-model="input"
                type="text"
                autocomplete="off"
                autocapitalize="off"
                spellcheck="false"
                class="w-full bg-transparent text-[16px] sm:text-sm text-white placeholder:text-gray-500 focus:outline-none"
                placeholder="输入命令（如 cd code / hot / help）…"
                @keydown.up.prevent="browseHistory(-1)"
                @keydown.down.prevent="browseHistory(1)"
                @keydown.tab.prevent="acceptCompletion"
                @keydown.ctrl.l.prevent="clearEntries"
                @keydown.ctrl.o.prevent="toggleAllCollapse"
              />
              <span
                v-if="ghostHint"
                class="hidden sm:inline pointer-events-none absolute left-0 top-0 text-gray-500/70 select-none whitespace-pre"
              >
                <span class="invisible">{{ input }}</span>{{ ghostHint }}
              </span>
            </div>
            <button
              type="submit"
              class="px-3 py-1.5 sm:py-1 rounded bg-[#42b883] hover:bg-[#33a06f] text-white text-xs font-sans transition-colors shrink-0"
            >
              执行
            </button>
          </form>

          <!-- Live Usage Hint Row -->
          <div v-if="hintEnabled && liveUsageHint" class="text-[11px] font-mono text-gray-400 flex items-center gap-2 truncate">
            <span class="text-[#4fc08d]">提示:</span>
            <span class="text-gray-300">{{ liveUsageHint.usage }}</span>
            <span class="text-gray-500">— {{ liveUsageHint.summary }}</span>
          </div>
        </div>
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
import { expandHistoryBang, parseCommandLine, parseDurationSeconds } from '@/utils/terminal/parser'
import { COMMAND_DOCS, allCommandKeywords, findCommandDoc } from '@/utils/terminal/registry'
import {
  VFS_STORAGE_KEY,
  VirtualSpace,
  defaultUserNodes,
  formatPromptPath,
  resolveVfsPath,
  type VfsNode,
} from '@/utils/terminal/vfs'
import TerminalHotBox from '@/components/terminal/TerminalHotBox.vue'
import TerminalTimingBox from '@/components/terminal/TerminalTimingBox.vue'
import TerminalVarbookBox from '@/components/terminal/TerminalVarbookBox.vue'
import TerminalShortcutBox from '@/components/terminal/TerminalShortcutBox.vue'

type Tone = 'normal' | 'accent' | 'muted' | 'error' | 'warn'
type BgMode = 'grid' | 'emerald' | 'midnight' | 'slate'

interface OutputLine {
  text: string
  tone?: Tone
  to?: string
  href?: string
  selfTab?: boolean
  runCmd?: string
  meta?: string
}

interface WidgetSpec {
  kind: 'hot' | 'timing' | 'varbook' | 'shortcut'
  seconds?: number
  label?: string
  query?: string
}

interface TerminalEntry {
  id: number
  promptPath?: string
  command?: string
  collapsed: boolean
  lines: OutputLine[]
  widget?: WidgetSpec
}

interface TodoItem {
  id: number
  text: string
  done: boolean
}

const TODOS_KEY = 'vilivala_terminal_todos'
const BG_KEY = 'vilivala_terminal_bg'
const BG_IMG_KEY = 'vilivala_terminal_bg_img'
const HINT_KEY = 'vilivala_terminal_hint'

const SEARCH_ENGINES: Record<string, { label: string; build: (q: string) => string }> = {
  google: { label: 'Google', build: (q) => `https://www.google.com/search?q=${q}` },
  g: { label: 'Google', build: (q) => `https://www.google.com/search?q=${q}` },
  baidu: { label: '百度', build: (q) => `https://www.baidu.com/s?wd=${q}` },
  bd: { label: '百度', build: (q) => `https://www.baidu.com/s?wd=${q}` },
  bing: { label: 'Bing', build: (q) => `https://www.bing.com/search?q=${q}` },
  github: { label: 'GitHub', build: (q) => `https://github.com/search?q=${q}` },
  gh: { label: 'GitHub', build: (q) => `https://github.com/search?q=${q}` },
  bilibili: { label: 'Bilibili', build: (q) => `https://search.bilibili.com/all?keyword=${q}` },
  bili: { label: 'Bilibili', build: (q) => `https://search.bilibili.com/all?keyword=${q}` },
  zhihu: { label: '知乎', build: (q) => `https://www.zhihu.com/search?type=content&q=${q}` },
  douban: { label: '豆瓣', build: (q) => `https://www.douban.com/search?q=${q}` },
  mdn: { label: 'MDN', build: (q) => `https://developer.mozilla.org/zh-CN/search?q=${q}` },
}

const bgPresets: { key: BgMode; label: string }[] = [
  { key: 'grid', label: '网格' },
  { key: 'emerald', label: '翡翠' },
  { key: 'midnight', label: '午夜' },
  { key: 'slate', label: '极简' },
]

const starterChips = ['help', 'hot', 'cd code', 'ls', 'timing 25m', 'varbook 用户配置', 'shortcut']
const allCmdNames = allCommandKeywords()

const router = useRouter()
const theme = useThemeStore()

const space = new VirtualSpace(loadSavedVfs())
const vfsVersion = ref(0)
const cwd = ref('/')

const inputRef = ref<HTMLInputElement | null>(null)
const outputRef = ref<HTMLElement | null>(null)
const input = ref('')
const entries = ref<TerminalEntry[]>([])
const history = ref<string[]>([])
const historyIdx = ref(-1)

const bgMode = ref<BgMode>((localStorage.getItem(BG_KEY) as BgMode) || 'grid')
const bgImage = ref<string>(localStorage.getItem(BG_IMG_KEY) || '')
const hintEnabled = ref(localStorage.getItem(HINT_KEY) !== 'off')
const todos = ref<TodoItem[]>(loadTodos())
const clockTime = ref('')
const clockDate = ref('')
const siteStats = ref({ posts: 0, categories: 0 })

let seq = 0
let clockTimer: number | undefined

const promptPath = computed(() => formatPromptPath(cwd.value))

const dockLinks = computed(() => {
  // Touch vfsVersion so Vue tracks mutations to VirtualSpace
  void vfsVersion.value
  return space.allLinks()
})

const anyExpanded = computed(() => entries.value.some((e) => canCollapse(e) && !e.collapsed))

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

const bgWrapperStyle = computed(() => {
  if (!bgImage.value) return undefined
  return {
    backgroundImage: `url("${bgImage.value.replace(/"/g, '')}")`,
    backgroundSize: 'cover',
    backgroundPosition: 'center',
  }
})

const liveUsageHint = computed(() => {
  const trimmed = input.value.trim()
  if (!trimmed) return null
  const firstWord = trimmed.split(/\s+/)[0]
  return findCommandDoc(firstWord) || null
})

const ghostHint = computed(() => {
  const val = input.value
  if (!val) return ''
  // 1) First token -> command completion
  if (!val.includes(' ')) {
    const match = allCmdNames.find((c) => c.startsWith(val.toLowerCase()) && c !== val.toLowerCase())
    return match ? match.slice(val.length) : ''
  }
  // 2) Path argument completion for cd / ls / rm / mv / cp
  const m = val.match(/^(cd|ls|rm|mv|cp)\s+([^\s]*)$/i)
  if (m) {
    const partial = m[2]
    void vfsVersion.value
    const children = space.listChildren(cwd.value, false)
    const hit = children.find(
      (n) => n.name.toLowerCase().startsWith(partial.toLowerCase()) && n.name.toLowerCase() !== partial.toLowerCase()
    )
    return hit ? hit.name.slice(partial.length) : ''
  }
  return ''
})

function loadSavedVfs(): VfsNode[] {
  try {
    const raw = localStorage.getItem(VFS_STORAGE_KEY)
    if (raw) {
      const parsed = JSON.parse(raw)
      if (Array.isArray(parsed) && parsed.length) return parsed
    }
  } catch {
    /* ignore */
  }
  return defaultUserNodes()
}

function persistVfs() {
  vfsVersion.value++
  localStorage.setItem(VFS_STORAGE_KEY, JSON.stringify(space.exportUserNodes()))
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

function setBgPreset(mode: BgMode) {
  bgMode.value = mode
  bgImage.value = ''
  localStorage.setItem(BG_KEY, mode)
  localStorage.removeItem(BG_IMG_KEY)
}

function setBgImageUrl(url: string) {
  bgImage.value = url
  localStorage.setItem(BG_IMG_KEY, url)
}

function canCollapse(entry: TerminalEntry): boolean {
  return !!entry.command && (entry.lines.length > 2 || !!entry.widget)
}

function toggleAllCollapse() {
  const targetCollapsed = anyExpanded.value
  for (const e of entries.value) {
    if (canCollapse(e)) e.collapsed = targetCollapsed
  }
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

function pushEntry(command: string | undefined, lines: OutputLine[], widget?: WidgetSpec) {
  entries.value.push({
    id: ++seq,
    promptPath: promptPath.value,
    command,
    collapsed: false,
    lines,
    widget,
  })
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
    const suffix = input.value.includes(' ') ? '' : ' '
    input.value += ghostHint.value + suffix
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

function navigateUrl(url: string, selfTab: boolean) {
  if (url.startsWith('/')) {
    router.push(url)
    return
  }
  if (selfTab) {
    window.location.href = url
  } else {
    window.open(url, '_blank', 'noopener,noreferrer')
  }
}

function runWebSearch(engineKey: string, query: string, selfTab: boolean): OutputLine[] {
  const engine = SEARCH_ENGINES[engineKey]
  if (!engine) {
    return [{ text: `不支持的搜索源：${engineKey}（可选 google|baidu|bing|github|bili|zhihu|douban|mdn）`, tone: 'error' }]
  }
  if (!query) {
    return [{ text: `用法：${engineKey} <关键词> [-s]`, tone: 'warn' }]
  }
  const url = engine.build(encodeURIComponent(query))
  navigateUrl(url, selfTab)
  return [
    {
      text: `已在${selfTab ? '当前页' : '新标签页'}打开 ${engine.label} 搜索：${query}`,
      href: url,
      selfTab,
    },
  ]
}

async function executeCommand(rawInput: string) {
  const trimmed = rawInput.trim()
  if (!trimmed) return

  // 1) Expand !! or !n
  const bang = expandHistoryBang(trimmed, history.value)
  if (bang.error) {
    pushEntry(trimmed, [{ text: bang.error, tone: 'error' }])
    return
  }
  const effectiveLine = bang.expanded
  history.value.push(effectiveLine)
  historyIdx.value = -1

  const parsed = parseCommandLine(effectiveLine)
  const { cmd, args, rest, flags } = parsed

  // ---- Clear / Home ----
  if (cmd === 'clear' || cmd === 'cl') {
    clearEntries()
    return
  }

  if (cmd === 'home' || cmd === 'exit') {
    pushEntry(effectiveLine, [{ text: '正在返回博客首页…', tone: 'accent' }])
    router.push('/')
    return
  }

  // ---- Help / Man ----
  if (cmd === 'help' || cmd === 'man') {
    if (args[0]) {
      const doc = findCommandDoc(args[0])
      if (!doc) {
        pushEntry(effectiveLine, [{ text: `未找到命令「${args[0]}」的说明，输入 help 查看总览。`, tone: 'warn' }])
        return
      }
      const lines: OutputLine[] = [
        { text: `${doc.name}${doc.aliases?.length ? ` (别名: ${doc.aliases.join(', ')})` : ''} — ${doc.summary}`, tone: 'accent' },
        { text: `  用法: ${doc.usage}` },
      ]
      if (doc.options?.length) {
        lines.push({ text: '  选项:', tone: 'muted' })
        doc.options.forEach((o) => lines.push({ text: `    ${o}` }))
      }
      if (doc.examples?.length) {
        lines.push({ text: '  示例（点击直接运行）:', tone: 'muted' })
        doc.examples.forEach((ex) => lines.push({ text: `    $ ${ex}`, runCmd: ex }))
      }
      pushEntry(effectiveLine, lines)
      return
    }

    const groups: { key: (typeof COMMAND_DOCS)[number]['category']; title: string }[] = [
      { key: 'space', title: '📁 目录与书签空间（支持把博客分类当文件夹逛）' },
      { key: 'blog', title: '📝 博客阅读与检索' },
      { key: 'web', title: '🌐 全网搜索与直达（支持 -s 当前页打开）' },
      { key: 'tools', title: '🧰 极客微工具卡片' },
      { key: 'system', title: '⚙️ 外观与系统控制' },
    ]
    const lines: OutputLine[] = [
      { text: '提示：输入 help <命令名> 可查看单条命令详情（例如 help search / help cd）', tone: 'muted' },
    ]
    for (const g of groups) {
      lines.push({ text: g.title, tone: 'accent' })
      for (const d of COMMAND_DOCS.filter((x) => x.category === g.key)) {
        lines.push({
          text: `  ${d.usage.padEnd(36, ' ')} ${d.summary}`,
        })
      }
    }
    pushEntry(effectiveLine, lines)
    return
  }

  // ---- Interactive Widgets: hot / timing / varbook / shortcut ----
  if (cmd === 'hot') {
    pushEntry(effectiveLine, [], { kind: 'hot' })
    return
  }

  if (cmd === 'shortcut' || cmd === 'keys') {
    pushEntry(effectiveLine, [], { kind: 'shortcut' })
    return
  }

  if (cmd === 'timing' || cmd === 'timer' || cmd === 'pomodoro') {
    const durArg = args[0] || (cmd === 'pomodoro' ? '25m' : '')
    const secs = parseDurationSeconds(durArg)
    if (!secs) {
      pushEntry(effectiveLine, [{ text: '用法：timing <时长，如 25m / 90s / 60> [标题]', tone: 'warn' }])
      return
    }
    const label = args.slice(1).join(' ') || undefined
    pushEntry(effectiveLine, [], { kind: 'timing', seconds: secs, label })
    return
  }

  if (cmd === 'varbook' || cmd === 'var') {
    if (!rest) {
      pushEntry(effectiveLine, [{ text: '用法：varbook <中文开发词或英文词组>，例如 varbook 用户配置', tone: 'warn' }])
      return
    }
    pushEntry(effectiveLine, [], { kind: 'varbook', query: rest })
    return
  }

  // ---- Virtual Space: pwd / cd / mkdir / touch / mv / cp / rm / ls ----
  if (cmd === 'pwd') {
    pushEntry(effectiveLine, [{ text: cwd.value, tone: 'accent' }])
    return
  }

  if (cmd === 'cd') {
    const targetPath = resolveVfsPath(cwd.value, args[0] || '/')
    const node = space.findInsensitive(targetPath)
    if (!node || !space.isDirectory(node)) {
      pushEntry(effectiveLine, [{ text: `目录不存在：${args[0] || targetPath}`, tone: 'error' }])
      return
    }
    cwd.value = node.path
    const hintMsg =
      node.kind === 'category'
        ? `已进入博客分类「${node.name}」（共 ${node.totalCount ?? 0} 篇），输入 ls 查看子目录与文章`
        : `当前目录：${node.path}`
    pushEntry(effectiveLine, [{ text: hintMsg, tone: 'muted' }])
    return
  }

  if (cmd === 'mkdir') {
    if (!args[0]) {
      pushEntry(effectiveLine, [{ text: '用法：mkdir <文件夹名>', tone: 'warn' }])
      return
    }
    const res = space.mkdir(resolveVfsPath(cwd.value, args[0]))
    if (!res.ok) {
      pushEntry(effectiveLine, [{ text: res.error || '创建失败', tone: 'error' }])
      return
    }
    persistVfs()
    pushEntry(effectiveLine, [{ text: `已创建文件夹：${args[0]}`, tone: 'accent' }])
    return
  }

  if (cmd === 'touch' || cmd === 'add') {
    if (args.length < 2) {
      pushEntry(effectiveLine, [{ text: '用法：touch <书签名> <网址>', tone: 'warn' }])
      return
    }
    const url = normalizeUrl(args[1])
    const res = space.addLink(resolveVfsPath(cwd.value, args[0]), url)
    if (!res.ok) {
      pushEntry(effectiveLine, [{ text: res.error || '添加失败', tone: 'error' }])
      return
    }
    persistVfs()
    pushEntry(effectiveLine, [{ text: `已保存书签：${args[0]} → ${url}`, tone: 'accent' }])
    return
  }

  if (cmd === 'mv' || cmd === 'move') {
    if (args.length < 2) {
      pushEntry(effectiveLine, [{ text: '用法：mv <源书签> <目标路径>', tone: 'warn' }])
      return
    }
    const res = space.move(resolveVfsPath(cwd.value, args[0]), resolveVfsPath(cwd.value, args[1]))
    if (!res.ok) {
      pushEntry(effectiveLine, [{ text: res.error || '移动失败', tone: 'error' }])
      return
    }
    persistVfs()
    pushEntry(effectiveLine, [{ text: `已移动 ${args[0]} → ${args[1]}`, tone: 'accent' }])
    return
  }

  if (cmd === 'cp' || cmd === 'copy') {
    if (args.length < 2) {
      pushEntry(effectiveLine, [{ text: '用法：cp <源书签> <目标路径>', tone: 'warn' }])
      return
    }
    const res = space.copy(resolveVfsPath(cwd.value, args[0]), resolveVfsPath(cwd.value, args[1]))
    if (!res.ok) {
      pushEntry(effectiveLine, [{ text: res.error || '复制失败', tone: 'error' }])
      return
    }
    persistVfs()
    pushEntry(effectiveLine, [{ text: `已复制 ${args[0]} → ${args[1]}`, tone: 'accent' }])
    return
  }

  if (cmd === 'rm' || cmd === 'del') {
    if (!args[0]) {
      pushEntry(effectiveLine, [{ text: '用法：rm <书签或目录> [-r]', tone: 'warn' }])
      return
    }
    const res = space.remove(resolveVfsPath(cwd.value, args[0]), flags.recursive)
    if (!res.ok) {
      pushEntry(effectiveLine, [{ text: res.error || '删除失败', tone: 'error' }])
      return
    }
    persistVfs()
    pushEntry(effectiveLine, [{ text: `已删除：${args[0]}`, tone: 'accent' }])
    return
  }

  if (cmd === 'ls' || cmd === 'list' || cmd === 'dir') {
    const special = (args[0] || '').toLowerCase()
    try {
      if (special === 'posts') {
        const res = await getPublishedPosts({ page: 1, page_size: 10 })
        pushEntry(effectiveLine, [
          { text: `最新 ${res.data.list.length} / ${res.data.total} 篇文章：`, tone: 'muted' },
          ...res.data.list.map((p) => ({
            text: `  📄 ${p.is_top ? '[顶] ' : ''}${p.title}`,
            meta: `${p.slug} · ${formatDate(p.created_at)}`,
            to: `/posts/${p.slug}`,
          })),
        ])
        return
      }
      if (special === 'categories') {
        const res = await getCategories(false)
        pushEntry(
          effectiveLine,
          res.data.map((c) => ({
            text: `  📂 ${c.path}`,
            meta: `${c.total_count} 篇`,
            to: `/categories/${c.slug}`,
          }))
        )
        return
      }
      if (special === 'pages') {
        const res = await getPages()
        pushEntry(
          effectiveLine,
          res.data.map((p) => ({
            text: `  📑 ${p.title}`,
            meta: `/pages/${p.slug}`,
            to: `/pages/${p.slug}`,
          }))
        )
        return
      }

      const targetDir = args[0] ? resolveVfsPath(cwd.value, args[0]) : cwd.value
      const dirNode = space.findInsensitive(targetDir)
      if (!dirNode || !space.isDirectory(dirNode)) {
        pushEntry(effectiveLine, [{ text: `未找到目录：${args[0] || targetDir}`, tone: 'error' }])
        return
      }

      const children = space.listChildren(dirNode.path, flags.recursive)
      const lines: OutputLine[] = []

      for (const item of children) {
        const displayName = flags.recursive ? item.path : item.name
        if (item.kind === 'category') {
          lines.push({
            text: `  📂 ${displayName}/`,
            meta: `博客分类 · ${item.totalCount ?? 0} 篇（点击 cd 进入）`,
            runCmd: `cd ${item.path}`,
          })
        } else if (item.kind === 'dir') {
          lines.push({
            text: `  📁 ${displayName}/`,
            meta: '自定义书签夹（点击 cd 进入）',
            runCmd: `cd ${item.path}`,
          })
        } else if (item.kind === 'link') {
          lines.push({
            text: `  ${item.icon || '🔗'} ${displayName}`,
            meta: item.url,
            to: item.url?.startsWith('/') ? item.url : undefined,
            href: item.url?.startsWith('/') ? undefined : item.url,
            selfTab: flags.self,
          })
        }
      }

      // If listing a mounted blog category directory, also list its published articles!
      if (dirNode.kind === 'category' && dirNode.categoryId) {
        const res = await getPublishedPosts({ page: 1, page_size: 12, category_id: dirNode.categoryId })
        if (res.data.list.length) {
          lines.push({
            text: `  --- 分类「${dirNode.name}」下的文章（展示前 ${res.data.list.length} / ${res.data.total} 篇）---`,
            tone: 'muted',
          })
          for (const p of res.data.list) {
            lines.push({
              text: `  📄 ${p.title}`,
              meta: `${p.slug} · ${formatDate(p.created_at)}`,
              to: `/posts/${p.slug}`,
            })
          }
        }
      }

      pushEntry(
        effectiveLine,
        lines.length ? lines : [{ text: '（空目录）可使用 touch <名称> <网址> 添加书签', tone: 'muted' }]
      )
    } catch {
      pushEntry(effectiveLine, [{ text: '读取列表失败，请检查后端连接', tone: 'error' }])
    }
    return
  }

  // ---- Blog internal search & post preview ----
  if (cmd === 's' || cmd === 'find' || (cmd === 'search' && flags.from === 'blog')) {
    if (!rest) {
      pushEntry(effectiveLine, [{ text: '用法：s <关键词>', tone: 'warn' }])
      return
    }
    try {
      const res = await searchPosts(rest, 10)
      if (!res.data.list.length) {
        pushEntry(effectiveLine, [{ text: `未找到包含「${rest}」的博客文章`, tone: 'muted' }])
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
      pushEntry(effectiveLine, lines)
    } catch {
      pushEntry(effectiveLine, [{ text: '检索失败，请稍后重试', tone: 'error' }])
    }
    return
  }

  if (cmd === 'open' || cmd === 'cat') {
    const slug = args[0]
    if (!slug) {
      pushEntry(effectiveLine, [{ text: `用法：${cmd} <文章slug>`, tone: 'warn' }])
      return
    }
    if (cmd === 'open') {
      pushEntry(effectiveLine, [{ text: `正在打开 /posts/${slug}`, to: `/posts/${slug}` }])
      router.push(`/posts/${slug}`)
      return
    }
    try {
      const res = await getPostBySlug(slug)
      const p = res.data
      const preview = (p.summary || p.content || '').slice(0, 260)
      pushEntry(effectiveLine, [
        { text: `# ${p.title}`, tone: 'accent' },
        { text: `发布于 ${formatDate(p.created_at)} · ${p.word_count} 字 · 预计 ${p.reading_time || 1} 分钟`, tone: 'muted' },
        { text: preview + (preview.length >= 260 ? '…' : '') },
        { text: `→ 点击阅读完整文章 (/posts/${p.slug})`, to: `/posts/${p.slug}` },
      ])
    } catch {
      pushEntry(effectiveLine, [{ text: `未找到 slug 为「${slug}」的文章`, tone: 'error' }])
    }
    return
  }

  // ---- Multi-engine web search ----
  if (cmd === 'search') {
    const engine = flags.from || 'google'
    pushEntry(effectiveLine, runWebSearch(engine, rest, flags.self))
    return
  }

  if (cmd in SEARCH_ENGINES) {
    pushEntry(effectiveLine, runWebSearch(cmd, rest, flags.self))
    return
  }

  // ---- Goto & Bookmarks ----
  if (cmd === 'goto' || cmd === 'visit') {
    if (!rest) {
      pushEntry(effectiveLine, [{ text: '用法：goto <书签名或网址> [-s]', tone: 'warn' }])
      return
    }
    const found = space.allLinks().find((b) => b.name.toLowerCase() === rest.toLowerCase())
    const targetUrl = found?.url ? found.url : normalizeUrl(rest)
    navigateUrl(targetUrl, flags.self)
    pushEntry(effectiveLine, [
      {
        text: `正在前往 ${targetUrl}`,
        to: targetUrl.startsWith('/') ? targetUrl : undefined,
        href: targetUrl.startsWith('/') ? undefined : targetUrl,
        selfTab: flags.self,
      },
    ])
    return
  }

  if (cmd === 'bm' || cmd === 'bookmarks') {
    const sub = (args[0] || '').toLowerCase()
    if (sub === 'add' && args.length >= 3) {
      const url = normalizeUrl(args[2])
      const res = space.addLink(resolveVfsPath(cwd.value, args[1]), url)
      if (!res.ok) {
        pushEntry(effectiveLine, [{ text: res.error || '添加失败', tone: 'error' }])
        return
      }
      persistVfs()
      pushEntry(effectiveLine, [{ text: `已添加书签：${args[1]} → ${url}`, tone: 'accent' }])
      return
    }
    if (sub === 'rm' && args[1]) {
      const found = space.allLinks().find((b) => b.name.toLowerCase() === args[1].toLowerCase())
      if (!found) {
        pushEntry(effectiveLine, [{ text: `未找到书签：${args[1]}`, tone: 'warn' }])
        return
      }
      space.remove(found.path, false)
      persistVfs()
      pushEntry(effectiveLine, [{ text: `已删除书签：${found.name}`, tone: 'accent' }])
      return
    }
    pushEntry(
      effectiveLine,
      space.allLinks().map((b) => ({
        text: `  ${b.icon || '🔗'} ${b.name}`,
        meta: `${b.path} → ${b.url}`,
        to: b.url?.startsWith('/') ? b.url : undefined,
        href: b.url?.startsWith('/') ? undefined : b.url,
      }))
    )
    return
  }

  // ---- Todo ----
  if (cmd === 'todo') {
    const sub = (args[0] || '').toLowerCase()
    if (sub === 'add' && args.length >= 2) {
      const text = args.slice(1).join(' ')
      const nextId = todos.value.reduce((m, t) => Math.max(m, t.id), 0) + 1
      todos.value.push({ id: nextId, text, done: false })
      saveTodos()
      pushEntry(effectiveLine, [{ text: `已添加待办 #${nextId}：${text}`, tone: 'accent' }])
      return
    }
    if (sub === 'done' && args[1]) {
      const item = todos.value.find((t) => t.id === Number(args[1]))
      if (!item) {
        pushEntry(effectiveLine, [{ text: `未找到待办 #${args[1]}`, tone: 'warn' }])
        return
      }
      item.done = !item.done
      saveTodos()
      pushEntry(effectiveLine, [{ text: `待办 #${item.id} 已标记为 ${item.done ? '完成' : '未完成'}`, tone: 'accent' }])
      return
    }
    if (sub === 'rm' && args[1]) {
      const before = todos.value.length
      todos.value = todos.value.filter((t) => t.id !== Number(args[1]))
      saveTodos()
      pushEntry(effectiveLine, [
        {
          text: todos.value.length < before ? `已删除待办 #${args[1]}` : `未找到待办 #${args[1]}`,
          tone: todos.value.length < before ? 'accent' : 'warn',
        },
      ])
      return
    }
    if (!todos.value.length) {
      pushEntry(effectiveLine, [{ text: '暂无待办。使用 todo add <事项> 添加一条吧！', tone: 'muted' }])
      return
    }
    pushEntry(
      effectiveLine,
      todos.value.map((t) => ({
        text: `  [${t.done ? '✓' : ' '}] #${t.id} ${t.text}`,
        tone: t.done ? 'muted' : 'normal',
      }))
    )
    return
  }

  // ---- System & Appearance: whoami / date / history / bg / theme / hint / reset ----
  if (cmd === 'date') {
    pushEntry(effectiveLine, [{ text: new Date().toLocaleString('zh-CN', { hour12: false }), tone: 'accent' }])
    return
  }

  if (cmd === 'history' || cmd === 'h') {
    const lines: OutputLine[] = history.value.slice(-20).map((h, i) => ({
      text: `  ${String(i + 1).padStart(2, ' ')}  ${h}`,
      meta: `输入 !${i + 1} 重跑`,
      runCmd: h,
    }))
    pushEntry(effectiveLine, lines.length ? lines : [{ text: '暂无历史记录', tone: 'muted' }])
    return
  }

  if (cmd === 'whoami' || cmd === 'about' || cmd === 'info') {
    pushEntry(effectiveLine, [
      { text: `${siteConfig.blogger.name} — ${siteConfig.blogger.slogan}`, tone: 'accent' },
      { text: `站点：${siteConfig.title}（${siteConfig.hero.title}）` },
      { text: `统计：${siteStats.value.posts} 篇已发布文章 · ${siteStats.value.categories} 个分类`, tone: 'muted' },
      { text: '→ 打开完整「关于」页面', to: '/pages/about' },
    ])
    return
  }

  if (cmd === 'bg' || cmd === 'background') {
    const arg = (args[0] || '').trim()
    if (!arg) {
      pushEntry(effectiveLine, [
        {
          text: `当前背景：${bgImage.value || bgMode.value}。用法：bg [grid|emerald|midnight|slate|random|<图片URL>]`,
          tone: 'muted',
        },
      ])
      return
    }
    const lower = arg.toLowerCase()
    if (['grid', 'emerald', 'midnight', 'slate'].includes(lower)) {
      setBgPreset(lower as BgMode)
      pushEntry(effectiveLine, [{ text: `已切换终端背景为预设「${lower}」`, tone: 'accent' }])
      return
    }
    if (lower === 'random') {
      const choices: BgMode[] = ['grid', 'emerald', 'midnight', 'slate']
      const others = choices.filter((c) => c !== bgMode.value)
      const picked = others[Math.floor(Math.random() * others.length)]
      setBgPreset(picked)
      pushEntry(effectiveLine, [{ text: `已随机切换背景为「${picked}」`, tone: 'accent' }])
      return
    }
    if (/^(https?:\/\/|\/img\/)/i.test(arg)) {
      setBgImageUrl(arg)
      pushEntry(effectiveLine, [{ text: `已设置自定义壁纸：${arg}`, tone: 'accent' }])
      return
    }
    pushEntry(effectiveLine, [{ text: '无效背景参数，请输入预设名、random 或以 https:// 开头的图片网址', tone: 'error' }])
    return
  }

  if (cmd === 'theme') {
    const target = args[0]?.toLowerCase() as ThemeMode | undefined
    if (!target || !['system', 'light', 'dark'].includes(target)) {
      pushEntry(effectiveLine, [{ text: `当前全站主题：${theme.mode}。可选：system | light | dark`, tone: 'muted' }])
      return
    }
    theme.setMode(target)
    pushEntry(effectiveLine, [{ text: `已切换全站主题为 ${target}`, tone: 'accent' }])
    return
  }

  if (cmd === 'hint') {
    const val = (args[0] || '').toLowerCase()
    if (val === 'on' || val === 'off') {
      hintEnabled.value = val === 'on'
      localStorage.setItem(HINT_KEY, val)
    } else {
      hintEnabled.value = !hintEnabled.value
      localStorage.setItem(HINT_KEY, hintEnabled.value ? 'on' : 'off')
    }
    pushEntry(effectiveLine, [
      { text: `实时命令语法提示已${hintEnabled.value ? '开启' : '关闭'}`, tone: 'accent' },
    ])
    return
  }

  if (cmd === 'reset') {
    setBgPreset('grid')
    hintEnabled.value = true
    localStorage.setItem(HINT_KEY, 'on')
    space.resetUserNodes(defaultUserNodes())
    persistVfs()
    cwd.value = '/'
    if (flags.all) {
      todos.value = []
      saveTodos()
    }
    pushEntry(effectiveLine, [
      {
        text: `已恢复默认壁纸、书签空间与提示开关${flags.all ? '（并已清空待办）' : '（待办保留，加 --all 可一并清空）'}`,
        tone: 'accent',
      },
    ])
    return
  }

  pushEntry(effectiveLine, [
    { text: `未知命令：${cmd}。输入 help 查看所有指令，或输入 g ${effectiveLine} 搜索全网。`, tone: 'error' },
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

function onGlobalKeydown(e: KeyboardEvent) {
  if (e.ctrlKey && e.key.toLowerCase() === 'o') {
    e.preventDefault()
    toggleAllCollapse()
    return
  }
  const active = document.activeElement
  const isTypingElsewhere =
    active instanceof HTMLInputElement ||
    active instanceof HTMLTextAreaElement ||
    (active instanceof HTMLElement && active.isContentEditable)
  if (!isTypingElsewhere && !e.ctrlKey && !e.metaKey && !e.altKey && /^[a-zA-Z0-9/!]$/.test(e.key)) {
    focusInput()
  }
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
    description: `${siteConfig.title} 极客终端浏览器起始页，支持分类目录漫游、热文榜、番茄钟、命名助手与多源检索。`,
    path: '/terminal',
  })
  tickClock()
  clockTimer = window.setInterval(tickClock, 1000)
  window.addEventListener('keydown', onGlobalKeydown)

  pushEntry(undefined, [
    { text: `Welcome to ${siteConfig.title} Geek Start Page (v2.1.2)`, tone: 'accent' },
    { text: `${siteConfig.hero.title} — ${siteConfig.hero.tagline}`, tone: 'muted' },
    { text: '试试点击上方快捷按钮（hot / cd code / timing 25m / varbook 用户配置），或输入 help 查看指南。', tone: 'muted' },
  ])

  try {
    const [p, treeRes, flatRes] = await Promise.all([
      getPublishedPosts({ page: 1, page_size: 1 }),
      getCategories(true),
      getCategories(false),
    ])
    space.setMountedCategories(treeRes.data)
    vfsVersion.value++
    siteStats.value = { posts: p.data.total, categories: flatRes.data.length }
  } catch {
    /* ignore */
  }
  focusInput()
})

onBeforeUnmount(() => {
  if (clockTimer) window.clearInterval(clockTimer)
  window.removeEventListener('keydown', onGlobalKeydown)
})
</script>

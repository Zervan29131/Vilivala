<template>
  <header class="sticky top-0 z-40 nav-blur border-b border-gray-200/80 dark:border-white/10">
    <div class="max-w-7xl mx-auto px-3.5 sm:px-4 h-14 flex items-center justify-between gap-3 sm:gap-4">
      <!-- Brand -->
      <router-link to="/" class="flex items-center gap-2.5 shrink-0 group min-w-0">
        <img :src="siteConfig.logo" :alt="site.config.site_title" class="w-8 h-8 rounded-full object-cover shrink-0" />
        <span class="text-[16px] sm:text-[17px] font-bold tracking-tight text-[#213547] dark:text-white group-hover:text-[#42b883] transition-colors truncate">
          {{ site.config.site_title }}
        </span>
      </router-link>

      <!-- Search (tablet & desktop) -->
      <form class="hidden md:flex items-center flex-1 max-w-xs" @submit.prevent="submitSearch">
        <div class="relative w-full">
          <input
            v-model="kw"
            class="w-full pl-9 pr-14 py-1.5 text-sm rounded-full bg-gray-100/90 dark:bg-white/10 text-[#213547] dark:text-gray-100 placeholder:text-gray-400 border border-transparent focus:border-[#42b883] focus:bg-white dark:focus:bg-[#242424] focus:outline-none transition"
            placeholder="搜索文章、关键词…"
          />
          <svg class="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-4.35-4.35m0 0A7.5 7.5 0 1010.5 18a7.5 7.5 0 006.15-3.35z" />
          </svg>
          <button
            type="button"
            class="absolute right-2 top-1/2 -translate-y-1/2 px-1.5 py-0.5 text-[10px] font-mono rounded bg-white/80 dark:bg-white/10 border border-gray-200 dark:border-white/15 text-gray-400 hover:text-[#42b883] hover:border-[#42b883] transition-colors"
            title="打开全文检索浮层 (Ctrl+K / ⌘K)"
            @click="$emit('open-search')"
          >
            ⌘K
          </button>
        </div>
      </form>

      <!-- Desktop Nav (lg+) -->
      <nav class="hidden lg:flex items-center gap-1 text-sm">
        <router-link to="/" class="nav-item" exact-active-class="nav-active">首页</router-link>

        <!-- Dynamic 3-level category dropdowns -->
        <div v-for="root in topCategories" :key="root.id" class="relative group">
          <router-link
            :to="`/categories/${root.slug}`"
            class="nav-item inline-flex items-center gap-1"
            :class="{ 'nav-active': route.path.startsWith(`/categories/${root.slug}`) }"
          >
            {{ root.name }}
            <svg v-if="root.children?.length" class="w-3.5 h-3.5 opacity-60 transition-transform group-hover:rotate-180" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </router-link>

          <div
            v-if="root.children?.length"
            class="invisible opacity-0 group-hover:visible group-hover:opacity-100 transition-all duration-150 absolute left-0 top-full pt-2 z-50 min-w-[200px]"
          >
            <div class="card p-2 shadow-lg space-y-0.5">
              <div v-for="child in root.children" :key="child.id">
                <router-link
                  :to="`/categories/${child.slug}`"
                  class="flex items-center justify-between px-3 py-1.5 rounded-lg text-sm text-gray-700 dark:text-gray-200 hover:bg-[#42b883]/10 hover:text-[#42b883] transition-colors"
                >
                  <span>{{ child.name }}</span>
                  <span class="text-xs text-gray-400">{{ child.total_count }}</span>
                </router-link>
                <div v-if="child.children?.length" class="pl-3 mt-0.5 space-y-0.5 border-l border-gray-200 dark:border-white/10 ml-3">
                  <router-link
                    v-for="grand in child.children"
                    :key="grand.id"
                    :to="`/categories/${grand.slug}`"
                    class="flex items-center justify-between px-2.5 py-1 rounded-md text-xs text-gray-500 dark:text-gray-400 hover:bg-[#42b883]/10 hover:text-[#42b883] transition-colors"
                  >
                    <span>{{ grand.name }}</span>
                    <span class="opacity-70">{{ grand.total_count }}</span>
                  </router-link>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 索引 Dropdown -->
        <div class="relative group">
          <button type="button" class="nav-item inline-flex items-center gap-1">
            索引
            <svg class="w-3.5 h-3.5 opacity-60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </button>
          <div class="invisible opacity-0 group-hover:visible group-hover:opacity-100 transition-all duration-150 absolute right-0 top-full pt-2 z-50 w-36">
            <div class="card p-1.5 shadow-lg space-y-0.5">
              <router-link to="/archives" class="drop-link">时间归档</router-link>
              <router-link to="/categories" class="drop-link">分类与标签</router-link>
            </div>
          </div>
        </div>

        <!-- 更多 Dropdown -->
        <div class="relative group">
          <button type="button" class="nav-item inline-flex items-center gap-1">
            更多
            <svg class="w-3.5 h-3.5 opacity-60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </button>
          <div class="invisible opacity-0 group-hover:visible group-hover:opacity-100 transition-all duration-150 absolute right-0 top-full pt-2 z-50 w-36">
            <div class="card p-1.5 shadow-lg space-y-0.5">
              <router-link to="/pages/about" class="drop-link">关于</router-link>
              <router-link to="/pages/links" class="drop-link">友情链接</router-link>
              <router-link to="/pages/changelog" class="drop-link">更新日志</router-link>
              <router-link to="/pages/sponsor" class="drop-link">赞助支持</router-link>
              <a href="/rss.xml" target="_blank" rel="noopener" class="drop-link">RSS 订阅</a>
            </div>
          </div>
        </div>

        <router-link
          to="/terminal"
          class="nav-item font-mono text-xs inline-flex items-center gap-1"
          title="终端风格浏览器起始页"
        >
          <span class="text-[#42b883]">&gt;_</span>
          <span>终端</span>
        </router-link>

        <router-link to="/admin" class="nav-item">后台</router-link>

        <!-- Theme Switcher -->
        <button
          type="button"
          class="ml-1 px-2.5 py-1.5 rounded-lg text-xs text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-white/10 transition-colors inline-flex items-center gap-1"
          :title="`当前主题：${themeLabel}（点击切换）`"
          @click="theme.cycle()"
        >
          <span>{{ themeIcon }}</span>
          <span class="hidden xl:inline">{{ themeLabel }}</span>
        </button>
      </nav>

      <!-- Mobile & Tablet Actions (< lg) -->
      <div class="flex lg:hidden items-center gap-1 shrink-0">
        <button
          type="button"
          class="md:hidden w-9 h-9 rounded-lg flex items-center justify-center text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-white/10"
          aria-label="打开搜索"
          @click="$emit('open-search')"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-4.35-4.35m0 0A7.5 7.5 0 1010.5 18a7.5 7.5 0 006.15-3.35z" />
          </svg>
        </button>
        <button
          type="button"
          class="w-9 h-9 rounded-lg flex items-center justify-center text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-white/10"
          :title="themeLabel"
          aria-label="切换主题"
          @click="theme.cycle()"
        >
          {{ themeIcon }}
        </button>
        <button
          type="button"
          class="w-9 h-9 rounded-lg flex items-center justify-center text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-white/10"
          :aria-label="mobileOpen ? '关闭菜单' : '打开菜单'"
          :aria-expanded="mobileOpen"
          @click="mobileOpen = !mobileOpen"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path v-if="!mobileOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Mobile / Tablet Overlay Drawer -->
    <teleport to="body">
      <div
        v-if="mobileOpen"
        class="lg:hidden fixed inset-0 top-14 z-40 bg-black/40 backdrop-blur-[1px]"
        @click.self="mobileOpen = false"
      >
        <nav
          class="bg-white dark:bg-[#1e1e1e] border-b border-gray-200 dark:border-white/10 shadow-xl px-4 py-4 space-y-3 max-h-[calc(100dvh-3.5rem)] overflow-y-auto overscroll-contain pb-safe"
          aria-label="移动端导航"
        >
          <form class="md:hidden" @submit.prevent="submitSearch">
            <input
              v-model="kw"
              class="input-field !text-[16px]"
              placeholder="搜索文章、关键词…"
            />
          </form>

          <div class="flex flex-col space-y-1 text-sm">
            <router-link to="/" class="drop-link" @click="mobileOpen = false">首页</router-link>

            <!-- Accordion 3-level Categories -->
            <div v-for="root in topCategories" :key="root.id" class="space-y-0.5">
              <div class="flex items-center justify-between">
                <router-link
                  :to="`/categories/${root.slug}`"
                  class="drop-link flex-1"
                  @click="mobileOpen = false"
                >
                  {{ root.name }} <span class="text-xs text-gray-400">({{ root.total_count }})</span>
                </router-link>
                <button
                  v-if="root.children?.length"
                  type="button"
                  class="px-3 py-2 text-xs text-gray-400 hover:text-[#42b883]"
                  :aria-label="`展开 ${root.name} 子分类`"
                  @click="toggleCatExpand(root.id)"
                >
                  <span class="inline-block transition-transform" :class="{ 'rotate-180': expandedCats.has(root.id) }">▾</span>
                </button>
              </div>

              <div
                v-if="root.children?.length && expandedCats.has(root.id)"
                class="pl-4 ml-3 border-l border-gray-200 dark:border-white/10 space-y-0.5"
              >
                <div v-for="child in root.children" :key="child.id">
                  <router-link
                    :to="`/categories/${child.slug}`"
                    class="drop-link !py-1.5 !text-xs flex items-center justify-between"
                    @click="mobileOpen = false"
                  >
                    <span>{{ child.name }}</span>
                    <span class="text-gray-400">{{ child.total_count }}</span>
                  </router-link>
                  <div v-if="child.children?.length" class="pl-3 ml-2 border-l border-gray-200/70 dark:border-white/10">
                    <router-link
                      v-for="grand in child.children"
                      :key="grand.id"
                      :to="`/categories/${grand.slug}`"
                      class="drop-link !py-1 !text-xs text-gray-500 dark:text-gray-400 flex items-center justify-between"
                      @click="mobileOpen = false"
                    >
                      <span>{{ grand.name }}</span>
                      <span class="opacity-70">{{ grand.total_count }}</span>
                    </router-link>
                  </div>
                </div>
              </div>
            </div>

            <div class="pt-1 border-t border-gray-100 dark:border-white/10" />
            <router-link to="/archives" class="drop-link" @click="mobileOpen = false">时间归档</router-link>
            <router-link to="/categories" class="drop-link" @click="mobileOpen = false">分类与标签</router-link>
            <router-link to="/terminal" class="drop-link font-mono" @click="mobileOpen = false">&gt;_ 终端起始页</router-link>
            <router-link to="/pages/about" class="drop-link" @click="mobileOpen = false">关于</router-link>
            <router-link to="/pages/links" class="drop-link" @click="mobileOpen = false">友情链接</router-link>
            <router-link to="/pages/changelog" class="drop-link" @click="mobileOpen = false">更新日志</router-link>
            <router-link to="/pages/sponsor" class="drop-link" @click="mobileOpen = false">赞助支持</router-link>
            <router-link to="/admin" class="drop-link text-[#42b883]" @click="mobileOpen = false">后台管理</router-link>
          </div>
        </nav>
      </div>
    </teleport>
  </header>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getCategories } from '@/api'
import type { Category } from '@/api/types'
import { siteConfig } from '@/site.config'
import { useSiteStore } from '@/stores/site'
import { useThemeStore } from '@/stores/theme'

defineEmits<{ 'open-search': [] }>()

const route = useRoute()
const router = useRouter()
const site = useSiteStore()
const theme = useThemeStore()

const topCategories = ref<Category[]>([])
const kw = ref('')
const mobileOpen = ref(false)
const expandedCats = ref<Set<number>>(new Set())

const themeIcon = computed(() =>
  theme.mode === 'dark' ? '🌙' : theme.mode === 'light' ? '☀️' : '🖥️'
)
const themeLabel = computed(() =>
  theme.mode === 'dark' ? '深色' : theme.mode === 'light' ? '浅色' : '跟随系统'
)

function toggleCatExpand(id: number) {
  const next = new Set(expandedCats.value)
  if (next.has(id)) next.delete(id)
  else next.add(id)
  expandedCats.value = next
}

function submitSearch() {
  mobileOpen.value = false
  router.push({ path: '/', query: { keyword: kw.value.trim() || undefined } })
}

function onKeydown(e: KeyboardEvent) {
  if (mobileOpen.value && e.key === 'Escape') {
    mobileOpen.value = false
  }
}

watch(
  () => route.query.keyword,
  (v) => {
    kw.value = (v as string) || ''
  },
  { immediate: true }
)

watch(
  () => route.fullPath,
  () => {
    mobileOpen.value = false
  }
)

watch(mobileOpen, (open) => {
  document.body.style.overflow = open ? 'hidden' : ''
})

onMounted(async () => {
  window.addEventListener('keydown', onKeydown)
  try {
    const res = await getCategories(true)
    const order = ['code', 'toread', 'fun']
    topCategories.value = [...res.data].sort((a, b) => {
      const ia = order.indexOf(a.slug.toLowerCase())
      const ib = order.indexOf(b.slug.toLowerCase())
      return (ia === -1 ? 99 : ia) - (ib === -1 ? 99 : ib)
    })
  } catch {
    /* ignore */
  }
})

onBeforeUnmount(() => {
  document.body.style.overflow = ''
  window.removeEventListener('keydown', onKeydown)
})
</script>

<style scoped>
.nav-blur {
  @apply bg-white/95 dark:bg-[#1a1a1a]/95;
}
@supports ((-webkit-backdrop-filter: blur(12px)) or (backdrop-filter: blur(12px))) {
  .nav-blur {
    @apply bg-white/85 dark:bg-[#1a1a1a]/85;
    -webkit-backdrop-filter: blur(12px);
    backdrop-filter: blur(12px);
  }
}
.nav-item {
  @apply px-3 py-1.5 rounded-lg text-gray-600 dark:text-gray-300 hover:text-[#42b883] dark:hover:text-[#4fc08d] transition-colors;
}
.nav-active {
  @apply text-[#42b883] dark:text-[#4fc08d] font-medium;
}
.drop-link {
  @apply block px-3 py-2 rounded-lg text-sm text-gray-700 dark:text-gray-200 hover:bg-[#42b883]/10 hover:text-[#42b883] transition-colors;
}
</style>

<template>
  <div class="min-h-dvh-safe flex flex-col md:flex-row bg-[#f6f6f7] dark:bg-[#1a1a1a]">
    <!-- Mobile Sticky Top Header (< md) -->
    <header class="md:hidden sticky top-0 z-40 h-13 px-4 py-2.5 bg-white/95 dark:bg-[#242424]/95 border-b border-gray-200/80 dark:border-white/10 flex items-center justify-between gap-3">
      <div class="flex items-center gap-2 min-w-0">
        <img :src="siteConfig.logo" alt="Logo" class="w-6 h-6 rounded-full shrink-0" />
        <span class="font-bold text-sm text-[#213547] dark:text-white truncate">后台 · {{ currentLabel }}</span>
      </div>
      <div class="flex items-center gap-1 shrink-0">
        <button
          type="button"
          class="w-9 h-9 rounded-lg flex items-center justify-center text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-white/10"
          aria-label="切换主题"
          @click="theme.cycle()"
        >
          {{ theme.mode === 'dark' ? '🌙' : theme.mode === 'light' ? '☀️' : '🖥️' }}
        </button>
        <button
          type="button"
          class="w-9 h-9 rounded-lg flex items-center justify-center text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-white/10"
          :aria-label="drawerOpen ? '关闭管理菜单' : '打开管理菜单'"
          @click="drawerOpen = !drawerOpen"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path v-if="!drawerOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </header>

    <!-- Mobile Drawer (< md) -->
    <div
      v-if="drawerOpen"
      class="md:hidden fixed inset-0 top-[52px] z-40 bg-black/40"
      @click.self="drawerOpen = false"
    >
      <nav class="bg-white dark:bg-[#242424] border-b border-gray-200 dark:border-white/10 p-3 space-y-1 text-sm shadow-xl">
        <router-link
          v-for="l in links"
          :key="l.to"
          :to="l.to"
          class="alink"
          active-class="alink-active"
          @click="drawerOpen = false"
        >
          {{ l.label }}
        </router-link>
        <div class="pt-2 mt-2 border-t border-gray-100 dark:border-white/10 flex items-center justify-between px-3 py-1.5 text-xs text-gray-400">
          <span>当前登录：{{ auth.user?.nickname || auth.user?.username }}</span>
          <router-link to="/" class="text-[#42b883]" @click="drawerOpen = false">返回前台 →</router-link>
        </div>
        <button class="alink w-full text-left text-red-500 hover:!bg-red-500/10" @click="logout">
          退出登录
        </button>
      </nav>
    </div>

    <!-- Tablet & Desktop Sidebar (md+) -->
    <aside class="hidden md:flex md:w-56 md:shrink-0 md:sticky md:top-0 md:h-screen bg-white dark:bg-[#242424] border-r border-gray-200/80 dark:border-white/10 p-4 flex-col">
      <div class="font-bold text-base mb-6 flex items-center gap-2 text-[#213547] dark:text-white">
        <img :src="siteConfig.logo" alt="Logo" class="w-6 h-6 rounded-full" />
        <span>管理后台</span>
      </div>
      <nav class="flex flex-col gap-1 flex-1 text-sm overflow-y-auto">
        <router-link v-for="l in links" :key="l.to" :to="l.to" class="alink" active-class="alink-active">
          {{ l.label }}
        </router-link>
        <router-link to="/" class="alink">← 返回前台</router-link>
      </nav>
      <div class="mt-6 pt-4 border-t border-gray-200/80 dark:border-white/10 text-sm flex flex-col gap-1">
        <button type="button" class="alink text-left" @click="theme.cycle()">
          {{ theme.mode === 'dark' ? '🌙 深色' : theme.mode === 'light' ? '☀️ 浅色' : '🖥️ 跟随系统' }}
        </button>
        <div class="px-3 py-1 text-xs text-gray-400 truncate">
          {{ auth.user?.nickname || auth.user?.username }}
        </div>
        <button class="alink w-full text-left text-red-500 hover:!bg-red-500/10" @click="logout">退出登录</button>
      </div>
    </aside>

    <main class="flex-1 p-3.5 sm:p-6 lg:p-8 min-w-0 pb-safe">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { siteConfig } from '@/site.config'

const auth = useAuthStore()
const theme = useThemeStore()
const route = useRoute()
const router = useRouter()
const drawerOpen = ref(false)

const links = [
  { to: '/admin/dashboard', label: '数据概览' },
  { to: '/admin/posts', label: '文章管理' },
  { to: '/admin/categories', label: '分类管理' },
  { to: '/admin/tags', label: '标签管理' },
  { to: '/admin/pages', label: '独立页面' },
  { to: '/admin/settings', label: '站点设置' },
]

const currentLabel = computed(() => {
  const hit = links.find((l) => route.path.startsWith(l.to))
  return hit ? hit.label : '文章编辑'
})

watch(
  () => route.fullPath,
  () => {
    drawerOpen.value = false
  }
)

function logout() {
  drawerOpen.value = false
  auth.logout()
  router.push('/admin/login')
}
</script>

<style scoped>
.alink {
  @apply block px-3 py-2 rounded-lg text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-white/5 transition-colors whitespace-nowrap;
}
.alink-active {
  @apply bg-[#42b883] text-white hover:!bg-[#33a06f];
}
</style>

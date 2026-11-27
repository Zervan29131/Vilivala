<template>
  <div class="min-h-screen flex flex-col md:flex-row bg-[#f6f6f7] dark:bg-[#1a1a1a]">
    <aside class="md:w-56 bg-white dark:bg-[#242424] border-b md:border-b-0 md:border-r border-gray-200/80 dark:border-white/10 p-4 flex md:flex-col gap-2 md:gap-0 items-center md:items-stretch overflow-x-auto">
      <div class="font-bold text-base md:mb-6 shrink-0 mr-4 md:mr-0 flex items-center gap-2 text-[#213547] dark:text-white">
        <img :src="siteConfig.logo" alt="Logo" class="w-6 h-6 rounded-full" />
        <span>管理后台</span>
      </div>
      <nav class="flex md:flex-col gap-1 flex-1 text-sm">
        <router-link v-for="l in links" :key="l.to" :to="l.to" class="alink" active-class="alink-active">
          {{ l.label }}
        </router-link>
        <router-link to="/" class="alink">← 返回前台</router-link>
      </nav>
      <div class="md:mt-6 md:pt-4 md:border-t border-gray-200/80 dark:border-white/10 text-sm shrink-0 flex md:flex-col gap-1">
        <button type="button" class="alink text-left" @click="theme.cycle()">
          {{ theme.mode === 'dark' ? '🌙 深色' : theme.mode === 'light' ? '☀️ 浅色' : '🖥️ 跟随系统' }}
        </button>
        <div class="hidden md:block px-3 py-1 text-xs text-gray-400">
          {{ auth.user?.nickname || auth.user?.username }}
        </div>
        <button class="alink w-full text-left text-red-500 hover:!bg-red-500/10" @click="logout">退出登录</button>
      </div>
    </aside>
    <main class="flex-1 p-4 sm:p-8 min-w-0">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { siteConfig } from '@/site.config'

const auth = useAuthStore()
const theme = useThemeStore()
const router = useRouter()

const links = [
  { to: '/admin/dashboard', label: '数据概览' },
  { to: '/admin/posts', label: '文章管理' },
  { to: '/admin/categories', label: '分类管理' },
  { to: '/admin/tags', label: '标签管理' },
  { to: '/admin/pages', label: '独立页面' },
  { to: '/admin/settings', label: '站点设置' },
]

function logout() {
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

<template>
  <div class="min-h-screen flex flex-col">
    <NavBar />
    <main class="flex-1">
      <router-view />
    </main>
    <footer class="border-t border-gray-200/70 dark:border-white/10 py-8 mt-12">
      <div class="max-w-7xl mx-auto px-4 flex flex-col sm:flex-row items-center justify-between gap-3 text-xs text-gray-400">
        <div>© {{ new Date().getFullYear() }} {{ siteConfig.title }} · {{ siteConfig.description }}</div>
        <div class="flex items-center gap-4">
          <router-link to="/pages/about" class="hover:text-[#42b883] transition-colors">关于</router-link>
          <router-link to="/pages/links" class="hover:text-[#42b883] transition-colors">友链</router-link>
          <router-link to="/pages/changelog" class="hover:text-[#42b883] transition-colors">日志</router-link>
          <router-link to="/pages/sponsor" class="hover:text-[#42b883] transition-colors">支持</router-link>
          <a :href="siteConfig.githubUrl" target="_blank" rel="noopener" class="hover:text-[#42b883] transition-colors">GitHub</a>
        </div>
      </div>
    </footer>

    <!-- Floating bottom-right theme button (like original Vilivala) -->
    <button
      type="button"
      class="fixed bottom-5 right-5 z-30 w-10 h-10 rounded-full bg-white dark:bg-[#2a2a2a] border border-gray-200 dark:border-white/15 shadow-md flex items-center justify-center text-base hover:scale-105 hover:border-[#42b883] transition"
      :title="`切换主题（当前：${theme.mode}）`"
      @click="theme.cycle()"
    >
      {{ theme.mode === 'dark' ? '🌙' : theme.mode === 'light' ? '☀️' : '🖥️' }}
    </button>
  </div>
</template>

<script setup lang="ts">
import NavBar from '@/components/NavBar.vue'
import { siteConfig } from '@/site.config'
import { useThemeStore } from '@/stores/theme'

const theme = useThemeStore()
</script>

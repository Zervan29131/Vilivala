<template>
  <div class="min-h-screen flex flex-col">
    <ReadingProgress v-if="showReadingProgress" />
    <NavBar @open-search="searchOpen = true" />
    <main class="flex-1">
      <router-view />
    </main>
    <footer class="border-t border-gray-200/70 dark:border-white/10 py-8 mt-12">
      <div class="max-w-7xl mx-auto px-4 space-y-2.5 text-xs text-gray-400">
        <div class="flex flex-col sm:flex-row items-center justify-between gap-3">
          <div class="flex flex-wrap items-center justify-center sm:justify-start gap-x-2 gap-y-1">
            <span>© {{ yearRange }} {{ site.config.site_title }} · {{ site.config.site_description }}</span>
            <template v-if="site.config.stats?.published_posts">
              <span>·</span>
              <span>共 {{ site.config.stats.published_posts }} 篇文章 · {{ site.config.stats.total_views }} 次阅读</span>
            </template>
          </div>
          <div class="flex flex-wrap items-center justify-center gap-4">
            <router-link to="/terminal" class="hover:text-[#42b883] transition-colors font-mono">&gt;_ 终端</router-link>
            <router-link to="/pages/about" class="hover:text-[#42b883] transition-colors">关于</router-link>
            <router-link to="/pages/links" class="hover:text-[#42b883] transition-colors">友链</router-link>
            <router-link to="/pages/changelog" class="hover:text-[#42b883] transition-colors">日志</router-link>
            <router-link to="/pages/sponsor" class="hover:text-[#42b883] transition-colors">支持</router-link>
            <a href="/rss.xml" target="_blank" rel="noopener" class="hover:text-[#42b883] transition-colors">RSS</a>
            <a :href="siteConfig.githubUrl" target="_blank" rel="noopener" class="hover:text-[#42b883] transition-colors">GitHub</a>
          </div>
        </div>

        <!-- Optional ICP / Gongan filing row (configured in Admin -> 站点设置) -->
        <div
          v-if="site.config.icp_number || site.config.gongan_number"
          class="flex flex-wrap items-center justify-center sm:justify-start gap-4 text-[11px] text-gray-400/90"
        >
          <a
            v-if="site.config.icp_number"
            :href="site.config.icp_url || 'https://beian.miit.gov.cn/'"
            target="_blank"
            rel="noopener noreferrer nofollow"
            class="hover:text-[#42b883] transition-colors"
          >
            {{ site.config.icp_number }}
          </a>
          <a
            v-if="site.config.gongan_number"
            :href="site.config.gongan_url || 'https://www.beian.gov.cn/'"
            target="_blank"
            rel="noopener noreferrer nofollow"
            class="hover:text-[#42b883] transition-colors"
          >
            {{ site.config.gongan_number }}
          </a>
        </div>
      </div>
    </footer>

    <!-- Floating Back-to-top + Theme button -->
    <BackToTop />
    <button
      type="button"
      class="fixed bottom-5 right-5 z-30 w-10 h-10 rounded-full bg-white dark:bg-[#2a2a2a] border border-gray-200 dark:border-white/15 shadow-md flex items-center justify-center text-base hover:scale-105 hover:border-[#42b883] transition"
      :title="`切换主题（当前：${theme.mode}）`"
      aria-label="切换主题"
      @click="theme.cycle()"
    >
      {{ theme.mode === 'dark' ? '🌙' : theme.mode === 'light' ? '☀️' : '🖥️' }}
    </button>

    <!-- Full-text Search Modal (Ctrl/Cmd + K) -->
    <SearchModal :open="searchOpen" @close="searchOpen = false" />
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import NavBar from '@/components/NavBar.vue'
import BackToTop from '@/components/BackToTop.vue'
import ReadingProgress from '@/components/ReadingProgress.vue'
import SearchModal from '@/components/SearchModal.vue'
import { siteConfig } from '@/site.config'
import { useSiteStore } from '@/stores/site'
import { useThemeStore } from '@/stores/theme'

const route = useRoute()
const site = useSiteStore()
const theme = useThemeStore()
const searchOpen = ref(false)

const currentYear = new Date().getFullYear()
const yearRange = computed(() => {
  const since = Number(site.config.since_year) || currentYear
  return since < currentYear ? `${since}–${currentYear}` : `${currentYear}`
})

const showReadingProgress = computed(
  () => route.path.startsWith('/posts/') || route.path.startsWith('/pages/')
)

function onGlobalKeydown(e: KeyboardEvent) {
  if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault()
    searchOpen.value = !searchOpen.value
  }
}

onMounted(() => {
  site.ensureLoaded()
  window.addEventListener('keydown', onGlobalKeydown)
})
onBeforeUnmount(() => window.removeEventListener('keydown', onGlobalKeydown))
</script>

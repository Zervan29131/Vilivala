<template>
  <div class="max-w-3xl mx-auto px-4 py-20 text-center">
    <p class="text-6xl sm:text-7xl font-extrabold tracking-tight text-[#42b883] font-mono select-none">
      404
    </p>
    <h1 class="mt-4 text-2xl sm:text-3xl font-bold text-[#213547] dark:text-white">
      页面不存在
    </h1>
    <p class="mt-3 text-sm text-gray-500 dark:text-gray-400 max-w-md mx-auto leading-relaxed">
      你访问的路径 <code class="px-1.5 py-0.5 rounded bg-gray-100 dark:bg-white/10 font-mono text-xs">{{ route.fullPath }}</code> 可能已被移动、删除或输入有误。
    </p>
    <div class="mt-8 flex flex-wrap items-center justify-center gap-3">
      <router-link to="/" class="btn-primary">返回博客首页</router-link>
      <router-link to="/archives" class="btn-secondary">浏览时间归档</router-link>
      <router-link to="/terminal" class="btn-secondary font-mono">&gt;_ 终端起始页</router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { updatePageMeta } from '@/utils/seo'

const route = useRoute()

onMounted(() => {
  updatePageMeta({
    title: '页面不存在 (404)',
    description: '你访问的页面不存在或已被移除。',
    path: route.path,
  })
})
</script>

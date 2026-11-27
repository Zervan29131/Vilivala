<template>
  <div class="max-w-4xl mx-auto px-4 py-8">
    <div class="flex items-baseline justify-between mb-8">
      <h1 class="text-2xl font-bold text-[#213547] dark:text-white">时间归档</h1>
      <span class="text-sm text-gray-400">共 {{ totalPosts }} 篇文章</span>
    </div>

    <div v-if="loading" class="space-y-4">
      <div v-for="i in 3" :key="i" class="card h-32 animate-pulse bg-gray-200/60 dark:bg-white/5" />
    </div>

    <div v-else class="space-y-8">
      <section v-for="y in archives" :key="y.year" class="card p-6">
        <div class="flex items-baseline justify-between pb-3 mb-4 border-b border-gray-200/80 dark:border-white/10">
          <h2 class="text-xl font-bold text-[#213547] dark:text-white">{{ y.year }} 年</h2>
          <span class="text-xs text-gray-400">{{ y.count }} 篇</span>
        </div>

        <div class="space-y-5">
          <div v-for="m in y.months" :key="m.month">
            <h3 class="text-sm font-semibold text-[#42b883] dark:text-[#4fc08d] mb-2">
              {{ String(m.month).padStart(2, '0') }} 月
            </h3>
            <ul class="space-y-2 border-l-2 border-gray-200 dark:border-white/10 pl-4 ml-1.5">
              <li v-for="p in m.posts" :key="p.id" class="flex items-baseline justify-between gap-4 text-sm">
                <router-link
                  :to="`/posts/${p.slug}`"
                  class="text-[#213547] dark:text-gray-200 hover:text-[#42b883] dark:hover:text-[#4fc08d] transition-colors truncate"
                >
                  {{ p.title }}
                </router-link>
                <span class="text-xs text-gray-400 shrink-0 font-mono">
                  {{ formatDate(p.created_at).slice(5) }}
                </span>
              </li>
            </ul>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { getArchives } from '@/api'
import type { ArchiveYear } from '@/api/types'
import { formatDate } from '@/utils/format'

const archives = ref<ArchiveYear[]>([])
const loading = ref(true)

const totalPosts = computed(() => archives.value.reduce((s, y) => s + y.count, 0))

onMounted(async () => {
  try {
    archives.value = (await getArchives()).data
  } finally {
    loading.value = false
  }
})
</script>

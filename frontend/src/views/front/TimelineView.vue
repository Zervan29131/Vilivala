<template>
  <div class="max-w-4xl mx-auto px-4 py-8">
    <div class="flex items-baseline justify-between mb-8">
      <h1 class="text-2xl font-bold text-[#213547] dark:text-white">时间归档</h1>
      <span class="text-sm text-gray-400">共 {{ totalPosts }} 篇文章</span>
    </div>

    <LoadingState v-if="loading" :rows="3" />
    <ErrorState v-else-if="errorMsg" :message="errorMsg" @retry="load" />
    <EmptyState v-else-if="!archives.length" message="暂无归档文章" />

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
import LoadingState from '@/components/LoadingState.vue'
import ErrorState from '@/components/ErrorState.vue'
import EmptyState from '@/components/EmptyState.vue'

const archives = ref<ArchiveYear[]>([])
const loading = ref(true)
const errorMsg = ref('')

const totalPosts = computed(() => archives.value.reduce((s, y) => s + y.count, 0))

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    archives.value = (await getArchives()).data
  } catch (err) {
    errorMsg.value = err instanceof Error ? err.message : '加载归档失败'
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

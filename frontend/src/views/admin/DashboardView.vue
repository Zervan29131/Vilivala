<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between flex-wrap gap-3">
      <h1 class="text-xl font-bold text-[#213547] dark:text-white">数据概览</h1>
      <div class="flex items-center gap-2">
        <router-link to="/admin/posts/new" class="btn-primary">+ 写新文章</router-link>
        <router-link to="/admin/settings" class="btn-secondary">站点设置</router-link>
      </div>
    </div>

    <LoadingState v-if="loading" :rows="3" height-class="h-28" />

    <template v-else-if="data">
      <!-- KPI Cards -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="card p-5">
          <div class="text-xs text-gray-400">已发布 / 总文章</div>
          <div class="mt-2 flex items-baseline gap-2">
            <span class="text-2xl sm:text-3xl font-bold text-[#42b883]">{{ data.published_posts }}</span>
            <span class="text-sm text-gray-400">/ {{ data.total_posts }}</span>
          </div>
          <div class="mt-1 text-xs text-gray-400">草稿箱：{{ data.draft_posts }} 篇</div>
        </div>

        <div class="card p-5">
          <div class="text-xs text-gray-400">全站累计字数</div>
          <div class="mt-2 text-2xl sm:text-3xl font-bold text-[#213547] dark:text-white">
            {{ formatCompactWords(data.total_words) }}
          </div>
          <div class="mt-1 text-xs text-gray-400">约 {{ Math.max(1, Math.round(data.total_words / 400)) }} 分钟阅读体量</div>
        </div>

        <div class="card p-5">
          <div class="text-xs text-gray-400">累计阅读次数</div>
          <div class="mt-2 text-2xl sm:text-3xl font-bold text-[#213547] dark:text-white">
            {{ data.total_views.toLocaleString() }}
          </div>
          <div class="mt-1 text-xs text-gray-400">篇均 {{ avgViews }} 次阅读</div>
        </div>

        <div class="card p-5">
          <div class="text-xs text-gray-400">内容索引结构</div>
          <div class="mt-2 text-2xl sm:text-3xl font-bold text-[#213547] dark:text-white">
            {{ data.categories }} <span class="text-sm font-normal text-gray-400">分类</span>
          </div>
          <div class="mt-1 text-xs text-gray-400">{{ data.tags }} 个标签 · {{ data.pages }} 个独立页</div>
        </div>
      </div>

      <!-- Category Distribution Bars -->
      <div v-if="data.category_distribution?.length" class="card p-5">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-sm font-semibold text-[#213547] dark:text-white">顶级板块文章占比</h2>
          <router-link to="/admin/categories" class="text-xs text-[#42b883] hover:underline">管理分类 →</router-link>
        </div>
        <div class="space-y-3">
          <div v-for="c in data.category_distribution" :key="c.id" class="space-y-1">
            <div class="flex items-center justify-between text-xs">
              <span class="font-medium text-[#213547] dark:text-gray-200">{{ c.name }}</span>
              <span class="text-gray-400 font-mono">{{ c.total_count }} 篇 ({{ pct(c.total_count) }}%)</span>
            </div>
            <div class="h-2 rounded-full bg-gray-100 dark:bg-white/10 overflow-hidden">
              <div
                class="h-full rounded-full bg-[#42b883] transition-all duration-300"
                :style="{ width: `${Math.max(3, pct(c.total_count))}%` }"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Top 10 Posts & Recent Updates -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div class="card p-5">
          <div class="flex items-center justify-between mb-3">
            <h2 class="text-sm font-semibold text-[#213547] dark:text-white">🔥 阅读量 Top 10</h2>
            <router-link to="/admin/posts" class="text-xs text-[#42b883] hover:underline">全部文章 →</router-link>
          </div>
          <ul class="divide-y divide-gray-100 dark:divide-white/10 text-sm">
            <li v-for="(p, idx) in data.top_posts" :key="p.id" class="py-2.5 flex items-center justify-between gap-3">
              <div class="min-w-0 flex items-center gap-2.5">
                <span
                  class="w-5 h-5 rounded text-[11px] font-mono flex items-center justify-center shrink-0"
                  :class="idx < 3 ? 'bg-[#42b883]/15 text-[#42b883] font-bold' : 'bg-gray-100 dark:bg-white/10 text-gray-400'"
                >
                  {{ idx + 1 }}
                </span>
                <a
                  :href="`/posts/${p.slug}`"
                  target="_blank"
                  rel="noopener"
                  class="truncate text-[#213547] dark:text-gray-200 hover:text-[#42b883]"
                >
                  {{ p.title }}
                </a>
              </div>
              <span class="text-xs text-gray-400 font-mono shrink-0">{{ p.view_count }} 阅读</span>
            </li>
          </ul>
        </div>

        <div class="card p-5">
          <div class="flex items-center justify-between mb-3">
            <h2 class="text-sm font-semibold text-[#213547] dark:text-white">🕒 最近更新</h2>
          </div>
          <ul class="divide-y divide-gray-100 dark:divide-white/10 text-sm">
            <li v-for="p in data.recent_posts" :key="p.id" class="py-2.5 flex items-center justify-between gap-3">
              <div class="min-w-0 flex items-center gap-2">
                <span
                  class="badge shrink-0"
                  :class="
                    p.status === 'published'
                      ? 'bg-emerald-500/15 text-emerald-600 dark:text-emerald-400'
                      : 'bg-amber-500/15 text-amber-600 dark:text-amber-400'
                  "
                >
                  {{ p.status === 'published' ? '已发' : '草稿' }}
                </span>
                <span class="truncate text-[#213547] dark:text-gray-200">{{ p.title }}</span>
              </div>
              <div class="flex items-center gap-2.5 shrink-0 text-xs">
                <span class="text-gray-400 font-mono hidden sm:inline">{{ formatDate(p.updated_at) }}</span>
                <router-link :to="`/posts/${p.slug}`" target="_blank" class="text-gray-500 hover:text-[#42b883]">
                  预览
                </router-link>
                <router-link :to="`/admin/posts/${p.id}/edit`" class="text-[#42b883] hover:underline">
                  编辑
                </router-link>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { adminGetDashboard } from '@/api'
import type { DashboardData } from '@/api/types'
import LoadingState from '@/components/LoadingState.vue'
import { formatDate } from '@/utils/format'

const data = ref<DashboardData | null>(null)
const loading = ref(true)

const avgViews = computed(() => {
  if (!data.value?.published_posts) return 0
  return Math.round(data.value.total_views / data.value.published_posts)
})

function formatCompactWords(n: number): string {
  if (n >= 10000) return `${(n / 10000).toFixed(1)} 万`
  return n.toLocaleString()
}

function pct(count: number): number {
  const base = data.value?.published_posts || 1
  return Math.min(100, Math.round((count / base) * 100))
}

onMounted(async () => {
  try {
    data.value = (await adminGetDashboard()).data
  } catch {
    /* handled by interceptor */
  } finally {
    loading.value = false
  }
})
</script>

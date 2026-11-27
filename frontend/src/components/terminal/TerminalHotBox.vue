<template>
  <div class="my-1.5 rounded-xl border border-white/10 bg-white/[0.04] p-3.5 max-w-2xl" @click.stop>
    <div class="flex items-center justify-between mb-2.5">
      <span class="text-xs font-semibold text-[#4fc08d] flex items-center gap-1.5">
        <span>🔥</span>
        <span>本站热门文章 Top 10</span>
      </span>
      <span class="text-[11px] text-gray-500">点击标题直达阅读</span>
    </div>

    <div v-if="loading" class="py-4 text-xs text-gray-400">正在加载热榜…</div>
    <div v-else-if="!posts.length" class="py-3 text-xs text-gray-400">暂无文章数据</div>
    <ul v-else class="space-y-1.5">
      <li
        v-for="(p, idx) in posts"
        :key="p.id"
        class="flex items-center justify-between gap-3 px-2.5 py-1.5 rounded-lg hover:bg-white/[0.06] transition-colors"
      >
        <div class="min-w-0 flex items-center gap-2.5">
          <span
            class="w-5 h-5 rounded text-[11px] font-bold flex items-center justify-center shrink-0"
            :class="
              idx === 0
                ? 'bg-amber-500/25 text-amber-300'
                : idx === 1
                  ? 'bg-slate-300/20 text-slate-200'
                  : idx === 2
                    ? 'bg-orange-500/20 text-orange-300'
                    : 'bg-white/10 text-gray-400'
            "
          >
            {{ idx + 1 }}
          </span>
          <router-link :to="`/posts/${p.slug}`" class="truncate text-xs sm:text-[13px] text-gray-100 hover:text-[#4fc08d]">
            {{ p.title }}
          </router-link>
        </div>
        <div class="flex items-center gap-2 shrink-0 text-[11px] text-gray-400">
          <span v-if="p.category?.name" class="hidden sm:inline px-1.5 py-0.5 rounded bg-white/5 text-gray-400">
            {{ p.category.name }}
          </span>
          <span class="text-[#4fc08d]/90">{{ p.view_count }} 阅读</span>
        </div>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getPublishedPosts } from '@/api'
import type { Post } from '@/api/types'

const posts = ref<Post[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await getPublishedPosts({ page: 1, page_size: 10, sort: 'views' })
    posts.value = res.data.list
  } catch {
    posts.value = []
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <article class="card p-5 sm:p-6 hover:border-[#42b883]/60 hover:shadow-md transition-all duration-200 group">
    <div class="flex flex-wrap items-center gap-2 text-xs text-gray-400 mb-2">
      <span
        v-if="post.is_top"
        class="badge bg-[#42b883]/15 text-[#42b883] dark:text-[#4fc08d] font-semibold"
      >
        置顶
      </span>
      <span>{{ formatDate(post.created_at) }}</span>
      <span>·</span>
      <span>约 {{ post.reading_time || 1 }} 分钟</span>
      <span>·</span>
      <span>{{ post.view_count }} 阅读</span>
      <template v-if="post.category?.path?.length">
        <span>·</span>
        <span class="inline-flex items-center gap-1">
          <template v-for="(node, idx) in post.category.path" :key="node.id">
            <span v-if="idx > 0" class="opacity-50">/</span>
            <router-link
              :to="`/categories/${node.slug}`"
              class="text-[#42b883] dark:text-[#4fc08d] hover:underline"
            >
              {{ node.name }}
            </router-link>
          </template>
        </span>
      </template>
    </div>

    <router-link :to="`/posts/${post.slug}`" class="block">
      <h2 class="text-lg sm:text-xl font-semibold text-[#213547] dark:text-white group-hover:text-[#42b883] dark:group-hover:text-[#4fc08d] transition-colors mb-2">
        {{ post.title }}
      </h2>
      <p v-if="post.summary" class="text-sm text-gray-500 dark:text-gray-400 leading-6 line-clamp-2">
        {{ post.summary }}
      </p>
    </router-link>

    <div v-if="post.tags?.length" class="mt-3.5 flex flex-wrap gap-1.5">
      <router-link
        v-for="t in post.tags"
        :key="t.id"
        :to="`/tags/${t.slug}`"
        class="badge bg-gray-100 dark:bg-white/10 text-gray-500 dark:text-gray-400 hover:bg-[#42b883]/15 hover:text-[#42b883] transition-colors"
      >
        # {{ t.name }}
      </router-link>
    </div>
  </article>
</template>

<script setup lang="ts">
import type { Post } from '@/api/types'
import { formatDate } from '@/utils/format'

defineProps<{ post: Post }>()
</script>

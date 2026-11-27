<template>
  <div class="max-w-5xl mx-auto px-4 py-8">
    <h1 class="text-2xl font-bold text-[#213547] dark:text-white mb-6">分类与标签</h1>

    <div class="grid md:grid-cols-3 gap-5 mb-12">
      <div v-for="root in categories" :key="root.id" class="card p-5">
        <router-link
          :to="`/categories/${root.slug}`"
          class="flex items-baseline justify-between pb-3 mb-3 border-b border-gray-200/80 dark:border-white/10 group"
        >
          <span class="text-lg font-bold text-[#213547] dark:text-white group-hover:text-[#42b883] transition-colors">
            {{ root.name }}
          </span>
          <span class="text-xs text-gray-400">{{ root.total_count }} 篇</span>
        </router-link>

        <ul class="space-y-1.5 text-sm">
          <li v-for="child in root.children" :key="child.id">
            <router-link
              :to="`/categories/${child.slug}`"
              class="flex items-center justify-between py-1 text-gray-600 dark:text-gray-300 hover:text-[#42b883] transition-colors"
            >
              <span>{{ child.name }}</span>
              <span class="text-xs text-gray-400">{{ child.total_count }}</span>
            </router-link>
            <ul v-if="child.children?.length" class="pl-4 space-y-1 border-l border-gray-200 dark:border-white/10 ml-1">
              <li v-for="grand in child.children" :key="grand.id">
                <router-link
                  :to="`/categories/${grand.slug}`"
                  class="flex items-center justify-between py-0.5 text-xs text-gray-500 dark:text-gray-400 hover:text-[#42b883] transition-colors"
                >
                  <span>{{ grand.name }}</span>
                  <span>{{ grand.total_count }}</span>
                </router-link>
              </li>
            </ul>
          </li>
        </ul>
      </div>
    </div>

    <h2 class="text-xl font-bold text-[#213547] dark:text-white mb-4">全部标签（{{ tags.length }}）</h2>
    <div class="card p-6 flex flex-wrap gap-2">
      <router-link
        v-for="t in tags"
        :key="t.id"
        :to="`/tags/${t.slug}`"
        class="badge !text-sm !px-3 !py-1 bg-gray-100 dark:bg-white/10 text-gray-600 dark:text-gray-300 hover:bg-[#42b883] hover:text-white transition-colors"
      >
        # {{ t.name }} <span class="ml-1 opacity-60">{{ t.post_count }}</span>
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getCategories, getTags } from '@/api'
import type { Category, Tag } from '@/api/types'

const categories = ref<Category[]>([])
const tags = ref<Tag[]>([])

onMounted(async () => {
  const [c, t] = await Promise.all([getCategories(true), getTags()])
  const order = ['code', 'toread', 'fun']
  categories.value = [...c.data].sort((a, b) => {
    const ia = order.indexOf(a.slug.toLowerCase())
    const ib = order.indexOf(b.slug.toLowerCase())
    return (ia === -1 ? 99 : ia) - (ib === -1 ? 99 : ib)
  })
  tags.value = [...t.data].sort((a, b) => (b.post_count || 0) - (a.post_count || 0))
})
</script>

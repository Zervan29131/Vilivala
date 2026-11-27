<template>
  <aside class="space-y-5">
    <BloggerCard />

    <!-- 3-level Category Tree -->
    <div class="card p-5">
      <div class="flex items-center justify-between mb-3">
        <h3 class="text-sm font-semibold text-[#213547] dark:text-white">分类目录</h3>
        <router-link to="/categories" class="text-xs text-[#42b883] hover:underline">全部</router-link>
      </div>
      <ul class="space-y-1">
        <li v-for="root in categories" :key="root.id">
          <router-link
            :to="`/categories/${root.slug}`"
            class="cat-row font-medium"
            :class="{ 'cat-active': route.params.slug === root.slug }"
          >
            <span>{{ root.name }}</span>
            <span class="text-xs text-gray-400">{{ root.total_count }}</span>
          </router-link>
          <ul v-if="root.children?.length" class="pl-3 mt-0.5 space-y-0.5 border-l border-gray-200 dark:border-white/10 ml-2.5">
            <li v-for="child in root.children" :key="child.id">
              <router-link
                :to="`/categories/${child.slug}`"
                class="cat-row text-[13px]"
                :class="{ 'cat-active': route.params.slug === child.slug }"
              >
                <span>{{ child.name }}</span>
                <span class="text-xs text-gray-400">{{ child.total_count }}</span>
              </router-link>
              <ul v-if="child.children?.length" class="pl-3 mt-0.5 space-y-0.5 border-l border-gray-200 dark:border-white/10 ml-2.5">
                <li v-for="grand in child.children" :key="grand.id">
                  <router-link
                    :to="`/categories/${grand.slug}`"
                    class="cat-row text-xs"
                    :class="{ 'cat-active': route.params.slug === grand.slug }"
                  >
                    <span>{{ grand.name }}</span>
                    <span class="text-xs text-gray-400">{{ grand.total_count }}</span>
                  </router-link>
                </li>
              </ul>
            </li>
          </ul>
        </li>
        <li v-if="!categories.length" class="text-xs text-gray-400">暂无分类</li>
      </ul>
    </div>

    <!-- Tag Cloud -->
    <div class="card p-5">
      <div class="flex items-center justify-between mb-3">
        <h3 class="text-sm font-semibold text-[#213547] dark:text-white">热门标签</h3>
        <span class="text-xs text-gray-400">{{ tags.length }}</span>
      </div>
      <div class="flex flex-wrap gap-1.5">
        <router-link
          v-for="t in visibleTags"
          :key="t.id"
          :to="`/tags/${t.slug}`"
          class="badge bg-gray-100 dark:bg-white/10 text-gray-600 dark:text-gray-300 hover:bg-[#42b883] hover:text-white transition-colors"
        >
          # {{ t.name }}
        </router-link>
        <span v-if="!tags.length" class="text-xs text-gray-400">暂无标签</span>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { getCategories, getTags } from '@/api'
import type { Category, Tag } from '@/api/types'
import BloggerCard from '@/components/BloggerCard.vue'

const route = useRoute()
const categories = ref<Category[]>([])
const tags = ref<Tag[]>([])

const visibleTags = computed(() =>
  [...tags.value]
    .sort((a, b) => (b.post_count || 0) - (a.post_count || 0))
    .slice(0, 24)
)

onMounted(async () => {
  const [c, t] = await Promise.all([getCategories(true), getTags()])
  const order = ['code', 'toread', 'fun']
  categories.value = [...c.data].sort((a, b) => {
    const ia = order.indexOf(a.slug.toLowerCase())
    const ib = order.indexOf(b.slug.toLowerCase())
    return (ia === -1 ? 99 : ia) - (ib === -1 ? 99 : ib)
  })
  tags.value = t.data
})
</script>

<style scoped>
.cat-row {
  @apply flex items-center justify-between px-2 py-1 rounded-md text-sm text-gray-600 dark:text-gray-300 hover:bg-[#42b883]/10 hover:text-[#42b883] transition-colors;
}
.cat-active {
  @apply bg-[#42b883]/10 text-[#42b883] font-medium;
}
</style>

<template>
  <div class="max-w-7xl mx-auto px-4 py-8">
    <div class="grid grid-cols-1 lg:grid-cols-[1fr_280px] gap-8">
      <section class="min-w-0">
        <!-- Breadcrumb / Header -->
        <div class="mb-6">
          <div class="flex items-center flex-wrap gap-1.5 text-xs text-gray-400 mb-2">
            <router-link to="/" class="hover:text-[#42b883]">首页</router-link>
            <span>/</span>
            <router-link to="/categories" class="hover:text-[#42b883]">
              {{ filterType === 'category' ? '分类' : '标签' }}
            </router-link>
            <template v-if="breadcrumb.length">
              <template v-for="node in breadcrumb" :key="node.id">
                <span>/</span>
                <router-link :to="`/categories/${node.slug}`" class="hover:text-[#42b883]">
                  {{ node.name }}
                </router-link>
              </template>
            </template>
          </div>
          <div class="flex items-baseline gap-3">
            <h1 class="text-2xl font-bold text-[#213547] dark:text-white">
              {{ filterType === 'tag' ? `# ${heading}` : heading }}
            </h1>
            <span class="text-sm text-gray-400">共 {{ total }} 篇</span>
          </div>
        </div>

        <!-- Child Categories Cards (if viewing a parent category) -->
        <div v-if="childCategories.length" class="grid grid-cols-2 sm:grid-cols-3 gap-3 mb-6">
          <router-link
            v-for="c in childCategories"
            :key="c.id"
            :to="`/categories/${c.slug}`"
            class="card px-4 py-3 hover:border-[#42b883] transition-colors flex items-center justify-between"
          >
            <span class="text-sm font-medium text-[#213547] dark:text-gray-100 truncate">{{ c.name }}</span>
            <span class="text-xs text-gray-400 shrink-0 ml-2">{{ c.total_count }}</span>
          </router-link>
        </div>

        <div v-if="loading" class="space-y-4">
          <div v-for="i in 3" :key="i" class="card h-32 animate-pulse bg-gray-200/60 dark:bg-white/5" />
        </div>
        <div v-else-if="posts.length" class="space-y-4">
          <PostCard v-for="p in posts" :key="p.id" :post="p" />
        </div>
        <div v-else class="card p-12 text-center text-gray-400 text-sm">该目录下暂无文章</div>

        <Pagination :page="page" :page-size="pageSize" :total="total" @change="goPage" />
      </section>

      <SideBar />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getCategories, getPublishedPosts, getTags } from '@/api'
import type { Category, Post } from '@/api/types'
import PostCard from '@/components/PostCard.vue'
import Pagination from '@/components/Pagination.vue'
import SideBar from '@/components/SideBar.vue'

const props = defineProps<{ filterType: 'category' | 'tag'; filterSlug: string }>()

const route = useRoute()
const router = useRouter()
const posts = ref<Post[]>([])
const total = ref(0)
const pageSize = 10
const loading = ref(false)
const heading = ref('')
const breadcrumb = ref<Category[]>([])
const childCategories = ref<Category[]>([])

const page = computed(() => Math.max(1, Number(route.query.page) || 1))

function findInTree(nodes: Category[], slug: string, chain: Category[] = []): { node: Category; chain: Category[] } | null {
  for (const n of nodes) {
    const next = [...chain, n]
    if (n.slug === slug) return { node: n, chain: next }
    if (n.children?.length) {
      const hit = findInTree(n.children, slug, next)
      if (hit) return hit
    }
  }
  return null
}

async function load() {
  loading.value = true
  breadcrumb.value = []
  childCategories.value = []
  try {
    const params: Record<string, unknown> = {
      page: page.value,
      page_size: pageSize,
    }
    if (props.filterType === 'category') {
      const tree = (await getCategories(true)).data
      const hit = findInTree(tree, props.filterSlug)
      if (!hit) {
        heading.value = props.filterSlug
        posts.value = []
        total.value = 0
        return
      }
      heading.value = hit.node.name
      breadcrumb.value = hit.chain
      childCategories.value = hit.node.children || []
      params.category_id = hit.node.id
    } else {
      const tags = (await getTags()).data
      const hit = tags.find((x) => x.slug === props.filterSlug)
      if (!hit) {
        heading.value = props.filterSlug
        posts.value = []
        total.value = 0
        return
      }
      heading.value = hit.name
      params.tag_id = hit.id
    }
    const res = await getPublishedPosts(params)
    posts.value = res.data.list
    total.value = res.data.total
  } finally {
    loading.value = false
  }
}

function goPage(n: number) {
  router.push({ query: { ...route.query, page: n === 1 ? undefined : n } })
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

watch(() => [route.fullPath, props.filterType, props.filterSlug], load, { immediate: true })
</script>

<template>
  <div>
    <!-- Hero Banner (hidden when actively searching) -->
    <section
      v-if="!route.query.keyword && page === 1"
      class="hero-grid text-white py-14 sm:py-20 px-4 border-b border-white/10"
    >
      <div class="max-w-5xl mx-auto text-center">
        <img
          :src="siteConfig.hero.image"
          alt="Hero Logo"
          class="h-14 sm:h-16 mx-auto mb-6 object-contain select-none"
        />
        <h1 class="text-3xl sm:text-5xl font-bold tracking-tight text-white">
          {{ site.config.hero_title }}
        </h1>
        <p class="mt-4 text-sm sm:text-base text-gray-300/90 max-w-xl mx-auto min-h-[1.6rem]">
          <span>{{ typedText }}</span>
          <span v-if="typewriterActive" class="inline-block w-[2px] h-4 bg-[#4fc08d] ml-0.5 align-middle animate-pulse" aria-hidden="true" />
        </p>
        <div class="mt-7">
          <router-link
            :to="siteConfig.hero.actionLink"
            class="inline-flex items-center px-6 py-2.5 rounded-md bg-black/80 hover:bg-[#42b883] text-white text-sm font-medium border border-white/15 hover:border-[#42b883] transition-colors duration-200"
          >
            {{ siteConfig.hero.actionText }}
          </router-link>
        </div>

        <!-- 3 Section Cards -->
        <div class="mt-14 grid grid-cols-1 sm:grid-cols-3 gap-6 max-w-4xl mx-auto">
          <router-link
            v-for="f in siteConfig.features"
            :key="f.slug"
            :to="`/categories/${f.slug}`"
            class="group flex flex-col items-center p-5 rounded-xl bg-white/[0.03] hover:bg-white/[0.08] border border-white/[0.06] hover:border-[#42b883]/50 transition-all duration-200 hover:-translate-y-1"
          >
            <img
              :src="f.imgUrl"
              :alt="f.title"
              class="w-24 h-24 sm:w-28 sm:h-28 object-contain opacity-90 group-hover:opacity-100 group-hover:scale-105 transition duration-200"
            />
            <div class="mt-4 text-lg font-semibold text-white group-hover:text-[#4fc08d] transition-colors">
              {{ f.title }}
            </div>
            <div class="mt-1 text-xs text-gray-400">{{ f.details }}</div>
          </router-link>
        </div>
      </div>
    </section>

    <!-- Main Posts + Sidebar -->
    <div class="max-w-7xl mx-auto px-4 py-8">
      <div class="grid grid-cols-1 lg:grid-cols-[1fr_280px] gap-8">
        <section class="min-w-0">
          <!-- Search Status Bar -->
          <div v-if="route.query.keyword" class="mb-5 flex items-center justify-between card px-4 py-3 text-sm">
            <span>
              搜索 “<strong class="text-[#42b883]">{{ route.query.keyword }}</strong>” 共找到
              <strong>{{ total }}</strong> 篇文章
            </span>
            <router-link to="/" class="text-xs text-gray-400 hover:text-[#42b883]">清除搜索</router-link>
          </div>

          <LoadingState v-if="loading" :rows="4" />
          <ErrorState v-else-if="errorMsg" :message="errorMsg" @retry="load" />
          <div v-else-if="posts.length" class="space-y-4">
            <PostCard v-for="p in posts" :key="p.id" :post="p" />
          </div>
          <EmptyState v-else message="没有找到相关文章" />

          <Pagination :page="page" :page-size="pageSize" :total="total" @change="goPage" />
        </section>

        <SideBar />
      </div>
    </div>

    <!-- Bottom-left Waving Panda Widget (from legacy Vilivala) -->
    <img
      :src="siteConfig.pandaWidget"
      alt=""
      class="hidden xl:block fixed bottom-0 left-2 w-[120px] h-[106px] object-contain opacity-85 hover:opacity-100 transition-opacity pointer-events-none select-none z-20"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getPublishedPosts } from '@/api'
import type { Post } from '@/api/types'
import { siteConfig } from '@/site.config'
import { useSiteStore } from '@/stores/site'
import PostCard from '@/components/PostCard.vue'
import Pagination from '@/components/Pagination.vue'
import SideBar from '@/components/SideBar.vue'
import LoadingState from '@/components/LoadingState.vue'
import ErrorState from '@/components/ErrorState.vue'
import EmptyState from '@/components/EmptyState.vue'

const route = useRoute()
const router = useRouter()
const site = useSiteStore()
const posts = ref<Post[]>([])
const total = ref(0)
const pageSize = 10
const loading = ref(false)
const errorMsg = ref('')

const page = computed(() => Math.max(1, Number(route.query.page) || 1))

// Lightweight native multi-string typewriter (inspired by Bilibala, zero CDN deps)
const typedText = ref(site.config.hero_tagline)
const typewriterActive = ref(false)
let twTimer: number | undefined

function startTypewriter() {
  window.clearTimeout(twTimer)
  const list = site.config.subtitles?.filter(Boolean) || [site.config.hero_tagline]
  const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  if (reduced || list.length <= 1) {
    typedText.value = list[0] || site.config.hero_tagline
    typewriterActive.value = false
    return
  }
  typewriterActive.value = true
  let strIdx = 0
  let charIdx = list[0].length
  let deleting = false
  typedText.value = list[0]

  const step = () => {
    const current = list[strIdx % list.length]
    if (!deleting) {
      charIdx++
      typedText.value = current.slice(0, charIdx)
      if (charIdx >= current.length) {
        deleting = true
        twTimer = window.setTimeout(step, 2600)
        return
      }
      twTimer = window.setTimeout(step, 90)
    } else {
      charIdx--
      typedText.value = current.slice(0, charIdx)
      if (charIdx <= 0) {
        deleting = false
        strIdx = (strIdx + 1) % list.length
        twTimer = window.setTimeout(step, 350)
        return
      }
      twTimer = window.setTimeout(step, 45)
    }
  }
  twTimer = window.setTimeout(() => {
    deleting = true
    step()
  }, 2600)
}

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    const res = await getPublishedPosts({
      page: page.value,
      page_size: pageSize,
      keyword: (route.query.keyword as string) || undefined,
    })
    posts.value = res.data.list
    total.value = res.data.total
  } catch (err) {
    errorMsg.value = err instanceof Error ? err.message : '加载文章失败'
  } finally {
    loading.value = false
  }
}

function goPage(n: number) {
  router.push({ query: { ...route.query, page: n === 1 ? undefined : n } })
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

watch(() => route.fullPath, load, { immediate: true })
watch(() => site.config.subtitles, startTypewriter)
onMounted(startTypewriter)
onBeforeUnmount(() => window.clearTimeout(twTimer))
</script>

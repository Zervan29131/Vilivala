<template>
  <div
    v-if="site.config.copyright_enabled"
    class="mt-6 rounded-xl border border-gray-200/80 dark:border-white/10 bg-gray-50/80 dark:bg-white/[0.03] px-5 py-4 text-xs text-gray-600 dark:text-gray-300 space-y-1.5"
  >
    <div class="flex flex-wrap items-center justify-between gap-2">
      <div>
        <span class="font-semibold text-[#213547] dark:text-gray-100">本文作者：</span>
        <span>{{ site.config.blogger_name }}</span>
      </div>
      <span class="text-gray-400 font-mono">发布于 {{ formatDate(createdAt) }}</span>
    </div>
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="min-w-0 truncate">
        <span class="font-semibold text-[#213547] dark:text-gray-100">本文链接：</span>
        <a :href="fullUrl" class="text-[#42b883] dark:text-[#4fc08d] hover:underline font-mono">{{ fullUrl }}</a>
      </div>
      <button
        type="button"
        class="px-2 py-0.5 rounded border border-gray-300 dark:border-white/15 text-[11px] hover:border-[#42b883] hover:text-[#42b883] transition-colors shrink-0"
        @click="copyUrl"
      >
        {{ copied ? '已复制' : '复制链接' }}
      </button>
    </div>
    <div>
      <span class="font-semibold text-[#213547] dark:text-gray-100">版权声明：</span>
      <span class="text-gray-500 dark:text-gray-400">{{ site.config.copyright_notice }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useSiteStore } from '@/stores/site'
import { copyText } from '@/utils/clipboard'
import { formatDate } from '@/utils/format'
import { toast } from '@/utils/toast'

const props = defineProps<{
  slug: string
  createdAt: string
}>()

const site = useSiteStore()
const copied = ref(false)

const fullUrl = computed(() => `${window.location.origin}/posts/${props.slug}`)

async function copyUrl() {
  if (await copyText(fullUrl.value)) {
    copied.value = true
    toast.success('文章链接已复制')
    setTimeout(() => (copied.value = false), 1800)
  } else {
    toast.error('复制失败')
  }
}
</script>

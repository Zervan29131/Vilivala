<template>
  <div v-if="site.config.notice_enabled" class="card p-4">
    <div class="flex items-center justify-between mb-2">
      <div class="flex items-center gap-1.5 text-xs font-semibold text-[#213547] dark:text-white">
        <span class="text-[#42b883]">📢</span>
        <span>{{ site.config.notice_mode === 'hitokoto' ? '每日一言' : '站点公告' }}</span>
      </div>
      <button
        v-if="site.config.notice_mode === 'hitokoto'"
        type="button"
        class="text-[11px] text-gray-400 hover:text-[#42b883] transition-colors"
        :disabled="fetching"
        @click="loadHitokoto"
      >
        ↻ 换一句
      </button>
    </div>
    <p class="text-xs text-gray-600 dark:text-gray-300 leading-relaxed">
      {{ displayText }}
    </p>
    <p v-if="hitokotoFrom && site.config.notice_mode === 'hitokoto'" class="mt-1.5 text-[11px] text-gray-400 text-right">
      —— {{ hitokotoFrom }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useSiteStore } from '@/stores/site'

const site = useSiteStore()
const hitokotoText = ref('')
const hitokotoFrom = ref('')
const fetching = ref(false)

const displayText = computed(() => {
  if (site.config.notice_mode === 'hitokoto' && hitokotoText.value) {
    return hitokotoText.value
  }
  return site.config.notice_text
})

async function loadHitokoto() {
  fetching.value = true
  const ctrl = new AbortController()
  const timer = window.setTimeout(() => ctrl.abort(), 3000)
  try {
    const resp = await fetch('https://v1.hitokoto.cn/?encode=json', { signal: ctrl.signal })
    if (resp.ok) {
      const data = await resp.json()
      if (data?.hitokoto) {
        hitokotoText.value = data.hitokoto
        hitokotoFrom.value = data.from_who || data.from || ''
      }
    }
  } catch {
    /* fallback to notice_text gracefully */
  } finally {
    window.clearTimeout(timer)
    fetching.value = false
  }
}

watch(
  () => site.config.notice_mode,
  (mode) => {
    if (mode === 'hitokoto' && !hitokotoText.value) loadHitokoto()
  },
  { immediate: true }
)
</script>

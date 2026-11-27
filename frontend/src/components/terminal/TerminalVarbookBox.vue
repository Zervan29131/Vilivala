<template>
  <div class="my-1.5 rounded-xl border border-white/10 bg-white/[0.04] p-3.5 max-w-xl" @click.stop>
    <div class="flex items-center justify-between mb-2.5">
      <span class="text-xs font-semibold text-[#4fc08d]">
        🔤 变量命名建议（词根：{{ parsed.words.join(' · ') || query }}）
      </span>
      <span class="text-[11px] text-gray-500">点击任意行复制</span>
    </div>

    <div v-if="!parsed.variants.length" class="text-xs text-amber-300">
      未能提取有效词根，请输入常见中文开发词汇（如「用户配置」）或英文词组。
    </div>
    <div v-else class="divide-y divide-white/10">
      <button
        v-for="v in parsed.variants"
        :key="v.style"
        type="button"
        class="w-full py-2 px-2.5 flex items-center justify-between gap-4 text-left rounded hover:bg-white/[0.07] transition-colors group"
        @click="copyVal(v.value)"
      >
        <span class="text-xs text-gray-400 shrink-0">{{ v.label }}</span>
        <span class="flex items-center gap-2 min-w-0">
          <code class="text-[13px] text-white group-hover:text-[#4fc08d] truncate">{{ v.value }}</code>
          <span class="text-[11px] text-gray-500 group-hover:text-[#4fc08d] shrink-0">
            {{ copiedVal === v.value ? '已复制 ✓' : '复制' }}
          </span>
        </span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { buildNamingVariants } from '@/utils/terminal/varbook'
import { toast } from '@/utils/toast'

const props = defineProps<{ query: string }>()
const copiedVal = ref('')

const parsed = computed(() => buildNamingVariants(props.query))

async function copyVal(val: string) {
  try {
    await navigator.clipboard.writeText(val)
    copiedVal.value = val
    toast.success(`已复制：${val}`)
    setTimeout(() => {
      if (copiedVal.value === val) copiedVal.value = ''
    }, 1500)
  } catch {
    toast.error('复制失败')
  }
}
</script>

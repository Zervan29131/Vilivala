<template>
  <nav v-if="pages > 1" class="mt-8" aria-label="分页导航">
    <!-- Compact Mobile Pagination (< 640px) -->
    <div class="flex sm:hidden items-center justify-between gap-3 text-sm">
      <button
        type="button"
        class="btn-secondary min-h-[38px] !px-3.5 disabled:opacity-40"
        :disabled="page <= 1"
        @click="$emit('change', page - 1)"
      >
        上一页
      </button>
      <span class="text-xs text-gray-500 dark:text-gray-400 font-mono">
        第 {{ page }} / {{ pages }} 页
      </span>
      <button
        type="button"
        class="btn-secondary min-h-[38px] !px-3.5 disabled:opacity-40"
        :disabled="page >= pages"
        @click="$emit('change', page + 1)"
      >
        下一页
      </button>
    </div>

    <!-- Tablet & Desktop Numbered Pagination (sm+) -->
    <div class="hidden sm:flex items-center justify-center flex-wrap gap-1.5 text-sm">
      <button
        type="button"
        class="btn-secondary min-h-[36px] !px-3 !py-1.5 disabled:opacity-40"
        :disabled="page <= 1"
        @click="$emit('change', page - 1)"
      >
        上一页
      </button>
      <template v-for="(item, idx) in visiblePages" :key="idx">
        <span v-if="item === '...'" class="px-2 text-gray-400">…</span>
        <button
          v-else
          type="button"
          class="min-w-[36px] h-[36px] px-2.5 rounded-lg border transition-colors"
          :class="
            item === page
              ? 'bg-[#42b883] text-white border-[#42b883] font-medium'
              : 'bg-white dark:bg-[#242424] border-gray-200 dark:border-white/10 text-gray-600 dark:text-gray-300 hover:border-[#42b883]'
          "
          @click="$emit('change', item)"
        >
          {{ item }}
        </button>
      </template>
      <button
        type="button"
        class="btn-secondary min-h-[36px] !px-3 !py-1.5 disabled:opacity-40"
        :disabled="page >= pages"
        @click="$emit('change', page + 1)"
      >
        下一页
      </button>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ page: number; pageSize: number; total: number }>()
defineEmits<{ (e: 'change', page: number): void }>()

const pages = computed(() => Math.max(1, Math.ceil(props.total / props.pageSize)))

const visiblePages = computed<(number | '...')[]>(() => {
  const total = pages.value
  const cur = props.page
  if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1)
  const out: (number | '...')[] = [1]
  if (cur > 3) out.push('...')
  for (let i = Math.max(2, cur - 1); i <= Math.min(total - 1, cur + 1); i++) {
    out.push(i)
  }
  if (cur < total - 2) out.push('...')
  out.push(total)
  return out
})
</script>

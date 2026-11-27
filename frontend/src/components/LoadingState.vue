<template>
  <div class="space-y-4" role="status" aria-live="polite" :aria-label="label">
    <div
      v-for="i in rows"
      :key="i"
      class="card animate-pulse bg-gray-200/60 dark:bg-white/5"
      :class="heightClass"
    />
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    rows?: number
    heightClass?: string
    label?: string
  }>(),
  {
    rows: 3,
    heightClass: 'h-32',
    label: '正在加载内容…',
  }
)
</script>

<template>
  <div class="card p-12 text-center space-y-2">
    <p class="text-sm text-gray-500 dark:text-gray-400">{{ message }}</p>
    <slot />
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    message?: string
  }>(),
  {
    message: '暂无内容',
  }
)
</script>

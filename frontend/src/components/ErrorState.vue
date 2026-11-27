<template>
  <div class="card p-10 text-center space-y-3" role="alert">
    <div class="text-2xl select-none">⚠️</div>
    <p class="text-sm text-gray-600 dark:text-gray-300">{{ message }}</p>
    <div class="flex items-center justify-center gap-3 pt-1">
      <button v-if="retryable" type="button" class="btn-primary" @click="$emit('retry')">重试</button>
      <router-link v-if="showHome" to="/" class="btn-secondary">返回首页</router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    message?: string
    retryable?: boolean
    showHome?: boolean
  }>(),
  {
    message: '加载失败，请稍后重试',
    retryable: true,
    showHome: false,
  }
)

defineEmits<{ retry: [] }>()
</script>

<template>
  <router-view />
  <div class="fixed bottom-4 right-4 z-50 space-y-2">
    <transition-group name="toast">
      <div
        v-for="t in toasts"
        :key="t.id"
        class="px-4 py-2.5 rounded-lg shadow-lg text-sm text-white"
        :class="{
          'bg-emerald-600': t.type === 'success',
          'bg-red-600': t.type === 'error',
          'bg-gray-800': t.type === 'info',
        }"
      >
        {{ t.text }}
      </div>
    </transition-group>
  </div>
</template>

<script setup lang="ts">
import { toasts } from '@/utils/toast'
</script>

<style>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.25s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(12px);
}
</style>

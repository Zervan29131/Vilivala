<template>
  <transition name="fade-up">
    <button
      v-if="visible"
      type="button"
      class="fixed floating-safe-stack right-3.5 sm:right-5 z-30 w-10 h-10 rounded-full bg-white dark:bg-[#2a2a2a] border border-gray-200 dark:border-white/15 shadow-md flex items-center justify-center text-gray-600 dark:text-gray-200 hover:scale-105 hover:border-[#42b883] hover:text-[#42b883] transition"
      title="回到顶部"
      aria-label="回到顶部"
      @click="scrollToTop"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
      </svg>
    </button>
  </transition>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'

const visible = ref(false)

function onScroll() {
  visible.value = window.scrollY > 360
}

function scrollToTop() {
  const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  window.scrollTo({ top: 0, behavior: reduced ? 'auto' : 'smooth' })
}

onMounted(() => {
  onScroll()
  window.addEventListener('scroll', onScroll, { passive: true })
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', onScroll)
})
</script>

<style scoped>
.fade-up-enter-active,
.fade-up-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.fade-up-enter-from,
.fade-up-leave-to {
  opacity: 0;
  transform: translateY(8px);
}
</style>

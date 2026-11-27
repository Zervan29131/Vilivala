<template>
  <div
    class="fixed top-0 left-0 right-0 z-50 h-[2.5px] pointer-events-none bg-transparent"
    role="progressbar"
    aria-label="阅读进度"
    :aria-valuenow="Math.round(progress)"
    aria-valuemin="0"
    aria-valuemax="100"
  >
    <div
      class="h-full bg-[#42b883] dark:bg-[#4fc08d] transition-[width] duration-75 ease-out"
      :style="{ width: `${progress}%` }"
    />
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'

const progress = ref(0)
let raf = 0

function update() {
  raf = 0
  const doc = document.documentElement
  const scrollable = doc.scrollHeight - window.innerHeight
  if (scrollable <= 0) {
    progress.value = 0
    return
  }
  progress.value = Math.min(100, Math.max(0, (window.scrollY / scrollable) * 100))
}

function onScroll() {
  if (!raf) raf = window.requestAnimationFrame(update)
}

onMounted(() => {
  update()
  window.addEventListener('scroll', onScroll, { passive: true })
  window.addEventListener('resize', onScroll, { passive: true })
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', onScroll)
  window.removeEventListener('resize', onScroll)
  if (raf) window.cancelAnimationFrame(raf)
})
</script>

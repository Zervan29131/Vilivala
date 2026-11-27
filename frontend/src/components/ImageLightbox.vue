<template>
  <teleport to="body">
    <transition name="lb-fade">
      <div
        v-if="open"
        class="fixed inset-0 z-[70] bg-black/85 backdrop-blur-sm flex flex-col items-center justify-center p-4 sm:p-8 select-none"
        role="dialog"
        aria-modal="true"
        aria-label="图片预览"
        @click.self="close"
      >
        <!-- Top controls -->
        <div class="fixed top-4 right-4 flex items-center gap-2 z-10">
          <span v-if="images.length > 1" class="px-2.5 py-1 rounded-full bg-white/10 text-white/80 text-xs font-mono">
            {{ index + 1 }} / {{ images.length }}
          </span>
          <button
            type="button"
            class="w-9 h-9 rounded-full bg-white/10 hover:bg-white/20 text-white flex items-center justify-center transition"
            :title="zoomed ? '还原大小' : '放大查看'"
            @click="zoomed = !zoomed"
          >
            {{ zoomed ? '−' : '+' }}
          </button>
          <button
            type="button"
            class="w-9 h-9 rounded-full bg-white/10 hover:bg-white/20 text-white flex items-center justify-center transition"
            title="关闭 (Esc)"
            aria-label="关闭预览"
            @click="close"
          >
            ✕
          </button>
        </div>

        <!-- Prev / Next -->
        <button
          v-if="images.length > 1"
          type="button"
          class="fixed left-3 sm:left-6 top-1/2 -translate-y-1/2 w-10 h-10 rounded-full bg-white/10 hover:bg-white/25 text-white flex items-center justify-center transition"
          aria-label="上一张"
          @click.stop="prev"
        >
          ‹
        </button>
        <button
          v-if="images.length > 1"
          type="button"
          class="fixed right-3 sm:right-6 top-1/2 -translate-y-1/2 w-10 h-10 rounded-full bg-white/10 hover:bg-white/25 text-white flex items-center justify-center transition"
          aria-label="下一张"
          @click.stop="next"
        >
          ›
        </button>

        <!-- Image -->
        <div class="max-w-full max-h-[85vh] overflow-auto flex items-center justify-center">
          <img
            v-if="current"
            :src="current.src"
            :alt="current.alt"
            class="rounded-lg shadow-2xl transition-transform duration-200 cursor-zoom-in"
            :class="zoomed ? 'max-w-none scale-125 cursor-zoom-out' : 'max-w-[92vw] max-h-[82vh] object-contain'"
            @click.stop="zoomed = !zoomed"
          />
        </div>

        <p v-if="current?.alt" class="mt-3 text-xs sm:text-sm text-gray-300 text-center max-w-2xl truncate">
          {{ current.alt }}
        </p>
      </div>
    </transition>
  </teleport>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'

export interface LightboxImage {
  src: string
  alt: string
}

const props = defineProps<{
  open: boolean
  images: LightboxImage[]
  startIndex?: number
}>()

const emit = defineEmits<{ close: [] }>()

const index = ref(0)
const zoomed = ref(false)

const current = computed(() => props.images[index.value] || null)

function close() {
  zoomed.value = false
  emit('close')
}

function prev() {
  if (!props.images.length) return
  zoomed.value = false
  index.value = (index.value - 1 + props.images.length) % props.images.length
}

function next() {
  if (!props.images.length) return
  zoomed.value = false
  index.value = (index.value + 1) % props.images.length
}

function onKeydown(e: KeyboardEvent) {
  if (!props.open) return
  if (e.key === 'Escape') close()
  else if (e.key === 'ArrowLeft') prev()
  else if (e.key === 'ArrowRight') next()
}

watch(
  () => [props.open, props.startIndex] as const,
  ([isOpen, start]) => {
    if (isOpen) {
      index.value = Math.max(0, Math.min(start ?? 0, Math.max(0, props.images.length - 1)))
      zoomed.value = false
    }
  }
)

onMounted(() => window.addEventListener('keydown', onKeydown))
onBeforeUnmount(() => window.removeEventListener('keydown', onKeydown))
</script>

<style scoped>
.lb-fade-enter-active,
.lb-fade-leave-active {
  transition: opacity 0.18s ease;
}
.lb-fade-enter-from,
.lb-fade-leave-to {
  opacity: 0;
}
</style>

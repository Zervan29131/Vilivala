<template>
  <div class="my-1.5 rounded-xl border border-white/10 bg-white/[0.04] p-3.5 max-w-md space-y-2.5" @click.stop>
    <div class="flex items-center justify-between gap-2">
      <div class="flex items-center gap-2 text-xs font-semibold text-gray-200">
        <span>⏱️</span>
        <span>{{ label || '专注倒计时' }}</span>
      </div>
      <span
        class="px-2 py-0.5 rounded text-[11px]"
        :class="
          finished
            ? 'bg-[#42b883]/25 text-[#4fc08d]'
            : running
              ? 'bg-sky-500/20 text-sky-300'
              : 'bg-amber-500/20 text-amber-300'
        "
      >
        {{ finished ? '🎉 已完成' : running ? '进行中' : '已暂停' }}
      </span>
    </div>

    <div class="flex items-baseline justify-between">
      <div class="text-2xl font-bold tracking-wider text-white">
        {{ formattedRemaining }}
      </div>
      <div class="text-[11px] text-gray-400">总时长 {{ formattedTotal }}</div>
    </div>

    <!-- Progress Bar -->
    <div class="h-1.5 rounded-full bg-white/10 overflow-hidden">
      <div
        class="h-full rounded-full bg-[#42b883] transition-all duration-300"
        :style="{ width: `${progressPct}%` }"
      />
    </div>

    <div class="flex items-center gap-2 pt-0.5 text-xs font-sans">
      <button
        v-if="!finished"
        type="button"
        class="px-3 py-1 rounded bg-[#42b883] hover:bg-[#33a06f] text-white transition-colors"
        @click="running = !running"
      >
        {{ running ? '暂停' : '继续' }}
      </button>
      <button
        type="button"
        class="px-3 py-1 rounded bg-white/10 hover:bg-white/20 text-gray-200 transition-colors"
        @click="resetTimer"
      >
        重置
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { toast } from '@/utils/toast'

const props = defineProps<{
  seconds: number
  label?: string
}>()

const total = Math.max(1, props.seconds)
const remaining = ref(total)
const running = ref(true)
let timer: number | undefined

const finished = computed(() => remaining.value <= 0)
const progressPct = computed(() => Math.min(100, Math.round(((total - remaining.value) / total) * 100)))

function fmt(sec: number): string {
  const s = Math.max(0, sec)
  const h = Math.floor(s / 3600)
  const m = Math.floor((s % 3600) / 60)
  const r = s % 60
  if (h > 0) {
    return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}:${String(r).padStart(2, '0')}`
  }
  return `${String(m).padStart(2, '0')}:${String(r).padStart(2, '0')}`
}

const formattedRemaining = computed(() => fmt(remaining.value))
const formattedTotal = computed(() => fmt(total))

function resetTimer() {
  remaining.value = total
  running.value = true
}

onMounted(() => {
  timer = window.setInterval(() => {
    if (!running.value || remaining.value <= 0) return
    remaining.value--
    if (remaining.value <= 0) {
      running.value = false
      toast.success(`⏱️ 倒计时完成：${props.label || formattedTotal.value}`)
    }
  }, 1000)
})

onBeforeUnmount(() => {
  if (timer) window.clearInterval(timer)
})
</script>

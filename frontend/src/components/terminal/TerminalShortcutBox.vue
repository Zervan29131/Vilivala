<template>
  <div class="my-1.5 rounded-xl border border-white/10 bg-white/[0.04] p-3.5 max-w-xl" @click.stop>
    <div class="text-xs font-semibold text-[#4fc08d] mb-2.5">⌨️ 终端快捷键速查</div>
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 text-xs">
      <div
        v-for="item in shortcuts"
        :key="item.keys"
        class="flex items-center justify-between gap-2 px-2.5 py-1.5 rounded bg-white/[0.03]"
      >
        <span class="text-gray-300">{{ item.desc }}</span>
        <kbd class="px-1.5 py-0.5 rounded bg-white/10 text-[#4fc08d] text-[11px] shrink-0">{{ item.keys }}</kbd>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const shortcuts = [
  { keys: 'Tab', desc: '补全命令或目录名' },
  { keys: '↑ / ↓', desc: '翻阅历史命令' },
  { keys: 'Ctrl + L', desc: '清空当前屏幕' },
  { keys: 'Ctrl + O', desc: '折叠 / 展开全部输出' },
  { keys: '!! / !3', desc: '重跑上一条或第 N 条命令' },
  { keys: 'a – z', desc: '任意字母自动聚焦输入框' },
  { keys: '⌘K / Ctrl + K', desc: '唤起全站文章搜索框' },
]
</script>

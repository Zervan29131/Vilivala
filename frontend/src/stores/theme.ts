import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export type ThemeMode = 'system' | 'light' | 'dark'

export const useThemeStore = defineStore('theme', () => {
  const stored = (localStorage.getItem('theme') as ThemeMode) || 'system'
  const mode = ref<ThemeMode>(['system', 'light', 'dark'].includes(stored) ? stored : 'system')
  const systemDark = ref(window.matchMedia('(prefers-color-scheme: dark)').matches)

  const isDark = computed(() =>
    mode.value === 'dark' || (mode.value === 'system' && systemDark.value)
  )

  function apply() {
    const dark = isDark.value
    document.documentElement.classList.toggle('dark', dark)
    document.documentElement.style.colorScheme = dark ? 'dark' : 'light'
    let meta = document.head.querySelector('meta[name="theme-color"]') as HTMLMetaElement | null
    if (!meta) {
      meta = document.createElement('meta')
      meta.name = 'theme-color'
      document.head.appendChild(meta)
    }
    meta.content = dark ? '#1a1a1a' : '#f6f6f7'
  }

  function setMode(next: ThemeMode) {
    mode.value = next
    localStorage.setItem('theme', next)
    apply()
  }

  function cycle() {
    const order: ThemeMode[] = ['system', 'light', 'dark']
    const next = order[(order.indexOf(mode.value) + 1) % order.length]
    setMode(next)
  }

  // Listen for OS theme changes
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
    systemDark.value = e.matches
    if (mode.value === 'system') apply()
  })

  apply()

  return { mode, isDark, setMode, cycle }
})

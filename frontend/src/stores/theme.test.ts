import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useThemeStore } from './theme'

function mockStorage() {
  const store = new Map<string, string>()
  const storage = {
    getItem: (k: string) => (store.has(k) ? store.get(k)! : null),
    setItem: (k: string, v: string) => {
      store.set(k, String(v))
    },
    removeItem: (k: string) => {
      store.delete(k)
    },
    clear: () => {
      store.clear()
    },
    get length() {
      return store.size
    },
    key: (i: number) => Array.from(store.keys())[i] ?? null,
  }
  Object.defineProperty(globalThis, 'localStorage', { value: storage, writable: true, configurable: true })
  Object.defineProperty(window, 'localStorage', { value: storage, writable: true, configurable: true })
}

function mockMatchMedia(matches: boolean) {
  Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: vi.fn().mockImplementation((query: string) => ({
      matches,
      media: query,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
    })),
  })
}

describe('theme store', () => {
  beforeEach(() => {
    mockStorage()
    document.documentElement.className = ''
    document.head.innerHTML = ''
    mockMatchMedia(false)
    setActivePinia(createPinia())
  })

  it('cycles system -> light -> dark and syncs html class, colorScheme, and theme-color', () => {
    const store = useThemeStore()
    expect(store.mode).toBe('system')
    expect(store.isDark).toBe(false)

    store.cycle()
    expect(store.mode).toBe('light')
    expect(localStorage.getItem('theme')).toBe('light')
    expect(document.documentElement.classList.contains('dark')).toBe(false)

    store.cycle()
    expect(store.mode).toBe('dark')
    expect(store.isDark).toBe(true)
    expect(document.documentElement.classList.contains('dark')).toBe(true)
    expect(document.documentElement.style.colorScheme).toBe('dark')
    const meta = document.head.querySelector('meta[name="theme-color"]') as HTMLMetaElement | null
    expect(meta?.content).toBe('#1a1a1a')
  })
})

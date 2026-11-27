import { describe, expect, it } from 'vitest'
import { formatDate, makeSlug } from './format'

describe('format utilities', () => {
  it('formats ISO date strings to YYYY-MM-DD', () => {
    expect(formatDate('2026-11-27T12:30:00+08:00')).toBe('2026-11-27')
    expect(formatDate('invalid')).toBe('')
  })

  it('creates clean slugs and falls back for non-ASCII strings', () => {
    expect(makeSlug('Hello, Vue 3!')).toBe('hello-vue-3')
    expect(makeSlug('中文标题', 'post')).toMatch(/^post-[a-z0-9]+$/)
  })
})

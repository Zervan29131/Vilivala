import { describe, expect, it, vi } from 'vitest'
import { copyText } from './clipboard'

describe('cross-browser clipboard helper', () => {
  it('uses navigator.clipboard.writeText when available', async () => {
    const writeText = vi.fn().mockResolvedValue(undefined)
    Object.defineProperty(navigator, 'clipboard', {
      value: { writeText },
      configurable: true,
    })
    expect(await copyText('hello')).toBe(true)
    expect(writeText).toHaveBeenCalledWith('hello')
  })

  it('falls back to textarea + execCommand when navigator.clipboard is unavailable (HTTP)', async () => {
    Object.defineProperty(navigator, 'clipboard', {
      value: undefined,
      configurable: true,
    })
    document.execCommand = vi.fn().mockReturnValue(true)
    expect(await copyText('fallback-text')).toBe(true)
    expect(document.execCommand).toHaveBeenCalledWith('copy')
  })
})

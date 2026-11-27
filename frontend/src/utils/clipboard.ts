/**
 * Cross-browser clipboard copy helper.
 * Uses Async Clipboard API in HTTPS/localhost, and falls back to an offscreen
 * textarea + execCommand('copy') for HTTP intranet deployments & older mobile WebViews.
 */
export async function copyText(text: string): Promise<boolean> {
  if (typeof navigator !== 'undefined' && navigator.clipboard?.writeText) {
    try {
      await navigator.clipboard.writeText(text)
      return true
    } catch {
      // Fall through to execCommand fallback
    }
  }

  if (typeof document === 'undefined') return false
  const ta = document.createElement('textarea')
  ta.value = text
  ta.setAttribute('readonly', '')
  ta.style.position = 'fixed'
  ta.style.top = '-9999px'
  ta.style.left = '-9999px'
  ta.style.opacity = '0'
  document.body.appendChild(ta)

  try {
    ta.focus()
    ta.select()
    ta.setSelectionRange(0, text.length) // Required by iOS Safari
    return document.execCommand('copy')
  } catch {
    return false
  } finally {
    document.body.removeChild(ta)
  }
}

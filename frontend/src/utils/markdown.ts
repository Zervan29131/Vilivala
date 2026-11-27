import { Marked } from 'marked'
import hljs from 'highlight.js/lib/common'
import DOMPurify from 'dompurify'

export interface TocItem {
  id: string
  text: string
  level: number
}

export interface RenderOptions {
  /** If provided, strips a leading `# <title>` heading when it duplicates the page title. */
  stripLeadingTitle?: string
}

function escapeHtml(s: string): string {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')
}

/** Build a readable, Unicode-aware heading slug (Chinese + Latin). */
export function slugifyHeading(text: string): string {
  const cleaned = text
    .replace(/<[^>]+>/g, '')
    .trim()
    .toLowerCase()
    .replace(/[\s/\\]+/g, '-')
    .replace(/[^\p{L}\p{N}_-]+/gu, '')
    .replace(/-+/g, '-')
    .replace(/^-+|-+$/g, '')
  return cleaned || 'section'
}

const BARE_EMAIL_RE = /^[^\s@:]+@[^\s@]+\.[^\s@]+$/

function stripDuplicateH1(src: string, title?: string): string {
  if (!title) return src
  const trimmed = src.replace(/^\uFEFF/, '')
  const match = trimmed.match(/^\s*#\s+(.+?)\s*(?:\r?\n|$)/)
  if (!match) return src
  const norm = (s: string) => s.replace(/[`*_~]/g, '').trim().toLowerCase()
  if (norm(match[1]) === norm(title)) {
    return trimmed.slice(match[0].length).replace(/^\r?\n/, '')
  }
  return src
}

/** Render Markdown to sanitized HTML with highlighted code blocks, lazy images, and a Unicode TOC. */
export function renderMarkdown(src: string, options: RenderOptions = {}): { html: string; toc: TocItem[] } {
  const source = stripDuplicateH1(src || '', options.stripLeadingTitle)
  const toc: TocItem[] = []
  const usedIds = new Map<string, number>()

  const nextId = (rawText: string) => {
    const base = slugifyHeading(rawText)
    const seen = usedIds.get(base) || 0
    usedIds.set(base, seen + 1)
    return seen === 0 ? base : `${base}-${seen + 1}`
  }

  const marked = new Marked()
  marked.use({
    renderer: {
      heading(text: string, level: number) {
        const plain = text.replace(/<[^>]+>/g, '').trim()
        const id = nextId(plain)
        if (level <= 3) toc.push({ id, text: plain, level })
        return `<h${level} id="${escapeHtml(id)}">${text}</h${level}>`
      },
      link(href: string, title: string | null | undefined, text: string) {
        let safeHref = (href || '').trim()
        if (BARE_EMAIL_RE.test(safeHref)) {
          safeHref = `mailto:${safeHref}`
        }
        const isExternal = /^https?:\/\//i.test(safeHref)
        const titleAttr = title ? ` title="${escapeHtml(title)}"` : ''
        const targetAttr = isExternal ? ' target="_blank" rel="noopener noreferrer"' : ''
        return `<a href="${escapeHtml(safeHref)}"${titleAttr}${targetAttr}>${text}</a>`
      },
      image(href: string, title: string | null | undefined, text: string) {
        const srcAttr = escapeHtml((href || '').trim())
        const altAttr = escapeHtml(text || '')
        const titleAttr = title ? ` title="${escapeHtml(title)}"` : ''
        return `<img src="${srcAttr}" alt="${altAttr}"${titleAttr} loading="lazy" decoding="async" data-lightbox="true" />`
      },
      code(code: string, infostring: string | undefined) {
        const lang = (infostring || '').split(/\s+/)[0]
        const valid = lang && hljs.getLanguage(lang)
        const body = valid ? hljs.highlight(code, { language: lang }).value : escapeHtml(code)
        return (
          `<div class="code-block"><div class="code-head"><span>${escapeHtml(lang || 'text')}</span>` +
          `<button type="button" class="copy-btn">复制</button></div>` +
          `<pre><code class="hljs">${body}</code></pre></div>`
        )
      },
    },
  })

  const rawHtml = marked.parse(source, { async: false, gfm: true }) as string
  const html = DOMPurify.sanitize(rawHtml, {
    ADD_TAGS: ['button'],
    ADD_ATTR: ['target', 'rel', 'loading', 'decoding', 'data-lightbox', 'type'],
  })
  return { html, toc }
}

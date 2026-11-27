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

interface FriendCardItem {
  name: string
  desc: string
  avatar: string
  link: string
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

function normalizeSelfAsset(url: string): string {
  return url.trim().replace(/^https?:\/\/zervan\.cn\/img\//i, '/img/')
}

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

const CALLOUT_LABELS: Record<string, { kind: string; title: string }> = {
  note: { kind: 'info', title: '说明' },
  info: { kind: 'info', title: '信息' },
  tip: { kind: 'tip', title: '提示' },
  important: { kind: 'tip', title: '重要' },
  warning: { kind: 'warning', title: '注意' },
  caution: { kind: 'danger', title: '警告' },
  danger: { kind: 'danger', title: '危险' },
}

function unquoteYamlVal(v: string): string {
  const s = v.trim()
  if ((s.startsWith("'") && s.endsWith("'")) || (s.startsWith('"') && s.endsWith('"'))) {
    return s.slice(1, -1).trim()
  }
  return s
}

function parseCardListYaml(raw: string): FriendCardItem[] {
  const cleaned = raw.replace(/^\s*```(?:yaml|yml)?\s*\r?\n?/i, '').replace(/\r?\n?\s*```\s*$/i, '')
  const items: FriendCardItem[] = []
  let current: FriendCardItem | null = null

  for (const line of cleaned.split(/\r?\n/)) {
    const trimmed = line.trim()
    if (!trimmed || trimmed.startsWith('#')) continue
    const itemStart = trimmed.match(/^-\s*([a-zA-Z_]+)\s*:\s*(.*)$/)
    if (itemStart) {
      if (current && current.name && current.link) items.push(current)
      current = { name: '', desc: '', avatar: '', link: '' }
      const [, k, v] = itemStart
      if (k in current) current[k as keyof FriendCardItem] = unquoteYamlVal(v)
      continue
    }
    const kv = trimmed.match(/^([a-zA-Z_]+)\s*:\s*(.*)$/)
    if (kv && current) {
      const [, k, v] = kv
      if (k === 'name' || k === 'desc' || k === 'avatar' || k === 'link') {
        current[k] = unquoteYamlVal(v)
      }
    }
  }
  if (current && current.name && current.link) items.push(current)
  return items.filter((it) => /^(https?:\/\/|\/)/i.test(it.link))
}

function renderCardListHtml(items: FriendCardItem[]): string {
  if (!items.length) return ''
  const cards = items
    .map((it) => {
      const initial = escapeHtml((it.name.trim()[0] || '?').toUpperCase())
      const avatar = it.avatar ? normalizeSelfAsset(it.avatar) : ''
      const avatarHtml = avatar
        ? `<img src="${escapeHtml(avatar)}" alt="${escapeHtml(it.name)}" loading="lazy" decoding="async" class="friend-avatar" data-friend-avatar="${initial}" />`
        : `<span class="friend-initial">${initial}</span>`
      return (
        `<a href="${escapeHtml(it.link.trim())}" target="_blank" rel="noopener noreferrer" class="friend-card">` +
        `<span class="friend-avatar-wrap">${avatarHtml}</span>` +
        `<span class="friend-meta">` +
        `<strong class="friend-name">${escapeHtml(it.name)}</strong>` +
        `<span class="friend-desc">${escapeHtml(it.desc || it.link)}</span>` +
        `</span></a>`
      )
    })
    .join('')
  return `<div class="friend-grid">${cards}</div>`
}

/** Pre-process VuePress/VitePress `::: cardList` and `::: tip|info|warning|danger` containers. */
function preprocessContainers(src: string): string {
  // 1. ::: cardList [cols] ... :::
  let out = src.replace(
    /^:::\s*cardList(?:\s+\d+)?\s*\r?\n([\s\S]*?)^:::\s*$/gm,
    (_m, body: string) => renderCardListHtml(parseCardListYaml(body))
  )

  // 2. ::: tip|info|warning|danger|note [custom title] ... :::
  out = out.replace(
    /^:::\s*(tip|info|warning|danger|note|important|caution)(?:[ \t]+([^\r\n]+))?\r?\n([\s\S]*?)^:::\s*$/gim,
    (_m, rawType: string, customTitle: string | undefined, body: string) => {
      const meta = CALLOUT_LABELS[rawType.toLowerCase()] || CALLOUT_LABELS.info
      const heading = escapeHtml((customTitle || '').trim() || meta.title)
      return `<div class="callout callout-${meta.kind}"><div class="callout-title">${heading}</div>\n\n${body.trim()}\n\n</div>`
    }
  )
  return out
}

/** Attach delegated broken-image fallback inside a rendered `.md-body` container. */
export function attachMarkdownImageFallbacks(container: HTMLElement | null) {
  if (!container) return
  container.querySelectorAll('img').forEach((img) => {
    if (img.dataset.fallbackBound === '1') return
    img.dataset.fallbackBound = '1'
    const onFail = () => {
      // Friend card avatar fallback -> clean initial badge
      if (img.classList.contains('friend-avatar')) {
        const span = document.createElement('span')
        span.className = 'friend-initial'
        span.textContent = img.dataset.friendAvatar || '?'
        img.replaceWith(span)
        return
      }
      const src = img.getAttribute('src') || ''
      const box = document.createElement('div')
      box.className = 'md-img-broken'
      const text = document.createElement('span')
      text.textContent = `🖼️ 外链图片暂不可用${img.alt ? `（${img.alt}）` : ''}`
      box.appendChild(text)
      if (/^https?:\/\//i.test(src)) {
        const link = document.createElement('a')
        link.href = src
        link.target = '_blank'
        link.rel = 'noopener noreferrer'
        link.className = 'md-img-broken-link'
        link.textContent = '尝试打开原图 ↗'
        box.appendChild(link)
      }
      img.replaceWith(box)
    }
    if (img.complete && img.naturalWidth === 0) {
      onFail()
    } else {
      img.addEventListener('error', onFail, { once: true })
    }
  })
}

/** Render Markdown to sanitized HTML with Callouts, Friend Cards, code blocks, lazy images, and Unicode TOC. */
export function renderMarkdown(src: string, options: RenderOptions = {}): { html: string; toc: TocItem[] } {
  const stripped = stripDuplicateH1(src || '', options.stripLeadingTitle)
  const source = preprocessContainers(stripped)
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
        const escId = escapeHtml(id)
        const anchor = `<a class="header-anchor" href="#${escId}" data-anchor="${escId}" aria-label="复制章节链接">#</a>`
        return `<h${level} id="${escId}" class="md-heading">${text}${anchor}</h${level}>`
      },
      blockquote(quote: string) {
        // Support GitHub-style alerts: > [!TIP] / > [!WARNING] / > [!NOTE] / > [!IMPORTANT] / > [!CAUTION]
        const alertMatch = quote.match(
          /^\s*<p>\s*\[!(NOTE|TIP|IMPORTANT|WARNING|CAUTION)\](?:[ \t]*([^\n<]*))?(?:\n|<br\s*\/?>)?([\s\S]*)/i
        )
        if (alertMatch) {
          const [, rawKind, inlineTitle, rest] = alertMatch
          const meta = CALLOUT_LABELS[rawKind.toLowerCase()] || CALLOUT_LABELS.info
          const title = escapeHtml((inlineTitle || '').trim() || meta.title)
          const body = rest.trim() ? `<p>${rest}` : ''
          return `<div class="callout callout-${meta.kind}"><div class="callout-title">${title}</div>${body}</div>`
        }
        return `<blockquote>${quote}</blockquote>`
      },
      link(href: string, title: string | null | undefined, text: string) {
        let safeHref = (href || '').trim()
        if (BARE_EMAIL_RE.test(safeHref)) {
          safeHref = `mailto:${safeHref}`
        }
        const isExternal = /^https?:\/\//i.test(safeHref)
        const titleAttr = title ? ` title="${escapeHtml(title)}"` : ''
        const targetAttr = isExternal ? ' target="_blank" rel="noopener noreferrer" class="ext-link"' : ''
        return `<a href="${escapeHtml(safeHref)}"${titleAttr}${targetAttr}>${text}</a>`
      },
      image(href: string, title: string | null | undefined, text: string) {
        const normalized = normalizeSelfAsset(href || '')
        const srcAttr = escapeHtml(normalized)
        const cleanAlt = (text || '').trim()
        const altAttr = escapeHtml(cleanAlt)
        const titleAttr = title ? ` title="${escapeHtml(title)}"` : ''
        const imgTag = `<img src="${srcAttr}" alt="${altAttr}"${titleAttr} loading="lazy" decoding="async" data-lightbox="true" />`
        if (cleanAlt && !/^(img|image|photo|pic|screenshot|\d+)$/i.test(cleanAlt)) {
          return `<figure class="md-figure">${imgTag}<figcaption class="md-figcaption">${altAttr}</figcaption></figure>`
        }
        return imgTag
      },
      code(code: string, infostring: string | undefined) {
        const lang = (infostring || '').split(/\s+/)[0]
        const valid = lang && hljs.getLanguage(lang)
        const body = valid ? hljs.highlight(code, { language: lang }).value : escapeHtml(code)
        return (
          `<div class="code-block"><div class="code-head"><span>${escapeHtml(lang || 'text')}</span>` +
          `<button type="button" class="copy-btn">复制</button></div>` +
          `<pre><code class="hljs language-${escapeHtml(lang || 'text')}">${body}</code></pre></div>`
        )
      },
    },
  })

  const rawHtml = marked.parse(source, { async: false, gfm: true }) as string
  const html = DOMPurify.sanitize(rawHtml, {
    ADD_TAGS: ['button', 'figure', 'figcaption'],
    ADD_ATTR: [
      'target',
      'rel',
      'loading',
      'decoding',
      'data-lightbox',
      'data-anchor',
      'data-friend-avatar',
      'type',
    ],
  })
  return { html, toc }
}

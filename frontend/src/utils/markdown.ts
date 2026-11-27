import { Marked } from 'marked'
import hljs from 'highlight.js/lib/common'

export interface TocItem {
  id: string
  text: string
  level: number
}

function escapeHtml(s: string) {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
}

/** Render Markdown to HTML with highlighted code blocks and a table of contents. */
export function renderMarkdown(src: string): { html: string; toc: TocItem[] } {
  const toc: TocItem[] = []
  const marked = new Marked()
  marked.use({
    renderer: {
      heading(text: string, level: number) {
        const id = `h-${toc.length}`
        if (level <= 3) toc.push({ id, text: text.replace(/<[^>]+>/g, ''), level })
        return `<h${level} id="${id}">${text}</h${level}>`
      },
      code(code: string, infostring: string | undefined) {
        const lang = (infostring || '').split(/\s+/)[0]
        const valid = lang && hljs.getLanguage(lang)
        const body = valid
          ? hljs.highlight(code, { language: lang }).value
          : escapeHtml(code)
        return (
          `<div class="code-block"><div class="code-head"><span>${escapeHtml(lang || 'text')}</span>` +
          `<button type="button" class="copy-btn">复制</button></div>` +
          `<pre><code class="hljs">${body}</code></pre></div>`
        )
      },
    },
  })
  const html = marked.parse(src, { async: false, gfm: true }) as string
  return { html, toc }
}

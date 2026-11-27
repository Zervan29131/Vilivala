import { describe, expect, it } from 'vitest'
import { renderMarkdown, slugifyHeading } from './markdown'

describe('markdown utilities', () => {
  it('generates readable Unicode heading slugs and deduplicates repeats', () => {
    expect(slugifyHeading('快速开始 & 安装')).toBe('快速开始-安装')
    const { toc, html } = renderMarkdown('## 简介\n\n## 简介\n\n### Hello World!')
    expect(toc.map((t) => t.id)).toEqual(['简介', '简介-2', 'hello-world'])
    expect(html).toContain('id="简介-2"')
  })

  it('sanitizes dangerous HTML / XSS payloads', () => {
    const { html } = renderMarkdown('Hello <script>alert(1)</script><img src=x onerror="alert(2)">')
    expect(html).not.toContain('<script')
    expect(html).not.toContain('onerror')
  })

  it('normalizes bare email links to mailto without double-prefixing, and adds lazy loading on images', () => {
    const { html } = renderMarkdown(
      '联系 [me](user@example.com) 或 [already](mailto:user@example.com) 和 ![配图](https://example.com/a.png)'
    )
    expect(html).toContain('href="mailto:user@example.com"')
    expect(html).not.toContain('mailto:mailto:')
    expect(html).toContain('loading="lazy"')
    expect(html).toContain('data-lightbox="true"')
  })

  it('strips a duplicate leading H1 when it matches the post title', () => {
    const { html, toc } = renderMarkdown('# Git实际操作\n\n## 常用命令\n\n内容', {
      stripLeadingTitle: 'Git实际操作',
    })
    expect(html).not.toContain('<h1')
    expect(toc).toHaveLength(1)
    expect(toc[0].text).toBe('常用命令')
  })

  it('renders highlighted code block with copy button', () => {
    const { html } = renderMarkdown('```ts\nconst x = 1\n```')
    expect(html).toContain('class="copy-btn"')
    expect(html).toContain('class="hljs"')
  })
})

import { describe, expect, it } from 'vitest'
import { renderMarkdown, slugifyHeading } from './markdown'

describe('markdown utilities', () => {
  it('generates readable Unicode heading slugs, deduplicates repeats, and adds # anchor', () => {
    expect(slugifyHeading('快速开始 & 安装')).toBe('快速开始-安装')
    const { toc, html } = renderMarkdown('## 简介\n\n## 简介\n\n### Hello World!')
    expect(toc.map((t) => t.id)).toEqual(['简介', '简介-2', 'hello-world'])
    expect(html).toContain('id="简介-2"')
    expect(html).toContain('class="header-anchor"')
  })

  it('sanitizes dangerous HTML / XSS payloads', () => {
    const { html } = renderMarkdown('Hello <script>alert(1)</script><img src=x onerror="alert(2)">')
    expect(html).not.toContain('<script')
    expect(html).not.toContain('onerror')
  })

  it('normalizes bare email links, marks external links, and wraps descriptive images in figure', () => {
    const { html } = renderMarkdown(
      '联系 [me](user@example.com) 或 [already](mailto:user@example.com) 访问 [官网](https://vuejs.org) 和 ![架构示意图](https://zervan.cn/img/console.jpg)'
    )
    expect(html).toContain('href="mailto:user@example.com"')
    expect(html).not.toContain('mailto:mailto:')
    expect(html).toContain('class="ext-link"')
    expect(html).toContain('src="/img/console.jpg"')
    expect(html).toContain('<figure class="md-figure">')
    expect(html).toContain('<figcaption class="md-figcaption">架构示意图</figcaption>')
  })

  it('renders GitHub [!TIP] and VuePress ::: warning callouts', () => {
    const { html } = renderMarkdown(
      '> [!TIP]\n> 推荐使用 pnpm\n\n::: warning 版本要求\n需要 Node 18+\n:::'
    )
    expect(html).toContain('callout callout-tip')
    expect(html).toContain('callout callout-warning')
    expect(html).toContain('版本要求')
  })

  it('parses VuePress ::: cardList YAML into friend link cards', () => {
    const { html } = renderMarkdown(
      '::: cardList\n```yaml\n- name: Zervan\n  desc: Be the change\n  avatar: https://zervan.cn/img/console.jpg\n  link: https://zervan.cn/\n```\n:::'
    )
    expect(html).toContain('class="friend-grid"')
    expect(html).toContain('class="friend-card"')
    expect(html).toContain('Zervan')
    expect(html).toContain('src="/img/console.jpg"')
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
    const { html } = renderMarkdown('```diff\n+added\n-removed\n```')
    expect(html).toContain('class="copy-btn"')
    expect(html).toContain('hljs-addition')
    expect(html).toContain('hljs-deletion')
  })
})

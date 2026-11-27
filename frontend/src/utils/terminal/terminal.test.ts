import { describe, expect, it } from 'vitest'
import { expandHistoryBang, parseCommandLine, parseDurationSeconds } from './parser'
import { VirtualSpace, formatPromptPath, resolveVfsPath } from './vfs'
import { buildNamingVariants } from './varbook'
import type { Category } from '@/api/types'

describe('terminal parser', () => {
  it('parses quoted arguments and flags (-s, -r, -f)', () => {
    const p = parseCommandLine('search "Vue 3 教程" -f zhihu -s')
    expect(p.cmd).toBe('search')
    expect(p.args).toEqual(['Vue 3 教程'])
    expect(p.flags.from).toBe('zhihu')
    expect(p.flags.self).toBe(true)
  })

  it('expands !! and !n from command history', () => {
    const hist = ['ls', 'cd code', 'hot']
    expect(expandHistoryBang('!!', hist).expanded).toBe('hot')
    expect(expandHistoryBang('!2', hist).expanded).toBe('cd code')
    expect(expandHistoryBang('!9', hist).error).toBeTruthy()
  })

  it('parses duration strings into seconds', () => {
    expect(parseDurationSeconds('25m')).toBe(1500)
    expect(parseDurationSeconds('90s')).toBe(90)
    expect(parseDurationSeconds('45')).toBe(45)
    expect(parseDurationSeconds('abc')).toBeNull()
  })
})

describe('terminal virtual space (VFS)', () => {
  it('mounts read-only blog categories and supports user bookmark folder operations', () => {
    const space = new VirtualSpace([])
    const mockTree: Category[] = [
      {
        id: 1,
        name: 'code',
        slug: 'code',
        parent_id: 0,
        depth: 1,
        sort: 0,
        post_count: 5,
        total_count: 36,
        path: 'code',
        created_at: '',
        children: [
          {
            id: 2,
            name: '博客开发指南',
            slug: 'code-blog',
            parent_id: 1,
            depth: 2,
            sort: 0,
            post_count: 10,
            total_count: 10,
            path: 'code / 博客开发指南',
            created_at: '',
          },
        ],
      },
    ]
    space.setMountedCategories(mockTree)

    // Path helpers
    expect(resolveVfsPath('/code', '..')).toBe('/')
    expect(formatPromptPath('/code/博客开发指南')).toBe('~/code/博客开发指南')

    // Mounted category is browsable & protected from deletion
    expect(space.findInsensitive('/CODE')?.path).toBe('/code')
    expect(space.remove('/code', true).ok).toBe(false)

    // User folder & bookmark CRUD
    expect(space.mkdir('/工具').ok).toBe(true)
    expect(space.addLink('/工具/V2EX', 'https://v2ex.com').ok).toBe(true)
    expect(space.remove('/工具', false).ok).toBe(false) // non-empty without -r
    expect(space.copy('/工具/V2EX', '/V2EX备份').ok).toBe(true)
    expect(space.move('/V2EX备份', '/V2EX新名').ok).toBe(true)
    expect(space.getNode('/V2EX新名')?.url).toBe('https://v2ex.com')
    expect(space.remove('/工具', true).ok).toBe(true)
  })
})

describe('terminal varbook naming helper', () => {
  it('converts Chinese dev phrases and English words into 5 naming conventions', () => {
    const cn = buildNamingVariants('用户配置')
    expect(cn.variants.map((v) => v.value)).toEqual([
      'userConfig',
      'UserConfig',
      'user_config',
      'USER_CONFIG',
      'user-config',
    ])

    const en = buildNamingVariants('Post Detail View')
    expect(en.variants[0].value).toBe('postDetailView')
    expect(en.variants[2].value).toBe('post_detail_view')
  })
})

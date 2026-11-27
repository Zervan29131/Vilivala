import type { Category } from '@/api/types'

export interface VfsNode {
  name: string
  path: string
  kind: 'dir' | 'link' | 'category'
  url?: string
  icon?: string
  readonly?: boolean
  categoryId?: number
  categorySlug?: string
  totalCount?: number
}

export const VFS_STORAGE_KEY = 'vilivala_terminal_vfs_v1'

export function defaultUserNodes(): VfsNode[] {
  return [
    { name: '博客首页', path: '/博客首页', kind: 'link', url: '/', icon: '🏠' },
    { name: 'GitHub', path: '/GitHub', kind: 'link', url: 'https://github.com', icon: '🐙' },
    { name: '知乎', path: '/知乎', kind: 'link', url: 'https://www.zhihu.com', icon: '📘' },
    { name: 'Bilibili', path: '/Bilibili', kind: 'link', url: 'https://www.bilibili.com', icon: '📺' },
    { name: 'dev', path: '/dev', kind: 'dir', icon: '📁' },
    { name: 'Vue文档', path: '/dev/Vue文档', kind: 'link', url: 'https://cn.vuejs.org', icon: '💚' },
    { name: 'MDN', path: '/dev/MDN', kind: 'link', url: 'https://developer.mozilla.org/zh-CN/', icon: '📚' },
    { name: 'Go文档', path: '/dev/Go文档', kind: 'link', url: 'https://go.dev/doc/', icon: '🐹' },
  ]
}

/** Normalize any relative or absolute virtual path against cwd. */
export function resolveVfsPath(cwd: string, target = ''): string {
  const raw = target.trim()
  if (!raw || raw === '.') return normalizePath(cwd)
  if (raw === '~' || raw === '/') return '/'
  const joined = raw.startsWith('~/')
    ? '/' + raw.slice(2)
    : raw.startsWith('/')
      ? raw
      : (cwd === '/' ? '' : cwd) + '/' + raw
  return normalizePath(joined)
}

function normalizePath(p: string): string {
  const parts = p.split('/').filter(Boolean)
  const stack: string[] = []
  for (const seg of parts) {
    if (seg === '.') continue
    if (seg === '..') stack.pop()
    else stack.push(seg)
  }
  return '/' + stack.join('/')
}

export function parentPath(p: string): string {
  const norm = normalizePath(p)
  if (norm === '/') return '/'
  const idx = norm.lastIndexOf('/')
  return idx <= 0 ? '/' : norm.slice(0, idx)
}

export function baseName(p: string): string {
  const norm = normalizePath(p)
  if (norm === '/') return '/'
  return norm.slice(norm.lastIndexOf('/') + 1)
}

export function formatPromptPath(cwd: string): string {
  const norm = normalizePath(cwd)
  return norm === '/' ? '~' : `~${norm}`
}

/** Convert blog category tree into read-only mounted VfsNodes (`/code`, `/code/教程`, ...). */
export function mountCategoryTree(tree: Category[], parent = ''): VfsNode[] {
  const out: VfsNode[] = []
  for (const cat of tree) {
    const safeName = cat.name.replace(/\//g, '-')
    const path = `${parent}/${safeName}`
    out.push({
      name: safeName,
      path,
      kind: 'category',
      readonly: true,
      icon: '📂',
      categoryId: cat.id,
      categorySlug: cat.slug,
      totalCount: cat.total_count,
    })
    if (cat.children?.length) {
      out.push(...mountCategoryTree(cat.children, path))
    }
  }
  return out
}

export class VirtualSpace {
  private userNodes = new Map<string, VfsNode>()
  private mountedNodes = new Map<string, VfsNode>()

  constructor(initialUserNodes: VfsNode[] = defaultUserNodes()) {
    this.resetUserNodes(initialUserNodes)
  }

  resetUserNodes(nodes: VfsNode[] = defaultUserNodes()) {
    this.userNodes.clear()
    for (const n of nodes) {
      this.userNodes.set(normalizePath(n.path), { ...n, path: normalizePath(n.path) })
    }
  }

  setMountedCategories(tree: Category[]) {
    this.mountedNodes.clear()
    for (const n of mountCategoryTree(tree)) {
      this.mountedNodes.set(n.path, n)
    }
  }

  exportUserNodes(): VfsNode[] {
    return Array.from(this.userNodes.values())
  }

  /** Top-level + nested user links for the top browser quick-launch dock. */
  allLinks(): VfsNode[] {
    return Array.from(this.userNodes.values()).filter((n) => n.kind === 'link')
  }

  getNode(absPath: string): VfsNode | null {
    const norm = normalizePath(absPath)
    if (norm === '/') {
      return { name: '/', path: '/', kind: 'dir', readonly: true }
    }
    return this.mountedNodes.get(norm) || this.userNodes.get(norm) || null
  }

  /** Case-insensitive lookup inside a directory so `cd toread` matches `/ToRead`. */
  findInsensitive(absPath: string): VfsNode | null {
    const exact = this.getNode(absPath)
    if (exact) return exact
    const lower = normalizePath(absPath).toLowerCase()
    for (const n of this.mountedNodes.values()) {
      if (n.path.toLowerCase() === lower) return n
    }
    for (const n of this.userNodes.values()) {
      if (n.path.toLowerCase() === lower) return n
    }
    return null
  }

  isDirectory(node: VfsNode | null): boolean {
    return !!node && (node.kind === 'dir' || node.kind === 'category')
  }

  listChildren(dirPath: string, recursive = false): VfsNode[] {
    const norm = normalizePath(dirPath)
    const prefix = norm === '/' ? '/' : norm + '/'
    const baseDepth = norm === '/' ? 0 : norm.split('/').filter(Boolean).length

    const merged = new Map<string, VfsNode>()
    for (const n of this.mountedNodes.values()) merged.set(n.path, n)
    for (const n of this.userNodes.values()) {
      if (!merged.has(n.path)) merged.set(n.path, n)
    }

    const res: VfsNode[] = []
    for (const n of merged.values()) {
      if (n.path === norm || !n.path.startsWith(prefix)) continue
      const depth = n.path.split('/').filter(Boolean).length
      if (recursive || depth === baseDepth + 1) {
        res.push(n)
      }
    }
    return res.sort((a, b) => {
      const aDir = this.isDirectory(a) ? 0 : 1
      const bDir = this.isDirectory(b) ? 0 : 1
      if (aDir !== bDir) return aDir - bDir
      return a.path.localeCompare(b.path, 'zh-CN')
    })
  }

  mkdir(absPath: string): { ok: boolean; error?: string } {
    const norm = normalizePath(absPath)
    if (norm === '/') return { ok: false, error: '根目录已存在' }
    if (this.getNode(norm)) return { ok: false, error: `路径已存在：${norm}` }
    const pNode = this.getNode(parentPath(norm))
    if (!pNode || pNode.kind !== 'dir') {
      return { ok: false, error: '父目录不存在或为只读博客分类目录' }
    }
    this.userNodes.set(norm, { name: baseName(norm), path: norm, kind: 'dir', icon: '📁' })
    return { ok: true }
  }

  addLink(absPath: string, url: string, icon = '🔖'): { ok: boolean; error?: string } {
    const norm = normalizePath(absPath)
    if (norm === '/') return { ok: false, error: '请指定书签名称' }
    const existing = this.getNode(norm)
    if (existing?.readonly) return { ok: false, error: '该路径为只读博客分类，不能覆盖' }
    if (existing?.kind === 'dir') return { ok: false, error: '同名文件夹已存在' }
    const pNode = this.getNode(parentPath(norm))
    if (!pNode || pNode.kind !== 'dir') {
      return { ok: false, error: '目标父目录不存在或为只读分类目录' }
    }
    this.userNodes.set(norm, { name: baseName(norm), path: norm, kind: 'link', url, icon })
    return { ok: true }
  }

  remove(absPath: string, recursive = false): { ok: boolean; error?: string } {
    const norm = normalizePath(absPath)
    if (norm === '/') return { ok: false, error: '不能删除根目录' }
    const node = this.getNode(norm)
    if (!node) return { ok: false, error: `未找到：${norm}` }
    if (node.readonly) return { ok: false, error: `「${node.name}」为博客分类目录，不可删除` }
    const children = this.listChildren(norm, true)
    if (children.length > 0 && !recursive) {
      return { ok: false, error: `目录非空，请添加 -r 参数递归删除（rm -r ${node.name}）` }
    }
    this.userNodes.delete(norm)
    for (const c of children) {
      if (!c.readonly) this.userNodes.delete(c.path)
    }
    return { ok: true }
  }

  copy(srcAbs: string, dstAbs: string): { ok: boolean; error?: string } {
    const sNorm = normalizePath(srcAbs)
    let dNorm = normalizePath(dstAbs)
    const src = this.getNode(sNorm)
    if (!src) return { ok: false, error: `源路径不存在：${sNorm}` }
    if (src.kind !== 'link') return { ok: false, error: '仅支持复制书签条目' }

    const dstExisting = this.getNode(dNorm)
    if (dstExisting && dstExisting.kind === 'dir') {
      dNorm = normalizePath(`${dNorm}/${src.name}`)
    }
    return this.addLink(dNorm, src.url || '', src.icon)
  }

  move(srcAbs: string, dstAbs: string): { ok: boolean; error?: string } {
    const sNorm = normalizePath(srcAbs)
    const src = this.getNode(sNorm)
    if (!src) return { ok: false, error: `源路径不存在：${sNorm}` }
    if (src.readonly) return { ok: false, error: '博客分类目录不可移动' }
    const cp = this.copy(sNorm, dstAbs)
    if (!cp.ok) return cp
    this.userNodes.delete(sNorm)
    return { ok: true }
  }
}

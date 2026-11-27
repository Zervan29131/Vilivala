// Shared TypeScript types mirroring docs/API.md (v1.1)

export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface PaginatedResponse<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

export interface PathNode {
  id: number
  name: string
  slug: string
}

export interface CategoryRef {
  id: number
  name: string
  slug: string
  path: PathNode[]
}

export interface Category {
  id: number
  name: string
  slug: string
  parent_id: number
  depth: number
  sort: number
  post_count: number
  total_count: number
  path: string
  created_at: string
  children?: Category[]
}

export interface Tag {
  id: number
  name: string
  slug: string
  post_count?: number
  created_at: string
}

export interface PostLink {
  title: string
  slug: string
}

export interface RelatedPost {
  id: number
  title: string
  slug: string
  summary: string
  reading_time: number
  created_at: string
}

export interface SearchResultItem {
  id: number
  title: string
  slug: string
  snippet: string
  category_name: string
  created_at: string
}

export interface SearchResponse {
  list: SearchResultItem[]
  total: number
}

export interface Post {
  id: number
  title: string
  slug: string
  summary: string
  content?: string
  cover_image: string
  status: 'draft' | 'published'
  view_count: number
  is_top: boolean
  word_count: number
  reading_time: number
  category_id?: number
  category?: CategoryRef
  tags?: Tag[]
  prev?: PostLink | null
  next?: PostLink | null
  related?: RelatedPost[]
  created_at: string
  updated_at: string
}

export interface ArchivePostItem {
  id: number
  title: string
  slug: string
  created_at: string
}

export interface ArchiveMonth {
  month: number
  posts: ArchivePostItem[]
}

export interface ArchiveYear {
  year: number
  count: number
  months: ArchiveMonth[]
}

export interface PageSummary {
  slug: string
  title: string
  updated_at: string
}

export interface PageDetail extends PageSummary {
  content: string
}

export interface UploadResult {
  url: string
  size: number
  content_type: string
}

export interface User {
  id: number
  username: string
  nickname: string
  avatar: string
}

export interface LoginResponse {
  token: string
  expires_in: number
  user: User
}

// Request payloads
export interface LoginPayload {
  username: string
  password: string
}

export interface PostPayload {
  title: string
  slug: string
  summary?: string
  content: string
  cover_image?: string
  status: 'draft' | 'published'
  category_id?: number
  tag_ids?: number[]
  is_top?: boolean
  created_at?: string
}

export interface CategoryPayload {
  name: string
  slug: string
  parent_id?: number
  sort?: number
}

export interface TagPayload {
  name: string
  slug: string
}

export interface PagePayload {
  title: string
  content: string
}

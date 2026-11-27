import http from './request'
import type {
  PaginatedResponse,
  Post,
  Category,
  Tag,
  ArchiveYear,
  PageSummary,
  PageDetail,
  SearchResponse,
  SiteConfigData,
  DashboardData,
  UploadResult,
  LoginPayload,
  LoginResponse,
  PostPayload,
  CategoryPayload,
  TagPayload,
  PagePayload,
} from './types'

// ---- Public ----
export const getPublishedPosts = (params?: {
  page?: number
  page_size?: number
  category_id?: number
  tag_id?: number
  keyword?: string
}) => http.get<PaginatedResponse<Post>>('/posts', params)

export const getPostBySlug = (slug: string) => http.get<Post>(`/posts/${slug}`)
export const searchPosts = (q: string, limit = 12) =>
  http.get<SearchResponse>('/search', { q, limit })
export const getCategories = (tree = false) =>
  http.get<Category[]>('/categories', tree ? { tree: 1 } : undefined)
export const getTags = () => http.get<Tag[]>('/tags')
export const getArchives = () => http.get<ArchiveYear[]>('/archives')
export const getPages = () => http.get<PageSummary[]>('/pages')
export const getPageBySlug = (slug: string) => http.get<PageDetail>(`/pages/${slug}`)
export const getSiteConfig = () => http.get<SiteConfigData>('/site-config')

// ---- Auth ----
export const login = (payload: LoginPayload) => http.post<LoginResponse>('/auth/login', payload)

// ---- Admin: dashboard & site config ----
export const adminGetDashboard = () => http.get<DashboardData>('/admin/dashboard')
export const adminUpdateSiteConfig = (p: Omit<SiteConfigData, 'stats'>) =>
  http.put<SiteConfigData>('/admin/site-config', p)

// ---- Admin: posts ----
export const adminGetPosts = (params?: {
  page?: number
  page_size?: number
  status?: string
  keyword?: string
}) => http.get<PaginatedResponse<Post>>('/admin/posts', params)

export const adminGetPost = (id: number) => http.get<Post>(`/admin/posts/${id}`)
export const adminCreatePost = (p: PostPayload) => http.post<Post>('/admin/posts', p)
export const adminUpdatePost = (id: number, p: Partial<PostPayload>) =>
  http.put<null>(`/admin/posts/${id}`, p)
export const adminDeletePost = (id: number) => http.delete<null>(`/admin/posts/${id}`)
export const adminUpdatePostStatus = (id: number, status: 'draft' | 'published') =>
  http.patch<null>(`/admin/posts/${id}/status`, { status })

// ---- Admin: categories ----
export const adminCreateCategory = (p: CategoryPayload) => http.post<Category>('/admin/categories', p)
export const adminUpdateCategory = (id: number, p: CategoryPayload) =>
  http.put<null>(`/admin/categories/${id}`, p)
export const adminDeleteCategory = (id: number) => http.delete<null>(`/admin/categories/${id}`)

// ---- Admin: tags ----
export const adminCreateTag = (p: TagPayload) => http.post<Tag>('/admin/tags', p)
export const adminUpdateTag = (id: number, p: TagPayload) => http.put<null>(`/admin/tags/${id}`, p)
export const adminDeleteTag = (id: number) => http.delete<null>(`/admin/tags/${id}`)

// ---- Admin: pages ----
export const adminGetPages = () => http.get<PageSummary[]>('/admin/pages')
export const adminPutPage = (slug: string, p: PagePayload) =>
  http.put<null>(`/admin/pages/${slug}`, p)
export const adminDeletePage = (slug: string) => http.delete<null>(`/admin/pages/${slug}`)

// ---- Admin: upload ----
export const adminUploadImage = (file: File) => {
  const form = new FormData()
  form.append('file', file)
  return http.upload<UploadResult>('/admin/upload', form)
}

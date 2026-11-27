import axios, { type AxiosRequestConfig } from 'axios'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'
import { toast } from '@/utils/toast'
import type { ApiResponse } from './types'

const instance = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: { 'Content-Type': 'application/json' },
})

function handleUnauthorized() {
  const auth = useAuthStore()
  if (auth.token) {
    auth.logout()
    toast.error('登录已过期，请重新登录')
  }
  if (router.currentRoute.value.path !== '/admin/login') {
    router.replace('/admin/login')
  }
}

// Request: attach JWT
instance.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.token) config.headers.Authorization = `Bearer ${auth.token}`
  return config
})

// Response: unwrap body, unify errors
instance.interceptors.response.use(
  (response) => {
    const res = response.data as ApiResponse
    if (res.code !== 0) {
      toast.error(res.message || '请求失败')
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    return res as never
  },
  (error) => {
    const status = error.response?.status
    const msg: string = error.response?.data?.message || error.message || '网络错误'
    if (status === 401) {
      // wrong credentials on the login page: just show the message
      if (error.config?.url?.includes('/auth/login')) {
        toast.error(msg)
      } else {
        handleUnauthorized()
      }
    } else {
      toast.error(msg)
    }
    return Promise.reject(new Error(msg))
  }
)

/** Typed helpers: the interceptor already unwraps axios' response to the API body. */
const http = {
  get: <T>(url: string, params?: object) =>
    instance.get(url, { params } as AxiosRequestConfig) as unknown as Promise<ApiResponse<T>>,
  post: <T>(url: string, data?: object) =>
    instance.post(url, data) as unknown as Promise<ApiResponse<T>>,
  put: <T>(url: string, data?: object) =>
    instance.put(url, data) as unknown as Promise<ApiResponse<T>>,
  patch: <T>(url: string, data?: object) =>
    instance.patch(url, data) as unknown as Promise<ApiResponse<T>>,
  delete: <T>(url: string) =>
    instance.delete(url) as unknown as Promise<ApiResponse<T>>,
  upload: <T>(url: string, formData: FormData) =>
    instance.post(url, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      timeout: 30000,
    }) as unknown as Promise<ApiResponse<T>>,
}

export default http

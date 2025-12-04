import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'

// 创建Axios实例
const service = axios.create({
  baseURL: '/api', // 代理前缀，对应vite.config.js
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json;charset=utf-8'
  }
})

// 请求拦截器：携带JWT令牌
service.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器：统一处理错误
service.interceptors.response.use(
  (response) => {
    const res = response.data
    // 后端返回code!=200则视为错误
    if (res.code !== 200) {
      ElMessage.error(res.msg || '请求失败')
      return Promise.reject(res)
    }
    return res
  },
  (error) => {
    // 401：令牌过期/未登录
    if (error.response?.status === 401) {
      const userStore = useUserStore()
      userStore.clearUser() // 清空用户状态
      ElMessageBox.confirm('登录状态已过期，请重新登录', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        window.location.href = '/login' // 跳转到登录页
      })
    }
    ElMessage.error(error.message || '服务器错误')
    return Promise.reject(error)
  }
)

export default service
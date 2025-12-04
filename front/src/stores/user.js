import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  // 状态：令牌、用户信息
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))

  // 保存用户信息（登录后调用）
  const setUser = (t, info) => {
    token.value = t
    userInfo.value = info
    // 持久化到本地存储
    localStorage.setItem('token', t)
    localStorage.setItem('userInfo', JSON.stringify(info))
  }

  // 清空用户信息（退出登录/令牌过期）
  const clearUser = () => {
    token.value = ''
    userInfo.value = {}
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  }

  return { token, userInfo, setUser, clearUser }
})
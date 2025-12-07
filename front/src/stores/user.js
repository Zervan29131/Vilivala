import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: {
      id: 0,
      username: '',
      avatar: '',
      role: ''
    }
  }),
  actions: {
    // 设置用户信息
    setUser(token, userInfo) {
      this.token = token
      this.userInfo = userInfo
      localStorage.setItem('token', token)
    },
    // 清空用户信息
    clearUser() {
      this.token = ''
      this.userInfo = { id: 0, username: '', avatar: '', role: '' }
      localStorage.removeItem('token')
    }
  }
})
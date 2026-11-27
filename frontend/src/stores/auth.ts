import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/api/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('blog_token') ?? '')
  const user = ref<User | null>(
    JSON.parse(localStorage.getItem('blog_user') ?? 'null')
  )

  const isLoggedIn = computed(() => !!token.value)

  function setAuth(newToken: string, newUser: User) {
    token.value = newToken
    user.value = newUser
    localStorage.setItem('blog_token', newToken)
    localStorage.setItem('blog_user', JSON.stringify(newUser))
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('blog_token')
    localStorage.removeItem('blog_user')
  }

  return { token, user, isLoggedIn, setAuth, logout }
})

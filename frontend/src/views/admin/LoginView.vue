<template>
  <div class="min-h-screen flex items-center justify-center bg-[#f6f6f7] dark:bg-[#1a1a1a] px-4">
    <form class="card p-8 w-full max-w-sm shadow-md" @submit.prevent="submit">
      <div class="text-center mb-6">
        <img :src="siteConfig.logo" alt="Logo" class="w-12 h-12 rounded-full mx-auto mb-3" />
        <h1 class="text-xl font-bold text-[#213547] dark:text-white">后台登录</h1>
        <p class="text-xs text-gray-400 mt-1">{{ siteConfig.title }}</p>
      </div>
      <label class="block text-sm text-gray-600 dark:text-gray-300 mb-1">用户名</label>
      <input v-model="form.username" class="input-field mb-4" autocomplete="username" />
      <label class="block text-sm text-gray-600 dark:text-gray-300 mb-1">密码</label>
      <input v-model="form.password" type="password" class="input-field mb-6" autocomplete="current-password" />
      <button class="btn-primary w-full justify-center" :disabled="loading || !form.username || !form.password">
        {{ loading ? '登录中…' : '登录' }}
      </button>
      <router-link to="/" class="block text-center text-xs text-gray-400 mt-4 hover:text-[#42b883]">
        ← 返回前台
      </router-link>
    </form>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { login } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { siteConfig } from '@/site.config'
import { toast } from '@/utils/toast'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const loading = ref(false)
const form = reactive({ username: '', password: '' })

async function submit() {
  loading.value = true
  try {
    const res = await login(form)
    auth.setAuth(res.data.token, res.data.user)
    toast.success('登录成功')
    router.replace((route.query.redirect as string) || '/admin/posts')
  } catch {
    /* error toast shown by interceptor */
  } finally {
    loading.value = false
  }
}
</script>

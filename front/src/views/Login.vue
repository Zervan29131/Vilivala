<template>
  <div class="login-container">
    <div class="login-card">
      <h2 class="login-title">用户登录</h2>
      <form @submit.prevent="handleLogin">
        <div class="form-item">
          <label class="form-label">用户名</label>
          <input
            v-model="form.username"
            type="text"
            class="form-input"
            placeholder="请输入用户名"
            required
          />
        </div>
        <div class="form-item">
          <label class="form-label">密码</label>
          <input
            v-model="form.password"
            type="password"
            class="form-input"
            placeholder="请输入密码"
            required
          />
        </div>
        <div class="form-error" v-if="errorMsg">{{ errorMsg }}</div>
        <button type="submit" class="btn-submit" :disabled="loading">
          <span v-if="loading">登录中...</span>
          <span v-else>登录</span>
        </button>
      </form>
      <div class="login-link">
        还没有账号？<router-link to="/register">立即注册</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import axios from 'axios'

const router = useRouter()
const route = useRoute()
const loading = ref(false)
const errorMsg = ref('')
const form = ref({
  username: '',
  password: ''
})

// 登录逻辑 + 跳转处理
const handleLogin = async () => {
  try {
    loading.value = true
    errorMsg.value = ''
    // 校验表单
    if (!form.value.username || !form.value.password) {
      errorMsg.value = '用户名和密码不能为空'
      return
    }
    // 请求后端登录接口
    const res = await axios.post('/api/v1/user/login', form.value)
    if (res.data.code === 200) {
      // 存储token
      localStorage.setItem('token', res.data.data?.token || '')
      // 跳转：优先跳转会来源页，否则跳首页
      const redirect = route.query.redirect || '/'
      router.push(redirect)
    } else {
      errorMsg.value = res.data.msg || '登录失败'
    }
  } catch (err) {
    errorMsg.value = '服务器错误，请稍后重试'
    console.error('登录失败：', err)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
}
.login-card {
  width: 400px;
  padding: 30px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}
.login-title {
  text-align: center;
  margin-bottom: 20px;
  color: #333;
}
.form-item {
  margin-bottom: 16px;
}
.form-label {
  display: block;
  margin-bottom: 6px;
  color: #666;
  font-size: 14px;
}
.form-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
}
.form-input:focus {
  outline: none;
  border-color: #409eff;
}
.form-error {
  color: #f56c6c;
  font-size: 13px;
  margin-bottom: 10px;
  text-align: center;
}
.btn-submit {
  width: 100%;
  padding: 10px;
  background-color: #409eff;
  color: #fff;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.2s;
}
.btn-submit:disabled {
  background-color: #a0cfff;
  cursor: not-allowed;
}
.btn-submit:hover:not(:disabled) {
  background-color: #337ecc;
}
.login-link {
  margin-top: 16px;
  text-align: center;
  font-size: 14px;
  color: #666;
}
.login-link a {
  color: #409eff;
  text-decoration: none;
}
.login-link a:hover {
  text-decoration: underline;
}
</style>
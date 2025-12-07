<template>
  <div class="register-container">
    <div class="register-card">
      <h2 class="register-title">用户注册</h2>
      <form @submit.prevent="handleRegister">
        <div class="form-item">
          <label class="form-label">用户名</label>
          <input
            v-model="form.username"
            type="text"
            class="form-input"
            placeholder="请输入用户名（2-20位）"
            required
          />
        </div>
        <div class="form-item">
          <label class="form-label">密码</label>
          <input
            v-model="form.password"
            type="password"
            class="form-input"
            placeholder="请输入密码（6位以上）"
            required
          />
        </div>
        <div class="form-item">
          <label class="form-label">确认密码</label>
          <input
            v-model="form.confirmPwd"
            type="password"
            class="form-input"
            placeholder="请再次输入密码"
            required
          />
        </div>
        <div class="form-error" v-if="errorMsg">{{ errorMsg }}</div>
        <button type="submit" class="btn-submit" :disabled="loading">
          <span v-if="loading">注册中...</span>
          <span v-else>注册</span>
        </button>
      </form>
      <div class="register-link">
        已有账号？<router-link to="/login">立即登录</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()
const loading = ref(false)
const errorMsg = ref('')
const form = ref({
  username: '',
  password: '',
  confirmPwd: ''
})

// 注册逻辑 + 跳转处理
const handleRegister = async () => {
  try {
    loading.value = true
    errorMsg.value = ''
    // 表单校验
    if (!form.value.username || form.value.username.length < 2 || form.value.username.length > 20) {
      errorMsg.value = '用户名长度需在2-20位之间'
      return
    }
    if (!form.value.password || form.value.password.length < 6) {
      errorMsg.value = '密码长度需大于6位'
      return
    }
    if (form.value.password !== form.value.confirmPwd) {
      errorMsg.value = '两次输入的密码不一致'
      return
    }
    // 请求后端注册接口
    const res = await axios.post('/api/v1/user/register', {
      username: form.value.username,
      password: form.value.password,
      avatar: '' // 默认为空
    })
    if (res.data.code === 200) {
      // 注册成功后跳转到登录页
      router.push({ path: '/login', query: { tip: '注册成功，请登录' } })
    } else {
      errorMsg.value = res.data.msg || '注册失败'
    }
  } catch (err) {
    errorMsg.value = '服务器错误，请稍后重试'
    console.error('注册失败：', err)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
}
.register-card {
  width: 400px;
  padding: 30px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}
.register-title {
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
  background-color: #67c23a;
  color: #fff;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.2s;
}
.btn-submit:disabled {
  background-color: #b3e19d;
  cursor: not-allowed;
}
.btn-submit:hover:not(:disabled) {
  background-color: #52af29;
}
.register-link {
  margin-top: 16px;
  text-align: center;
  font-size: 14px;
  color: #666;
}
.register-link a {
  color: #409eff;
  text-decoration: none;
}
.register-link a:hover {
  text-decoration: underline;
}
</style>
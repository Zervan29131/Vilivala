<template>
  <div class="layout-container">
    <!-- 头部导航 -->
    <header class="layout-header">
      <div class="header-content">
        <router-link to="/" class="logo">Vilivala博客</router-link>
        <nav class="nav-menu">
          <router-link to="/" class="nav-item" active-class="active">首页</router-link>
          <router-link to="/publish" class="nav-item" active-class="active" v-if="hasToken">发布文章</router-link>
        </nav>
        <div class="user-actions">
          <router-link to="/login" class="btn-login" v-if="!hasToken">登录</router-link>
          <router-link to="/register" class="btn-register" v-if="!hasToken">注册</router-link>
          <div class="user-info" v-if="hasToken">
            <span class="username">{{ username }}</span>
            <button class="btn-logout" @click="handleLogout">退出</button>
          </div>
        </div>
      </div>
    </header>

    <!-- 主要内容 -->
    <main class="layout-main">
      <slot />
    </main>

    <!-- 页脚 -->
    <footer class="layout-footer">
      <div class="footer-content">
        <p>© 2025 Vilivala博客 - 所有权利保留</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()
const username = ref('')
const hasToken = computed(() => !!localStorage.getItem('token'))

// 获取当前用户信息
const getUserInfo = async () => {
  if (!hasToken.value) return
  try {
    const res = await axios.get('/api/v1/user/info')
    if (res.data.code === 200) {
      username.value = res.data.data.username
    }
  } catch (err) {
    console.error('获取用户信息失败：', err)
    localStorage.removeItem('token')
    router.push('/login')
  }
}

// 退出登录
const handleLogout = () => {
  localStorage.removeItem('token')
  router.push('/login')
}

onMounted(() => {
  getUserInfo()
})
</script>

<style scoped>
.layout-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}
.layout-header {
  background-color: #fff;
  border-bottom: 1px solid #e5e7eb;
  position: sticky;
  top: 0;
  z-index: 100;
}
.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.logo {
  font-size: 20px;
  font-weight: bold;
  color: #409eff;
  text-decoration: none;
}
.nav-menu {
  display: flex;
  gap: 24px;
}
.nav-item {
  color: #333;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.2s;
}
.nav-item:hover, .nav-item.active {
  color: #409eff;
}
.user-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}
.btn-login {
  padding: 6px 12px;
  color: #409eff;
  border: 1px solid #409eff;
  border-radius: 4px;
  text-decoration: none;
  font-size: 14px;
  transition: all 0.2s;
}
.btn-login:hover {
  background-color: #f0f8ff;
}
.btn-register {
  padding: 6px 12px;
  background-color: #409eff;
  color: #fff;
  border-radius: 4px;
  text-decoration: none;
  font-size: 14px;
  transition: background-color 0.2s;
}
.btn-register:hover {
  background-color: #337ecc;
}
.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}
.username {
  font-size: 14px;
  color: #333;
}
.btn-logout {
  padding: 4px 8px;
  background-color: #f56c6c;
  color: #fff;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: background-color 0.2s;
}
.btn-logout:hover {
  background-color: #e45656;
}
.layout-main {
  flex: 1;
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}
.layout-footer {
  background-color: #f9f9f9;
  padding: 20px 0;
  text-align: center;
  border-top: 1px solid #e5e7eb;
}
.footer-content {
  font-size: 12px;
  color: #666;
}
</style>
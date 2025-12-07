<template>
  <!-- 仅头部导航，无其他布局 -->
  <el-header class="nav-header" style="background-color: #fff; border-bottom: 1px solid #e6e6e6">
    <div class="nav-container">
      <!-- Logo -->
      <el-link href="/" class="nav-logo" style="font-size: 20px; font-weight: bold; color: #409eff">
        Vilivala博客
      </el-link>

      <!-- 导航菜单 -->
      <el-menu :default-active="activePath" mode="horizontal" background-color="transparent" class="nav-menu">
        <el-menu-item index="/" @click="handleMenuClick('/')">
          <el-icon><House /></el-icon>
          <span>首页</span>
        </el-menu-item>

        <!-- 未登录 -->
        <template v-if="!userStore.token">
          <el-menu-item index="/login" @click="handleMenuClick('/login')">
            <el-icon><User /></el-icon>
            <span>登录</span>
          </el-menu-item>
          <el-menu-item index="/register" @click="handleMenuClick('/register')">
            <el-icon><UserFilled /></el-icon>
            <span>注册</span>
          </el-menu-item>
        </template>

        <!-- 已登录 -->
        <template v-else>
          <el-menu-item index="/publish" @click="handleMenuClick('/publish')">
            <el-icon><Edit /></el-icon>
            <span>发布文章</span>
          </el-menu-item>
          <el-sub-menu index="user-center" popper-class="user-submenu">
            <template #title>
              <el-avatar :src="userStore.userInfo.avatar || ''" size="small" class="user-avatar">
                {{ userStore.userInfo.username?.slice(0, 1) || 'U' }}
              </el-avatar>
              <span class="user-name">{{ userStore.userInfo.username || '未知用户' }}</span>
            </template>
            <el-menu-item index="/profile" @click="handleMenuClick('/profile')">
              <el-icon><User /></el-icon>
              <span>个人中心</span>
            </el-menu-item>
            <el-menu-item @click="handleLogout">
              <el-icon><SwitchButton /></el-icon>
              <span>退出登录</span>
            </el-menu-item>
          </el-sub-menu>
        </template>
      </el-menu>
    </div>
  </el-header>
</template>

<script setup>
import { computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { House, User, UserFilled, Edit, SwitchButton } from '@element-plus/icons-vue'
import { getUserInfo } from '@/api/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 当前激活路由
const activePath = computed(() => route.path)

// 菜单跳转
const handleMenuClick = (path) => {
  if (path === route.path) return
  router.push(path)
}

// 退出登录
const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定退出登录？', '提示', { type: 'warning' })
    userStore.clearUser()
    ElMessage.success('退出成功')
    router.push('/login')
  } catch (err) {
    ElMessage.info('已取消退出')
  }
}

// 拉取用户信息
const fetchUserInfo = async () => {
  if (userStore.token && !userStore.userInfo.id) {
    try {
      const res = await getUserInfo()
      if (res.code === 200) userStore.setUser(userStore.token, res.data)
      else throw new Error(res.msg)
    } catch (err) {
      console.error('获取用户信息失败：', err)
      userStore.clearUser()
      router.push('/login')
    }
  }
}

// 监听路由变化
onMounted(() => {
  fetchUserInfo()
  router.afterEach(() => (activePath.value = route.path))
})

onUnmounted(() => {
  router.afterEach(() => (activePath.value = route.path))
})
</script>

<style scoped>
.nav-header {
  padding: 0 !important;
  height: 60px !important;
  line-height: 60px;
  position: sticky;
  top: 0;
  z-index: 999;
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}

.nav-menu {
  border-bottom: none !important;
}

.user-avatar {
  margin-right: 8px;
}

.user-name {
  font-size: 14px;
  color: #333;
}

/* 响应式适配 */
@media (max-width: 768px) {
  .nav-logo {
    font-size: 16px !important;
  }
  .user-name {
    display: none;
  }
}

:deep(.el-sub-menu__title) {
  padding: 0 10px !important;
  height: 60px !important;
  line-height: 60px !important;
}

:deep(.el-menu-item) {
  padding: 0 15px !important;
}
</style>
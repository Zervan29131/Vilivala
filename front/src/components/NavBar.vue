<template>
  <el-header style="background-color: #fff; border-bottom: 1px solid #e6e6e6">
    <div class="container" style="display: flex; justify-content: space-between; align-items: center">
      <!-- 左侧logo/标题 -->
      <el-link href="/" style="font-size: 20px; font-weight: bold; color: #409eff">
        我的博客
      </el-link>

      <!-- 右侧导航 -->
      <el-menu :default-active="activePath" mode="horizontal" background-color="transparent">
        <el-menu-item index="/">
          <el-icon><House /></el-icon>
          <span>首页</span>
        </el-menu-item>

        <!-- 登录/注册（未登录） -->
        <template v-if="!userStore.token">
          <el-menu-item index="/login">
            <el-icon><User /></el-icon>
            <span>登录</span>
          </el-menu-item>
          <el-menu-item index="/register">
            <el-icon><UserFilled /></el-icon>
            <span>注册</span>
          </el-menu-item>
        </template>

        <!-- 已登录菜单 -->
        <template v-else>
          <el-menu-item index="/article/edit">
            <el-icon><Edit /></el-icon>
            <span>发布文章</span>
          </el-menu-item>
          <el-sub-menu index="profile" title="">
            <template #title>
              <el-avatar :src="userStore.userInfo.avatar || ''" size="small">
                {{ userStore.userInfo.username?.slice(0, 1) }}
              </el-avatar>
              <span style="margin-left: 5px">{{ userStore.userInfo.username }}</span>
            </template>
            <el-menu-item index="/profile">
              <el-icon><User /></el-icon>
              个人中心
            </el-menu-item>
            <el-menu-item @click="logout">
              <el-icon><SwitchButton /></el-icon>
              退出登录
            </el-menu-item>
          </el-sub-menu>
        </template>
      </el-menu>
    </div>
  </el-header>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { House, User, UserFilled, Edit, SwitchButton } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 当前激活的路由
const activePath = computed(() => route.path)

// 退出登录
const logout = () => {
  userStore.clearUser()
  ElMessage.success('退出登录成功')
  router.push('/login')
}

// 页面加载时获取用户信息（已登录状态）
onMounted(() => {
  if (userStore.token && !userStore.userInfo.id) {
    import('@/api/user').then(({ getUserInfo }) => {
      getUserInfo().then(res => {
        userStore.setUser(userStore.token, res.data)
      })
    })
  }
})
</script>

<style scoped>
.el-header {
  padding: 0;
}
.el-menu {
  border-bottom: none;
}
</style>
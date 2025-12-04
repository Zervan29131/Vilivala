<template>
  <el-container style="min-height: 100vh">
    <!-- 导航栏 -->
    <NavBar />
    <!-- 主内容区 -->
    <el-main>
      <div class="container">
        <router-view />
      </div>
    </el-main>
    <!-- 页脚 -->
    <el-footer style="text-align: center; padding: 20px; border-top: 1px solid #e6e6e6">
      <p>我的博客 ©2025 版权所有</p>
    </el-footer>
  </el-container>
</template>

<script setup>
import NavBar from '@/components/NavBar.vue'
</script>

<style scoped>
.el-container {
  padding: 0;
}

.el-main {
  padding: 0;
}
</style>
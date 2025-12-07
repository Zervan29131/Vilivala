<template>
  <el-container style="min-height: 100vh">
    <!-- 导航栏 -->
    <!-- <Layout /> -->
    <!-- 主内容区 -->
    <el-main>
      <div class="container">
        <router-view />
      </div>
    </el-main>

  </el-container>
</template>

<script setup>
// import NavBar from '@/components/NavBar.vue'
// import Layout from './components/Layout.vue';
</script>

<style scoped>
.el-container {
  padding: 0;
}

.el-main {
  padding: 0;
}
</style>
<<<<<<< HEAD
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
=======
<script setup>
import HelloWorld from './components/HelloWorld.vue'
>>>>>>> parent of 48535c3 (add)
</script>

<template>
  <div>
    <a href="https://vite.dev" target="_blank">
      <img src="/vite.svg" class="logo" alt="Vite logo" />
    </a>
    <a href="https://vuejs.org/" target="_blank">
      <img src="./assets/vue.svg" class="logo vue" alt="Vue logo" />
    </a>
  </div>
  <HelloWorld msg="Vite + Vue" />
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>

import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import ArticleDetail from '../views/ArticleDetail.vue'
import PublishArticle from '../views/PublishArticle.vue'

// 路由守卫：未登录拦截（仅拦截需要登录的页面）
const requireAuth = (to, from, next) => {
  const token = localStorage.getItem('token')
  if (token) {
    next()
  } else {
    // 未登录跳转到登录页，记录来源页
    next({ path: '/login', query: { redirect: to.fullPath } })
  }
}

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: { title: '首页 - 博客系统' }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { title: '登录 - 博客系统' },
    // 已登录时，访问登录页跳首页
    beforeEnter: (to, from, next) => {
      if (localStorage.getItem('token')) {
        next('/')
      } else {
        next()
      }
    }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { title: '注册 - 博客系统' },
    // 已登录时，访问注册页跳首页
    beforeEnter: (to, from, next) => {
      if (localStorage.getItem('token')) {
        next('/')
      } else {
        next()
      }
    }
  },
  {
    path: '/article/:id',
    name: 'ArticleDetail',
    component: ArticleDetail,
    meta: { title: '文章详情 - 博客系统' }
  },
  {
    path: '/publish',
    name: 'PublishArticle',
    component: PublishArticle,
    meta: { title: '发布文章 - 博客系统', requiresAuth: true },
    beforeEnter: requireAuth // 需要登录才能访问
  },
  // 404路由（必须放最后）
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// 全局路由守卫：设置页面标题
router.beforeEach((to, from, next) => {
  document.title = to.meta.title || '博客系统'
  next()
})

export default router
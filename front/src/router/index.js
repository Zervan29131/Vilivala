import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

// 导入页面组件
const Login = () => import('@/pages/auth/Login.vue')
const Register = () => import('@/pages/auth/Register.vue')
const Home = () => import('@/pages/home/Index.vue')
const ArticleDetail = () => import('@/pages/article/Detail.vue')
const ArticleEdit = () => import('@/pages/article/Edit.vue')
const Profile = () => import('@/pages/user/Profile.vue')

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
      meta: { title: '首页' }
    },
    {
      path: '/login',
      name: 'login',
      component: Login,
      meta: { title: '登录' }
    },
    {
      path: '/register',
      name: 'register',
      component: Register,
      meta: { title: '注册' }
    },
    {
      path: '/article/:id',
      name: 'articleDetail',
      component: ArticleDetail,
      meta: { title: '文章详情' }
    },
    {
      path: '/article/edit',
      name: 'articleEdit',
      component: ArticleEdit,
      meta: { title: '发布文章', requireAuth: true } // 需要登录
    },
    {
      path: '/article/edit/:id',
      name: 'articleEdit',
      component: ArticleEdit,
      meta: { title: '编辑文章', requireAuth: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: Profile,
      meta: { title: '个人中心', requireAuth: true }
    },
    // 404页面
    {
      path: '/:pathMatch(.*)*',
      redirect: '/'
    }
  ]
})

// 路由守卫：验证登录状态
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title || '我的博客'
  
  const userStore = useUserStore()
  // 需要登录的页面，检查令牌
  if (to.meta.requireAuth && !userStore.token) {
    next('/login') // 未登录跳转到登录页
  } else {
    next()
  }
})

export default router
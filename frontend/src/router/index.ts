import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { siteConfig } from '@/site.config'
import { updatePageMeta } from '@/utils/seo'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    // ---- Front-facing routes ----
    {
      path: '/',
      component: () => import('@/components/FrontLayout.vue'),
      children: [
        { path: '', component: () => import('@/views/front/HomeView.vue') },
        { path: 'posts/:slug', component: () => import('@/views/front/PostDetailView.vue') },
        { path: 'archives', component: () => import('@/views/front/TimelineView.vue') },
        { path: 'categories', component: () => import('@/views/front/CategoriesView.vue') },
        {
          path: 'categories/:slug',
          component: () => import('@/views/front/ArchiveView.vue'),
          props: (route) => ({ filterType: 'category', filterSlug: route.params.slug }),
        },
        {
          path: 'tags/:slug',
          component: () => import('@/views/front/ArchiveView.vue'),
          props: (route) => ({ filterType: 'tag', filterSlug: route.params.slug }),
        },
        { path: 'pages/:slug', component: () => import('@/views/front/PageView.vue') },
        { path: ':pathMatch(.*)*', component: () => import('@/views/front/NotFoundView.vue') },
      ],
    },
    // ---- Terminal browser start page ----
    {
      path: '/terminal',
      component: () => import('@/views/front/TerminalView.vue'),
    },
    // ---- Admin routes ----
    {
      path: '/admin/login',
      component: () => import('@/views/admin/LoginView.vue'),
    },
    {
      path: '/admin',
      component: () => import('@/views/admin/AdminLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '', redirect: '/admin/dashboard' },
        { path: 'dashboard', component: () => import('@/views/admin/DashboardView.vue') },
        { path: 'posts', component: () => import('@/views/admin/PostsView.vue') },
        { path: 'posts/new', component: () => import('@/views/admin/PostEditorView.vue') },
        { path: 'posts/:id/edit', component: () => import('@/views/admin/PostEditorView.vue') },
        { path: 'categories', component: () => import('@/views/admin/CategoriesView.vue') },
        { path: 'tags', component: () => import('@/views/admin/TagsView.vue') },
        { path: 'pages', component: () => import('@/views/admin/PagesView.vue') },
        { path: 'settings', component: () => import('@/views/admin/SiteSettingsView.vue') },
      ],
    },
  ],
  scrollBehavior(_to, _from, savedPosition) {
    if (savedPosition) return savedPosition
    return { top: 0 }
  },
})

router.beforeEach((to, _from, next) => {
  if (to.meta.requiresAuth) {
    const authStore = useAuthStore()
    if (!authStore.isLoggedIn) {
      next({ path: '/admin/login', query: { redirect: to.fullPath } })
      return
    }
  }
  next()
})

router.afterEach((to) => {
  if (!to.path.startsWith('/posts/') && !to.path.startsWith('/pages/') && to.path !== '/terminal') {
    const titleMap: Record<string, string> = {
      '/archives': '时间归档',
      '/categories': '分类与标签',
      '/admin/login': '后台登录',
    }
    updatePageMeta({
      title: titleMap[to.path],
      description: `${siteConfig.hero.title} · ${siteConfig.description}`,
      path: to.path,
    })
  }
})

export default router

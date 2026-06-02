import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import MainView from '../views/MainView.vue'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/',
      name: 'loginview',
      component: LoginView
    },
    {
      path: '/main',
      name: 'main',
      component: MainView,
      meta: { requiresAuth: true }
    },
    {
      path: '/first-set-auth',
      name: 'FirstSetAuth',
      component: () => import('@/views/FirstSetAuth.vue'),
    }
  ],
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  // **如果目标路由需要登录，但用户未登录，则跳转到 /login**
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    console.log('未登录，跳转到登录页')
    return next('/login')
  }

  // **如果用户已登录，且访问 login 页，则跳转到首页**
  if (authStore.isLoggedIn && to.path === '/login') {
    console.log('已登录，跳转到首页')
    return next('/main')
  }

  next() // 继续导航
})

export default router

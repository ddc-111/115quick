import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      redirect: '/dashboard'
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: () => import('@/views/Dashboard.vue')
    },
    {
      path: '/token',
      name: 'TokenConfig',
      component: () => import('@/views/TokenConfig.vue')
    },
    {
      path: '/downloads',
      name: 'DownloadList',
      component: () => import('@/views/DownloadList.vue')
    },
    {
      path: '/history',
      name: 'TaskHistory',
      component: () => import('@/views/TaskHistory.vue')
    },
    {
      path: '/settings',
      name: 'Settings',
      component: () => import('@/views/Settings.vue')
    }
  ]
})

export default router

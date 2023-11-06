import { createRouter, createWebHashHistory } from 'vue-router'
import { useCounterStore } from '@/stores/counter'

export const routers = [
  {
    path: '/404',
    name: '404',
    component: () => import('../views/404.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    name: 'home',
    component: () => import('../views/HomeView.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/payment/:id',
    name: 'payment',
    component: () => import('../views/PaymentView.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/LoginView.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/admin',
    name: 'dashboard',
    component: () => import('../views/admin/Dashboard.vue'),
    meta: {
      requiresAuth: true,
      title: '仪表盘',
      icon: 'dashboard'
    }
  },
  {
    path: '/admin/settings',
    name: 'settings',
    component: () => import('../views/admin/Settings.vue'),
    meta: {
      requiresAuth: true,
      title: '设置',
      icon: 'setting-1',
    }
  },
  {
    path: '/admin/orders',
    name: 'orders',
    component: () => import('../views/admin/OrdersView.vue'),
    meta: {
      requiresAuth: true,
      title: '订单管理',
      icon: 'table',
    }
  },
  {
    path: '/admin/qrcode',
    name: 'qrcode',
    component: () => import('../views/admin/QrcodeView.vue'),
    meta: {
      requiresAuth: true,
      title: '收款码管理',
      icon: 'qrcode',
    }
  },
  {
    path: '/admin/payments',
    name: 'payments',
    component: () => import('../views/admin/PaymentsView.vue'),
    meta: {
      requiresAuth: true,
      title: '支付记录',
      icon: 'list',
    }
  },
]

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: routers
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const counterStore = useCounterStore()
  // 如果访问的页面需要登录，但是用户没有登录，就跳转到登录页面
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!counterStore.isLogin) {
      next({
        path: '/login'
      })
    } else {
      next()
    }
  } else {
    next()
  }
})


export default router

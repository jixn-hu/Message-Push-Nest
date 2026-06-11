import { createRouter, createWebHistory, createWebHashHistory } from 'vue-router'
// import LoginInex from '../components/Login.vue'
import { CONSTANT } from '../constant'
import axios from 'axios'
import config from '../../config.js'

const router = createRouter({
  // 使用 HTML5 History 模式，确保 URL 变化反映在浏览器地址栏中
  history: createWebHashHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component:() => import('../components/Login.vue')
    },
    {
      path: '/',
      name: 'index',
      component: () => import('../components/Index.vue'),
      children: [
        {
          // 默认子路由，显示 Dashboard
          path: '',
          name: 'dashboard',
          component: () => import('../components/pages/dashboard/Dashboard.vue')
        },
        {
          path: 'sendlogs',
          name: 'sendlogs',
          component: () => import('../components/pages/sendLogs/SendLogs.vue')
        },
        {
          path: 'hostedmessage',
          name: 'hostedmessage',
          component: () => import('../components/pages/hostedMessage/HostedMessage.vue')
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('../components/pages/settings/Settings.vue')
        },
        {
          path: 'sendways',
          name: 'sendways',
          component: () => import('../components/pages/sendWays/SendWays.vue')
        },
        {
          path: 'templates',
          name: 'templates',
          component: () => import('../components/pages/messageTemplate/MessageTemplate.vue')
        }
      ]
    },
    // {
    //   path: '/settings',
    //   name: 'settings',
    //   component: () => import('../views/tabsTools/settings/settings.vue')
    // },
    // {
    //   path: '/hostedMessage',
    //   name: 'hostedmessage',
    //   component: () => import('../views/tabsTools/hostedMessage/hostedMessage.vue')
    // },
    {
      path: '/:catchAll(.*)',
      name: '404',
      component: () => import('../components/404.vue')
    },
  ]
})

// 自动登录：无 token 时自动调用 /auth 获取 token
const autoLogin = async () => {
  // 可通过 VITE_AUTO_LOGIN_ACCOUNT / VITE_AUTO_LOGIN_PASSWORD 环境变量配置
  const account = import.meta.env.VITE_AUTO_LOGIN_ACCOUNT || 'admin'
  const password = import.meta.env.VITE_AUTO_LOGIN_PASSWORD || '123456'
  try {
    const baseURL = config.apiUrl + config.pathPrefix
    const res = await axios.post(baseURL + '/auth', { username: account, passwd: password })
    if (res.data && res.data.code === 200 && res.data.data && res.data.data.token) {
      localStorage.setItem(CONSTANT.STORE_TOKEN_NAME, res.data.data.token)
      if (res.data.data.role) {
        localStorage.setItem('__message_nest_role__', res.data.data.role)
      }
      return true
    }
  } catch (e) {
    console.warn('自动登录失败，跳转登录页面')
  }
  return false
}

// 登录失效重定向到登录页面
router.beforeEach(async (to, from, next) => {
  const token = localStorage.getItem(CONSTANT.STORE_TOKEN_NAME);
  let isAuthenticated = Boolean(token && token.trim() !== '');

  // 404页面不需要登录验证
  if (to.name === '404') {
    next();
    return;
  }

  // 无 token 时尝试自动登录
  if (!isAuthenticated && to.path !== '/login') {
    const success = await autoLogin()
    if (success) {
      isAuthenticated = true
    }
  }

  // 如果没有token且不是访问登录页，跳转到登录页
  if (!isAuthenticated && to.path !== '/login') {
    next('/login');
  }
  // 如果有token且访问登录页，跳转到首页
  else if (isAuthenticated && to.path === '/login') {
    next('/');
  }
  // 普通用户不能访问渠道和设置页
  else if (isAuthenticated) {
    const role = localStorage.getItem('__message_nest_role__') || 'admin';
    if (role !== 'admin' && (to.path === '/sendways' || to.path === '/settings')) {
      next('/');
    } else {
      next();
    }
  }
  // 其他情况正常访问
  else {
    next();
  }
});

export default router

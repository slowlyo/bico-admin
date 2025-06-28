import { defineConfig } from '@umijs/max';

export default defineConfig({
  antd: {
    // 配置 antd 的 locale 为中文
    configProvider: {},
    // 暗色主题
    dark: false,
  },
  access: {},
  model: {},
  initialState: {},
  request: {
    dataField: 'data',
  },
  layout: {
    title: 'Bico Admin',
    locale: false, // 关闭国际化
  },
  routes: [
    {
      path: '/',
      redirect: '/dashboard',
    },
    {
      name: '仪表板',
      path: '/dashboard',
      component: './Home',
      icon: 'DashboardOutlined',
    },
    {
      name: '用户管理',
      path: '/users',
      component: './Users',
      icon: 'UserOutlined',
      access: 'canManageUsers',
    },
    {
      path: '/profile',
      component: './Profile',
      hideInMenu: true,
    },
    {
      path: '/login',
      component: './Login',
      layout: false,
    },
  ],
  npmClient: 'pnpm',
  hash: true,
  // 代理配置
  proxy: {
    '/admin': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
})
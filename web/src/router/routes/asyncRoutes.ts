import { RoutesAlias } from '../routesAlias'
import { AppRouteRecord } from '@/types/router'

/**
 * 菜单列表、异步路由
 *
 * 支持两种模式:
 * 前端静态配置 - 直接使用本文件中定义的路由配置
 * 后端动态配置 - 后端返回菜单数据，前端解析生成路由
 *
 * 菜单标题（title）:
 * 可以是 i18n 的 key，也可以是字符串，比如：'用户列表'
 *
 * RoutesAlias.Layout 指向的是布局组件，后端返回的菜单数据中，component 字段需要指向 /index/index
 * 路由元数据（meta）：异步路由在 asyncRoutes 中配置，静态路由在 staticRoutes 中配置
 */
export const asyncRoutes: AppRouteRecord[] = [
  {
    name: 'Dashboard',
    path: '/dashboard',
    component: RoutesAlias.Layout,
    meta: {
      title: '首页',
      icon: '&#xe721;'
      // 首页不需要权限控制，所有用户都可以访问
    },
    children: [
      {
        path: 'console',
        name: 'Console',
        component: RoutesAlias.Dashboard,
        meta: {
          title: '控制台',
          keepAlive: false,
          fixedTab: true
          // 控制台不需要权限控制
        }
      }
    ]
  },
  {
    path: '/system',
    name: 'System',
    component: RoutesAlias.Layout,
    meta: {
      title: '系统管理',
      icon: '&#xe7b9;',
      // 系统管理模块需要任意一个系统权限
      permissions: ['system.admin_user:list', 'system.role:list']
    },
    children: [
      {
        path: 'admin-user',
        name: 'AdminUsers',
        component: RoutesAlias.AdminUsers,
        meta: {
          title: '管理员管理',
          keepAlive: true,
          // 管理员管理页面需要查看权限
          permissions: ['system.admin_user:list']
        }
      },
      {
        path: 'admin-role',
        name: 'AdminRoles',
        component: RoutesAlias.AdminRoles,
        meta: {
          title: '角色管理',
          keepAlive: true,
          // 角色管理页面需要查看权限
          permissions: ['system.role:list']
        }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: RoutesAlias.Profile,
        meta: {
          title: '个人中心',
          keepAlive: true,
          isHide: true
          // 个人中心不需要权限控制，所有用户都可以访问
        }
      }
    ]
  },

]

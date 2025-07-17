/**
 * 快速入口配置
 * 包含：应用列表、快速链接等配置
 */
import { RoutesAlias } from '@/router/routesAlias'
import type { FastEnterConfig } from '@/types/config'

const fastEnterConfig: FastEnterConfig = {
  // 显示条件（屏幕宽度）
  minWidth: 1200,
  // 应用列表
  applications: [
    {
      name: '工作台',
      description: '系统概览与数据统计',
      icon: '&#xe721;',
      iconColor: '#377dff',
      path: RoutesAlias.Dashboard,
      enabled: true,
      order: 1
    },
  ],
  // 快速链接
  quickLinks: [
    {
      name: '登录',
      path: RoutesAlias.Login,
      enabled: true,
      order: 1
    },
    {
      name: '个人中心',
      path: RoutesAlias.Profile,
      enabled: true,
      order: 2
    }
  ]
}

export default Object.freeze(fastEnterConfig)

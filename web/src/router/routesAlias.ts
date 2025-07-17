/**
 * 路由别名，方便快速找到页面，同时可以用作路由跳转
 */
export enum RoutesAlias {
  // 布局和认证
  Layout = '/index/index', // 布局容器
  Login = '/auth/login', // 登录

  // 异常页面
  Exception403 = '/exception/403', // 403
  Exception404 = '/exception/404', // 404
  Exception500 = '/exception/500', // 500

  // 结果页面
  Success = '/result/success', // 成功
  Fail = '/result/fail', // 失败

  // 仪表板
  Dashboard = '/dashboard/console', // 工作台

  // 系统管理
  AdminUsers = '/system/admin-user', // 管理员管理
  AdminRoles = '/system/admin-role', // 角色管理
  Profile = '/system/profile' // 个人中心
}

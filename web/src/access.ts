import { InitialState } from './app';

export default (initialState: InitialState) => {
  // 在这里按照初始化数据定义项目中的权限，统一管理
  // 参考文档 https://umijs.org/docs/max/access

  const { currentUser, permissions = [] } = initialState || {};

  // 通用权限检查函数
  const hasPermission = (permission: string) => {
    return permissions.some(p => p === permission);
  };

  // 检查多个权限中的任意一个
  const hasAnyPermission = (permissionList: string[]) => {
    return permissionList.some(permission => hasPermission(permission));
  };

  // 检查是否已登录
  const isLogin = !!currentUser;

  return {
    // 基础权限
    isLogin,
    canSeeAdmin: isLogin,

    // 管理员管理权限
    canViewAdminUsers: hasPermission('system.admin_user:list'),
    canCreateAdminUser: hasPermission('system.admin_user:create'),
    canEditAdminUser: hasPermission('system.admin_user:update'),
    canDeleteAdminUser: hasPermission('system.admin_user:delete'),
    canResetAdminPassword: hasPermission('system.admin_user:reset_password'),

    // 角色管理权限
    canViewRoles: hasPermission('system.role:list'),
    canCreateRole: hasPermission('system.role:create'),
    canEditRole: hasPermission('system.role:update'),
    canDeleteRole: hasPermission('system.role:delete'),

    // 通用权限检查函数，供组件使用
    hasPermission,
    hasAnyPermission,
  };
};

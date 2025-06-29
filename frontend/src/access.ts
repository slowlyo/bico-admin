import { hasPermission, PERMISSIONS } from '@/constants/permissions';

export default (initialState: any) => {
  // 在这里按照初始化数据定义项目中的权限，统一管理
  // 参考文档 https://umijs.org/docs/max/access
  const { currentUser } = initialState ?? {};
  const userRole = currentUser?.role || '';

  return {
    // 是否可以访问管理员功能
    canSeeAdmin: hasPermission(userRole, PERMISSIONS.SYSTEM.VIEW),
    // 是否可以管理用户
    canManageUsers: hasPermission(userRole, PERMISSIONS.USER.VIEW),
    // 是否可以管理角色
    canManageRoles: hasPermission(userRole, PERMISSIONS.ROLE.VIEW),
    // 是否已登录
    isLogin: !!currentUser,

    // 具体权限检查
    canCreateUser: hasPermission(userRole, PERMISSIONS.USER.CREATE),
    canUpdateUser: hasPermission(userRole, PERMISSIONS.USER.UPDATE),
    canDeleteUser: hasPermission(userRole, PERMISSIONS.USER.DELETE),
    canManageUserStatus: hasPermission(userRole, PERMISSIONS.USER.MANAGE_STATUS),
    canResetPassword: hasPermission(userRole, PERMISSIONS.USER.RESET_PASSWORD),

    canCreateRole: hasPermission(userRole, PERMISSIONS.ROLE.CREATE),
    canUpdateRole: hasPermission(userRole, PERMISSIONS.ROLE.UPDATE),
    canDeleteRole: hasPermission(userRole, PERMISSIONS.ROLE.DELETE),
    canAssignPermissions: hasPermission(userRole, PERMISSIONS.ROLE.ASSIGN_PERMISSIONS),
  };
};

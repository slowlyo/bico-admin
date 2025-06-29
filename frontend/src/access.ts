import { hasPermission, PERMISSIONS } from '@/constants/permissions';

export default (initialState: any) => {
  // 在这里按照初始化数据定义项目中的权限，统一管理
  // 参考文档 https://umijs.org/docs/max/access
  const { currentUser, userPermissions } = initialState ?? {};
  const permissions = userPermissions || [];

  return {
    // 是否可以访问管理员功能
    canSeeAdmin: hasPermission(permissions, PERMISSIONS.SYSTEM.VIEW),
    // 是否可以管理用户
    canManageUsers: hasPermission(permissions, PERMISSIONS.USER.VIEW),
    // 是否可以管理角色
    canManageRoles: hasPermission(permissions, PERMISSIONS.ROLE.VIEW),
    // 是否已登录
    isLogin: !!currentUser,

    // 个人资料权限豁免 - 所有登录用户都可以访问
    canViewProfile: !!currentUser,
    canUpdateProfile: !!currentUser,
    canChangePassword: !!currentUser,

    // 具体权限检查
    canCreateUser: hasPermission(permissions, PERMISSIONS.USER.CREATE),
    canUpdateUser: hasPermission(permissions, PERMISSIONS.USER.UPDATE),
    canDeleteUser: hasPermission(permissions, PERMISSIONS.USER.DELETE),
    canManageUserStatus: hasPermission(permissions, PERMISSIONS.USER.MANAGE_STATUS),
    canResetPassword: hasPermission(permissions, PERMISSIONS.USER.RESET_PASSWORD),

    canCreateRole: hasPermission(permissions, PERMISSIONS.ROLE.CREATE),
    canUpdateRole: hasPermission(permissions, PERMISSIONS.ROLE.UPDATE),
    canDeleteRole: hasPermission(permissions, PERMISSIONS.ROLE.DELETE),
    canAssignPermissions: hasPermission(permissions, PERMISSIONS.ROLE.ASSIGN_PERMISSIONS),
  };
};

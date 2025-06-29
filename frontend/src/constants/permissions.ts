// 权限配置常量
export const PERMISSIONS = {
  // 系统管理权限
  SYSTEM: {
    VIEW: 'system:view',
    MANAGE: 'system:manage',
  },
  
  // 用户管理权限
  USER: {
    VIEW: 'user:view',
    CREATE: 'user:create',
    UPDATE: 'user:update',
    DELETE: 'user:delete',
    MANAGE_STATUS: 'user:manage_status',
    RESET_PASSWORD: 'user:reset_password',
  },
  
  // 角色管理权限
  ROLE: {
    VIEW: 'role:view',
    CREATE: 'role:create',
    UPDATE: 'role:update',
    DELETE: 'role:delete',
    ASSIGN_PERMISSIONS: 'role:assign_permissions',
  },
  
  // 个人资料权限
  PROFILE: {
    VIEW: 'profile:view',
    UPDATE: 'profile:update',
    CHANGE_PASSWORD: 'profile:change_password',
  },
} as const;

// 角色权限映射
export const ROLE_PERMISSIONS = {
  admin: [
    // 系统管理
    PERMISSIONS.SYSTEM.VIEW,
    PERMISSIONS.SYSTEM.MANAGE,
    
    // 用户管理
    PERMISSIONS.USER.VIEW,
    PERMISSIONS.USER.CREATE,
    PERMISSIONS.USER.UPDATE,
    PERMISSIONS.USER.DELETE,
    PERMISSIONS.USER.MANAGE_STATUS,
    PERMISSIONS.USER.RESET_PASSWORD,
    
    // 角色管理
    PERMISSIONS.ROLE.VIEW,
    PERMISSIONS.ROLE.CREATE,
    PERMISSIONS.ROLE.UPDATE,
    PERMISSIONS.ROLE.DELETE,
    PERMISSIONS.ROLE.ASSIGN_PERMISSIONS,
    
    // 个人资料
    PERMISSIONS.PROFILE.VIEW,
    PERMISSIONS.PROFILE.UPDATE,
    PERMISSIONS.PROFILE.CHANGE_PASSWORD,
  ],
  
  manager: [
    // 用户管理（部分权限）
    PERMISSIONS.USER.VIEW,
    PERMISSIONS.USER.CREATE,
    PERMISSIONS.USER.UPDATE,
    PERMISSIONS.USER.MANAGE_STATUS,
    
    // 个人资料
    PERMISSIONS.PROFILE.VIEW,
    PERMISSIONS.PROFILE.UPDATE,
    PERMISSIONS.PROFILE.CHANGE_PASSWORD,
  ],
  
  user: [
    // 个人资料
    PERMISSIONS.PROFILE.VIEW,
    PERMISSIONS.PROFILE.UPDATE,
    PERMISSIONS.PROFILE.CHANGE_PASSWORD,
  ],
} as const;

// 权限检查函数
export const hasPermission = (userRole: string, permission: string): boolean => {
  const rolePermissions = ROLE_PERMISSIONS[userRole as keyof typeof ROLE_PERMISSIONS];
  return rolePermissions ? rolePermissions.includes(permission) : false;
};

// 检查多个权限（需要全部满足）
export const hasAllPermissions = (userRole: string, permissions: string[]): boolean => {
  return permissions.every(permission => hasPermission(userRole, permission));
};

// 检查多个权限（满足其中一个即可）
export const hasAnyPermission = (userRole: string, permissions: string[]): boolean => {
  return permissions.some(permission => hasPermission(userRole, permission));
};

// 获取用户所有权限
export const getUserPermissions = (userRole: string): string[] => {
  return ROLE_PERMISSIONS[userRole as keyof typeof ROLE_PERMISSIONS] || [];
};

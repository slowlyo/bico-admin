// 权限配置常量 - 重新设计的层级结构
export const PERMISSIONS = {
  // 系统管理权限
  SYSTEM: {
    VIEW: 'system:view',
    MANAGE: 'system:manage',
    DASHBOARD: 'system:dashboard',
    SETTINGS: 'system:settings',
  },

  // 用户管理权限
  USER: {
    VIEW: 'user:view',
    CREATE: 'user:create',
    UPDATE: 'user:update',
    DELETE: 'user:delete',
    MANAGE_STATUS: 'user:manage_status',
    RESET_PASSWORD: 'user:reset_password',
    EXPORT: 'user:export',
    IMPORT: 'user:import',
  },

  // 角色管理权限
  ROLE: {
    VIEW: 'role:view',
    CREATE: 'role:create',
    UPDATE: 'role:update',
    DELETE: 'role:delete',
    ASSIGN_PERMISSIONS: 'role:assign_permissions',
    MANAGE_STATUS: 'role:manage_status',
  },

  // 权限管理权限
  PERMISSION: {
    VIEW: 'permission:view',
    CREATE: 'permission:create',
    UPDATE: 'permission:update',
    DELETE: 'permission:delete',
    MANAGE_TREE: 'permission:manage_tree',
  },

  // 个人资料权限（所有用户都有，无需权限检查）
  PROFILE: {
    VIEW: 'profile:view',
    UPDATE: 'profile:update',
    CHANGE_PASSWORD: 'profile:change_password',
    UPLOAD_AVATAR: 'profile:upload_avatar',
  },

  // 内容管理权限（预留）
  CONTENT: {
    VIEW: 'content:view',
    CREATE: 'content:create',
    UPDATE: 'content:update',
    DELETE: 'content:delete',
    PUBLISH: 'content:publish',
  },

  // 日志管理权限（预留）
  LOG: {
    VIEW: 'log:view',
    EXPORT: 'log:export',
    DELETE: 'log:delete',
  },
} as const;

// 角色权限映射 - 优化后的权限分配
export const ROLE_PERMISSIONS = {
  // 超级管理员 - 拥有所有权限
  admin: [
    // 系统管理
    PERMISSIONS.SYSTEM.VIEW,
    PERMISSIONS.SYSTEM.MANAGE,
    PERMISSIONS.SYSTEM.DASHBOARD,
    PERMISSIONS.SYSTEM.SETTINGS,

    // 用户管理
    PERMISSIONS.USER.VIEW,
    PERMISSIONS.USER.CREATE,
    PERMISSIONS.USER.UPDATE,
    PERMISSIONS.USER.DELETE,
    PERMISSIONS.USER.MANAGE_STATUS,
    PERMISSIONS.USER.RESET_PASSWORD,
    PERMISSIONS.USER.EXPORT,
    PERMISSIONS.USER.IMPORT,

    // 角色管理
    PERMISSIONS.ROLE.VIEW,
    PERMISSIONS.ROLE.CREATE,
    PERMISSIONS.ROLE.UPDATE,
    PERMISSIONS.ROLE.DELETE,
    PERMISSIONS.ROLE.ASSIGN_PERMISSIONS,
    PERMISSIONS.ROLE.MANAGE_STATUS,

    // 权限管理
    PERMISSIONS.PERMISSION.VIEW,
    PERMISSIONS.PERMISSION.CREATE,
    PERMISSIONS.PERMISSION.UPDATE,
    PERMISSIONS.PERMISSION.DELETE,
    PERMISSIONS.PERMISSION.MANAGE_TREE,

    // 内容管理
    PERMISSIONS.CONTENT.VIEW,
    PERMISSIONS.CONTENT.CREATE,
    PERMISSIONS.CONTENT.UPDATE,
    PERMISSIONS.CONTENT.DELETE,
    PERMISSIONS.CONTENT.PUBLISH,

    // 日志管理
    PERMISSIONS.LOG.VIEW,
    PERMISSIONS.LOG.EXPORT,
    PERMISSIONS.LOG.DELETE,

    // 个人资料（所有用户都有）
    PERMISSIONS.PROFILE.VIEW,
    PERMISSIONS.PROFILE.UPDATE,
    PERMISSIONS.PROFILE.CHANGE_PASSWORD,
    PERMISSIONS.PROFILE.UPLOAD_AVATAR,
  ],

  // 管理者 - 部分管理权限
  manager: [
    // 系统管理（仅查看）
    PERMISSIONS.SYSTEM.VIEW,
    PERMISSIONS.SYSTEM.DASHBOARD,

    // 用户管理（部分权限）
    PERMISSIONS.USER.VIEW,
    PERMISSIONS.USER.CREATE,
    PERMISSIONS.USER.UPDATE,
    PERMISSIONS.USER.MANAGE_STATUS,
    PERMISSIONS.USER.EXPORT,

    // 角色管理（仅查看）
    PERMISSIONS.ROLE.VIEW,

    // 内容管理
    PERMISSIONS.CONTENT.VIEW,
    PERMISSIONS.CONTENT.CREATE,
    PERMISSIONS.CONTENT.UPDATE,
    PERMISSIONS.CONTENT.PUBLISH,

    // 个人资料
    PERMISSIONS.PROFILE.VIEW,
    PERMISSIONS.PROFILE.UPDATE,
    PERMISSIONS.PROFILE.CHANGE_PASSWORD,
    PERMISSIONS.PROFILE.UPLOAD_AVATAR,
  ],

  // 普通用户 - 基础权限
  user: [
    // 系统管理（仅仪表板）
    PERMISSIONS.SYSTEM.DASHBOARD,

    // 内容管理（仅查看）
    PERMISSIONS.CONTENT.VIEW,

    // 个人资料
    PERMISSIONS.PROFILE.VIEW,
    PERMISSIONS.PROFILE.UPDATE,
    PERMISSIONS.PROFILE.CHANGE_PASSWORD,
    PERMISSIONS.PROFILE.UPLOAD_AVATAR,
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

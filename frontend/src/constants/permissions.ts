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

  // 注意：个人资料相关功能无需权限验证，所有登录用户都可以访问
  // 已移除 PROFILE 相关权限定义

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

// 注意：角色权限映射已移除，现在通过API动态获取用户权限
// 权限常量保留在前端，便于开发时引用和类型检查

// 权限检查函数 - 现在通过API获取用户权限进行检查
// 这些函数需要配合全局状态管理使用

/**
 * 检查用户是否有指定权限
 * @param userPermissions 用户权限列表（从API获取）
 * @param permission 要检查的权限
 */
export const hasPermission = (userPermissions: string[], permission: string): boolean => {
  return userPermissions.includes(permission);
};

/**
 * 检查多个权限（需要全部满足）
 * @param userPermissions 用户权限列表（从API获取）
 * @param permissions 要检查的权限列表
 */
export const hasAllPermissions = (userPermissions: string[], permissions: string[]): boolean => {
  return permissions.every(permission => hasPermission(userPermissions, permission));
};

/**
 * 检查多个权限（满足任意一个即可）
 * @param userPermissions 用户权限列表（从API获取）
 * @param permissions 要检查的权限列表
 */
export const hasAnyPermission = (userPermissions: string[], permissions: string[]): boolean => {
  return permissions.some(permission => hasPermission(userPermissions, permission));
};

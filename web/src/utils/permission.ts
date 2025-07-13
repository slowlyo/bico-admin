// 权限工具函数

/**
 * 检查用户是否有指定权限
 * @param userPermissions 用户权限列表
 * @param permission 需要检查的权限
 * @returns 是否有权限
 */
export function hasPermission(userPermissions: string[], permission: string): boolean {
  return userPermissions.includes(permission);
}

/**
 * 检查用户是否有任意一个权限
 * @param userPermissions 用户权限列表
 * @param permissions 需要检查的权限列表
 * @returns 是否有任意一个权限
 */
export function hasAnyPermission(userPermissions: string[], permissions: string[]): boolean {
  return permissions.some(permission => userPermissions.includes(permission));
}

/**
 * 检查用户是否有所有权限
 * @param userPermissions 用户权限列表
 * @param permissions 需要检查的权限列表
 * @returns 是否有所有权限
 */
export function hasAllPermissions(userPermissions: string[], permissions: string[]): boolean {
  return permissions.every(permission => userPermissions.includes(permission));
}

/**
 * 权限常量定义
 */
export const PERMISSIONS = {
  // 管理员管理
  ADMIN_USER_LIST: 'admin_user:list',
  ADMIN_USER_CREATE: 'admin_user:create',
  ADMIN_USER_UPDATE: 'admin_user:update',
  ADMIN_USER_DELETE: 'admin_user:delete',
  ADMIN_USER_RESET_PASSWORD: 'admin_user:reset_password',

  // 角色管理
  ROLE_LIST: 'role:list',
  ROLE_CREATE: 'role:create',
  ROLE_UPDATE: 'role:update',
  ROLE_DELETE: 'role:delete',
} as const;

/**
 * 权限级别常量
 */
export const PERMISSION_LEVELS = {
  VIEW: 1,    // 查看权限
  ACTION: 2,  // 操作权限
  MANAGE: 3,  // 管理权限
  SUPER: 4,   // 超级权限
} as const;

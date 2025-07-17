import { storeToRefs } from 'pinia'
import { useUserStore } from '@/store/modules/user'

const userStore = useUserStore()

/**
 * 权限检查工具
 * 用法：
 * const { hasAuth } = useAuth()
 * hasAuth('system.admin_user:list') // 检查完整权限代码
 *
 * const { hasAuth } = useAuth('system.admin_user')
 * hasAuth('list') // 自动拼接为 'system.admin_user:list'
 */
export const useAuth = (prefix?: string) => {
  const { permissions } = storeToRefs(userStore)

  /**
   * 构建完整的权限代码
   * @param code 权限代码或操作名
   * @returns 完整的权限代码
   */
  const buildPermissionCode = (code: string): string => {
    if (!prefix || code.includes(':')) {
      return code
    }
    return `${prefix}:${code}`
  }

  /**
   * 检查是否拥有指定权限
   * @param permissionCode 权限代码或操作名
   * @returns 是否有权限
   */
  const hasAuth = (permissionCode: string): boolean => {
    const fullCode = buildPermissionCode(permissionCode)
    return permissions.value.includes(fullCode)
  }

  /**
   * 检查是否拥有多个权限中的任意一个
   * @param permissionCodes 权限代码或操作名数组
   * @returns 是否有任意一个权限
   */
  const hasAnyAuth = (permissionCodes: string[]): boolean => {
    return permissionCodes.some(code => hasAuth(code))
  }

  /**
   * 检查是否拥有所有指定权限
   * @param permissionCodes 权限代码或操作名数组
   * @returns 是否拥有所有权限
   */
  const hasAllAuth = (permissionCodes: string[]): boolean => {
    return permissionCodes.every(code => hasAuth(code))
  }

  return {
    hasAuth,
    hasAnyAuth,
    hasAllAuth
  }
}

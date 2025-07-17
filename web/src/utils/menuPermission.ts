import { AppRouteRecord } from '@/types/router'

/**
 * 菜单权限过滤工具
 * 根据用户权限过滤菜单项
 */

/**
 * 检查用户是否有权限访问某个菜单
 * @param requiredPermissions 菜单需要的权限列表
 * @param userPermissions 用户拥有的权限列表
 * @returns 是否有权限
 */
export function hasMenuPermission(
  requiredPermissions: string[] | undefined,
  userPermissions: string[]
): boolean {
  // 如果菜单没有权限要求，默认允许访问
  if (!requiredPermissions || requiredPermissions.length === 0) {
    return true
  }

  // 检查用户是否拥有任意一个所需权限
  return requiredPermissions.some(permission => 
    userPermissions.includes(permission)
  )
}

/**
 * 递归过滤菜单项
 * @param menuItems 菜单项列表
 * @param userPermissions 用户权限列表
 * @returns 过滤后的菜单项列表
 */
export function filterMenuItems(
  menuItems: AppRouteRecord[],
  userPermissions: string[]
): AppRouteRecord[] {
  return menuItems
    .filter(item => {
      // 检查当前菜单项的权限
      const hasPermission = hasMenuPermission(item.meta?.permissions, userPermissions)
      
      if (!hasPermission) {
        return false
      }

      // 如果有子菜单，递归过滤子菜单
      if (item.children && item.children.length > 0) {
        const filteredChildren = filterMenuItems(item.children, userPermissions)
        // 如果所有子菜单都被过滤掉了，则隐藏父菜单
        return filteredChildren.length > 0
      }

      // 叶子节点且有权限，保留
      return true
    })
    .map(item => ({
      ...item,
      children: item.children ? filterMenuItems(item.children, userPermissions) : undefined
    }))
}

/**
 * 检查路由是否需要权限验证
 * @param route 路由配置
 * @param userPermissions 用户权限列表
 * @returns 是否有权限访问
 */
export function hasRoutePermission(
  route: AppRouteRecord,
  userPermissions: string[]
): boolean {
  return hasMenuPermission(route.meta?.permissions, userPermissions)
}

/**
 * 获取用户可访问的第一个菜单路径
 * @param menuItems 菜单项列表
 * @param userPermissions 用户权限列表
 * @returns 第一个可访问的菜单路径
 */
export function getFirstAccessibleMenuPath(
  menuItems: AppRouteRecord[],
  userPermissions: string[]
): string | null {
  const filteredMenus = filterMenuItems(menuItems, userPermissions)

  for (const menu of filteredMenus) {
    if (menu.children && menu.children.length > 0) {
      // 递归查找子菜单中的第一个可访问路径
      const childPath = getFirstAccessibleMenuPath(menu.children, userPermissions)
      if (childPath) {
        return childPath
      }
    } else if (menu.path && !menu.meta?.isHide) {
      // 找到第一个可访问的叶子节点
      return menu.path
    }
  }

  return null
}

/**
 * 调试菜单权限过滤结果
 * @param menuItems 原始菜单项列表
 * @param userPermissions 用户权限列表
 */
export function debugMenuPermissions(
  menuItems: AppRouteRecord[],
  userPermissions: string[]
): void {
  console.group('🔐 菜单权限过滤调试')
  console.log('用户权限:', userPermissions)
  console.log('原始菜单数量:', menuItems.length)

  const filteredMenus = filterMenuItems(menuItems, userPermissions)
  console.log('过滤后菜单数量:', filteredMenus.length)

  console.log('过滤后的菜单结构:')
  filteredMenus.forEach(menu => {
    console.log(`- ${menu.meta?.title} (${menu.path})`)
    if (menu.children) {
      menu.children.forEach(child => {
        console.log(`  - ${child.meta?.title} (${child.path})`)
      })
    }
  })

  console.groupEnd()
}

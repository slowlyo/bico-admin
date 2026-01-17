/**
 * useAuth - 权限验证管理
 *
 * 提供统一的权限验证功能，支持前端和后端两种权限模式。
 * 用于控制页面按钮、操作等功能的显示和访问权限。
 *
 * ## 主要功能
 *
 * 1. 权限检查 - 检查用户是否拥有指定的权限标识
 * 2. 双模式支持 - 自动适配前端模式和后端模式的权限验证
 * 3. 前端模式 - 从用户信息中获取按钮权限列表（如 ['add', 'edit', 'delete']）
 * 4. 后端模式 - 从路由 meta 配置中获取权限列表（如 [{ authMark: 'add' }]）
 *
 * ## 使用示例
 *
 * ```typescript
 * const { hasAuth } = useAuth()
 *
 * // 检查是否有新增权限
 * if (hasAuth('add')) {
 *   // 显示新增按钮
 * }
 *
 * // 在模板中使用
 * <el-button v-if="hasAuth('edit')">编辑</el-button>
 * <el-button v-if="hasAuth('delete')">删除</el-button>
 * ```
 *
 * @module useAuth
 * @author Art Design Pro Team
 */

import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useUserStore } from '@/store/modules/user'
import { useAppMode } from '@/hooks/core/useAppMode'
import type { AppRouteRecord } from '@/types/router'

type AuthItem = NonNullable<AppRouteRecord['meta']['authList']>[number]

export const useAuth = () => {
  const route = useRoute()
  const userStore = useUserStore()
  const { isFrontendMode } = useAppMode()
  const { info } = storeToRefs(userStore)

  /**
   * 检查是否拥有某权限标识（前后端模式通用）
   * @param auth 权限标识
   * @returns 是否有权限
   */
  const hasAuth = (auth: string): boolean => {
    if (!auth) return true

    // 超级管理员直接放行
    if (info.value?.username === 'admin') {
      return true
    }

    // 1. 优先校验用户信息中的 permissions 列表（对应后端返回的权限标识数组）
    const permissions = info.value?.permissions ?? []
    if (permissions.includes(auth)) {
      return true
    }

    // 2. 如果是前端模式且 permissions 中没有，则直接返回 false
    if (isFrontendMode.value) {
      return false
    }

    // 3. 后端模式：额外校验路由 meta 中配置的权限列表（authMark）
    const backendAuthList: AuthItem[] = Array.isArray(route?.meta?.authList)
      ? (route?.meta?.authList as AuthItem[])
      : []

    return backendAuthList.some((item) => item?.authMark === auth)
  }

  return {
    hasAuth
  }
}

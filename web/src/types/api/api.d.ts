/**
 * API 接口类型定义模块
 *
 * 提供所有后端接口的类型定义
 *
 * ## 主要功能
 *
 * - 通用类型（分页参数、响应结构等）
 * - 认证类型（登录、用户信息等）
 * - 系统管理类型（用户、角色等）
 * - 全局命名空间声明
 *
 * ## 使用场景
 *
 * - API 请求参数类型约束
 * - API 响应数据类型定义
 * - 接口文档类型同步
 *
 * ## 注意事项
 *
 * - 在 .vue 文件使用需要在 eslint.config.mjs 中配置 globals: { Api: 'readonly' }
 * - 使用全局命名空间，无需导入即可使用
 *
 * ## 使用方式
 *
 * ```typescript
 * const params: Api.Auth.LoginParams = { userName: 'admin', password: '123456' }
 * const response: Api.Auth.UserInfo = await fetchUserInfo()
 * ```
 *
 * @module types/api/api
 * @author Art Design Pro Team
 */

declare namespace Api {
  /** 通用类型 */
  namespace Common {
    /** 分页参数 (前端状态) */
    interface PaginationParams {
      /** 当前页码 */
      current: number
      /** 每页条数 */
      size: number
      /** 总条数 */
      total: number
    }

    /** 通用请求分页参数 */
    interface CommonSearchParams {
      /** 当前页码 */
      page?: number
      /** 每页条数 */
      pageSize?: number
      /** 排序字段 */
      sortField?: string
      /** 排序顺序 */
      sortOrder?: string
    }

    /** 分页响应基础结构 */
    interface PaginatedResponse<T = any> {
      /** 数据列表 */
      list: T[]
      /** 总条数 */
      total: number
      /** 当前页码 (可选) */
      page?: number
      /** 每页条数 (可选) */
      pageSize?: number
    }

    /** 启用状态 */
    type EnableStatus = '1' | '2'
  }

  /** 认证类型 */
  namespace Auth {
    /** 登录参数 */
    interface LoginParams {
      username: string
      password: string
      captchaId: string
      captchaCode: string
    }

    /** 验证码响应 */
    interface CaptchaResponse {
      id: string
      image: string
    }

    /** 登录响应 */
    interface LoginResponse {
      token: string
    }

    /** 用户信息 */
    interface UserInfo {
      id: number
      username: string
      name: string
      avatar: string
      enabled: boolean
      permissions: string[]
    }
  }

  /** 应用配置类型 */
  namespace App {
    interface Config {
      name: string
      logo: string
      debug: boolean
    }
  }

  /** 系统管理类型 */
  namespace SystemManage {
    /** 用户列表 */
    type UserList = Api.Common.PaginatedResponse<UserListItem>

    /** 用户列表项 */
    interface UserListItem {
      id: number
      username: string
      name: string
      avatar: string
      enabled: boolean
      roles?: RoleListItem[]
      created_at: string
      updated_at: string
    }

    /** 用户创建/更新参数 */
    interface UserParams {
      username?: string
      password?: string
      name?: string
      avatar?: string
      enabled?: boolean
      role_ids?: number[]
    }

    /** 用户搜索参数 */
    interface UserSearchParams extends Api.Common.CommonSearchParams {
      username?: string
      name?: string
      enabled?: boolean
    }

    /** 角色列表 */
    type RoleList = Api.Common.PaginatedResponse<RoleListItem>

    /** 角色列表项 */
    interface RoleListItem {
      id: number
      name: string
      slug: string
      code?: string
      description: string
      enabled: boolean
      permissions?: string[]
      created_at: string
      updated_at: string
    }

    /** 角色创建/更新参数 */
    interface RoleParams {
      name?: string
      slug?: string
      code?: string
      description?: string
      enabled?: boolean
    }

    /** 角色搜索参数 */
    interface RoleSearchParams extends Api.Common.CommonSearchParams {
      name?: string
      code?: string
      description?: string
      enabled?: boolean
    }

    /** 权限树节点 */
    interface Permission {
      key: string
      label: string
      children?: Permission[]
    }
  }
}

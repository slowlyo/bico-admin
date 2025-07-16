/**
 * namespace: Api
 *
 * 所有接口相关类型定义
 * 在.vue文件使用会报错，需要在 eslint.config.mjs 中配置 globals: { Api: 'readonly' }
 */
declare namespace Api {
  /** 基础类型 */
  namespace Http {
    /** 基础响应 */
    interface BaseResponse<T = any> {
      // 状态码
      code: number
      // 消息 (支持两种字段名)
      msg?: string
      message?: string
      // 数据
      data: T
    }
  }

  /** 通用类型 */
  namespace Common {
    /** 分页参数 */
    interface PaginatingParams {
      /** 当前页码 */
      current: number
      /** 每页条数 */
      size: number
      /** 总条数 */
      total: number
    }

    /** 通用搜索参数 */
    type PaginatingSearchParams = Pick<PaginatingParams, 'current' | 'size'>

    /** 启用状态 */
    type EnableStatus = '1' | '2'
  }

  /** 认证类型 */
  namespace Auth {
    /** 登录参数 */
    interface LoginParams {
      username: string
      password: string
      captcha: string
    }

    /** 登录响应 */
    interface LoginResponse {
      token: string
      expires_at: string
      user_info: User.UserInfo
      permissions: string[]
    }
  }

  /** 用户类型 */
  namespace User {
    /** 用户信息 */
    interface UserInfo {
      id: number
      username: string
      name: string
      avatar: string
      email: string
      phone: string
      status: number
      status_text: string
      last_login_at?: string
      remark: string
      can_delete: boolean
      can_disable: boolean
      roles: UserRole[]
      created_at: string
      updated_at: string
    }

    /** 用户角色 */
    interface UserRole {
      id: number
      name: string
      code: string
      description: string
    }

    /** 用户列表数据 */
    interface UserListData {
      records: UserListItem[]
      current: number
      size: number
      total: number
    }

    /** 用户列表项 */
    interface UserListItem {
      id: number
      avatar: string
      createBy: string
      createTime: string
      updateBy: string
      updateTime: string
      status: '1' | '2' | '3' | '4' // 1: 在线 2: 离线 3: 异常 4: 注销
      userName: string
      userGender: string
      nickName: string
      userPhone: string
      userEmail: string
      userRoles: string[]
    }
  }

  /** 管理员用户类型 */
  namespace AdminUser {
    /** 管理员用户信息 */
    interface AdminUserInfo {
      id: number
      username: string
      name: string
      avatar?: string
      email?: string
      phone?: string
      status: number
      status_text: string
      last_login_at?: string
      remark?: string
      can_delete: boolean
      can_disable: boolean
      roles: AdminUserRole[]
      created_at: string
      updated_at: string
    }

    /** 管理员用户角色 */
    interface AdminUserRole {
      id: number
      name: string
      code: string
      description?: string
    }

    /** 管理员用户列表数据 */
    interface AdminUserListData {
      list: AdminUserInfo[]
      total: number
      page: number
      page_size: number
      total_pages: number
    }

    /** 管理员用户创建请求 */
    interface AdminUserCreateRequest {
      username: string
      password: string
      name: string
      avatar?: string
      email?: string
      phone?: string
      remark?: string
      enabled: boolean
      role_ids?: number[]
    }

    /** 管理员用户更新请求 */
    interface AdminUserUpdateRequest {
      username: string
      password?: string
      name: string
      avatar?: string
      email?: string
      phone?: string
      remark?: string
      enabled: boolean
      role_ids?: number[]
    }

    /** 管理员用户列表查询参数 */
    interface AdminUserListParams extends Common.PaginatingParams {
      username?: string
      name?: string
      status?: number
      role_id?: number
      sort_by?: string
      sort_desc?: boolean
      page?: number
      page_size?: number
      [key: string]: unknown
    }

    /** 管理员用户状态更新请求 */
    interface AdminUserStatusRequest {
      status: number
    }
  }

  /** 角色类型 */
  namespace Role {
    /** 角色信息 */
    interface RoleInfo {
      id: number
      name: string
      code: string
      description?: string
      status: number
      status_text: string
      created_at: string
      updated_at: string
    }

    /** 角色选项 */
    interface RoleOption {
      id: number
      name: string
      code: string
    }
  }
}

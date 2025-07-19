import request from '@/utils/http'

// 管理员用户相关类型定义
export namespace AdminUserTypes {
  /** 管理员用户信息 */
  export interface AdminUserInfo {
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
  export interface AdminUserRole {
    id: number
    name: string
    code: string
    description?: string
  }

  /** 管理员用户列表数据 */
  export interface AdminUserListData {
    list: AdminUserInfo[]
    total: number
    page: number
    page_size: number
    total_pages: number
  }

  /** 管理员用户创建请求 */
  export interface AdminUserCreateRequest {
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
  export interface AdminUserUpdateRequest {
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
  export interface AdminUserListParams {
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
  export interface AdminUserStatusRequest {
    status: number
  }
}

// 角色相关类型定义（在adminUserApi中也会用到）
export namespace RoleTypes {
  /** 角色选项 */
  export interface RoleOption {
    id: number
    name: string
    code: string
    description?: string
  }
}

/**
 * 管理员用户API服务
 */
export class AdminUserService {
  /**
   * 获取管理员用户列表
   */
  static getAdminUserList(params: AdminUserTypes.AdminUserListParams) {
    return request.get<Api.Http.BaseResponse<AdminUserTypes.AdminUserListData>>({
      url: '/admin-api/admin-users',
      params
    })
  }

  /**
   * 根据ID获取管理员用户
   */
  static getAdminUserById(id: number) {
    return request.get<Api.Http.BaseResponse<AdminUserTypes.AdminUserInfo>>({
      url: `/admin-api/admin-users/${id}`
    })
  }

  /**
   * 创建管理员用户
   */
  static createAdminUser(data: AdminUserTypes.AdminUserCreateRequest) {
    return request.post<Api.Http.BaseResponse<AdminUserTypes.AdminUserInfo>>({
      url: '/admin-api/admin-users',
      data: data
    })
  }

  /**
   * 更新管理员用户
   */
  static updateAdminUser(id: number, data: AdminUserTypes.AdminUserUpdateRequest) {
    return request.put<Api.Http.BaseResponse<AdminUserTypes.AdminUserInfo>>({
      url: `/admin-api/admin-users/${id}`,
      data: data
    })
  }

  /**
   * 删除管理员用户
   */
  static deleteAdminUser(id: number) {
    return request.del<Api.Http.BaseResponse<null>>({
      url: `/admin-api/admin-users/${id}`
    })
  }

  /**
   * 更新管理员用户状态
   */
  static updateAdminUserStatus(id: number, status: number) {
    return request.request<Api.Http.BaseResponse<null>>({
      url: `/admin-api/admin-users/${id}/status`,
      method: 'PATCH',
      data: { status },
      showErrorMessage: false, // 禁用自动错误显示，手动处理
      returnFullResponse: true // 返回完整响应，包括message
    })
  }

  /**
   * 更新管理员用户状态（返回完整响应包括消息）
   */
  static updateAdminUserStatusWithMessage(id: number, status: number) {
    return request.request<Api.Http.BaseResponse<null>>({
      url: `/admin-api/admin-users/${id}/status`,
      method: 'PATCH',
      data: { status },
      showErrorMessage: false, // 禁用自动错误显示，手动处理
      returnFullResponse: true // 返回完整响应
    })
  }
}

/**
 * 角色API服务
 */
export class RoleService {
  /**
   * 获取活跃角色选项
   */
  static getActiveRoles() {
    return request.get<RoleTypes.RoleOption[]>({
      url: '/admin-api/roles/options'
    })
  }

  /**
   * 获取角色列表
   */
  static getRoleList(params?: any) {
    return request.get<Api.Http.BaseResponse<any>>({
      url: '/admin-api/roles',
      params
    })
  }
}

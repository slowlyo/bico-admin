import request from '@/utils/http'

/**
 * 管理员用户API服务
 */
export class AdminUserService {
  /**
   * 获取管理员用户列表
   */
  static getAdminUserList(params: Api.AdminUser.AdminUserListParams) {
    return request.get<Api.Http.BaseResponse<Api.AdminUser.AdminUserListData>>({
      url: '/admin-api/admin-users',
      params
    })
  }

  /**
   * 根据ID获取管理员用户
   */
  static getAdminUserById(id: number) {
    return request.get<Api.Http.BaseResponse<Api.AdminUser.AdminUserInfo>>({
      url: `/admin-api/admin-users/${id}`
    })
  }

  /**
   * 创建管理员用户
   */
  static createAdminUser(data: Api.AdminUser.AdminUserCreateRequest) {
    return request.post<Api.Http.BaseResponse<Api.AdminUser.AdminUserInfo>>({
      url: '/admin-api/admin-users',
      data: data
    })
  }

  /**
   * 更新管理员用户
   */
  static updateAdminUser(id: number, data: Api.AdminUser.AdminUserUpdateRequest) {
    return request.put<Api.Http.BaseResponse<Api.AdminUser.AdminUserInfo>>({
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
    return request.get<Api.Role.RoleOption[]>({
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

import request from '@/utils/http'

/**
 * 管理员角色API服务
 */
export class AdminRoleService {
  /**
   * 获取角色列表
   */
  static getRoleList(params: any) {
    return request.get<Api.Http.BaseResponse<Api.Role.RoleListData>>({
      url: '/admin-api/roles',
      params
    })
  }

  /**
   * 根据ID获取角色
   */
  static getRoleById(id: number) {
    return request.get<Api.Http.BaseResponse<Api.Role.RoleInfo>>({
      url: `/admin-api/roles/${id}`
    })
  }

  /**
   * 创建角色
   */
  static createRole(data: Api.Role.RoleCreateRequest) {
    return request.post<Api.Http.BaseResponse<Api.Role.RoleInfo>>({
      url: '/admin-api/roles',
      data: data
    })
  }

  /**
   * 更新角色
   */
  static updateRole(id: number, data: Api.Role.RoleUpdateRequest) {
    return request.put<Api.Http.BaseResponse<Api.Role.RoleInfo>>({
      url: `/admin-api/roles/${id}`,
      data: data
    })
  }

  /**
   * 删除角色
   */
  static deleteRole(id: number) {
    return request.del<Api.Http.BaseResponse<null>>({
      url: `/admin-api/roles/${id}`
    })
  }

  /**
   * 更新角色状态
   */
  static updateRoleStatus(id: number, status: number) {
    return request.request<Api.Http.BaseResponse<null>>({
      url: `/admin-api/roles/${id}/status`,
      method: 'PATCH',
      data: { status },
      showErrorMessage: false,
      returnFullResponse: true
    })
  }

  /**
   * 获取权限树
   */
  static getPermissionTree(roleId?: number) {
    return request.get<Api.Role.PermissionTreeNode[]>({
      url: '/admin-api/roles/permissions',
      params: roleId ? { role_id: roleId } : undefined
    })
  }

  /**
   * 更新角色权限
   */
  static updateRolePermissions(id: number, permissions: string[]) {
    return request.put<Api.Http.BaseResponse<null>>({
      url: `/admin-api/roles/${id}/permissions`,
      data: { permissions }
    })
  }

  /**
   * 获取活跃角色选项
   */
  static getActiveRoles() {
    return request.get<Api.Role.RoleOption[]>({
      url: '/admin-api/roles/options'
    })
  }
}

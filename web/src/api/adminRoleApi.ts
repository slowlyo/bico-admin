import request from '@/utils/http'

// 角色相关类型定义
export namespace RoleTypes {
  /** 角色信息 */
  export interface RoleInfo {
    id: number
    name: string
    code: string
    description?: string
    status: number
    status_text: string
    user_count: number
    can_edit: boolean
    can_delete: boolean
    permissions?: RolePermission[]
    created_at: string
    updated_at: string
  }

  /** 角色权限 */
  export interface RolePermission {
    permission_code: string
    permission_name: string
    module: string
    level: number
  }

  /** 角色选项 */
  export interface RoleOption {
    id: number
    name: string
    code: string
    description?: string
  }

  /** 角色列表请求参数 */
  export interface RoleListParams {
    page: number
    page_size: number
    name?: string
    code?: string
    status?: number
    sort_by?: string
    sort_desc?: boolean
    [key: string]: unknown
  }

  /** 角色列表数据 */
  export interface RoleListData {
    list: RoleInfo[]
    total: number
    page: number
    page_size: number
  }

  /** 角色创建请求 */
  export interface RoleCreateRequest {
    name: string
    code: string
    description?: string
    status: number
    permissions: string[]
  }

  /** 角色更新请求 */
  export interface RoleUpdateRequest {
    name: string
    description?: string
    status: number
    permissions: string[]
  }

  /** 权限树节点 */
  export interface PermissionTreeNode {
    key: string
    title: string
    type: string
    selected: boolean
    children?: PermissionTreeNode[]
  }
}

/**
 * 管理员角色API服务
 */
export class AdminRoleService {
  /**
   * 获取角色列表
   */
  static getRoleList(params: any) {
    return request.get<Api.Http.BaseResponse<RoleTypes.RoleListData>>({
      url: '/admin-api/roles',
      params
    })
  }

  /**
   * 根据ID获取角色
   */
  static getRoleById(id: number) {
    return request.get<Api.Http.BaseResponse<RoleTypes.RoleInfo>>({
      url: `/admin-api/roles/${id}`
    })
  }

  /**
   * 创建角色
   */
  static createRole(data: RoleTypes.RoleCreateRequest) {
    return request.post<Api.Http.BaseResponse<RoleTypes.RoleInfo>>({
      url: '/admin-api/roles',
      data: data
    })
  }

  /**
   * 更新角色
   */
  static updateRole(id: number, data: RoleTypes.RoleUpdateRequest) {
    return request.put<Api.Http.BaseResponse<RoleTypes.RoleInfo>>({
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
    return request.get<RoleTypes.PermissionTreeNode[]>({
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
    return request.get<RoleTypes.RoleOption[]>({
      url: '/admin-api/roles/options'
    })
  }
}

import request from '@/utils/http'
import { AppRouteRecord } from '@/types/router'

// --- 用户管理 ---

/** 获取用户列表 */
export function fetchGetUserList(params: Api.SystemManage.UserSearchParams) {
  return request.get<Api.SystemManage.UserList>({
    url: '/admin-api/admin-users',
    params
  })
}

/** 获取用户详情 */
export function fetchGetUserInfo(id: number) {
  return request.get<Api.SystemManage.UserListItem>({
    url: `/admin-api/admin-users/${id}`
  })
}

/** 创建用户 */
export function fetchCreateUser(data: Api.SystemManage.UserParams) {
  return request.post({
    url: '/admin-api/admin-users',
    data
  })
}

/** 更新用户 */
export function fetchUpdateUser(id: number, data: Api.SystemManage.UserParams) {
  return request.put({
    url: `/admin-api/admin-users/${id}`,
    data
  })
}

/** 删除用户 */
export function fetchDeleteUser(id: number) {
  return request.del({
    url: `/admin-api/admin-users/${id}`
  })
}

/** 上传文件 */
export function fetchUploadFile(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post<{ url: string }>({
    url: '/admin-api/upload',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// --- 角色管理 ---

/** 获取角色列表 */
export function fetchGetRoleList(params: Api.SystemManage.RoleSearchParams) {
  return request.get<Api.SystemManage.RoleList>({
    url: '/admin-api/admin-roles',
    params
  })
}

/** 获取所有角色（下拉选择） */
export function fetchGetAllRoles() {
  return request.get<Api.SystemManage.RoleListItem[]>({
    url: '/admin-api/admin-roles/all'
  })
}

/** 创建角色 */
export function fetchCreateRole(data: Api.SystemManage.RoleParams) {
  return request.post({
    url: '/admin-api/admin-roles',
    data
  })
}

/** 更新角色 */
export function fetchUpdateRole(id: number, data: Api.SystemManage.RoleParams) {
  return request.put({
    url: `/admin-api/admin-roles/${id}`,
    data
  })
}

/** 删除角色 */
export function fetchDeleteRole(id: number) {
  return request.del({
    url: `/admin-api/admin-roles/${id}`
  })
}

/** 获取角色权限（已分配的权限标识） */
export function fetchGetRolePermissions(id: number) {
  return request.get<{ permissions: string[] }>({
    url: `/admin-api/admin-roles/${id}/permissions`
  })
}

/** 获取所有权限树 */
export function fetchGetAllPermissions() {
  return request.get<Api.SystemManage.Permission[]>({
    url: '/admin-api/admin-roles/permissions'
  })
}

/** 更新角色权限 */
export function fetchUpdateRolePermissions(id: number, permissions: string[]) {
  return request.put({
    url: `/admin-api/admin-roles/${id}/permissions`,
    data: { permissions }
  })
}

// --- 菜单管理 ---

/** 获取菜单列表 */
export function fetchGetMenuList() {
  return request.get<AppRouteRecord[]>({
    url: '/admin-api/v3/system/menus/simple'
  })
}

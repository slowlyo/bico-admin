import { request } from '@umijs/max';

// 角色权限响应
export interface RolePermissionResponse {
  permission_code: string;
  permission_name: string;
  module: string;
  level: number;
}

// 角色数据类型
export interface Role {
  id: number;
  code: string;
  name: string;
  description?: string;
  status: number;
  status_text: string;
  permissions: RolePermissionResponse[];
  user_count: number;    // 拥有该角色的用户数量
  can_edit: boolean;     // 是否可编辑
  can_delete: boolean;   // 是否可删除
  created_at: string;
  updated_at: string;
}

// 角色创建请求
export interface RoleCreateRequest {
  code: string;
  name: string;
  description?: string;
  status: number;
  permissions?: string[];  // 权限代码列表
}

// 角色更新请求
export interface RoleUpdateRequest {
  name: string;
  description?: string;
  status: number;
  permissions?: string[];  // 权限代码列表
}

// 权限树节点
export interface PermissionTreeNode {
  module: string;
  name: string;
  permissions: PermissionTreeItem[];
}

// 权限树项目
export interface PermissionTreeItem {
  code: string;
  name: string;
  level: number;
  level_text: string;
  buttons: string[];
  apis: string[];
  selected: boolean;
}

// 角色权限更新请求
export interface RolePermissionUpdateRequest {
  permissions: string[];
}

// 角色列表请求参数
export interface RoleListRequest {
  page?: number;
  page_size?: number;
  code?: string;
  name?: string;
  status?: number;
}

// 分页响应
export interface PageResponse<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

// API响应
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

/**
 * 获取角色列表
 */
export async function getRoleList(params: RoleListRequest): Promise<ApiResponse<PageResponse<Role>>> {
  return request('/admin/roles', {
    method: 'GET',
    params,
  });
}

/**
 * 获取角色详情
 */
export async function getRoleById(id: number): Promise<ApiResponse<Role>> {
  return request(`/admin/roles/${id}`, {
    method: 'GET',
  });
}

/**
 * 创建角色
 */
export async function createRole(data: RoleCreateRequest): Promise<ApiResponse<Role>> {
  return request('/admin/roles', {
    method: 'POST',
    data,
  });
}

/**
 * 更新角色
 */
export async function updateRole(id: number, data: RoleUpdateRequest): Promise<ApiResponse<Role>> {
  return request(`/admin/roles/${id}`, {
    method: 'PUT',
    data,
  });
}

/**
 * 删除角色
 */
export async function deleteRole(id: number): Promise<ApiResponse<null>> {
  return request(`/admin/roles/${id}`, {
    method: 'DELETE',
  });
}

/**
 * 更新角色状态
 */
export async function updateRoleStatus(id: number, status: number): Promise<ApiResponse<null>> {
  return request(`/admin/roles/${id}/status`, {
    method: 'PATCH',
    data: { status },
  });
}

/**
 * 获取权限树
 */
export async function getPermissionTree(roleId?: number): Promise<ApiResponse<PermissionTreeNode[]>> {
  return request('/admin/roles/permissions', {
    method: 'GET',
    params: roleId ? { role_id: roleId } : {},
  });
}

/**
 * 更新角色权限
 */
export async function updateRolePermissions(id: number, data: RolePermissionUpdateRequest): Promise<ApiResponse<null>> {
  return request(`/admin/roles/${id}/permissions`, {
    method: 'PUT',
    data,
  });
}

// 角色分配请求
export interface RoleAssignRequest {
  user_id: number;
  role_ids: number[];
}

// 用户角色响应
export interface UserRoleResponse {
  user_id: number;
  username: string;
  name: string;
  roles: Role[];
  created_at: string;
}

/**
 * 分配角色给用户
 */
export async function assignRolesToUser(data: RoleAssignRequest): Promise<ApiResponse<null>> {
  return request('/admin/roles/assign', {
    method: 'POST',
    data,
  });
}

/**
 * 获取用户角色
 */
export async function getUserRoles(userId: number): Promise<ApiResponse<UserRoleResponse>> {
  return request(`/admin/roles/user/${userId}`, {
    method: 'GET',
  });
}

/**
 * 获取角色统计
 */
export async function getRoleStats(): Promise<ApiResponse<any>> {
  return request('/admin/roles/stats', {
    method: 'GET',
  });
}

import { request } from '@umijs/max';

// 角色数据类型
export interface Role {
  id: number;
  code: string;
  name: string;
  description?: string;
  status: number;
  status_text: string;
  created_at: string;
  updated_at: string;
}

// 角色创建请求
export interface RoleCreateRequest {
  code: string;
  name: string;
  description?: string;
  status?: number;
}

// 角色更新请求
export interface RoleUpdateRequest {
  code?: string;
  name?: string;
  description?: string;
  status?: number;
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

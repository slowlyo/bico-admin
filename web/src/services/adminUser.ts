import { request } from '@umijs/max';

// 管理员用户类型定义
export interface AdminUser {
  id: number;
  username: string;
  name: string;
  avatar?: string;
  email?: string;
  phone?: string;
  status: number;
  status_text: string;
  last_login_at?: string;
  last_login_ip?: string;
  login_count: number;
  remark?: string;
  created_at: string;
  updated_at: string;
}

export interface AdminUserCreateRequest {
  username: string;
  password: string;
  name: string;
  avatar?: string;
  enabled: boolean;
}

export interface AdminUserUpdateRequest {
  username: string;
  password?: string;
  name: string;
  avatar?: string;
  enabled: boolean;
}

export interface AdminUserListParams {
  page?: number;
  page_size?: number;
}

export interface AdminUserListResponse {
  list: AdminUser[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

// 获取管理员用户列表
export async function getAdminUserList(
  params: AdminUserListParams = {},
  options?: { [key: string]: any },
) {
  return request<ApiResponse<AdminUserListResponse>>('/admin/admin-users', {
    method: 'GET',
    params: {
      page: 1,
      page_size: 10,
      ...params,
    },
    ...(options || {}),
  });
}

// 根据ID获取管理员用户
export async function getAdminUserById(
  id: number,
  options?: { [key: string]: any },
) {
  return request<ApiResponse<AdminUser>>(`/admin/admin-users/${id}`, {
    method: 'GET',
    ...(options || {}),
  });
}

// 创建管理员用户
export async function createAdminUser(
  data: AdminUserCreateRequest,
  options?: { [key: string]: any },
) {
  return request<ApiResponse<AdminUser>>('/admin/admin-users', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}

// 更新管理员用户
export async function updateAdminUser(
  id: number,
  data: AdminUserUpdateRequest,
  options?: { [key: string]: any },
) {
  return request<ApiResponse<AdminUser>>(`/admin/admin-users/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}

// 删除管理员用户
export async function deleteAdminUser(
  id: number,
  options?: { [key: string]: any },
) {
  return request<ApiResponse<null>>(`/admin/admin-users/${id}`, {
    method: 'DELETE',
    ...(options || {}),
  });
}

// 更新管理员用户状态
export async function updateAdminUserStatus(
  id: number,
  status: number,
  options?: { [key: string]: any },
) {
  return request<ApiResponse<null>>(`/admin/admin-users/${id}/status`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    data: { status },
    ...(options || {}),
  });
}

import { request } from '@umijs/max';

export interface AdminUser {
  id: number;
  username: string;
  name: string;
  avatar: string;
  enabled: boolean;
  roles?: AdminRole[];
  created_at: string;
  updated_at: string;
}

export interface AdminRole {
  id: number;
  name: string;
  code: string;
  description: string;
  enabled: boolean;
}

export interface AdminUserListParams {
  page?: number;
  pageSize?: number;
  username?: string;
  name?: string;
  enabled?: boolean;
  sortField?: string;
  sortOrder?: string;
}

export interface AdminUserCreateParams {
  username: string;
  password: string;
  name?: string;
  avatar?: string;
  enabled?: boolean;
  roleIds?: number[];
}

export interface AdminUserUpdateParams {
  name?: string;
  avatar?: string;
  enabled?: boolean;
  roleIds?: number[];
}

export async function getAdminUserList(params: AdminUserListParams) {
  return request<API.Response<AdminUser[]> & { total?: number }>('/admin-api/admin-users', {
    method: 'GET',
    params,
  });
}

export async function getAdminUser(id: number) {
  return request<API.Response<AdminUser>>(`/admin-api/admin-users/${id}`, {
    method: 'GET',
  });
}

export async function createAdminUser(data: AdminUserCreateParams) {
  return request<API.Response<AdminUser>>('/admin-api/admin-users', {
    method: 'POST',
    data,
  });
}

export async function updateAdminUser(id: number, data: AdminUserUpdateParams) {
  return request<API.Response<AdminUser>>(`/admin-api/admin-users/${id}`, {
    method: 'PUT',
    data,
  });
}

export async function deleteAdminUser(id: number) {
  return request<API.Response<null>>(`/admin-api/admin-users/${id}`, {
    method: 'DELETE',
  });
}

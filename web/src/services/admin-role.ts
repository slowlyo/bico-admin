import { request } from '@umijs/max';

export interface AdminRole {
  id: number;
  name: string;
  code: string;
  description: string;
  enabled: boolean;
  permissions?: string[];
  created_at: string;
  updated_at: string;
}

export interface Permission {
  key: string;
  label: string;
  children?: Permission[];
}

export interface AdminRoleListParams {
  page?: number;
  pageSize?: number;
  name?: string;
  code?: string;
  enabled?: boolean;
  sortField?: string;
  sortOrder?: string;
}

export interface AdminRoleCreateParams {
  name: string;
  code: string;
  description?: string;
  enabled?: boolean;
  permissions?: string[];
}

export interface AdminRoleUpdateParams {
  name?: string;
  description?: string;
  enabled?: boolean;
}

export interface UpdatePermissionsParams {
  permissions: string[];
}

export async function getAdminRoleList(params: AdminRoleListParams) {
  return request<API.Response<AdminRole[]> & { total?: number }>('/admin-api/admin-roles', {
    method: 'GET',
    params,
  });
}

export async function getAllAdminRoles() {
  return request<API.Response<AdminRole[]>>('/admin-api/admin-roles/all', {
    method: 'GET',
  });
}

export async function getAdminRole(id: number) {
  return request<API.Response<AdminRole>>(`/admin-api/admin-roles/${id}`, {
    method: 'GET',
  });
}

export async function createAdminRole(data: AdminRoleCreateParams) {
  return request<API.Response<AdminRole>>('/admin-api/admin-roles', {
    method: 'POST',
    data,
  });
}

export async function updateAdminRole(id: number, data: AdminRoleUpdateParams) {
  return request<API.Response<AdminRole>>(`/admin-api/admin-roles/${id}`, {
    method: 'PUT',
    data,
  });
}

export async function deleteAdminRole(id: number) {
  return request<API.Response<null>>(`/admin-api/admin-roles/${id}`, {
    method: 'DELETE',
  });
}

export async function updateRolePermissions(id: number, data: UpdatePermissionsParams) {
  return request<API.Response<null>>(`/admin-api/admin-roles/${id}/permissions`, {
    method: 'PUT',
    data,
  });
}

export async function getAllPermissions() {
  return request<API.Response<Permission[]>>('/admin-api/admin-roles/permissions', {
    method: 'GET',
  });
}

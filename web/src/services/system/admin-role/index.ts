import { request } from '@umijs/max';
import { buildApiUrl } from '../../config';
import type {
  AdminRole,
  AdminRoleListParams,
  AdminRoleCreateParams,
  AdminRoleUpdateParams,
  UpdatePermissionsParams,
  Permission,
} from './types';

/**
 * 获取角色列表
 */
export async function getAdminRoleList(params: AdminRoleListParams) {
  return request<API.Response<AdminRole[]> & { total?: number }>(
    buildApiUrl('/admin-roles'),
    {
      method: 'GET',
      params,
    }
  );
}

/**
 * 获取所有角色（不分页）
 */
export async function getAllAdminRoles() {
  return request<API.Response<AdminRole[]>>(buildApiUrl('/admin-roles/all'), {
    method: 'GET',
  });
}

/**
 * 获取角色详情
 */
export async function getAdminRole(id: number) {
  return request<API.Response<AdminRole>>(buildApiUrl(`/admin-roles/${id}`), {
    method: 'GET',
  });
}

/**
 * 创建角色
 */
export async function createAdminRole(data: AdminRoleCreateParams) {
  return request<API.Response<AdminRole>>(buildApiUrl('/admin-roles'), {
    method: 'POST',
    data,
  });
}

/**
 * 更新角色
 */
export async function updateAdminRole(id: number, data: AdminRoleUpdateParams) {
  return request<API.Response<AdminRole>>(buildApiUrl(`/admin-roles/${id}`), {
    method: 'PUT',
    data,
  });
}

/**
 * 删除角色
 */
export async function deleteAdminRole(id: number) {
  return request<API.Response<null>>(buildApiUrl(`/admin-roles/${id}`), {
    method: 'DELETE',
  });
}

/**
 * 更新角色权限
 */
export async function updateRolePermissions(
  id: number,
  data: UpdatePermissionsParams
) {
  return request<API.Response<null>>(
    buildApiUrl(`/admin-roles/${id}/permissions`),
    {
      method: 'PUT',
      data,
    }
  );
}

/**
 * 获取所有权限
 */
export async function getAllPermissions() {
  return request<API.Response<Permission[]>>(
    buildApiUrl('/admin-roles/permissions'),
    {
      method: 'GET',
    }
  );
}

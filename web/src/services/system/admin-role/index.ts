/**
 * 角色管理 - 特殊接口
 * 标准 CRUD 使用 createCrudService('/admin-roles')
 */
import { request } from '@umijs/max';
import { buildApiUrl } from '../../config';
import type { AdminRole, Permission, UpdatePermissionsParams } from './types';

/** 获取所有角色（不分页，用于下拉选择） */
export async function getAllAdminRoles() {
  return request<API.Response<AdminRole[]>>(buildApiUrl('/admin-roles/all'), {
    method: 'GET',
  });
}

/** 获取所有权限树 */
export async function getAllPermissions() {
  return request<API.Response<Permission[]>>(buildApiUrl('/admin-roles/permissions'), {
    method: 'GET',
  });
}

/** 更新角色权限 */
export async function updateRolePermissions(id: number, data: UpdatePermissionsParams) {
  return request<API.Response<null>>(buildApiUrl(`/admin-roles/${id}/permissions`), {
    method: 'PUT',
    data,
  });
}

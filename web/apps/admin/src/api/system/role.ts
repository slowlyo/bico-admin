import { requestClient } from '#/api/request';

import type { PageResult } from './user';

export interface AdminRole {
  created_at: string;
  description: string;
  enabled: boolean;
  id: number;
  name: string;
  permissions?: string[];
  system: boolean;
}

export interface AdminRoleQuery {
  enabled?: boolean;
  name?: string;
  page: number;
  pageSize: number;
  sortField?: string;
  sortOrder?: string;
}

export interface AdminRolePayload {
  description: string;
  enabled: boolean;
  name: string;
}

export interface PermissionNode {
  children?: PermissionNode[];
  key: string;
  label: string;
}

// 角色分页列表。系统保留标记和权限码由列表后处理填上。
export function listAdminRolesApi(params: AdminRoleQuery) {
  return requestClient.get<PageResult<AdminRole>>('/admin-roles', { params });
}

// 下拉只用启用角色，不分页。
export function listAllAdminRolesApi() {
  return requestClient.get<AdminRole[]>('/admin-roles/all');
}

// 新建角色。权限不在这个表单里配。
export function createAdminRoleApi(data: AdminRolePayload) {
  return requestClient.post('/admin-roles', data);
}

// 更新角色名称、描述和启用状态。
export function updateAdminRoleApi(id: number, data: AdminRolePayload) {
  return requestClient.put(`/admin-roles/${id}`, data);
}

// 删除角色。系统保留角色由页面先拦住。
export function deleteAdminRoleApi(id: number) {
  return requestClient.delete(`/admin-roles/${id}`);
}

// 批量删除。调用方需要先去掉系统保留角色。
export function deleteAdminRolesApi(ids: number[]) {
  return requestClient.delete('/admin-roles/batch', { data: { ids } });
}

// 完整权限树，配置抽屉和保存时补父节点都用它。
export function listPermissionsApi() {
  return requestClient.get<PermissionNode[]>('/admin-roles/permissions');
}

// 保存角色权限。需要带上父节点，否则菜单权限会丢。
export function updateRolePermissionsApi(id: number, permissions: string[]) {
  return requestClient.put(`/admin-roles/${id}/permissions`, { permissions });
}

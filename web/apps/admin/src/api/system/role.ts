import type { definitions, paths } from '#/api/generated/schema';

import type { Defined, PageQuery } from '#/api/crud';

import { createCrudApi } from '#/api/crud';
import { requestClient } from '#/api/request';

// 以下类型来自后端 Swagger（api/generated/schema.d.ts），后端改字段后执行 make web-types 同步。
export type AdminRole = Defined<definitions['model.AdminRole']>;

export type AdminRoleQuery = NonNullable<
  paths['/admin-roles']['get']['parameters']['query']
> &
  PageQuery;

export type AdminRoleCreatePayload = definitions['handler.createRoleReq'];

export type AdminRoleUpdatePayload = definitions['handler.updateRoleReq'];

export interface PermissionNode {
  children?: PermissionNode[];
  key: string;
  label: string;
}

// 角色走标准 CRUD 接口。系统保留标记和权限码由列表后处理填上；权限不在创建表单里配。
export const adminRoleApi = createCrudApi<
  AdminRole,
  AdminRoleQuery,
  AdminRoleCreatePayload,
  AdminRoleUpdatePayload
>('/admin-roles');

// 下拉只用启用角色，不分页。
export function listAllAdminRolesApi() {
  return requestClient.get<AdminRole[]>('/admin-roles/all');
}

// 完整权限树，配置抽屉和保存时补父节点都用它。
export function listPermissionsApi() {
  return requestClient.get<PermissionNode[]>('/admin-roles/permissions');
}

// 保存角色权限。需要带上父节点，否则菜单权限会丢。
export function updateRolePermissionsApi(id: number, permissions: string[]) {
  return requestClient.put(`/admin-roles/${id}/permissions`, { permissions });
}

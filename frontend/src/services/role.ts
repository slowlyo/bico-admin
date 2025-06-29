import { request } from '@umijs/max';

export interface RoleItem {
  id: number;
  name: string;
  code: string;
  description?: string;
  status: number;
  created_at: string;
  updated_at: string;
  permissions?: PermissionItem[];
}

export interface PermissionItem {
  id: number;
  name: string;
  code: string;
  type: number;
  resource?: string;
  action?: string;
  description?: string;
  parent_id?: number;
  sort: number;
  status: number;
  created_at: string;
  updated_at: string;
  parent?: PermissionItem;
  children?: PermissionItem[];
}

export interface CreateRoleParams {
  name: string;
  code: string;
  description?: string;
  permission_ids?: number[];
}

export interface UpdateRoleParams {
  name?: string;
  code?: string;
  description?: string;
  status?: number;
  permission_ids?: number[];
}

export interface RoleListParams {
  current?: number;
  pageSize?: number;
  name?: string;
  code?: string;
  status?: number;
}

export interface RoleListResult {
  data: RoleItem[];
  total: number;
  success: boolean;
}

/**
 * 获取角色列表
 */
export async function getRoleList(
  params: RoleListParams,
): Promise<RoleListResult> {
  try {
    const response = await request('/admin/roles', {
      method: 'GET',
      params: {
        page: params.current || 1,
        page_size: params.pageSize || 10,
        search: params.name || params.code || '',
        status: params.status,
      },
    });

    // 后端现在直接返回Ant Design Pro标准格式
    return {
      data: response?.data || [],
      total: response?.total || 0,
      success: response?.success !== false,
    };
  } catch (error) {
    console.error('获取角色列表失败:', error);
    return {
      data: [],
      total: 0,
      success: false,
    };
  }
}

/**
 * 获取单个角色
 */
export async function getRole(id: number): Promise<{
  code: number;
  message: string;
  data: RoleItem;
}> {
  return request(`/admin/roles/${id}`, {
    method: 'GET',
  });
}

/**
 * 创建角色
 */
export async function createRole(params: CreateRoleParams): Promise<{
  code: number;
  message: string;
  data: RoleItem;
}> {
  return request('/admin/roles', {
    method: 'POST',
    data: params,
  });
}

/**
 * 更新角色
 */
export async function updateRole(
  id: number,
  params: UpdateRoleParams,
): Promise<{
  code: number;
  message: string;
  data: RoleItem;
}> {
  return request(`/admin/roles/${id}`, {
    method: 'PUT',
    data: params,
  });
}

/**
 * 删除角色
 */
export async function deleteRole(id: number): Promise<{
  code: number;
  message: string;
}> {
  return request(`/admin/roles/${id}`, {
    method: 'DELETE',
  });
}

/**
 * 批量删除角色
 */
export async function batchDeleteRoles(ids: number[]): Promise<{
  code: number;
  message: string;
}> {
  return request('/admin/roles/batch', {
    method: 'DELETE',
    data: { ids },
  });
}

/**
 * 更新角色状态
 */
export async function updateRoleStatus(
  id: number,
  status: number,
): Promise<{
  code: number;
  message: string;
  data: RoleItem;
}> {
  return request(`/admin/roles/${id}/status`, {
    method: 'PUT',
    data: { status },
  });
}

/**
 * 获取角色权限
 */
export async function getRolePermissions(id: number): Promise<{
  code: number;
  message: string;
  data: string[]; // 后端返回权限代码字符串数组
}> {
  try {
    const response = await request(`/admin/roles/${id}/permissions`, {
      method: 'GET',
    });
    return response;
  } catch (error) {
    console.error('获取角色权限失败:', error);
    // 返回空的权限列表
    return {
      code: 200,
      message: 'Success',
      data: [],
    };
  }
}

/**
 * 分配角色权限
 */
export async function assignRolePermissions(
  id: number,
  permissionCodes: string[],
): Promise<{
  code: number;
  message: string;
  data: RoleItem;
}> {
  return request(`/admin/roles/${id}/permissions`, {
    method: 'PUT',
    data: { permission_codes: permissionCodes },
  });
}

/**
 * 获取权限列表
 */
export async function getPermissionList(params?: {
  current?: number;
  pageSize?: number;
  type?: number;
}): Promise<{
  data: PermissionItem[];
  total: number;
  success: boolean;
}> {
  try {
    const response = await request('/admin/permissions', {
      method: 'GET',
      params: {
        page: params?.current || 1,
        page_size: params?.pageSize || 100, // 权限列表通常不分页
        type: params?.type,
      },
    });

    return {
      data: response?.data || [],
      total: response?.total || 0,
      success: response?.success !== false,
    };
  } catch (error) {
    console.error('获取权限列表失败:', error);
    return {
      data: [],
      total: 0,
      success: false,
    };
  }
}

/**
 * 获取权限树
 */
export async function getPermissionTree(): Promise<{
  code: number;
  message: string;
  data: PermissionItem[];
}> {
  try {
    const response = await request('/admin/permissions/tree', {
      method: 'GET',
    });
    return response;
  } catch (error) {
    console.error('获取权限树失败:', error);
    // 返回空的权限树结构
    return {
      code: 200,
      message: 'Success',
      data: [],
    };
  }
}

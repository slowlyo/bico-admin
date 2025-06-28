import { request } from '@umijs/max';

export interface UserItem {
  id: number;
  username: string;
  email: string;
  nickname?: string;
  phone?: string;
  status: number;
  created_at: string;
  updated_at: string;
  roles?: any[];
}

export interface CreateUserParams {
  username: string;
  email: string;
  password: string;
  nickname?: string;
  phone?: string;
  status?: number;
  role_ids?: number[];
}

export interface UpdateUserParams {
  username?: string;
  email?: string;
  nickname?: string;
  phone?: string;
  status?: number;
  role_ids?: number[];
}

export interface UserListParams {
  current?: number;
  pageSize?: number;
  username?: string;
  email?: string;
  status?: number;
}

export interface UserListResult {
  data: UserItem[];
  total: number;
  success: boolean;
}

/**
 * 获取用户列表
 */
export async function getUserList(
  params: UserListParams,
): Promise<UserListResult> {
  const response = await request('/admin/users', {
    method: 'GET',
    params: {
      page: params.current || 1,
      size: params.pageSize || 10,
      search: params.username || params.email || '',
      status: params.status,
    },
  });

  return {
    data: response.data || [],
    total: response.total || 0,
    success: response.code === 200,
  };
}

/**
 * 获取单个用户
 */
export async function getUser(id: number): Promise<{
  code: number;
  message: string;
  data: UserItem;
}> {
  return request(`/admin/users/${id}`, {
    method: 'GET',
  });
}

/**
 * 创建用户
 */
export async function createUser(params: CreateUserParams): Promise<{
  code: number;
  message: string;
  data: UserItem;
}> {
  return request('/admin/users', {
    method: 'POST',
    data: params,
  });
}

/**
 * 更新用户
 */
export async function updateUser(
  id: number,
  params: UpdateUserParams,
): Promise<{
  code: number;
  message: string;
  data: UserItem;
}> {
  return request(`/admin/users/${id}`, {
    method: 'PUT',
    data: params,
  });
}

/**
 * 删除用户
 */
export async function deleteUser(id: number): Promise<{
  code: number;
  message: string;
}> {
  return request(`/admin/users/${id}`, {
    method: 'DELETE',
  });
}

/**
 * 批量删除用户
 */
export async function batchDeleteUsers(ids: number[]): Promise<{
  code: number;
  message: string;
}> {
  return request('/admin/users/batch', {
    method: 'DELETE',
    data: { ids },
  });
}

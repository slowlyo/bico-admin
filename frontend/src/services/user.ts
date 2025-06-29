import { request } from '@umijs/max';

export interface UserItem {
  id: number;
  username: string;
  email: string;
  nickname?: string;
  phone?: string;
  status: number;
  role?: string;
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
  role?: string;
  role_ids?: number[];
}

export interface UpdateUserParams {
  username?: string;
  email?: string;
  nickname?: string;
  phone?: string;
  status?: number;
  role?: string;
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
  try {
    const response = await request('/admin/users', {
      method: 'GET',
      params: {
        page: params.current || 1,
        page_size: params.pageSize || 10,
        search: params.username || params.email || '',
        status: params.status,
      },
    });

    // 后端现在直接返回Ant Design Pro标准格式
    // { data: [...], total: 100, success: true, current: 1, pageSize: 10 }
    return {
      data: response?.data || [],
      total: response?.total || 0,
      success: response?.success !== false,
    };
  } catch (error) {
    console.error('获取用户列表失败:', error);
    return {
      data: [],
      total: 0,
      success: false,
    };
  }
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

/**
 * 更新用户状态
 */
export async function updateUserStatus(
  id: number,
  status: number,
): Promise<{
  code: number;
  message: string;
  data: UserItem;
}> {
  return request(`/admin/users/${id}/status`, {
    method: 'PUT',
    data: { status },
  });
}

/**
 * 重置用户密码
 */
export async function resetUserPassword(
  id: number,
  newPassword: string,
): Promise<{
  code: number;
  message: string;
}> {
  return request(`/admin/users/${id}/reset-password`, {
    method: 'PUT',
    data: { new_password: newPassword },
  });
}

/**
 * 修改用户密码
 */
export async function changeUserPassword(
  id: number,
  oldPassword: string,
  newPassword: string,
): Promise<{
  code: number;
  message: string;
}> {
  return request(`/admin/users/${id}/password`, {
    method: 'PUT',
    data: {
      old_password: oldPassword,
      new_password: newPassword,
    },
  });
}

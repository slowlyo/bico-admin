import { request } from '@umijs/max';
import { buildApiUrl } from '../../config';
import type {
  AdminUser,
  AdminUserListParams,
  AdminUserCreateParams,
  AdminUserUpdateParams,
} from './types';

/**
 * 获取管理员列表
 */
export async function getAdminUserList(params: AdminUserListParams) {
  return request<API.Response<AdminUser[]> & { total?: number }>(
    buildApiUrl('/admin-users'),
    {
      method: 'GET',
      params,
    }
  );
}

/**
 * 获取管理员详情
 */
export async function getAdminUser(id: number) {
  return request<API.Response<AdminUser>>(buildApiUrl(`/admin-users/${id}`), {
    method: 'GET',
  });
}

/**
 * 创建管理员
 */
export async function createAdminUser(data: AdminUserCreateParams) {
  return request<API.Response<AdminUser>>(buildApiUrl('/admin-users'), {
    method: 'POST',
    data,
  });
}

/**
 * 更新管理员
 */
export async function updateAdminUser(id: number, data: AdminUserUpdateParams) {
  return request<API.Response<AdminUser>>(buildApiUrl(`/admin-users/${id}`), {
    method: 'PUT',
    data,
  });
}

/**
 * 删除管理员
 */
export async function deleteAdminUser(id: number) {
  return request<API.Response<null>>(buildApiUrl(`/admin-users/${id}`), {
    method: 'DELETE',
  });
}

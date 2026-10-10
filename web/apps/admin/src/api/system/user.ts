import { requestClient } from '#/api/request';

export interface AdminUserRole {
  id: number;
  name: string;
}

export interface AdminUser {
  avatar: string;
  created_at: string;
  enabled: boolean;
  id: number;
  name: string;
  roles?: AdminUserRole[];
  username: string;
}

export interface AdminUserQuery {
  enabled?: boolean;
  name?: string;
  page: number;
  pageSize: number;
  role_ids?: string;
  sortField?: string;
  sortOrder?: string;
  username?: string;
}

export interface AdminUserPayload {
  avatar?: string;
  enabled: boolean;
  name: string;
  password?: string;
  role_ids: number[];
  username: string;
}

export interface PageResult<T> {
  list: T[];
  total: number;
}

// 用户列表走标准分页。角色筛选是逗号分隔的 id，和后台 form 绑定一致。
export function listAdminUsersApi(params: AdminUserQuery) {
  return requestClient.get<PageResult<AdminUser>>('/admin-users', { params });
}

// 新建用户。密码只在创建时必填，角色在同一请求里写入。
export function createAdminUserApi(data: AdminUserPayload) {
  return requestClient.post('/admin-users', data);
}

// 更新用户。密码留空时不传，避免空字符串触发最少 8 位校验。
export function updateAdminUserApi(id: number, data: AdminUserPayload) {
  return requestClient.put(`/admin-users/${id}`, data);
}

// 删除单个用户。
export function deleteAdminUserApi(id: number) {
  return requestClient.delete(`/admin-users/${id}`);
}

// 批量删除。后台按 ids 数组处理，系统用户没有单独豁免。
export function deleteAdminUsersApi(ids: number[]) {
  return requestClient.delete('/admin-users/batch', { data: { ids } });
}

// 只拿上传后的地址，不改当前登录用户的头像。
export function uploadUserAvatarApi(file: File) {
  return requestClient.upload<{ url: string }>('/auth/avatar', {
    avatar: file,
    file,
  });
}

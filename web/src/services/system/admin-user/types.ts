/**
 * 管理员管理类型定义
 */

export interface AdminRole {
  id: number;
  name: string;
  code: string;
  description: string;
  enabled: boolean;
}

export interface AdminUser {
  id: number;
  username: string;
  name: string;
  avatar: string;
  enabled: boolean;
  roles?: AdminRole[];
  created_at: string;
  updated_at: string;
}

export interface AdminUserListParams {
  page?: number;
  pageSize?: number;
  username?: string;
  name?: string;
  enabled?: boolean;
  sortField?: string;
  sortOrder?: string;
}

export interface AdminUserCreateParams {
  username: string;
  password: string;
  name?: string;
  avatar?: string;
  enabled?: boolean;
  roleIds?: number[];
}

export interface AdminUserUpdateParams {
  name?: string;
  avatar?: string;
  enabled?: boolean;
  roleIds?: number[];
}

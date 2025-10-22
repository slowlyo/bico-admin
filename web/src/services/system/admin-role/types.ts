/**
 * 角色管理类型定义
 */

export interface Permission {
  key: string;
  label: string;
  children?: Permission[];
}

export interface AdminRole {
  id: number;
  name: string;
  code: string;
  description: string;
  enabled: boolean;
  permissions?: string[];
  created_at: string;
  updated_at: string;
}

export interface AdminRoleListParams {
  page?: number;
  pageSize?: number;
  name?: string;
  code?: string;
  enabled?: boolean;
  sortField?: string;
  sortOrder?: string;
}

export interface AdminRoleCreateParams {
  name: string;
  code: string;
  description?: string;
  enabled?: boolean;
  permissions?: string[];
}

export interface AdminRoleUpdateParams {
  name?: string;
  description?: string;
  enabled?: boolean;
}

export interface UpdatePermissionsParams {
  permissions: string[];
}

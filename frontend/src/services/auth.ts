import { request } from '@umijs/max';

export interface LoginParams {
  username: string;
  password: string;
}

export interface LoginResult {
  code: number;
  message: string;
  data: {
    user: {
      id: number;
      username: string;
      email: string;
      nickname?: string;
      phone?: string;
      status: number;
      roles?: any[];
    };
    token: string;
  };
}

export interface UserProfile {
  id: number;
  username: string;
  email: string;
  nickname?: string;
  phone?: string;
  status: number;
  roles?: any[];
  created_at: string;
  updated_at: string;
}

export interface UpdateProfileParams {
  email?: string;
  nickname?: string;
  phone?: string;
}

export interface ChangePasswordParams {
  old_password: string;
  new_password: string;
  confirm_password: string;
}

/**
 * 用户登录
 */
export async function login(params: LoginParams, options?: any): Promise<LoginResult> {
  return request('/admin/auth/login', {
    method: 'POST',
    data: params,
    ...options,
  });
}

/**
 * 用户登出
 */
export async function logout(): Promise<{ code: number; message: string }> {
  return request('/admin/auth/logout', {
    method: 'POST',
  });
}

/**
 * 获取用户资料
 */
export async function getUserProfile(): Promise<{
  code: number;
  message: string;
  data: UserProfile;
}> {
  return request('/admin/auth/profile', {
    method: 'GET',
  });
}

/**
 * 更新用户资料
 */
export async function updateUserProfile(params: UpdateProfileParams): Promise<{
  code: number;
  message: string;
  data: UserProfile;
}> {
  return request('/admin/auth/profile', {
    method: 'PUT',
    data: params,
  });
}

/**
 * 修改密码
 */
export async function changePassword(params: ChangePasswordParams): Promise<{
  code: number;
  message: string;
}> {
  return request('/admin/auth/change-password', {
    method: 'POST',
    data: params,
  });
}

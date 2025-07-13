import { request } from '@umijs/max';

// 登录请求参数
export interface LoginRequest {
  username: string;
  password: string;
  captcha: string;
}

// 刷新token请求参数
export interface RefreshTokenRequest {
  refresh_token: string;
}

// 用户信息
export interface UserInfo {
  id: number;
  username: string;
  nickname: string;  // 后端返回的是nickname字段
  avatar?: string;
  email?: string;
  phone?: string;
  status: number;
  last_login_at?: string;
  last_login_ip?: string;
  login_count: number;
}

// 权限信息
export interface Permission {
  sign: string;
  name: string;
  type: string;
}

// 菜单信息
export interface Menu {
  sign: string;
  name: string;
  path?: string;
  icon?: string;
  parent_sign?: string;
  sort: number;
  children?: Menu[];
}

// 登录响应
export interface LoginResponse {
  token: string;
  expires_at: string;
  user_info: UserInfo;
  permissions: Permission[];
  menus: Menu[];
}

// API响应格式
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

// 登录
export async function login(data: LoginRequest, options?: { [key: string]: any }) {
  return request<ApiResponse<LoginResponse>>('/admin/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}

// 登出
export async function logout(options?: { [key: string]: any }) {
  return request<ApiResponse<null>>('/admin/auth/logout', {
    method: 'POST',
    ...(options || {}),
  });
}

// 刷新token
export async function refreshToken(data: RefreshTokenRequest, options?: { [key: string]: any }) {
  return request<ApiResponse<LoginResponse>>('/admin/auth/refresh', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}

// 获取用户信息
export async function getUserProfile(options?: { [key: string]: any }) {
  return request<ApiResponse<UserInfo>>('/admin/auth/profile', {
    method: 'GET',
    ...(options || {}),
  });
}

// 更新用户信息
export async function updateUserProfile(data: Partial<UserInfo>, options?: { [key: string]: any }) {
  return request<ApiResponse<UserInfo>>('/admin/auth/profile', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}

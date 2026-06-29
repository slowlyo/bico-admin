/**
 * 认证相关类型定义
 */

export interface LoginParams {
  username: string;
  password: string;
  captchaId: string;
  captchaCode: string;
}

export interface LoginResult {
  token: string;
}

export interface CurrentUser {
  id: number;
  username: string;
  name: string;
  avatar: string;
  roles?: string[];
  permissions?: string[];
}

export interface UpdateProfileParams {
  name?: string;
  avatar?: string;
}

export interface ChangePasswordParams {
  oldPassword: string;
  newPassword: string;
}

export interface CaptchaResult {
  id: string;
  image: string;
}

import { request } from '@umijs/max';
import { buildApiUrl } from '../config';
import type {
  LoginParams,
  LoginResult,
  CurrentUser,
  CaptchaResult,
} from './types';

/**
 * 登录
 */
export async function login(data: LoginParams) {
  return request<API.Response<LoginResult>>(buildApiUrl('/auth/login'), {
    method: 'POST',
    data,
  });
}

/**
 * 退出登录
 */
export async function logout() {
  return request<API.Response<any>>(buildApiUrl('/auth/logout'), {
    method: 'POST',
  });
}

/**
 * 获取当前用户信息
 */
export async function getCurrentUser() {
  return request<API.Response<CurrentUser>>(buildApiUrl('/auth/current-user'), {
    method: 'GET',
  });
}

/**
 * 获取验证码
 */
export async function getCaptcha() {
  return request<API.Response<CaptchaResult>>(buildApiUrl('/captcha'), {
    method: 'GET',
  });
}

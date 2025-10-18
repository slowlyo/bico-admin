import { request } from '@umijs/max';

/**
 * 登录接口
 * POST /admin-api/auth/login
 */
export async function login(body: API.LoginParams) {
  return request<API.Response<API.LoginResult>>('/admin-api/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
  });
}

/**
 * 退出登录接口
 * POST /admin-api/auth/logout
 */
export async function logout() {
  return request<API.Response<any>>('/admin-api/auth/logout', {
    method: 'POST',
  });
}

/**
 * 获取当前用户信息
 * GET /admin-api/auth/current-user
 */
export async function currentUser() {
  return request<API.Response<API.CurrentUser>>('/admin-api/auth/current-user', {
    method: 'GET',
  });
}

/**
 * 获取验证码
 * GET /admin-api/captcha
 */
export async function getCaptcha() {
  return request<API.Response<{ id: string; image: string }>>('/admin-api/captcha', {
    method: 'GET',
  });
}

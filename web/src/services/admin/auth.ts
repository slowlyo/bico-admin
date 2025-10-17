import { request } from '@umijs/max';

/**
 * 登录接口
 * POST /admin-api/login
 */
export async function login(body: API.LoginParams) {
  return request<API.Response<API.LoginResult>>('/admin-api/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
  });
}

/**
 * 退出登录接口
 * POST /admin-api/logout
 */
export async function logout() {
  return request<API.Response<any>>('/admin-api/logout', {
    method: 'POST',
  });
}

/**
 * 获取当前用户信息
 * GET /admin-api/current-user
 */
export async function currentUser() {
  return request<API.Response<API.CurrentUser>>('/admin-api/current-user', {
    method: 'GET',
  });
}

import { request } from '@umijs/max';

/**
 * 获取应用配置
 * GET /admin-api/app-config
 */
export async function getAppConfig() {
  return request<API.Response<API.AppConfig>>('/admin-api/app-config', {
    method: 'GET',
  });
}

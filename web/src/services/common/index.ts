import { request } from '@umijs/max';
import { buildApiUrl } from '../config';
import type { AppConfig } from './types';

/**
 * 获取应用配置
 */
export async function getAppConfig() {
  return request<API.Response<AppConfig>>(buildApiUrl('/app-config'), {
    method: 'GET',
  });
}

/**
 * API配置
 */

// API基础URL前缀
export const API_PREFIX = '/admin-api';

/**
 * 构建完整的API路径
 */
export function buildApiUrl(path: string): string {
  // 确保path以/开头
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;
  return `${API_PREFIX}${normalizedPath}`;
}

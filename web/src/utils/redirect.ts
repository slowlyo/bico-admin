/**
 * 跳转相关工具
 * 统一处理 redirect 的构造与校验，避免 hash 路由下取错路径
 */

import { history } from '@umijs/max';

export const LOGIN_PATH = '/auth/login';

/**
 * 获取当前真实路由（兼容 hash history）
 */
export function getCurrentPathWithSearch(): string {
  const { pathname, search } = history.location;
  return `${pathname}${search || ''}`;
}

/**
 * 校验 redirect，只允许站内相对路径，防止开放重定向
 */
export function getSafeRedirect(redirect?: string | null): string {
  if (!redirect) return '/';

  if (!redirect.startsWith('/')) return '/';

  // 禁止协议/域名形式（如 //evil.com 或 /\evil.com 等）
  if (redirect.startsWith('//')) return '/';
  if (redirect.includes('://')) return '/';

  return redirect;
}

/**
 * 从 url search 中读取并校验 redirect
 */
export function getSafeRedirectFromSearch(search?: string): string {
  const searchParams = new URLSearchParams(search || '');
  const redirect = searchParams.get('redirect');
  return getSafeRedirect(redirect);
}

/**
 * 构造登录页地址（带 redirect）
 */
export function buildLoginUrl(redirect?: string): string {
  const safeRedirect = getSafeRedirect(redirect);
  const searchParams = new URLSearchParams({
    redirect: safeRedirect,
  });
  return `${LOGIN_PATH}?${searchParams.toString()}`;
}

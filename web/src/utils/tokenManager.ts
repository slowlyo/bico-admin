// Token管理工具
import { refreshToken as refreshTokenApi } from '@/services/auth';
import { history } from '@umijs/max';
import { message } from 'antd';

// Token刷新状态
let isRefreshing = false;
let refreshPromise: Promise<string | null> | null = null;

/**
 * 检查token是否即将过期
 * @param token JWT token
 * @returns 是否即将过期（剩余时间少于30分钟）
 */
export function isTokenExpiringSoon(token: string): boolean {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    const expirationTime = payload.exp * 1000; // 转换为毫秒
    const currentTime = Date.now();
    const timeUntilExpiration = expirationTime - currentTime;
    
    // 如果剩余时间少于30分钟，认为即将过期
    return timeUntilExpiration < 30 * 60 * 1000;
  } catch (error) {
    console.error('解析token失败:', error);
    return true; // 解析失败时认为已过期
  }
}

/**
 * 检查token是否已过期
 * @param token JWT token
 * @returns 是否已过期
 */
export function isTokenExpired(token: string): boolean {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    const expirationTime = payload.exp * 1000;
    return Date.now() >= expirationTime;
  } catch (error) {
    console.error('解析token失败:', error);
    return true;
  }
}

/**
 * 获取token剩余时间（毫秒）
 * @param token JWT token
 * @returns 剩余时间（毫秒），如果已过期返回0
 */
export function getTokenRemainingTime(token: string): number {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    const expirationTime = payload.exp * 1000;
    const remainingTime = expirationTime - Date.now();
    return Math.max(0, remainingTime);
  } catch (error) {
    console.error('解析token失败:', error);
    return 0;
  }
}

/**
 * 刷新token
 * @returns 新的token，如果刷新失败返回null
 */
export async function refreshToken(): Promise<string | null> {
  // 如果正在刷新，返回现有的Promise
  if (isRefreshing && refreshPromise) {
    return refreshPromise;
  }

  isRefreshing = true;
  
  refreshPromise = (async () => {
    try {
      const currentToken = localStorage.getItem('token');
      if (!currentToken) {
        throw new Error('没有找到当前token');
      }

      const response = await refreshTokenApi({
        refresh_token: currentToken,
      });

      if (response.code === 200) {
        const { token, user_info, permissions } = response.data;
        
        // 更新本地存储
        localStorage.setItem('token', token);
        localStorage.setItem('userInfo', JSON.stringify(user_info));
        localStorage.setItem('permissions', JSON.stringify(permissions));
        
        console.log('Token刷新成功');
        return token;
      } else {
        throw new Error(response.message || 'Token刷新失败');
      }
    } catch (error) {
      console.error('Token刷新失败:', error);
      
      // 刷新失败，清除本地存储并跳转到登录页
      localStorage.removeItem('token');
      localStorage.removeItem('userInfo');
      localStorage.removeItem('permissions');
      
      // 只有在不是登录页时才跳转
      if (history.location.pathname !== '/login') {
        message.error('登录已过期，请重新登录');
        history.replace('/login');
      }
      
      return null;
    } finally {
      isRefreshing = false;
      refreshPromise = null;
    }
  })();

  return refreshPromise;
}

/**
 * 自动检查并刷新token
 * 如果token即将过期，自动刷新
 */
export async function autoRefreshToken(): Promise<void> {
  const token = localStorage.getItem('token');
  if (!token) {
    return;
  }

  // 如果token已过期，直接跳转到登录页
  if (isTokenExpired(token)) {
    localStorage.removeItem('token');
    localStorage.removeItem('userInfo');
    localStorage.removeItem('permissions');
    
    if (history.location.pathname !== '/login') {
      message.error('登录已过期，请重新登录');
      history.replace('/login');
    }
    return;
  }

  // 如果token即将过期，尝试刷新
  if (isTokenExpiringSoon(token)) {
    await refreshToken();
  }
}

/**
 * 启动token自动刷新定时器
 * 每5分钟检查一次token状态
 */
export function startTokenRefreshTimer(): () => void {
  const interval = setInterval(() => {
    autoRefreshToken();
  }, 5 * 60 * 1000); // 5分钟检查一次

  // 立即执行一次检查
  autoRefreshToken();

  // 返回清理函数
  return () => {
    clearInterval(interval);
  };
}

/**
 * 在请求前检查token状态
 * 如果token即将过期，先刷新再发送请求
 */
export async function ensureValidToken(): Promise<string | null> {
  const token = localStorage.getItem('token');
  if (!token) {
    return null;
  }

  // 如果token已过期，返回null
  if (isTokenExpired(token)) {
    return null;
  }

  // 如果token即将过期，尝试刷新
  if (isTokenExpiringSoon(token)) {
    const newToken = await refreshToken();
    return newToken || token; // 如果刷新失败，返回原token让后端处理
  }

  return token;
}

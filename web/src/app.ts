// 运行时配置
import { history } from '@umijs/max';
import { message } from 'antd';
import { UserInfo, logout as logoutApi } from '@/services/auth';

// 全局状态更新函数，用于在退出登录时同步更新状态
let globalSetInitialState: ((initialState: any) => void) | null = null;

export interface InitialState {
  currentUser?: UserInfo;
  permissions?: string[]; // 权限代码字符串数组
  fetchUserInfo?: () => Promise<{ userInfo?: UserInfo; permissions?: string[] }>;
  // 添加umi layout需要的字段
  name?: string;
  avatar?: string;
}

// 获取用户信息和权限
const fetchUserInfo = async (): Promise<{ userInfo?: UserInfo; permissions?: string[] }> => {
  try {
    const token = localStorage.getItem('token');
    if (!token) {
      return {};
    }

    // 调用后端API获取最新的用户信息和权限
    const response = await fetch('/admin/auth/profile', {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });

    if (response.ok) {
      const data = await response.json();
      if (data.code === 200) {
        // 更新本地存储
        localStorage.setItem('userInfo', JSON.stringify(data.data.user_info));
        localStorage.setItem('permissions', JSON.stringify(data.data.permissions));

        return {
          userInfo: data.data.user_info,
          permissions: data.data.permissions,
        };
      }
    }

    // 如果API调用失败，清除本地存储
    localStorage.removeItem('token');
    localStorage.removeItem('userInfo');
    localStorage.removeItem('permissions');
    return {};
  } catch (error) {
    console.error('获取用户信息失败:', error);
    // 发生错误时，尝试从本地存储读取作为降级方案
    try {
      const userInfoStr = localStorage.getItem('userInfo');
      const permissionsStr = localStorage.getItem('permissions');

      if (userInfoStr && permissionsStr) {
        return {
          userInfo: JSON.parse(userInfoStr),
          permissions: JSON.parse(permissionsStr),
        };
      }
    } catch (parseError) {
      console.error('解析本地存储失败:', parseError);
    }

    return {};
  }
};

// 全局初始化数据配置，用于 Layout 用户信息和权限初始化
export async function getInitialState(): Promise<InitialState> {
  const { location } = history;

  // 获取最新的用户信息和权限
  const { userInfo, permissions } = await fetchUserInfo();

  // 如果用户已登录但在登录页，静默跳转到首页
  if (userInfo && location.pathname === '/login') {
    // 使用 setTimeout 确保在下一个事件循环中执行，避免阻塞渲染
    setTimeout(() => {
      history.replace('/');
    }, 0);
  }

  // 如果用户已登录，返回用户信息
  if (userInfo) {
    return {
      currentUser: userInfo,
      permissions: permissions || [],
      fetchUserInfo,
      // 添加name和avatar字段供umi layout使用
      name: userInfo.nickname,
      avatar: userInfo.avatar,
    };
  }

  return {
    fetchUserInfo,
  };
}

// 登出处理
const handleLogout = async () => {
  try {
    await logoutApi();
  } catch (error) {
    console.error('登出失败:', error);
  } finally {
    // 先跳转到登录页，避免在当前页面清除状态时的闪烁
    history.replace('/login');

    // 延迟清除本地存储，确保页面已经跳转
    setTimeout(() => {
      localStorage.removeItem('token');
      localStorage.removeItem('userInfo');
      localStorage.removeItem('permissions');

      // 更新全局状态
      if (globalSetInitialState) {
        globalSetInitialState({
          fetchUserInfo,
        });
      }
    }, 50);

    message.success('已退出登录');
  }
};

// 请求拦截器
export const request = {
  // 请求拦截器
  requestInterceptors: [
    (config: any) => {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    },
  ],
  // 响应拦截器
  responseInterceptors: [
    (response: any) => {
      // 如果返回401，说明token过期，跳转到登录页
      if (response.status === 401) {
        // 先跳转到登录页
        history.replace('/login');

        // 延迟清除状态，避免闪烁
        setTimeout(() => {
          localStorage.removeItem('token');
          localStorage.removeItem('userInfo');
          localStorage.removeItem('permissions');

          // 更新全局状态
          if (globalSetInitialState) {
            globalSetInitialState({
              fetchUserInfo,
            });
          }
        }, 50);
      }
      return response;
    },
  ],
};

export const layout = (_initialState: InitialState, setInitialState: any) => {
  // 保存全局状态更新函数
  globalSetInitialState = setInitialState;

  return {
    logo: 'https://img.alicdn.com/tfs/TB1YHEpwUT1gK0jSZFhXXaAtVXa-28-27.svg',
    menu: {
      locale: false,
    },
    layout: 'mix',
    logout: handleLogout,
    onPageChange: () => {
      const { location } = history;
      // 如果没有登录，重定向到 login
      // 只检查localStorage中的token，避免状态更新延迟导致的问题
      const hasToken = localStorage.getItem('token');
      if (!hasToken && location.pathname !== '/login') {
        history.replace('/login');
      }
    },
  };
};

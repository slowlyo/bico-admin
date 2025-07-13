// 运行时配置
import { history } from '@umijs/max';
import { message } from 'antd';
import { UserInfo, Permission, Menu, logout as logoutApi } from '@/services/auth';

// 全局状态更新函数，用于在退出登录时同步更新状态
let globalSetInitialState: ((initialState: any) => void) | null = null;

export interface InitialState {
  currentUser?: UserInfo;
  permissions?: Permission[];
  menus?: Menu[];
  fetchUserInfo?: () => Promise<UserInfo | undefined>;
  // 添加umi layout需要的字段
  name?: string;
  avatar?: string;
}

// 获取用户信息
const fetchUserInfo = async (): Promise<UserInfo | undefined> => {
  try {
    const token = localStorage.getItem('token');
    if (!token) {
      return undefined;
    }

    const userInfoStr = localStorage.getItem('userInfo');
    if (userInfoStr) {
      return JSON.parse(userInfoStr);
    }

    return undefined;
  } catch (error) {
    console.error('获取用户信息失败:', error);
    return undefined;
  }
};

// 全局初始化数据配置，用于 Layout 用户信息和权限初始化
export async function getInitialState(): Promise<InitialState> {
  const { location } = history;

  // 同步获取用户信息，避免异步导致的状态不一致
  const currentUser = await fetchUserInfo();

  // 如果用户已登录但在登录页，静默跳转到首页
  if (currentUser && location.pathname === '/login') {
    // 使用 setTimeout 确保在下一个事件循环中执行，避免阻塞渲染
    setTimeout(() => {
      history.replace('/');
    }, 0);
  }

  // 如果用户已登录，返回用户信息
  if (currentUser) {
    const permissionsStr = localStorage.getItem('permissions');
    const menusStr = localStorage.getItem('menus');

    return {
      currentUser,
      permissions: permissionsStr ? JSON.parse(permissionsStr) : [],
      menus: menusStr ? JSON.parse(menusStr) : [],
      fetchUserInfo,
      // 添加name和avatar字段供umi layout使用
      name: currentUser.nickname,
      avatar: currentUser.avatar,
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
      localStorage.removeItem('menus');

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
          localStorage.removeItem('menus');

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

export const layout = (initialState: InitialState, setInitialState: any) => {
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

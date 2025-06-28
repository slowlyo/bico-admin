// 运行时配置
import { history } from '@umijs/max';
import { message } from 'antd';
import React from 'react';
import { getUserProfile } from '@/services/auth';
import { AvatarDropdown, AvatarName } from '@/components';

const loginPath = '/login';

/**
 * @see https://umijs.org/docs/api/runtime-config#getinitialstate
 */
export async function getInitialState(): Promise<{
  currentUser?: API.CurrentUser;
  loading?: boolean;
  fetchUserInfo?: () => Promise<API.CurrentUser | undefined>;
}> {
  const fetchUserInfo = async () => {
    try {
      const userInfo = await getUserProfile();
      if (userInfo.code === 200) {
        return userInfo.data;
      }
      return undefined;
    } catch (error) {
      // 如果获取用户信息失败，跳转到登录页
      history.push(loginPath);
      return undefined;
    }
  };

  // 如果不是登录页面，执行
  const { location } = history;
  if (location.pathname !== loginPath) {
    const currentUser = await fetchUserInfo();
    return {
      fetchUserInfo,
      currentUser,
    };
  }

  return {
    fetchUserInfo,
  };
}

// ProLayout 支持的api https://procomponents.ant.design/components/layout
export const layout = ({ initialState }: any) => {
  return {
    logo: 'https://gw.alipayobjects.com/zos/antfincdn/PmY%24TNNDBI/logo.svg',
    menu: {
      locale: false,
    },
    layout: 'mix',
    avatarProps: {
      src: 'https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png',
      title: <AvatarName />,
      render: (_: any, avatarChildren: any) => {
        return <AvatarDropdown>{avatarChildren}</AvatarDropdown>;
      },
    },
    menuHeaderRender: undefined,
    // 自定义页脚
    footerRender: () => {
      return (
        <div style={{ textAlign: 'center', color: '#999' }}>
          Bico Admin ©2024 Created by Bico Team
        </div>
      );
    },
    // 页面切换时的处理
    onPageChange: () => {
      const { location } = history;
      // 如果没有登录，重定向到 login
      if (!initialState?.currentUser && location.pathname !== loginPath) {
        history.push(loginPath);
      }
    },
  };
};

// 请求拦截器
export const request = {
  timeout: 10000,
  errorConfig: {
    errorThrower: (res: any) => {
      const { success, data, errorMessage } = res;
      if (!success) {
        const error: any = new Error(errorMessage);
        error.name = 'BizError';
        error.info = { errorMessage, data };
        throw error;
      }
    },
    errorHandler: (error: any, opts: any) => {
      if (opts?.skipErrorHandler) throw error;

      if (error.name === 'BizError') {
        message.error(error.info.errorMessage);
      } else if (error.response) {
        // 请求成功发出且服务器也响应了状态码，但状态代码超出了 2xx 的范围
        if (error.response.status === 401) {
          message.error('未授权，请重新登录');
          localStorage.removeItem('token');
          history.push('/login');
        } else {
          message.error(`请求错误 ${error.response.status}: ${error.response.statusText}`);
        }
      } else if (error.request) {
        // 请求已经成功发起，但没有收到响应
        message.error('网络错误，请检查您的网络连接');
      } else {
        // 发送请求时出了点问题
        message.error('请求配置错误');
      }
    },
  },
  requestInterceptors: [
    (config: any) => {
      // 添加认证头
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    },
  ],
  responseInterceptors: [
    (response: any) => {
      // 统一处理响应
      return response;
    },
  ],
};

import type { RequestOptions } from '@@/plugin-request/request';
import type { RequestConfig } from '@umijs/max';
import { message, notification } from 'antd';

// 错误处理方案： 错误类型
enum ErrorShowType {
  SILENT = 0,
  WARN_MESSAGE = 1,
  ERROR_MESSAGE = 2,
  NOTIFICATION = 3,
  REDIRECT = 9,
}
// 与后端约定的响应数据格式
interface ResponseStructure {
  success: boolean;
  data: any;
  errorCode?: number;
  errorMessage?: string;
  showType?: ErrorShowType;
}

/**
 * @name 错误处理
 * pro 自带的错误处理， 可以在这里做自己的改动
 * @doc https://umijs.org/docs/max/request#配置
 */
export const errorConfig: RequestConfig = {
  // 错误处理： umi@3 的错误处理方案。
  errorConfig: {
    // 禁用 errorThrower，使用响应拦截器处理
    errorThrower: () => {},
    // 错误接收及处理
    errorHandler: (error: any, opts: any) => {
      if (opts?.skipErrorHandler) throw error;
      
      // 业务错误不在这里显示消息，由业务代码自行处理
      if (error.name === 'BizError') {
        throw error;
      }
      
      // 网络错误统一处理
      if (error.response) {
        // HTTP 状态码错误
        const status = error.response.status;
        if (status === 401) {
          message.error('未授权，请重新登录');
        } else if (status === 403) {
          message.error('无权限访问');
        } else if (status === 404) {
          message.error('请求的资源不存在');
        } else if (status >= 500) {
          message.error('服务器错误，请稍后重试');
        } else {
          message.error(`请求错误: ${status}`);
        }
      } else if (error.request) {
        // 请求发出但没有收到响应
        message.error('网络错误，请检查网络连接');
      } else {
        // 其他错误
        message.error('请求失败，请重试');
      }
      
      throw error;
    },
  },

  // 请求拦截器
  requestInterceptors: [
    (config: RequestOptions) => {
      // 添加 token 到请求头
      const token = localStorage.getItem('token');
      if (token) {
        config.headers = {
          ...config.headers,
          Authorization: `Bearer ${token}`,
        };
      }
      return config;
    },
  ],

  // 响应拦截器
  responseInterceptors: [
    (response) => {
      // 拦截响应数据，进行个性化处理
      const { data } = response;

      // 处理后端统一响应格式
      if (data && typeof data === 'object' && 'code' in data) {
        const apiResponse = data as any;
        
        // code !== 0 表示业务错误
        if (apiResponse.code !== 0) {
          // 401 未授权（包括 token 失效、账户被禁用等），清除 token 并跳转登录
          if (apiResponse.code === 401 && window.location.pathname !== '/auth/login') {
            message.error(apiResponse.msg || '未授权，请重新登录');
            localStorage.removeItem('token');
            localStorage.removeItem('currentUser');
            setTimeout(() => {
              window.location.href = '/auth/login';
            }, 1000);
            const error: any = new Error(apiResponse.msg || '未授权');
            error.name = 'BizError';
            error.response = response;
            error.data = apiResponse;
            throw error;
          }
          
          // 403 权限不足，显示错误消息但不跳转
          if (apiResponse.code === 403) {
            message.error(apiResponse.msg || '无权访问');
          }
          
          // 抛出错误让业务代码可以捕获
          const error: any = new Error(apiResponse.msg || '请求失败');
          error.name = 'BizError';
          error.response = response;
          error.data = apiResponse;
          throw error;
        }
      }
      
      return response;
    },
  ],
};

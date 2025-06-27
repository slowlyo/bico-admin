import axios from 'axios';
import config from '../config';

// Token存储键
export const TOKEN_KEY = 'bico-admin-token';
export const USER_KEY = 'bico-admin-user';

// 创建共享的axios实例
export const axiosInstance = axios.create({
  baseURL: config.adminApiUrl,
});

// 设置请求拦截器，自动添加token
axiosInstance.interceptors.request.use(
  (config) => {
    const token = getToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 设置响应拦截器，处理401错误
axiosInstance.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response?.status === 401) {
      clearAuth();
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// API响应接口 - 匹配后端响应格式
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

// 用户信息接口
export interface UserInfo {
  id: number;
  username: string;
  email: string;
  nickname: string;
  avatar: string;
  phone: string;
  status: number;
  last_login_at: string | null;
  created_at: string;
  updated_at: string;
  roles: Array<{
    id: number;
    name: string;
    description: string;
  }>;
}

// 登录请求接口
export interface LoginRequest {
  username: string;
  password: string;
}

// 登录响应接口
export interface LoginResponse {
  user: UserInfo;
  token: string;
}

// 请求配置接口
export interface RequestConfig {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';
  headers?: Record<string, string>;
  body?: any;
  requireAuth?: boolean;
}

// 获取存储的token
export const getToken = (): string | null => {
  return localStorage.getItem(TOKEN_KEY);
};

// 设置token
export const setToken = (token: string): void => {
  localStorage.setItem(TOKEN_KEY, token);
};

// 移除token
export const removeToken = (): void => {
  localStorage.removeItem(TOKEN_KEY);
};

// 获取存储的用户信息
export const getUser = (): UserInfo | null => {
  const userStr = localStorage.getItem(USER_KEY);
  if (userStr) {
    try {
      return JSON.parse(userStr);
    } catch {
      return null;
    }
  }
  return null;
};

// 设置用户信息
export const setUser = (user: UserInfo): void => {
  localStorage.setItem(USER_KEY, JSON.stringify(user));
};

// 移除用户信息
export const removeUser = (): void => {
  localStorage.removeItem(USER_KEY);
};

// 清除所有认证信息
export const clearAuth = (): void => {
  removeToken();
  removeUser();
};

// HTTP请求工具类
class RequestUtil {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  // 发送请求
  async request<T = any>(
    url: string,
    config: RequestConfig = {}
  ): Promise<ApiResponse<T>> {
    const {
      method = 'GET',
      headers = {},
      body,
      requireAuth = true,
    } = config;

    // 构建完整URL
    const fullUrl = url.startsWith('http') ? url : `${this.baseURL}${url}`;

    // 构建请求头
    const requestHeaders: Record<string, string> = {
      'Content-Type': 'application/json',
      ...headers,
    };

    // 添加认证头
    if (requireAuth) {
      const token = getToken();
      if (token) {
        requestHeaders.Authorization = `Bearer ${token}`;
      }
    }

    // 构建请求配置
    const requestConfig: RequestInit = {
      method,
      headers: requestHeaders,
    };

    // 添加请求体
    if (body && method !== 'GET') {
      requestConfig.body = JSON.stringify(body);
    }

    try {
      const response = await fetch(fullUrl, requestConfig);
      
      // 检查响应状态
      if (!response.ok) {
        // 如果是401未授权，清除认证信息
        if (response.status === 401) {
          clearAuth();
          // 可以在这里触发重定向到登录页
          window.location.href = '/login';
        }
        
        // 尝试解析错误响应
        try {
          const errorData = await response.json();
          throw new Error(errorData.message || `HTTP ${response.status}`);
        } catch {
          throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
      }

      // 解析响应数据
      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Request failed:', error);
      throw error;
    }
  }

  // GET请求
  async get<T = any>(url: string, config?: Omit<RequestConfig, 'method'>): Promise<ApiResponse<T>> {
    return this.request<T>(url, { ...config, method: 'GET' });
  }

  // POST请求
  async post<T = any>(url: string, body?: any, config?: Omit<RequestConfig, 'method' | 'body'>): Promise<ApiResponse<T>> {
    return this.request<T>(url, { ...config, method: 'POST', body });
  }

  // PUT请求
  async put<T = any>(url: string, body?: any, config?: Omit<RequestConfig, 'method' | 'body'>): Promise<ApiResponse<T>> {
    return this.request<T>(url, { ...config, method: 'PUT', body });
  }

  // DELETE请求
  async delete<T = any>(url: string, config?: Omit<RequestConfig, 'method'>): Promise<ApiResponse<T>> {
    return this.request<T>(url, { ...config, method: 'DELETE' });
  }

  // PATCH请求
  async patch<T = any>(url: string, body?: any, config?: Omit<RequestConfig, 'method' | 'body'>): Promise<ApiResponse<T>> {
    return this.request<T>(url, { ...config, method: 'PATCH', body });
  }
}

// 创建请求实例 - 专门用于认证的API
export const request = new RequestUtil('http://localhost:8080');

// 认证相关API
export const authAPI = {
  // 用户登录
  login: async (data: LoginRequest): Promise<ApiResponse<LoginResponse>> => {
    const response = await axiosInstance.post('/auth/login', data);
    return response.data;
  },

  // 用户登出
  logout: async (): Promise<ApiResponse<null>> => {
    const response = await axiosInstance.post('/auth/logout');
    return response.data;
  },

  // 获取用户资料
  getProfile: async (): Promise<ApiResponse<UserInfo>> => {
    const response = await axiosInstance.get('/auth/profile');
    return response.data;
  },

  // 更新用户资料
  updateProfile: async (data: Partial<UserInfo>): Promise<ApiResponse<UserInfo>> => {
    const response = await axiosInstance.put('/auth/profile', data);
    return response.data;
  },

  // 修改密码
  changePassword: async (data: { old_password: string; new_password: string }): Promise<ApiResponse<null>> => {
    const response = await axiosInstance.post('/auth/change-password', data);
    return response.data;
  },
};

export default request;

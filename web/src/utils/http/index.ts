import axios, { InternalAxiosRequestConfig, AxiosRequestConfig, AxiosResponse } from 'axios'
import { useUserStore } from '@/store/modules/user'
import { ApiStatus } from './status'
import { HttpError, handleError, showError } from './error'
import { $t } from '@/locales'

// 常量定义
const REQUEST_TIMEOUT = 15000 // 请求超时时间(毫秒)
const MAX_RETRIES = 2 // 最大重试次数
const RETRY_DELAY = 1000 // 重试延迟时间(毫秒)

// 防止重复登出的标志
let isLoggingOut = false

// 防止重复显示认证错误提示
let lastAuthErrorTime = 0
const AUTH_ERROR_THROTTLE = 3000 // 3秒内不重复显示认证错误

// 扩展 AxiosRequestConfig 类型
interface ExtendedAxiosRequestConfig extends AxiosRequestConfig {
  showErrorMessage?: boolean
  returnFullResponse?: boolean // 是否返回完整响应（包括message）
  enableRetry?: boolean // 是否启用重试机制，默认关闭
  maxRetries?: number // 最大重试次数，默认使用全局配置
  _retry?: boolean // 标记请求是否已重试，避免无限循环
}

const { VITE_API_URL, VITE_WITH_CREDENTIALS } = import.meta.env

const axiosInstance = axios.create({
  timeout: REQUEST_TIMEOUT, // 请求超时时间(毫秒)
  baseURL: VITE_API_URL, // API地址
  withCredentials: VITE_WITH_CREDENTIALS === 'true', // 是否携带cookie，默认关闭
  transformRequest: [(data) => JSON.stringify(data)], // 请求数据转换为 JSON 字符串
  validateStatus: (status) => status >= 200 && status < 300, // 只接受 2xx 的状态码
  headers: {
    get: { 'Content-Type': 'application/x-www-form-urlencoded;charset=utf-8' },
    post: { 'Content-Type': 'application/json;charset=utf-8' }
  },
  transformResponse: [
    (data, headers) => {
      const contentType = headers['content-type']
      if (contentType && contentType.includes('application/json')) {
        try {
          return JSON.parse(data)
        } catch {
          return data
        }
      }
      return data
    }
  ]
})

// 处理认证错误（防重复提示）
function handleAuthError(message: string): void {
  const now = Date.now()
  if (now - lastAuthErrorTime > AUTH_ERROR_THROTTLE) {
    lastAuthErrorTime = now
    console.error('认证错误:', message)
    // 显示用户友好的错误提示（防重复）
    showError(new HttpError(message, ApiStatus.unauthorized), true)
  }
}

// 请求拦截器
axiosInstance.interceptors.request.use(
  async (request: InternalAxiosRequestConfig) => {
    // 确保token有效，如果即将过期会触发后台刷新（非阻塞）
    const { ensureValidToken } = await import('@/utils/tokenManager')
    const token = ensureValidToken()

    // 设置 token 和 请求头
    if (token) {
      request.headers.set('Authorization', `Bearer ${token}`)
      request.headers.set('Content-Type', 'application/json')
    }

    return request
  },
  (error) => {
    showError(new HttpError($t('httpMsg.requestConfigError'), ApiStatus.error))
    return Promise.reject(error)
  }
)

// 响应拦截器
axiosInstance.interceptors.response.use(
  (response: AxiosResponse<Api.Http.BaseResponse>) => {
    const { code, message, msg } = response.data
    const errorMessage = message || msg

    switch (code) {
      case ApiStatus.success:
        return response
      case ApiStatus.unauthorized:
        // 检查是否是登录请求，如果是登录请求则不执行登出操作
        const isLoginRequest = response.config.url?.includes('/auth/login')
        if (!isLoginRequest) {
          // 使用统一的认证错误处理（已包含用户提示）
          handleAuthError(errorMessage || $t('httpMsg.unauthorized'))
          // 异步执行登出操作，不阻塞当前响应处理
          logOut().catch(error => console.error('登出操作失败:', error))
        }
        // 抛出错误但不再显示额外的用户提示（已在handleAuthError中处理）
        throw new HttpError(errorMessage || $t('httpMsg.unauthorized'), ApiStatus.unauthorized)
      default:
        throw new HttpError(errorMessage || $t('httpMsg.requestFailed'), code)
    }
  },
  async (error) => {
    // 处理HTTP状态码401的情况
    if (error.response?.status === 401) {
      // 检查是否是登录或刷新token请求，避免无限循环
      const isAuthRequest = error.config?.url?.includes('/auth/login') ||
                           error.config?.url?.includes('/auth/refresh')

      if (!isAuthRequest) {
        try {
          // 尝试刷新token
          const { refreshToken } = await import('@/utils/tokenManager')
          const newToken = await refreshToken()

          if (newToken && error.config && !error.config._retry) {
            // 标记请求已重试，避免无限循环
            error.config._retry = true
            // 如果刷新成功，重新发送原请求
            error.config.headers.Authorization = `Bearer ${newToken}`
            return axiosInstance.request(error.config)
          } else {
            // 如果没有获取到新token，说明刷新失败
            handleAuthError('Token刷新失败，请重新登录')
            logOut().catch(error => console.error('登出操作失败:', error))
          }
        } catch (refreshError) {
          console.error('Token刷新失败:', refreshError)
          // 刷新失败，统一处理认证错误
          handleAuthError('认证失败，请重新登录')
          logOut().catch(error => console.error('登出操作失败:', error))
        }
      } else {
        // 认证请求失败，直接执行登出（不显示重复错误）
        logOut().catch(error => console.error('登出操作失败:', error))
      }
    }

    return Promise.reject(handleError(error))
  }
)

// 请求重试函数
async function retryRequest<T>(
  config: ExtendedAxiosRequestConfig,
  retries?: number
): Promise<T> {
  // 如果没有启用重试，直接执行请求
  if (!config.enableRetry) {
    return await request<T>(config)
  }

  const maxRetries = retries ?? config.maxRetries ?? MAX_RETRIES

  try {
    return await request<T>(config)
  } catch (error) {
    if (maxRetries > 0 && error instanceof HttpError && shouldRetry(error.code)) {
      await new Promise((resolve) => setTimeout(resolve, RETRY_DELAY))
      return retryRequest<T>(config, maxRetries - 1)
    }
    throw error
  }
}

// 判断是否需要重试（仅针对网络相关错误）
function shouldRetry(statusCode: number): boolean {
  return [
    ApiStatus.requestTimeout,
    ApiStatus.badGateway,
    ApiStatus.serviceUnavailable,
    ApiStatus.gatewayTimeout
  ].includes(statusCode)
}

// 请求函数
async function request<T = any>(config: ExtendedAxiosRequestConfig): Promise<T> {
  // 对 POST | PUT 请求特殊处理
  if (config.method?.toUpperCase() === 'POST' || config.method?.toUpperCase() === 'PUT') {
    if (config.params && !config.data) {
      config.data = config.params
      config.params = undefined
    }
  }

  try {
    const res = await axiosInstance.request<Api.Http.BaseResponse<T>>(config)
    // 根据配置决定返回完整响应还是只返回数据
    if (config.returnFullResponse) {
      return res.data as T // 返回完整响应 {code, message, data}
    }
    return res.data.data as T // 只返回数据部分
  } catch (error) {
    if (error instanceof HttpError) {
      // 对于认证错误，不再重复显示提示（已在拦截器中处理）
      const isAuthError = error.code === ApiStatus.unauthorized
      const showErrorMessage = config.showErrorMessage !== false && !isAuthError
      showError(error, showErrorMessage)
    }
    return Promise.reject(error)
  }
}

// API 方法集合
const api = {
  get<T>(config: ExtendedAxiosRequestConfig): Promise<T> {
    return retryRequest<T>({ ...config, method: 'GET' })
  },
  post<T>(config: ExtendedAxiosRequestConfig): Promise<T> {
    return retryRequest<T>({ ...config, method: 'POST' })
  },
  put<T>(config: ExtendedAxiosRequestConfig): Promise<T> {
    return retryRequest<T>({ ...config, method: 'PUT' })
  },
  del<T>(config: ExtendedAxiosRequestConfig): Promise<T> {
    return retryRequest<T>({ ...config, method: 'DELETE' })
  },
  request<T>(config: ExtendedAxiosRequestConfig): Promise<T> {
    return retryRequest<T>({ ...config })
  },
  // 带重试的请求方法
  getWithRetry<T>(config: ExtendedAxiosRequestConfig): Promise<T> {
    return retryRequest<T>({ ...config, method: 'GET', enableRetry: true })
  },
  postWithRetry<T>(config: ExtendedAxiosRequestConfig): Promise<T> {
    return retryRequest<T>({ ...config, method: 'POST', enableRetry: true })
  },
  putWithRetry<T>(config: ExtendedAxiosRequestConfig): Promise<T> {
    return retryRequest<T>({ ...config, method: 'PUT', enableRetry: true })
  },
  delWithRetry<T>(config: ExtendedAxiosRequestConfig): Promise<T> {
    return retryRequest<T>({ ...config, method: 'DELETE', enableRetry: true })
  }
}

// 退出登录函数
const logOut = async (): Promise<void> => {
  // 防止重复登出
  if (isLoggingOut) {
    return
  }

  isLoggingOut = true

  try {
    await useUserStore().logOut()
  } catch (error) {
    console.error('登出操作失败:', error)
    // 即使登出操作失败，也要确保跳转到登录页
    const { router } = await import('@/router')
    const { RoutesAlias } = await import('@/router/routesAlias')
    router.push(RoutesAlias.Login)
  } finally {
    // 重置标志，但延迟一点时间避免立即重复调用
    setTimeout(() => {
      isLoggingOut = false
    }, 1000)
  }
}

export default api

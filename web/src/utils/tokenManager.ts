// Token管理工具
// 优化版本：移除阻塞逻辑，采用非阻塞的后台刷新策略
import { UserService } from '@/api/usersApi'
import { useUserStore } from '@/store/modules/user'

// Token刷新状态
let isRefreshing = false
let refreshPromise: Promise<string | null> | null = null

/**
 * 检查token是否即将过期
 * @param token JWT token
 * @returns 是否即将过期（剩余时间少于10分钟）
 */
export function isTokenExpiringSoon(token: string): boolean {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    const expirationTime = payload.exp * 1000 // 转换为毫秒
    const currentTime = Date.now()
    const timeUntilExpiration = expirationTime - currentTime

    // 如果剩余时间少于10分钟，认为即将过期
    // 减少刷新频率，避免过于频繁的后台刷新
    return timeUntilExpiration < 10 * 60 * 1000
  } catch (error) {
    console.error('解析token失败:', error)
    return true // 解析失败时认为已过期
  }
}

/**
 * 检查token是否已过期
 * @param token JWT token
 * @returns 是否已过期
 */
export function isTokenExpired(token: string): boolean {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    const expirationTime = payload.exp * 1000
    return Date.now() >= expirationTime
  } catch (error) {
    console.error('解析token失败:', error)
    return true
  }
}

/**
 * 获取token剩余时间（毫秒）
 * @param token JWT token
 * @returns 剩余时间（毫秒），如果已过期返回0
 */
export function getTokenRemainingTime(token: string): number {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    const expirationTime = payload.exp * 1000
    const remainingTime = expirationTime - Date.now()
    return Math.max(0, remainingTime)
  } catch (error) {
    console.error('解析token失败:', error)
    return 0
  }
}

/**
 * 刷新token
 * @returns 新的token，如果刷新失败返回null
 */
export async function refreshToken(): Promise<string | null> {
  // 如果正在刷新，返回现有的Promise
  if (isRefreshing && refreshPromise) {
    return refreshPromise
  }

  isRefreshing = true
  
  refreshPromise = (async () => {
    try {
      const userStore = useUserStore()
      const currentToken = userStore.accessToken
      
      if (!currentToken) {
        throw new Error('没有找到当前token')
      }

      const response = await UserService.refreshToken({
        refresh_token: currentToken
      })

      if (response.token) {
        const { token, user_info, permissions } = response

        // 更新store中的数据
        userStore.setToken(token)
        userStore.setUserInfo(user_info)
        userStore.setPermissions(permissions)

        console.log('Token刷新成功')
        return token
      } else {
        console.warn('Token刷新失败: 服务器未返回新token')
        return null
      }
    } catch (error) {
      console.error('Token刷新失败:', error)

      // 刷新失败，静默处理，不显示用户提示
      // 让HTTP拦截器统一处理认证错误
      return null
    } finally {
      isRefreshing = false
      refreshPromise = null
    }
  })()

  return refreshPromise
}

/**
 * 自动检查并刷新token
 * 如果token即将过期，自动刷新（非阻塞）
 */
export function autoRefreshToken(): void {
  const userStore = useUserStore()
  const token = userStore.accessToken

  if (!token) {
    return
  }

  // 如果token已过期，清空token让HTTP拦截器处理
  if (isTokenExpired(token)) {
    userStore.setToken('')
    return
  }

  // 如果token即将过期，异步刷新（不阻塞）
  if (isTokenExpiringSoon(token)) {
    refreshToken().catch(error => {
      console.error('后台token刷新失败:', error)
    })
  }
}

/**
 * 启动token自动刷新定时器
 * 每3分钟检查一次token状态（优化频率）
 */
export function startTokenRefreshTimer(): () => void {
  const interval = setInterval(() => {
    autoRefreshToken()
  }, 3 * 60 * 1000) // 3分钟检查一次，配合10分钟的过期判断

  // 立即执行一次检查
  autoRefreshToken()

  // 返回清理函数
  return () => {
    clearInterval(interval)
  }
}

/**
 * 在请求前检查token状态（非阻塞版本）
 * 如果token即将过期，触发后台刷新但不等待
 */
export function ensureValidToken(): string | null {
  const userStore = useUserStore()
  const token = userStore.accessToken

  if (!token) {
    return null
  }

  // 如果token已过期，返回null
  if (isTokenExpired(token)) {
    return null
  }

  // 如果token即将过期，触发后台刷新（不阻塞当前请求）
  if (isTokenExpiringSoon(token)) {
    refreshToken().catch(error => {
      console.error('后台token刷新失败:', error)
    })
  }

  return token
}

/**
 * 确保token有效（阻塞版本，仅在必要时使用）
 * 如果token即将过期，会等待刷新完成
 * 主要用于响应拦截器中的token刷新重试逻辑
 */
export async function ensureValidTokenBlocking(): Promise<string | null> {
  const userStore = useUserStore()
  const token = userStore.accessToken

  if (!token) {
    return null
  }

  // 如果token已过期，返回null
  if (isTokenExpired(token)) {
    return null
  }

  // 如果token即将过期，等待刷新完成
  if (isTokenExpiringSoon(token)) {
    const newToken = await refreshToken()
    return newToken || token // 如果刷新失败，返回原token让后端处理
  }

  return token
}

import request from '@/utils/http'

/**
 * 获取验证码
 * @returns 验证码响应
 */
export function fetchCaptcha() {
  return request.get<Api.Auth.CaptchaResponse>({
    url: '/admin-api/captcha'
  })
}

/**
 * 获取应用配置
 * @returns 应用配置
 */
export function fetchAppConfig() {
  return request.get<Api.App.Config>({
    url: '/admin-api/app-config'
  })
}

/**
 * 登录
 * @param params 登录参数
 * @returns 登录响应
 */
export function fetchLogin(params: Api.Auth.LoginParams) {
  return request.post<Api.Auth.LoginResponse>({
    url: '/admin-api/auth/login',
    params
  })
}

/**
 * 退出登录
 * @returns 退出登录响应
 */
export function fetchLogout() {
  return request.post({
    url: '/admin-api/auth/logout'
  })
}

/**
 * 获取用户信息
 * @returns 用户信息
 */
export function fetchGetUserInfo() {
  return request.get<Api.Auth.UserInfo>({
    url: '/admin-api/auth/current-user'
  })
}

/**
 * 更新个人资料
 * @param data 个人资料参数
 */
export function fetchUpdateProfile(data: any) {
  return request.put({
    url: '/admin-api/auth/profile',
    data
  })
}

/**
 * 修改密码
 * @param data 密码参数
 */
export function fetchChangePassword(data: any) {
  return request.put({
    url: '/admin-api/auth/password',
    data
  })
}

/**
 * 上传头像
 * @param file 头像文件
 */
export function fetchUploadAvatar(file: File) {
  const formData = new FormData()
  formData.append('avatar', file)
  return request.post({
    url: '/admin-api/auth/avatar',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

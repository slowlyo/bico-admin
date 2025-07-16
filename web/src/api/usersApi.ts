import request from '@/utils/http'

export class UserService {
  // 登录
  static login(params: Api.Auth.LoginParams) {
    return request.post<Api.Auth.LoginResponse>({
      url: '/admin-api/auth/login',
      params,
      showErrorMessage: false // 禁用自动错误显示，由登录页面手动处理
    })
  }

  // 获取用户信息
  static getUserInfo() {
    return request.get<{ user_info: Api.User.UserInfo; permissions: string[] }>({
      url: '/admin-api/auth/profile'
    })
  }

  // 登出
  static logout() {
    return request.post<null>({
      url: '/admin-api/auth/logout'
    })
  }

  // 刷新token
  static refreshToken(params: { refresh_token: string }) {
    return request.post<Api.Auth.LoginResponse>({
      url: '/admin-api/auth/refresh',
      params
    })
  }

  // 获取用户列表
  static getUserList(params: Api.Common.PaginatingSearchParams) {
    return request.get<Api.User.UserListData>({
      url: '/api/user/list',
      params
    })
  }
}

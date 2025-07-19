import request from '@/utils/http'

// 认证相关类型定义
export namespace AuthTypes {
  /** 登录参数 */
  export interface LoginParams {
    username: string
    password: string
    captcha: string
  }

  /** 登录响应 */
  export interface LoginResponse {
    token: string
    expires_at: string
    user_info: UserTypes.UserInfo
    permissions: string[]
  }
}

// 用户相关类型定义
export namespace UserTypes {
  /** 用户信息 */
  export interface UserInfo {
    id: number
    username: string
    name: string
    avatar: string
    email: string
    phone: string
    status: number
    status_text: string
    last_login_at?: string
    remark: string
    can_delete: boolean
    can_disable: boolean
    roles: UserRole[]
    created_at: string
    updated_at: string
  }

  /** 用户角色 */
  export interface UserRole {
    id: number
    name: string
    code: string
    description: string
  }

  /** 用户列表数据 */
  export interface UserListData {
    records: UserListItem[]
    current: number
    size: number
    total: number
  }

  /** 个人信息更新请求 */
  export interface ProfileUpdateRequest {
    name: string
    avatar?: string
    email?: string
    phone?: string
  }

  /** 修改密码请求 */
  export interface ChangePasswordRequest {
    old_password: string
    new_password: string
  }

  /** 用户列表项 */
  export interface UserListItem {
    id: number
    avatar: string
    createBy: string
    createTime: string
    updateBy: string
    updateTime: string
    status: '1' | '2' | '3' | '4' // 1: 在线 2: 离线 3: 异常 4: 注销
    userName: string
    userGender: string
    nickName: string
    userPhone: string
    userEmail: string
    userRoles: string[]
  }
}

export class UserService {
  // 登录
  static login(params: AuthTypes.LoginParams) {
    return request.post<AuthTypes.LoginResponse>({
      url: '/admin-api/auth/login',
      data: params,
      showErrorMessage: false // 禁用自动错误显示，由登录页面手动处理
    })
  }

  // 获取用户信息
  static getUserInfo() {
    return request.get<{ user_info: UserTypes.UserInfo; permissions: string[] }>({
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
    return request.post<AuthTypes.LoginResponse>({
      url: '/admin-api/auth/refresh',
      data: params
    })
  }

  // 获取用户列表
  static getUserList(params: Api.Common.PaginatingSearchParams) {
    return request.get<UserTypes.UserListData>({
      url: '/api/user/list',
      params
    })
  }

  // 更新个人信息
  static updateProfile(data: UserTypes.ProfileUpdateRequest) {
    return request.put<Api.Http.BaseResponse<UserTypes.UserInfo>>({
      url: '/admin-api/profile',
      data: data
    })
  }

  // 修改密码
  static changePassword(data: UserTypes.ChangePasswordRequest) {
    return request.put<Api.Http.BaseResponse<null>>({
      url: '/admin-api/profile/password',
      data: data
    })
  }
}

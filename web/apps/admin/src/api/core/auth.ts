import { baseRequestClient, requestClient } from '#/api/request';

export namespace AuthApi {
  /** 登录接口参数 */
  export interface LoginParams {
    captchaCode: string;
    captchaId: string;
    password: string;
    username: string;
  }

  /** 登录接口返回值 */
  export interface LoginResult {
    token: string;
  }

  /** 图形验证码 */
  export interface CaptchaResult {
    id: string;
    image: string;
  }

  export interface RefreshTokenResult {
    data: string;
    status: number;
  }
}

// 登录。后台用 captchaId 对上验证码，成功只返回 token。
export async function loginApi(data: AuthApi.LoginParams) {
  return requestClient.post<AuthApi.LoginResult>('/auth/login', data);
}

// 获取登录验证码。image 是 base64 图片。
export async function getCaptchaApi() {
  return requestClient.get<AuthApi.CaptchaResult>('/captcha');
}

// 刷新 token。当前后台没有这个接口，保留函数是因为请求层在开启刷新时会调用它。
export async function refreshTokenApi() {
  return baseRequestClient.post<AuthApi.RefreshTokenResult>('/auth/refresh', {
    withCredentials: true,
  });
}

// 退出登录。必须走带 Token 的客户端，后台才能把当前 token 拉黑。
export async function logoutApi() {
  return requestClient.post('/auth/logout');
}

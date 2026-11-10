import {
  authenticateResponseInterceptor,
  isUnauthorizedResponse,
  RequestClient,
} from '@vben/request';
import { describe, expect, it, vi } from 'vitest';

describe('isUnauthorizedResponse', () => {
  it('识别 HTTP 401', () => {
    expect(
      isUnauthorizedResponse({ response: { status: 401, data: { code: 401 } } }),
    ).toBe(true);
  });

  it('识别 HTTP 200 且 body.code 为 401', () => {
    expect(
      isUnauthorizedResponse({
        status: 200,
        data: { code: 401, msg: '请先登录' },
      }),
    ).toBe(true);
  });

  it('忽略其他业务错误码', () => {
    expect(
      isUnauthorizedResponse({
        status: 200,
        data: { code: 500, msg: '服务器内部错误' },
      }),
    ).toBe(false);
  });
});

describe('authenticateResponseInterceptor', () => {
  it('业务码 401 时触发重新认证', async () => {
    const doReAuthenticate = vi.fn(async () => undefined);
    const client = new RequestClient();
    const interceptor = authenticateResponseInterceptor({
      client,
      doReAuthenticate,
      doRefreshToken: vi.fn(),
      enableRefreshToken: false,
      formatToken: (token) => `Bearer ${token}`,
    });

    const error = {
      config: { url: '/auth/current-user' },
      status: 200,
      data: { code: 401, msg: 'token 无效或已过期' },
    };

    await expect(interceptor.rejected!(error)).rejects.toBe(error);
    expect(doReAuthenticate).toHaveBeenCalledTimes(1);
  });
});

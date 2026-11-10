import { beforeEach, describe, expect, it, vi } from 'vitest';

const accessStore = vi.hoisted(() => {
  const store = {
    accessToken: 'stale-token' as null | string,
    refreshToken: 'stale-refresh' as null | string,
    setAccessToken: vi.fn((token: null | string) => {
      store.accessToken = token;
    }),
    setRefreshToken: vi.fn((token: null | string) => {
      store.refreshToken = token;
    }),
  };
  return store;
});

vi.mock('@vben/stores', () => ({
  useAccessStore: () => accessStore,
}));

const { clearPersistedAccessTokens } = await import(
  './clear-persisted-access-token'
);

describe('clearPersistedAccessTokens', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  it('清空当前与历史 namespace 下的持久化 token', () => {
    localStorage.setItem(
      'bico-admin-1.0.0-dev-core-access',
      JSON.stringify({ accessToken: 'old', refreshToken: 'old-r' }),
    );
    localStorage.setItem(
      'bico-admin-2.0.0-prod-core-access',
      JSON.stringify({ accessToken: 'older', refreshToken: 'older-r' }),
    );
    localStorage.setItem('other-app-core-access', JSON.stringify({ accessToken: 'x' }));

    clearPersistedAccessTokens('bico-admin');

    expect(accessStore.setAccessToken).toHaveBeenCalledWith(null);
    expect(JSON.parse(localStorage.getItem('bico-admin-1.0.0-dev-core-access')!)).toMatchObject({
      accessToken: null,
      refreshToken: null,
    });
    expect(JSON.parse(localStorage.getItem('bico-admin-2.0.0-prod-core-access')!)).toMatchObject({
      accessToken: null,
      refreshToken: null,
    });
    expect(JSON.parse(localStorage.getItem('other-app-core-access')!).accessToken).toBe('x');
  });
});

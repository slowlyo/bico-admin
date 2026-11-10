import { useAccessStore } from '@vben/stores';

const ACCESS_STORE_SUFFIX = '-core-access';

/**
 * 清除内存与 localStorage 中的 accessToken，包含历史版本 namespace 遗留的持久化数据。
 */
export function clearPersistedAccessTokens(appNamespace = import.meta.env.VITE_APP_NAMESPACE) {
  const accessStore = useAccessStore();
  accessStore.setAccessToken(null);
  accessStore.setRefreshToken(null);

  if (typeof localStorage === 'undefined' || !appNamespace) {
    return;
  }

  for (let index = 0; index < localStorage.length; index += 1) {
    const key = localStorage.key(index);
    if (!key?.startsWith(appNamespace) || !key.endsWith(ACCESS_STORE_SUFFIX)) {
      continue;
    }
    try {
      const raw = localStorage.getItem(key);
      if (!raw) {
        continue;
      }
      const parsed = JSON.parse(raw) as Record<string, unknown>;
      if (typeof parsed !== 'object' || parsed === null) {
        localStorage.removeItem(key);
        continue;
      }
      parsed.accessToken = null;
      parsed.refreshToken = null;
      localStorage.setItem(key, JSON.stringify(parsed));
    } catch {
      localStorage.removeItem(key);
    }
  }
}

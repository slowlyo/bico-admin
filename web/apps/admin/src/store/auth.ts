import type { UserInfo } from '@vben/types';

import { ref } from 'vue';
import { useRouter } from 'vue-router';

import { LOGIN_PATH } from '@vben/constants';
import { preferences } from '@vben/preferences';
import { resetAllStores, useAccessStore, useUserStore } from '@vben/stores';

import { notification } from 'antdv-next';
import { defineStore } from 'pinia';

import type { AuthApi } from '#/api';
import { getUserInfoApi, loginApi, logoutApi } from '#/api';
import { getAppConfigApi } from '#/api/common';
import { $t } from '#/locales';

export const useAuthStore = defineStore('auth', () => {
  const accessStore = useAccessStore();
  const userStore = useUserStore();
  const router = useRouter();

  const loginLoading = ref(false);
  // 退出请求若返回 401，拦截器会再次调用 logout，这个标记用来打断重入。
  let loggingOut = false;

  // 提交账号密码和验证码。token 写入后再拉当前用户，权限来自用户资料而不是单独的权限接口。
  async function authLogin(
    params: AuthApi.LoginParams,
    onSuccess?: () => Promise<void> | void,
  ) {
    let userInfo: null | UserInfo = null;
    try {
      loginLoading.value = true;
      const { token } = await loginApi(params);
      // 没有 token 时不能进入系统，交给登录页刷新验证码。
      if (!token) {
        throw new Error('登录失败');
      }

      accessStore.setAccessToken(token);
      userInfo = await fetchUserInfo();

      if (accessStore.loginExpired) {
        accessStore.setLoginExpired(false);
      } else {
        onSuccess
          ? await onSuccess?.()
          : await router.push(
              userInfo.homePath || preferences.app.defaultHomePath,
            );
      }

      if (userInfo?.realName) {
        notification.success({
          description: `${$t('authentication.loginSuccessDesc')}:${userInfo?.realName}`,
          duration: 3,
          title: $t('authentication.loginSuccess'),
        });
      }
    } finally {
      loginLoading.value = false;
    }

    return {
      userInfo,
    };
  }

  // 先通知后台拉黑 token，再清本地状态。后台失败也不留下本地登录态。
  async function logout(redirect: boolean = true) {
    if (loggingOut) {
      return;
    }
    loggingOut = true;
    try {
      await logoutApi();
    } catch {
      // 退出接口失败时仍然清本地登录态，避免卡在失效 token 上。
    } finally {
      loggingOut = false;
    }
    resetAllStores();
    accessStore.setLoginExpired(false);

    // 回登录页带上当前路由地址
    await router.replace({
      path: LOGIN_PATH,
      query: redirect
        ? {
            redirect: encodeURIComponent(router.currentRoute.value.fullPath),
          }
        : {},
    });
  }

  // 拉当前用户，并在调试模式补上开发者菜单权限。权限码同时写入按钮鉴权。
  async function fetchUserInfo() {
    const userInfo = await getUserInfoApi();
    const config = await getAppConfigApi();
    // 开发者工具不在角色权限树里，只跟服务端 debug 模式走。
    if (config.debug) {
      userInfo.roles = [...new Set([...(userInfo.roles ?? []), 'developer'])];
    }
    userStore.setUserInfo(userInfo);
    accessStore.setAccessCodes(userInfo.roles ?? []);
    return userInfo;
  }

  function $reset() {
    loginLoading.value = false;
  }

  return {
    $reset,
    authLogin,
    fetchUserInfo,
    loginLoading,
    logout,
  };
});

<script lang="ts" setup>
import type { VbenFormSchema } from '@vben/common-ui';

import { computed, markRaw, onMounted, provide, reactive } from 'vue';

import { AuthenticationLogin, z } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { getCaptchaApi } from '#/api';
import { useAuthStore } from '#/store';

import CaptchaInput from './captcha-input.vue';
import { CAPTCHA_KEY } from './captcha';

defineOptions({ name: 'Login' });

const authStore = useAuthStore();

// 验证码不放进表单 schema。schema 里一旦出现会变的对象，表单会整表重挂，验证码接口会被反复调用。
const captcha = reactive({ id: '', image: '' });

const formSchema = computed((): VbenFormSchema[] => {
  return [
    {
      component: 'VbenInput',
      componentProps: {
        placeholder: $t('authentication.usernameTip'),
      },
      fieldName: 'username',
      label: $t('authentication.username'),
      rules: z.string().min(1, { message: $t('authentication.usernameTip') }),
    },
    {
      component: 'VbenInputPassword',
      componentProps: {
        placeholder: $t('authentication.password'),
      },
      fieldName: 'password',
      label: $t('authentication.password'),
      rules: z.string().min(1, { message: $t('authentication.passwordTip') }),
    },
    {
      component: markRaw(CaptchaInput),
      fieldName: 'captchaCode',
      label: '验证码',
      rules: z.string().min(1, { message: '请输入验证码' }),
    },
  ];
});

// 旧验证码校验后即作废，所以进入页面和登录失败都要重新取。
async function loadCaptcha() {
  const data = await getCaptchaApi();
  captcha.id = data.id;
  captcha.image = data.image.startsWith('data:')
    ? data.image
    : `data:image/png;base64,${data.image}`;
}

provide(CAPTCHA_KEY, { loadCaptcha, state: captcha });

// 把验证码 id 和表单里的验证码一起提交。失败后换图，成功则由 store 跳转。
async function handleLogin(values: Record<string, any>) {
  try {
    await authStore.authLogin({
      captchaCode: values.captchaCode,
      captchaId: captcha.id,
      password: values.password,
      username: values.username,
    });
  } catch {
    await loadCaptcha();
  }
}

onMounted(() => {
  void loadCaptcha();
});
</script>

<template>
  <AuthenticationLogin
    :form-schema="formSchema"
    :loading="authStore.loginLoading"
    :show-code-login="false"
    :show-forget-password="false"
    :show-qrcode-login="false"
    :show-register="false"
    :show-third-party-login="false"
    @submit="handleLogin"
  />
</template>

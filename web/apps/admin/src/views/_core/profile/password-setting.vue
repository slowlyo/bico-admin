<script setup lang="ts">
import type { VbenFormSchema } from '#/adapter/form';

import { computed } from 'vue';
import { useRouter } from 'vue-router';

import { ProfilePasswordSetting, z } from '@vben/common-ui';
import { LOGIN_PATH } from '@vben/constants';
import { resetAllStores } from '@vben/stores';

import { message } from 'antdv-next';

import { changePasswordApi } from '#/api';

const router = useRouter();

const formSchema = computed((): VbenFormSchema[] => {
  return [
    {
      component: 'VbenInputPassword',
      componentProps: {
        placeholder: '请输入旧密码',
      },
      fieldName: 'oldPassword',
      label: '旧密码',
      rules: z.string().min(1, { message: '请输入旧密码' }),
    },
    {
      component: 'VbenInputPassword',
      componentProps: {
        passwordStrength: true,
        placeholder: '请输入新密码',
      },
      fieldName: 'newPassword',
      label: '新密码',
      // 后台 binding 要求新密码至少 8 位，短密码会直接被拒绝。
      rules: z.string().min(8, { message: '新密码至少 8 位' }),
    },
    {
      component: 'VbenInputPassword',
      componentProps: {
        passwordStrength: true,
        placeholder: '请再次输入新密码',
      },
      dependencies: {
        rules(values) {
          const { newPassword } = values;
          return z
            .string({ required_error: '请再次输入新密码' })
            .min(1, { message: '请再次输入新密码' })
            .refine((value) => value === newPassword, {
              message: '两次输入的密码不一致',
            });
        },
        triggerFields: ['newPassword'],
      },
      fieldName: 'confirmPassword',
      label: '确认密码',
    },
  ];
});

// 改密成功后当前 token 版本已作废。不再请求退出接口，避免再弹出「token 已失效」。
async function handleSubmit(values: Record<string, string | undefined>) {
  const oldPassword = values.oldPassword;
  const newPassword = values.newPassword;
  // 表单校验过了字段仍可能是空。空密码不能发给最少 8 位的接口。
  if (!oldPassword || !newPassword) {
    return;
  }
  await changePasswordApi({
    newPassword,
    oldPassword,
  });
  message.success('密码修改成功');
  resetAllStores();
  await router.replace({ path: LOGIN_PATH });
}
</script>
<template>
  <ProfilePasswordSetting
    class="w-1/3"
    :form-schema="formSchema"
    @submit="handleSubmit"
  />
</template>

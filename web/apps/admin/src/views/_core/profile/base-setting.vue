<script setup lang="ts">
import type { VbenFormSchema } from '#/adapter/form';

import { computed, onMounted, ref } from 'vue';

import { ProfileBaseSetting, z } from '@vben/common-ui';
import { preferences } from '@vben/preferences';
import { useUserStore } from '@vben/stores';

import { message } from 'antdv-next';

import { getUserInfoApi, updateProfileApi, uploadAvatarApi } from '#/api';
import AvatarUpload from '#/components/avatar-upload.vue';

const userStore = useUserStore();
const profileBaseSettingRef = ref();

const avatarSrc = computed(
  () => userStore.userInfo?.avatar || preferences.app.defaultAvatar,
);

const formSchema = computed((): VbenFormSchema[] => {
  return [
    {
      component: 'Input',
      fieldName: 'realName',
      label: '姓名',
      rules: z.string().min(1, { message: '请输入姓名' }),
    },
    {
      component: 'Input',
      fieldName: 'username',
      label: '用户名',
      rules: z.string().min(1, { message: '请输入用户名' }),
    },
  ];
});

// 只回填后台会保存的姓名和用户名，角色和简介没有对应接口。
async function loadProfile() {
  const data = await getUserInfoApi();
  userStore.setUserInfo(data);
  profileBaseSettingRef.value.getFormApi().setValues({
    realName: data.realName,
    username: data.username,
  });
}

// 上传接口不改用户字段。资料接口写回头像后，侧栏才能用返回的用户刷新。
async function saveAvatar(file: File) {
  const user = await uploadAvatarApi(file);
  userStore.setUserInfo(user);
  return user.avatar || '';
}

// 提交后用返回的用户覆盖本地资料，侧栏和顶栏不用再等下一次刷新。
async function handleSubmit(values: Record<string, string>) {
  const user = await updateProfileApi({
    name: values.realName,
    username: values.username,
  });
  userStore.setUserInfo(user);
  message.success('更新成功');
}

onMounted(() => {
  void loadProfile();
});
</script>
<template>
  <div>
    <div class="mb-6 flex justify-center">
      <AvatarUpload :model-value="avatarSrc" :on-upload="saveAvatar" :size="104" />
    </div>
    <ProfileBaseSetting
      ref="profileBaseSettingRef"
      :form-schema="formSchema"
      @submit="handleSubmit"
    />
  </div>
</template>
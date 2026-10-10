<script lang="ts" setup>
import type { AdminUser } from '#/api/system/user';

import { useVbenModal } from '@vben/common-ui';

import { message } from 'antdv-next';

import { useVbenForm, z } from '#/adapter/form';
import AvatarUpload from '#/components/avatar-upload.vue';
import { listAllAdminRolesApi } from '#/api/system/role';
import { createAdminUserApi, updateAdminUserApi, uploadUserAvatarApi } from '#/api/system/user';

const emit = defineEmits<{ success: [] }>();

const [Form, formApi] = useVbenForm({
  schema: [
    {
      component: 'Switch',
      dependencies: {
        show: false,
        triggerFields: ['username'],
      },
      fieldName: 'editing',
      label: '编辑',
    },
    {
      component: 'Input',
      componentProps: { class: 'w-full', placeholder: '请输入用户名' },
      fieldName: 'username',
      label: '用户名',
      rules: z.string().trim().min(1, '请输入用户名').max(64, '用户名不能超过64位'),
    },
    {
      component: 'InputPassword',
      componentProps: { placeholder: '请输入密码' },
      dependencies: {
        componentProps(values) {
          return {
            class: 'w-full',
            placeholder: values.editing ? '不修改请留空' : '请输入密码',
          };
        },
        required: (values) => !values.editing,
        triggerFields: ['editing'],
      },
      fieldName: 'password',
      label: '密码',
      rules: z
        .string()
        .min(8, '密码长度至少8位')
        .optional()
        .or(z.literal('')),
    },
    {
      component: 'Input',
      componentProps: { class: 'w-full', placeholder: '请输入姓名' },
      fieldName: 'name',
      label: '姓名',
    },
    {
      component: 'Input',
      fieldName: 'avatar',
      label: '头像',
    },
    {
      component: 'ApiSelect',
      componentProps: {
        api: listAllAdminRolesApi,
        class: 'w-full',
        labelField: 'name',
        mode: 'multiple',
        valueField: 'id',
      },
      fieldName: 'role_ids',
      label: '角色',
    },
    {
      component: 'Switch',
      defaultValue: true,
      fieldName: 'enabled',
      label: '状态',
    },
  ],
  showDefaultActions: false,
});

// 只拿上传地址。用户点确认时才随资料提交，不能顺手改掉当前登录用户。
async function uploadAvatar(file: File) {
  const uploaded = await uploadUserAvatarApi(file);
  return uploaded.url;
}

const [Modal, modalApi] = useVbenModal({
  async onConfirm() {
    const { valid } = await formApi.validate();
    if (!valid) {
      return;
    }
    const values = await formApi.getValues();
    const editing = Boolean(values.editing);
    // 编辑时密码留空表示不改，不能把空字符串传给最少 8 位校验。
    const payload = {
      avatar: values.avatar || '',
      enabled: Boolean(values.enabled),
      name: String(values.name || '').trim(),
      role_ids: (values.role_ids || []) as number[],
      username: String(values.username || '').trim(),
      ...(values.password ? { password: String(values.password) } : {}),
    };
    if (!editing && !payload.password) {
      message.error('请输入密码');
      return;
    }

    modalApi.lock();
    try {
      const data = modalApi.getData<AdminUser>();
      if (editing && data?.id) {
        await updateAdminUserApi(data.id, payload);
        message.success('更新成功');
      } else {
        await createAdminUserApi(payload);
        message.success('创建成功');
      }
      emit('success');
      modalApi.close();
    } finally {
      modalApi.unlock();
    }
  },
  async onOpenChange(isOpen) {
    if (!isOpen) {
      return;
    }
    const data = modalApi.getData<AdminUser>();
    const editing = Boolean(data?.id);
    modalApi.setState({ title: editing ? '编辑用户' : '新建用户' });
    await formApi.resetForm();
    await formApi.setValues({
      avatar: editing
        ? data?.avatar
        : `https://api.dicebear.com/9.x/thumbs/png?seed=${Math.random()}`,
      editing,
      enabled: data?.enabled ?? true,
      name: data?.name ?? '',
      password: '',
      role_ids: data?.roles?.map((role) => role.id) ?? [],
      username: data?.username ?? '',
    });
  },
});
</script>

<template>
  <Modal>
    <Form>
      <template #avatar="slotProps">
        <AvatarUpload
          :model-value="slotProps.value"
          :on-upload="uploadAvatar"
          @update:model-value="(url) => formApi.setFieldValue('avatar', url)"
        />
      </template>
    </Form>
  </Modal>
</template>

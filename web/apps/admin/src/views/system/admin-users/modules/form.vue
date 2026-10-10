<script lang="ts" setup>
import type { AdminUser } from '#/api/system/user';

import { message } from 'antdv-next';

import { z } from '#/adapter/form';
import { listAllAdminRolesApi } from '#/api/system/role';
import { adminUserApi, uploadUserAvatarApi } from '#/api/system/user';
import AvatarUpload from '#/components/avatar-upload.vue';
import { useCrudForm } from '#/hooks/use-crud-form';

const emit = defineEmits<{ success: [] }>();

// 把表单值整理成后端载荷；密码单独处理，因为创建必填、编辑可留空。
function toPayload(values: Record<string, any>) {
  return {
    avatar: values.avatar || '',
    enabled: Boolean(values.enabled),
    name: String(values.name || '').trim(),
    role_ids: (values.role_ids || []) as number[],
    username: String(values.username || '').trim(),
  };
}

// 只拿上传地址。用户点确认时才随资料提交，不能顺手改掉当前登录用户。
async function uploadAvatar(file: File) {
  const uploaded = await uploadUserAvatarApi(file);
  return uploaded.url;
}

const { Form, formApi, Modal } = useCrudForm<AdminUser>({
  create: async (values) => {
    const password = values.password ? String(values.password) : '';
    if (!password) {
      message.error('请输入密码');
      return false;
    }
    return adminUserApi.create({ ...toPayload(values), password });
  },
  entityName: '用户',
  onSuccess: () => emit('success'),
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
  toValues: (row) => ({
    avatar: row
      ? row.avatar
      : `https://api.dicebear.com/9.x/thumbs/png?seed=${Math.random()}`,
    editing: Boolean(row),
    enabled: row?.enabled ?? true,
    name: row?.name ?? '',
    password: '',
    role_ids: row?.roles?.map((role) => role.id) ?? [],
    username: row?.username ?? '',
  }),
  // 编辑时密码留空表示不改，不能把空字符串传给最少 8 位校验。
  update: (row, values) =>
    adminUserApi.update(row.id, {
      ...toPayload(values),
      ...(values.password ? { password: String(values.password) } : {}),
    }),
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

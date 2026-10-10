<script lang="ts" setup>
import type { AdminRole } from '#/api/system/role';

import { useVbenModal } from '@vben/common-ui';

import { message } from 'antdv-next';

import { useVbenForm, z } from '#/adapter/form';
import { createAdminRoleApi, updateAdminRoleApi } from '#/api/system/role';

const emit = defineEmits<{ success: [] }>();

const [Form, formApi] = useVbenForm({
  schema: [
    {
      component: 'Input',
      componentProps: { class: 'w-full', placeholder: '请输入角色名称' },
      fieldName: 'name',
      label: '角色名称',
      rules: z.string().trim().min(1, '请输入角色名称'),
    },
    {
      component: 'Textarea',
      componentProps: { class: 'w-full', placeholder: '请输入描述' },
      fieldName: 'description',
      label: '描述',
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

const [Modal, modalApi] = useVbenModal({
  async onConfirm() {
    const { valid } = await formApi.validate();
    if (!valid) {
      return;
    }
    const values = await formApi.getValues();
    const payload = {
      description: String(values.description || '').trim(),
      enabled: Boolean(values.enabled),
      name: String(values.name || '').trim(),
    };
    const data = modalApi.getData<AdminRole>();
    modalApi.lock();
    try {
      if (data?.id) {
        await updateAdminRoleApi(data.id, payload);
        message.success('更新成功');
      } else {
        await createAdminRoleApi(payload);
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
    const data = modalApi.getData<AdminRole>();
    modalApi.setState({ title: data?.id ? '编辑角色' : '新建角色' });
    await formApi.resetForm();
    await formApi.setValues({
      description: data?.description ?? '',
      enabled: data?.enabled ?? true,
      name: data?.name ?? '',
    });
  },
});
</script>

<template>
  <Modal>
    <Form />
  </Modal>
</template>

<script lang="ts" setup>
import type { AdminRole } from '#/api/system/role';

import { z } from '#/adapter/form';
import { adminRoleApi } from '#/api/system/role';
import { useCrudForm } from '#/hooks/use-crud-form';

const emit = defineEmits<{ success: [] }>();

// 把表单值整理成后端载荷，名称和描述去除首尾空格。
function toPayload(values: Record<string, any>) {
  return {
    description: String(values.description || '').trim(),
    enabled: Boolean(values.enabled),
    name: String(values.name || '').trim(),
  };
}

const { Form, Modal } = useCrudForm<AdminRole>({
  create: (values) => adminRoleApi.create(toPayload(values)),
  entityName: '角色',
  onSuccess: () => emit('success'),
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
  toValues: (row) => ({
    description: row?.description ?? '',
    enabled: row?.enabled ?? true,
    name: row?.name ?? '',
  }),
  update: (row, values) => adminRoleApi.update(row.id, toPayload(values)),
});
</script>

<template>
  <Modal>
    <Form />
  </Modal>
</template>

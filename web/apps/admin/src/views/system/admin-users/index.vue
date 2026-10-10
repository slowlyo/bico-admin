<script lang="ts" setup>
import type { AdminUser, AdminUserQuery } from '#/api/system/user';

import { Page, useVbenModal } from '@vben/common-ui';

import { Avatar, Button, Popconfirm, Tag } from 'antdv-next';

import { listAllAdminRolesApi } from '#/api/system/role';
import { adminUserApi } from '#/api/system/user';
import { useCrudGrid } from '#/hooks/use-crud-grid';

import UserForm from './modules/form.vue';

const [FormModal, formModalApi] = useVbenModal({
  connectedComponent: UserForm,
  destroyOnClose: true,
});

const { Grid, gridApi, remove, removeSelected } = useCrudGrid<
  AdminUser,
  AdminUserQuery
>({
  api: adminUserApi,
  entityName: '用户',
  formOptions: {
    schema: [
      {
        component: 'Input',
        fieldName: 'username',
        label: '用户名',
      },
      {
        component: 'Input',
        fieldName: 'name',
        label: '姓名',
      },
      {
        component: 'ApiSelect',
        componentProps: {
          api: listAllAdminRolesApi,
          labelField: 'name',
          mode: 'multiple',
          valueField: 'id',
        },
        fieldName: 'role_ids',
        label: '角色',
      },
      {
        component: 'Select',
        componentProps: {
          allowClear: true,
          options: [
            { label: '启用', value: 'true' },
            { label: '禁用', value: 'false' },
          ],
        },
        fieldName: 'enabled',
        label: '状态',
      },
    ],
    showCollapseButton: true,
  },
  gridOptions: {
    columns: [
      { type: 'checkbox', width: 48 },
      { field: 'id', title: 'ID', width: 80 },
      { field: 'username', minWidth: 140, title: '用户名' },
      { field: 'avatar', slots: { default: 'avatar' }, title: '头像', width: 80 },
      { field: 'name', minWidth: 140, title: '姓名' },
      { field: 'roles', minWidth: 200, slots: { default: 'roles' }, title: '角色' },
      { field: 'enabled', slots: { default: 'enabled' }, title: '状态', width: 100 },
      { field: 'created_at', sortable: true, title: '创建时间', width: 180 },
      {
        field: 'action',
        fixed: 'right',
        slots: { default: 'action' },
        title: '操作',
        width: 140,
      },
    ],
  },
});

// 新建不带行数据，表单据此走创建分支。
function openCreate() {
  formModalApi.setData({}).open();
}

// 编辑把当前行交给弹窗，角色 id 在表单里从 roles 取出。
function openEdit(row: Record<string, any>) {
  formModalApi.setData(row).open();
}
</script>

<template>
  <Page auto-content-height>
    <FormModal @success="gridApi.query()" />
    <Grid>
      <template #toolbar-tools>
        <Button
          v-access:code="'system:admin_user:create'"
          type="primary"
          @click="openCreate"
        >
          新建
        </Button>
        <Popconfirm title="确认删除选中用户？" @confirm="removeSelected">
          <Button
            v-access:code="'system:admin_user:delete'"
            class="ml-2"
            danger
          >
            批量删除
          </Button>
        </Popconfirm>
      </template>
      <template #avatar="{ row }">
        <Avatar :size="32" :src="row.avatar" />
      </template>
      <template #roles="{ row }">
        <Tag v-for="role in row.roles" :key="role.id">{{ role.name }}</Tag>
      </template>
      <template #enabled="{ row }">
        <Tag :color="row.enabled ? 'green' : 'red'">
          {{ row.enabled ? '启用' : '禁用' }}
        </Tag>
      </template>
      <template #action="{ row }">
        <Button
          v-access:code="'system:admin_user:edit'"
          size="small"
          type="link"
          @click="openEdit(row)"
        >
          编辑
        </Button>
        <Popconfirm title="确认删除该用户？" @confirm="remove(row)">
          <Button
            v-access:code="'system:admin_user:delete'"
            danger
            size="small"
            type="link"
          >
            删除
          </Button>
        </Popconfirm>
      </template>
    </Grid>
  </Page>
</template>

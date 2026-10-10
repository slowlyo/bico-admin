<script lang="ts" setup>
import type { AdminRole, AdminRoleQuery } from '#/api/system/role';

import { Page, useVbenDrawer, useVbenModal } from '@vben/common-ui';

import { Button, Popconfirm, Tag, Tooltip } from 'antdv-next';

import { adminRoleApi } from '#/api/system/role';
import { useCrudGrid } from '#/hooks/use-crud-grid';

import RoleForm from './modules/form.vue';
import RolePermission from './modules/permission.vue';

const [FormModal, formModalApi] = useVbenModal({
  connectedComponent: RoleForm,
  destroyOnClose: true,
});

const [PermissionDrawer, permissionDrawerApi] = useVbenDrawer({
  connectedComponent: RolePermission,
  destroyOnClose: true,
});

// 系统保留角色不可删，批量删除也据此跳过。
const canDelete = (row: AdminRole) => !row.system;

const { Grid, gridApi, remove, removeSelected } = useCrudGrid<
  AdminRole,
  AdminRoleQuery
>({
  api: adminRoleApi,
  canDelete,
  entityName: '角色',
  formOptions: {
    schema: [
      {
        component: 'Input',
        fieldName: 'name',
        label: '角色名称',
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
    showCollapseButton: false,
  },
  gridOptions: {
    checkboxConfig: {
      checkMethod: ({ row }: { row: AdminRole }) => canDelete(row),
      highlight: true,
    },
    columns: [
      { type: 'checkbox', width: 48 },
      { field: 'id', title: 'ID', width: 80 },
      { field: 'name', slots: { default: 'name' }, title: '角色名称', minWidth: 200 },
      { field: 'description', minWidth: 180, title: '描述' },
      { field: 'enabled', slots: { default: 'enabled' }, title: '状态', width: 100 },
      { field: 'created_at', sortable: true, title: '创建时间', width: 180 },
      {
        field: 'action',
        fixed: 'right',
        slots: { default: 'action' },
        title: '操作',
        width: 220,
      },
    ],
  },
});

// 新建角色不带权限，权限在单独抽屉里配。
function openCreate() {
  formModalApi.setData({}).open();
}

// 系统保留角色不能改名称和状态。
function openEdit(row: Record<string, any>) {
  if (row.system) {
    return;
  }
  formModalApi.setData(row).open();
}

// 超级管理员不打开抽屉，权限由后台固定为全部。
function openPermissions(row: Record<string, any>) {
  if (row.system) {
    return;
  }
  permissionDrawerApi.setData(row).open();
}
</script>

<template>
  <Page auto-content-height>
    <FormModal @success="gridApi.query()" />
    <PermissionDrawer @success="gridApi.query()" />
    <Grid>
      <template #toolbar-tools>
        <Button
          v-access:code="'system:admin_role:create'"
          type="primary"
          @click="openCreate"
        >
          新建
        </Button>
        <Popconfirm title="确认删除选中角色？" @confirm="removeSelected">
          <Button
            v-access:code="'system:admin_role:delete'"
            class="ml-2"
            danger
          >
            批量删除
          </Button>
        </Popconfirm>
      </template>
      <template #name="{ row }">
        <span class="inline-flex items-center gap-1">
          <span>{{ row.name }}</span>
          <Tag v-if="row.system">系统保留</Tag>
        </span>
      </template>
      <template #enabled="{ row }">
        <Tag :color="row.enabled ? 'green' : 'red'">
          {{ row.enabled ? '启用' : '禁用' }}
        </Tag>
      </template>
      <template #action="{ row }">
        <Tooltip v-if="row.system" title="超级管理员自动拥有全部权限">
          <Button
            v-access:code="'system:admin_role:permission'"
            disabled
            size="small"
            type="link"
          >
            配置权限
          </Button>
        </Tooltip>
        <Button
          v-else
          v-access:code="'system:admin_role:permission'"
          size="small"
          type="link"
          @click="openPermissions(row)"
        >
          配置权限
        </Button>
        <Tooltip v-if="row.system" title="系统保留角色不可修改">
          <Button
            v-access:code="'system:admin_role:edit'"
            disabled
            size="small"
            type="link"
          >
            编辑
          </Button>
        </Tooltip>
        <Button
          v-else
          v-access:code="'system:admin_role:edit'"
          size="small"
          type="link"
          @click="openEdit(row)"
        >
          编辑
        </Button>
        <Tooltip v-if="row.system" title="系统保留角色不可删除">
          <Button
            v-access:code="'system:admin_role:delete'"
            danger
            disabled
            size="small"
            type="link"
          >
            删除
          </Button>
        </Tooltip>
        <Popconfirm v-else title="确认删除该角色？" @confirm="remove(row)">
          <Button
            v-access:code="'system:admin_role:delete'"
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

<script lang="ts" setup>
import type { AdminRole } from '#/api/system/role';

import { Page, useVbenDrawer, useVbenModal } from '@vben/common-ui';

import { Button, Popconfirm, Tag, Tooltip, message } from 'antdv-next';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  deleteAdminRoleApi,
  deleteAdminRolesApi,
  listAdminRolesApi,
} from '#/api/system/role';

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

const [Grid, gridApi] = useVbenVxeGrid({
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
      checkMethod: ({ row }: { row: AdminRole }) => !row.system,
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
    height: 'auto',
    proxyConfig: {
      ajax: {
        query: async (
          {
            page,
            sort,
          }: {
            page: { currentPage: number; pageSize: number };
            sort?: { field?: string; order?: string };
          },
          formValues: Record<string, any>,
        ) => {
          const result = await listAdminRolesApi({
            enabled:
              formValues.enabled === 'true'
                ? true
                : formValues.enabled === 'false'
                  ? false
                  : undefined,
            name: formValues.name || undefined,
            page: page.currentPage,
            pageSize: page.pageSize,
            sortField: sort?.field || undefined,
            sortOrder:
              sort?.order === 'asc'
                ? 'ascend'
                : sort?.order === 'desc'
                  ? 'descend'
                  : undefined,
          });
          return { items: result.list ?? [], total: result.total ?? 0 };
        },
      },
    },
    rowConfig: { keyField: 'id' },
    toolbarConfig: { custom: true, refresh: true, search: true, zoom: true },
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

// 系统保留角色不可删。
async function removeRole(row: Record<string, any>) {
  if (row.system) {
    return;
  }
  await deleteAdminRoleApi(row.id);
  message.success('删除成功');
  await gridApi.query();
}

// 批量删除跳过系统角色，避免一次请求被保留角色拒绝。
async function removeSelected() {
  const rows = ((gridApi.grid?.getCheckboxRecords?.() ?? []) as AdminRole[]).filter(
    (row) => !row.system,
  );
  if (rows.length === 0) {
    message.warning('请选择角色');
    return;
  }
  await deleteAdminRolesApi(rows.map((row) => row.id));
  message.success('删除成功');
  await gridApi.query();
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
        <Popconfirm v-else title="确认删除该角色？" @confirm="removeRole(row)">
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

# 前端 CRUD 开发指南

> 技术栈：Vue 3 + Vben Admin + antdv-next + vxe-table。参考实现：`web/apps/admin/src/views/system/admin-users/`。

一个标准 CRUD 页面由三部分组成：

```
web/apps/admin/src/
├── api/system/user.ts                      # 接口层（createCrudApi，见 frontend-services.md）
├── hooks/use-crud-grid.ts                  # 列表页通用封装
├── hooks/use-crud-form.ts                  # 新建/编辑弹窗通用封装
├── views/system/admin-users/
│   ├── index.vue                           # 列表页：搜索 + 表格 + 操作
│   └── modules/form.vue                    # 新建/编辑弹窗
└── router/routes/modules/system.ts         # 路由与菜单
```

## 1. 路由

在 `router/routes/modules/` 下新增或修改路由文件，`authority` 对应后端权限码：

```ts
{
  name: 'AdminUsers',
  path: '/system/admin-users',
  component: () => import('#/views/system/admin-users/index.vue'),
  meta: {
    authority: ['system:admin_user:menu'],
    icon: 'lucide:users',
    title: '用户管理',
  },
}
```

当前 `accessMode` 为 `frontend`：前端按 `meta.authority` 与登录用户的权限列表过滤菜单，无需后端下发菜单。

## 2. 列表页（index.vue）

核心是 `useCrudGrid`（`hooks/use-crud-grid.ts`），封装了分页、服务端排序、搜索表单映射、单删与批量删除。页面只需声明搜索表单和列：

```vue
<script lang="ts" setup>
import { Page, useVbenModal } from '@vben/common-ui';
import { adminUserApi } from '#/api/system/user';
import { useCrudGrid } from '#/hooks/use-crud-grid';
import UserForm from './modules/form.vue';

const [FormModal, formModalApi] = useVbenModal({
  connectedComponent: UserForm,
  destroyOnClose: true,
});

const { Grid, gridApi, remove, removeSelected } = useCrudGrid<AdminUser, AdminUserQuery>({
  api: adminUserApi,
  entityName: '用户',
  formOptions: {
    schema: [{ component: 'Input', fieldName: 'username', label: '用户名' }],
  },
  gridOptions: {
    columns: [
      { type: 'checkbox', width: 48 },
      { field: 'username', title: '用户名', minWidth: 140 },
      { field: 'created_at', sortable: true, title: '创建时间', width: 180 },
    ],
  },
});
</script>
```

模板中：`<Popconfirm @confirm="remove(row)">` 删除单行，`@confirm="removeSelected"` 批量删除，成功后会自动刷新。

`useCrudGrid` 已处理：

- 分页：后端 `{ list, total }` 转为 vxe 的 `{ items, total }`。
- 排序：开启 `proxyConfig.sort` 与 `sortConfig.remote`，`asc/desc` 映射为 `ascend/descend`，参数名 `sortField`、`sortOrder`；可排序列须在后端 `SortableFields` 白名单内。
- 表单值：空值不传，`enabled` 的 `'true'/'false'` 转布尔，数组转逗号分隔。
- 默认配置：`rowConfig.keyField = 'id'`、复选框高亮、`height: 'auto'` 和工具栏；`gridOptions` 同名字段会覆盖默认值。
- 不可删除的行：传 `canDelete: (row) => boolean`，单删忽略这类行，批量删除自动跳过；复选框是否可选仍需在 `checkboxConfig.checkMethod` 里配置。

需要新建/编辑弹窗时，页面自己用 `formModalApi.setData(row).open()`，成功后调用 `gridApi.query()`。

## 3. 表单弹窗（modules/form.vue）

使用 `useCrudForm`（`hooks/use-crud-form.ts`），新建与编辑共用一个弹窗，它负责打开回填、标题、校验、提交锁定、成功提示和关闭：

```vue
<script lang="ts" setup>
import type { AdminRole } from '#/api/system/role';

import { z } from '#/adapter/form';
import { adminRoleApi } from '#/api/system/role';
import { useCrudForm } from '#/hooks/use-crud-form';

const emit = defineEmits<{ success: [] }>();

function toPayload(values: Record<string, any>) {
  return { name: String(values.name || '').trim(), enabled: Boolean(values.enabled) };
}

const { Form, Modal } = useCrudForm<AdminRole>({
  create: (values) => adminRoleApi.create(toPayload(values)),
  entityName: '角色',
  onSuccess: () => emit('success'),
  schema: [
    { component: 'Input', fieldName: 'name', label: '角色名称', rules: z.string().trim().min(1, '请输入角色名称') },
    { component: 'Switch', defaultValue: true, fieldName: 'enabled', label: '状态' },
  ],
  toValues: (row) => ({ name: row?.name ?? '', enabled: row?.enabled ?? true }),
  update: (row, values) => adminRoleApi.update(row.id, toPayload(values)),
});
</script>

<template>
  <Modal><Form /></Modal>
</template>
```

要点：

- 列表页 `setData({})` 表示新建，`setData(row)` 表示编辑；`toValues(row)` 在新建时 `row` 为空，在这里给默认值。
- `create` / `update` 返回 `false` 表示前置校验未通过（调用方自己提示原因），弹窗保持打开。
- 编辑时可选字段（如密码）留空不要提交，避免触发后端校验。
- 需要自定义字段插槽（如头像上传）时，拿返回的 `formApi`，在 `<Form>` 里写 `<template #字段名>`，参考 `admin-users/modules/form.vue`。

列表页通过 `@success="gridApi.query()"` 接收刷新信号：

```vue
<FormModal @success="gridApi.query()" />
```

## 4. 权限控制

按钮级权限使用 `v-access:code` 指令，权限码与后端 `crud.NewCRUDPerms` 生成的一致：

```vue
<Button v-access:code="'system:admin_user:create'" type="primary" @click="openCreate">
  新建
</Button>
```

权限码格式：`{namespace}:{module}:{menu|list|create|edit|delete}`。

## 5. 新增页面检查清单

1. 后端模块已注册，且权限已加入权限树。
2. `api/` 下用 `createCrudApi` 新增接口文件。
3. `views/` 下新增 `index.vue`（用 `useCrudGrid`）与 `modules/form.vue`。
4. `router/routes/modules/` 下新增路由，`authority` 使用 `:menu` 权限。
5. 在「角色管理」中为测试角色勾选新权限后验证。

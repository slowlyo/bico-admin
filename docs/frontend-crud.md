# 前端 CRUD 开发指南

> 使用 `CrudTable` 和 `createCrudService`，一个页面只需 ~100 行代码。

## 快速开始

### 最小示例

```tsx
import type { ProColumns } from '@ant-design/pro-components';
import { ProFormText, ProFormSwitch } from '@ant-design/pro-components';
import { CrudTable } from '@/components';
import { createCrudService } from '@/services/crud';

// 1. 定义类型
interface Article {
  id: number;
  title: string;
  content: string;
  enabled: boolean;
}

// 2. 创建服务（一行搞定 CRUD API）
const articleService = createCrudService<Article>('/articles');

// 3. 定义列
const columns: ProColumns<Article>[] = [
  { title: 'ID', dataIndex: 'id', width: 80, search: false },
  { title: '标题', dataIndex: 'title' },
  { title: '状态', dataIndex: 'enabled', valueType: 'switch' },
];

// 4. 导出页面
export default () => (
  <CrudTable<Article>
    title="文章"
    permissionPrefix="content:article"
    service={articleService}
    columns={columns}
    formContent={
      <>
        <ProFormText name="title" label="标题" rules={[{ required: true }]} />
        <ProFormText name="content" label="内容" />
        <ProFormSwitch name="enabled" label="状态" initialValue={true} />
      </>
    }
  />
);
```

**完成！** 自动支持：列表、搜索、分页、新建、编辑、删除、权限控制。

---

## 核心组件

### createCrudService

快速生成标准 CRUD API：

```ts
import { createCrudService } from '@/services/crud';

// 基本用法
const service = createCrudService<Article>('/articles');

// 生成的方法：
service.list(params)        // GET /articles
service.get(id)             // GET /articles/:id
service.create(data)        // POST /articles
service.update(id, data)    // PUT /articles/:id
service.delete(id)          // DELETE /articles/:id
```

### CrudTable

封装了 ProTable + 弹窗表单的完整 CRUD 页面：

```tsx
<CrudTable<T>
  // 必填
  title="文章"                          // 模块名称
  permissionPrefix="content:article"   // 权限前缀
  service={articleService}             // CRUD 服务
  columns={columns}                    // 列配置
  formContent={<FormFields />}         // 表单字段

  // 可选
  recordToValues={(r) => ({...})}      // 编辑时转换初始值
  transformParams={(p) => ({...})}     // 自定义请求参数
  rowKey="id"                          // 行 key，默认 "id"
  scrollX={1200}                       // 横向滚动宽度
  showCreate={true}                    // 是否显示新建按钮
  showDeleteConfirm={true}             // 是否显示删除确认
  toolBarExtra={[<Button />]}          // 额外工具栏按钮
  renderActions={(record, defaults) => ...}  // 自定义操作列
/>
```

### CrudModal

统一的创建/编辑弹窗（CrudTable 内部使用，也可单独使用）：

```tsx
<CrudModal<Article>
  title="文章"
  open={modalOpen}
  onOpenChange={setModalOpen}
  record={currentRow}                  // 有值=编辑，无值=创建
  onCreate={service.create}
  onUpdate={service.update}
  onSuccess={() => reload()}
  recordToValues={(r) => ({...})}
>
  <ProFormText name="title" label="标题" />
</CrudModal>
```

---

## 权限控制

CrudTable 自动处理权限，只需配置 `permissionPrefix`：

```tsx
permissionPrefix="system:admin_user"
```

自动生成的权限检查：
- `system:admin_user:create` - 新建按钮
- `system:admin_user:edit` - 编辑按钮
- `system:admin_user:delete` - 删除按钮

---

## 表单字段复用

创建和编辑共用同一套表单字段，可通过 `record` 判断模式：

```tsx
const FormContent: React.FC<{ record?: Article }> = ({ record }) => {
  const isEdit = !!record;
  
  return (
    <>
      {/* 仅创建时显示 */}
      {!isEdit && <ProFormText name="code" label="编码" />}
      
      {/* 创建和编辑都显示 */}
      <ProFormText name="title" label="标题" />
      <ProFormSwitch name="enabled" label="状态" />
    </>
  );
};
```

---

## 完整示例

用户管理页面（约 120 行 vs 原 450+ 行）：

```tsx
import type { ProColumns } from '@ant-design/pro-components';
import { ProFormText, ProFormSwitch, ProFormSelect } from '@ant-design/pro-components';
import { Avatar, Tag, Space, Upload, Button, message } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import React, { useState, useEffect } from 'react';
import { CrudTable } from '@/components';
import { createCrudService } from '@/services/crud';
import { getAllAdminRoles } from '@/services/system/admin-role';
import { buildApiUrl } from '@/services/config';

interface AdminUser {
  id: number;
  username: string;
  name: string;
  avatar: string;
  enabled: boolean;
  roles?: { id: number; name: string }[];
  created_at: string;
}

const userService = createCrudService<AdminUser>('/admin-users');

const columns: ProColumns<AdminUser>[] = [
  { title: 'ID', dataIndex: 'id', width: 80, search: false },
  { title: '用户名', dataIndex: 'username', width: 150 },
  {
    title: '头像',
    dataIndex: 'avatar',
    width: 80,
    search: false,
    render: (_, r) => <Avatar src={r.avatar} size={40} />,
  },
  { title: '姓名', dataIndex: 'name', width: 150 },
  {
    title: '角色',
    dataIndex: 'roles',
    search: false,
    render: (_, r) => r.roles?.map((role) => (
      <Tag key={role.id} color="blue">{role.name}</Tag>
    )),
  },
  {
    title: '状态',
    dataIndex: 'enabled',
    valueType: 'select',
    valueEnum: {
      true: { text: '启用', status: 'Success' },
      false: { text: '禁用', status: 'Default' },
    },
    render: (_, r) => (
      <Tag color={r.enabled ? 'green' : 'red'}>
        {r.enabled ? '启用' : '禁用'}
      </Tag>
    ),
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    valueType: 'dateTime',
    search: false,
    sorter: true,
  },
];

const FormContent: React.FC<{ record?: AdminUser }> = ({ record }) => {
  const isEdit = !!record;
  const [avatarUrl, setAvatarUrl] = useState('');

  useEffect(() => {
    setAvatarUrl(
      record?.avatar || 
      `https://api.dicebear.com/9.x/thumbs/png?seed=${Math.random()}`
    );
  }, [record]);

  return (
    <>
      {!isEdit && (
        <>
          <ProFormText name="username" label="用户名" rules={[{ required: true }]} />
          <ProFormText.Password name="password" label="密码" rules={[{ required: true }]} />
        </>
      )}
      <ProFormText name="name" label="姓名" />
      
      <div style={{ marginBottom: 24 }}>
        <div style={{ marginBottom: 8 }}>头像</div>
        <Space>
          <Avatar src={avatarUrl} size={64} />
          <Upload
            name="avatar"
            showUploadList={false}
            action={buildApiUrl('/auth/avatar')}
            headers={{ Authorization: `Bearer ${localStorage.getItem('token')}` }}
            onChange={(info) => {
              if (info.file.status === 'done') {
                setAvatarUrl(info.file.response?.data?.url);
                message.success('上传成功');
              }
            }}
          >
            <Button icon={<UploadOutlined />}>上传</Button>
          </Upload>
        </Space>
      </div>

      <ProFormSelect
        name="roleIds"
        label="角色"
        mode="multiple"
        request={async () => {
          const res = await getAllAdminRoles();
          return (res.data || []).map((r: any) => ({ label: r.name, value: r.id }));
        }}
      />
      <ProFormSwitch name="enabled" label="状态" initialValue={true} />
    </>
  );
};

export default function AdminUserList() {
  return (
    <CrudTable<AdminUser>
      title="用户"
      permissionPrefix="system:admin_user"
      service={userService}
      columns={columns}
      formContent={<FormContent />}
      recordToValues={(r) => ({
        name: r.name,
        enabled: r.enabled,
        roleIds: r.roles?.map((role) => role.id),
      })}
      transformParams={(params) => ({
        ...params,
        enabled: params.enabled === 'true' ? true : 
                 params.enabled === 'false' ? false : undefined,
      })}
    />
  );
}
```

---

## 对比

| 项目 | 传统方式 | CrudTable |
|------|----------|-----------|
| 文件数 | 4-5 个 | 1 个 |
| 代码行数 | 400-500 行 | ~100 行 |
| CreateForm | 单独组件 | 合并 |
| UpdateForm | 单独组件 | 合并 |
| 服务定义 | 5 个函数 | 1 行 |
| 权限控制 | 手动判断 | 自动 |

---

## 最佳实践

1. **类型优先** - 先定义好接口类型
2. **列配置抽离** - 列配置建议单独定义，便于复用
3. **表单组件化** - 复杂表单可抽成独立组件
4. **善用 recordToValues** - 处理编辑时的数据转换
5. **善用 transformParams** - 处理特殊的请求参数格式

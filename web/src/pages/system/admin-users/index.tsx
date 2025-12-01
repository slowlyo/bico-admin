/**
 * 用户管理 - 使用 CrudTable 重构
 */
import type { ProColumns } from '@ant-design/pro-components';
import { ProFormText, ProFormSwitch, ProFormSelect } from '@ant-design/pro-components';
import { Avatar, Tag, Space, Upload, Button, message } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import React, { useState, useEffect } from 'react';
import { CrudTable } from '@/components';
import { createCrudService } from '@/services/crud';
import { getAllAdminRoles } from '@/services/system/admin-role';
import { buildApiUrl } from '@/services/config';

// 类型定义
interface AdminUser {
  id: number;
  username: string;
  name: string;
  avatar: string;
  enabled: boolean;
  roles?: { id: number; name: string }[];
  created_at: string;
}

// CRUD 服务
const userService = createCrudService<AdminUser>('/admin-users');

// 列配置
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
    width: 200,
    render: (_, r) => <Space size={4}>{r.roles?.map((role) => <Tag key={role.id} color="blue">{role.name}</Tag>)}</Space>,
  },
  {
    title: '状态',
    dataIndex: 'enabled',
    width: 100,
    valueType: 'select',
    valueEnum: { true: { text: '启用', status: 'Success' }, false: { text: '禁用', status: 'Default' } },
    render: (_, r) => <Tag color={r.enabled ? 'green' : 'red'}>{r.enabled ? '启用' : '禁用'}</Tag>,
  },
  { title: '创建时间', dataIndex: 'created_at', valueType: 'dateTime', width: 180, search: false, sorter: true },
];

// 表单内容组件
const FormContent: React.FC<{ record?: AdminUser }> = ({ record }) => {
  const isEdit = !!record;
  const [avatarUrl, setAvatarUrl] = useState('');

  useEffect(() => {
    setAvatarUrl(record?.avatar || `https://api.dicebear.com/9.x/thumbs/png?seed=${Math.random()}`);
  }, [record]);

  return (
    <>
      {!isEdit && (
        <>
          <ProFormText name="username" label="用户名" placeholder="请输入用户名" rules={[{ required: true }]} />
          <ProFormText.Password name="password" label="密码" placeholder="请输入密码" rules={[{ required: true }]} />
        </>
      )}
      <ProFormText name="name" label="姓名" placeholder="请输入姓名" />
      <div style={{ marginBottom: 24 }}>
        <div style={{ marginBottom: 8, fontWeight: 500 }}>头像</div>
        <Space size={16}>
          <Avatar src={avatarUrl} size={64} />
          <Upload
            name="avatar"
            showUploadList={false}
            action={buildApiUrl('/auth/avatar')}
            headers={{ Authorization: `Bearer ${localStorage.getItem('token')}` }}
            onChange={(info) => {
              if (info.file.status === 'done') {
                const url = info.file.response?.data?.url;
                if (url) {
                  setAvatarUrl(url);
                  message.success('上传成功');
                }
              } else if (info.file.status === 'error') {
                message.error('上传失败');
              }
            }}
          >
            <Button icon={<UploadOutlined />}>上传头像</Button>
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
        enabled: params.enabled === 'true' ? true : params.enabled === 'false' ? false : undefined,
      })}
    />
  );
}

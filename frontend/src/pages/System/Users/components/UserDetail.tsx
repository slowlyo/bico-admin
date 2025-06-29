import React from 'react';
import { ProDescriptions } from '@ant-design/pro-components';
import { Card, Tag } from 'antd';

export interface UserDetailProps {
  user: {
    id: number;
    username: string;
    email: string;
    nickname?: string;
    phone?: string;
    status: number;
    role?: string;
    created_at: string;
    updated_at: string;
  };
  onClose: () => void;
}

const UserDetail: React.FC<UserDetailProps> = ({ user, onClose }) => {
  return (
    <div>
      <ProDescriptions
        title="用户详情"
        column={1}
        bordered
        dataSource={user}
        columns={[
          {
            title: 'ID',
            dataIndex: 'id',
            copyable: true,
          },
          {
            title: '用户名',
            dataIndex: 'username',
            copyable: true,
          },
          {
            title: '邮箱',
            dataIndex: 'email',
            copyable: true,
          },
          {
            title: '昵称',
            dataIndex: 'nickname',
            render: (text) => text || '-',
          },
          {
            title: '手机号',
            dataIndex: 'phone',
            copyable: true,
            render: (text) => text || '-',
          },
          {
            title: '角色',
            dataIndex: 'role',
            render: (role) => {
              const roleMap: Record<string, { text: string; color: string }> = {
                admin: { text: '管理员', color: 'red' },
                manager: { text: '管理者', color: 'orange' },
                user: { text: '普通用户', color: 'blue' },
              };
              const roleInfo = roleMap[role] || { text: role || '-', color: 'default' };
              return role ? (
                <Tag color={roleInfo.color}>{roleInfo.text}</Tag>
              ) : (
                '-'
              );
            },
          },
          {
            title: '状态',
            dataIndex: 'status',
            render: (status) => (
              <Tag color={status === 1 ? 'green' : 'red'}>
                {status === 1 ? '正常' : '禁用'}
              </Tag>
            ),
          },
          {
            title: '创建时间',
            dataIndex: 'created_at',
            valueType: 'dateTime',
          },
          {
            title: '更新时间',
            dataIndex: 'updated_at',
            valueType: 'dateTime',
          },
        ]}
      />

      <Card title="用户权限" style={{ marginTop: 16 }}>
        <div style={{ color: '#999', textAlign: 'center', padding: '20px 0' }}>
          用户权限由角色决定，请查看角色管理了解详细权限配置
        </div>
      </Card>
    </div>
  );
};

export default UserDetail;

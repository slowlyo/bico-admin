import React from 'react';
import { ProDescriptions } from '@ant-design/pro-components';
import { Button } from 'antd';

export interface UserDetailProps {
  user: {
    id: number;
    username: string;
    email: string;
    nickname?: string;
    phone?: string;
    status: number;
    created_at: string;
    updated_at: string;
  };
  onClose: () => void;
}

const UserDetail: React.FC<UserDetailProps> = ({ user, onClose }) => {
  return (
    <>
      <ProDescriptions
        column={1}
        title="用户详情"
        request={async () => ({
          data: user || {},
        })}
        params={{
          id: user?.id,
        }}
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
          },
          {
            title: '手机号',
            dataIndex: 'phone',
            copyable: true,
          },
          {
            title: '状态',
            dataIndex: 'status',
            valueEnum: {
              1: {
                text: '正常',
                status: 'Success',
              },
              0: {
                text: '禁用',
                status: 'Error',
              },
            },
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
      <div style={{ textAlign: 'center', marginTop: 24 }}>
        <Button onClick={onClose}>关闭</Button>
      </div>
    </>
  );
};

export default UserDetail;

import React, { useRef, useState } from 'react';
import { PlusOutlined } from '@ant-design/icons';
import {
  ActionType,
  PageContainer,
  ProTable,
  ProColumns,
} from '@ant-design/pro-components';
import { Button, Drawer, message, Popconfirm, Switch } from 'antd';
import { getUserList, deleteUser, updateUserStatus } from '@/services/user';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';
import UserDetail from './components/UserDetail';

export type UserItem = {
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

const UserList: React.FC = () => {
  const [createModalOpen, setCreateModalOpen] = useState<boolean>(false);
  const [updateModalOpen, setUpdateModalOpen] = useState<boolean>(false);
  const [showDetail, setShowDetail] = useState<boolean>(false);
  const actionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<UserItem>();




  /**
   * 删除用户
   */
  const handleRemove = async (record: UserItem) => {
    const hide = message.loading('正在删除');
    try {
      await deleteUser(record.id);
      hide();
      message.success('删除成功');
      actionRef.current?.reloadAndRest?.();
      return true;
    } catch (error) {
      hide();
      // 错误消息由全局错误处理器显示，这里不再重复显示
      return false;
    }
  };

  /**
   * 切换用户状态
   */
  const handleStatusChange = async (record: UserItem, checked: boolean) => {
    const hide = message.loading('正在更新状态');
    try {
      const newStatus = checked ? 1 : 0;
      await updateUserStatus(record.id, newStatus);
      hide();
      message.success('状态更新成功');
      actionRef.current?.reload?.();
    } catch (error) {
      hide();
      // 错误消息由全局错误处理器显示，这里不再重复显示
    }
  };

  const columns: ProColumns<UserItem>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      hideInSearch: true,
      width: 80,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      valueType: 'text',
      render: (dom, entity) => {
        return (
          <a
            onClick={() => {
              setCurrentRow(entity);
              setShowDetail(true);
            }}
          >
            {dom}
          </a>
        );
      },
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      valueType: 'text',
    },
    {
      title: '昵称',
      dataIndex: 'nickname',
      valueType: 'text',
    },
    {
      title: '手机号',
      dataIndex: 'phone',
      valueType: 'text',
    },
    {
      title: '角色',
      dataIndex: 'role',
      valueType: 'text',
      valueEnum: {
        'admin': {
          text: '管理员',
          status: 'Error',
        },
        'manager': {
          text: '管理者',
          status: 'Warning',
        },
        'user': {
          text: '普通用户',
          status: 'Default',
        },
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      hideInForm: true,
      render: (_, record) => (
        <Switch
          checked={record.status === 1}
          onChange={(checked) => handleStatusChange(record, checked)}
          checkedChildren="正常"
          unCheckedChildren="禁用"
        />
      ),
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
      sorter: true,
      dataIndex: 'created_at',
      valueType: 'dateTime',
      hideInSearch: true,
    },
    {
      title: '操作',
      dataIndex: 'option',
      valueType: 'option',
      render: (_, record) => [
        <a
          key="config"
          onClick={() => {
            setUpdateModalOpen(true);
            setCurrentRow(record);
          }}
        >
          编辑
        </a>,
        <Popconfirm
          key="delete"
          title="确定删除这个用户吗？"
          onConfirm={async () => {
            await handleRemove(record);
          }}
        >
          <a style={{ color: 'red' }}>删除</a>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<UserItem, API.PageParams>
        headerTitle="用户列表"
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [
          <Button
            type="primary"
            key="primary"
            onClick={() => {
              setCreateModalOpen(true);
            }}
          >
            <PlusOutlined /> 新建
          </Button>,
        ]}
        request={getUserList}
        columns={columns}
      />

      <CreateForm
        open={createModalOpen}
        onOpenChange={setCreateModalOpen}
        onFinish={async () => {
          setCreateModalOpen(false);
          if (actionRef.current) {
            actionRef.current.reload();
          }
        }}
      />

      <UpdateForm
        open={updateModalOpen}
        onOpenChange={setUpdateModalOpen}
        onFinish={async () => {
          setUpdateModalOpen(false);
          setCurrentRow(undefined);
          if (actionRef.current) {
            actionRef.current.reload();
          }
        }}
        values={currentRow || {}}
      />

      <Drawer
        width={600}
        open={showDetail}
        onClose={() => {
          setCurrentRow(undefined);
          setShowDetail(false);
        }}
        closable={false}
      >
        {currentRow?.username && (
          <UserDetail
            user={currentRow}
            onClose={() => {
              setCurrentRow(undefined);
              setShowDetail(false);
            }}
          />
        )}
      </Drawer>
    </PageContainer>
  );
};

export default UserList;

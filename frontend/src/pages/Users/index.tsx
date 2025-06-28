import React, { useRef, useState } from 'react';
import { PlusOutlined } from '@ant-design/icons';
import {
  ActionType,
  FooterToolbar,
  PageContainer,
  ProTable,
  ProColumns,
} from '@ant-design/pro-components';
import { Button, Drawer, message, Popconfirm } from 'antd';
import { getUserList, deleteUser } from '@/services/user';
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
  created_at: string;
  updated_at: string;
};

const UserList: React.FC = () => {
  const [createModalOpen, setCreateModalOpen] = useState<boolean>(false);
  const [updateModalOpen, setUpdateModalOpen] = useState<boolean>(false);
  const [showDetail, setShowDetail] = useState<boolean>(false);
  const actionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<UserItem>();
  const [selectedRowsState, setSelectedRows] = useState<UserItem[]>([]);

  /**
   * 删除用户
   */
  const handleRemove = async (selectedRows: UserItem[]) => {
    const hide = message.loading('正在删除');
    if (!selectedRows) return true;
    try {
      for (const row of selectedRows) {
        await deleteUser(row.id);
      }
      hide();
      message.success('删除成功，即将刷新');
      actionRef.current?.reloadAndRest?.();
      return true;
    } catch (error) {
      hide();
      message.error('删除失败，请重试');
      return false;
    }
  };

  const columns: ProColumns<UserItem>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      tip: '用户唯一标识',
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
      title: '用户名',
      dataIndex: 'username',
      valueType: 'text',
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
      title: '状态',
      dataIndex: 'status',
      hideInForm: true,
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
            await handleRemove([record]);
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
        rowSelection={{
          onChange: (_, selectedRows) => {
            setSelectedRows(selectedRows);
          },
        }}
      />
      {selectedRowsState?.length > 0 && (
        <FooterToolbar
          extra={
            <div>
              已选择{' '}
              <a style={{ fontWeight: 600 }}>{selectedRowsState.length}</a>{' '}
              项
            </div>
          }
        >
          <Button
            onClick={async () => {
              await handleRemove(selectedRowsState);
              setSelectedRows([]);
              actionRef.current?.reloadAndRest?.();
            }}
          >
            批量删除
          </Button>
        </FooterToolbar>
      )}

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

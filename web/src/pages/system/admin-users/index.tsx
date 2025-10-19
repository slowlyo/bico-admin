import { PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { ProTable } from '@ant-design/pro-components';
import { Button, message, Popconfirm, Space, Tag, Avatar } from 'antd';
import React, { useRef, useState } from 'react';
import { getAdminUserList, deleteAdminUser, type AdminUser } from '@/services/admin-user';
import { useAccess } from '@umijs/max';
import { PageContainer } from '@/components';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';

const AdminUserList: React.FC = () => {
  const actionRef = useRef<ActionType>(null);
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [updateModalVisible, setUpdateModalVisible] = useState(false);
  const [currentRow, setCurrentRow] = useState<AdminUser>();
  const access = useAccess() as Record<string, boolean>;

  const handleSuccess = () => {
    setCreateModalVisible(false);
    setUpdateModalVisible(false);
    setCurrentRow(undefined);
    actionRef.current?.reload();
  };

  const handleDelete = async (id: number) => {
    try {
      const res = await deleteAdminUser(id);
      if (res.code === 0) {
        message.success('删除成功');
        actionRef.current?.reload();
      } else {
        message.error(res.msg || '删除失败');
      }
    } catch (error) {
      message.error('删除失败');
    }
  };

  const columns: ProColumns<AdminUser>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
      search: false,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      width: 150,
    },
    {
      title: '头像',
      dataIndex: 'avatar',
      width: 80,
      search: false,
      render: (_, record) => (
        <Avatar src={record.avatar} size={40} />
      ),
    },
    {
      title: '姓名',
      dataIndex: 'name',
      width: 150,
    },
    {
      title: '角色',
      dataIndex: 'roles',
      search: false,
      width: 200,
      render: (_, record) => (
        <>
          {record.roles?.map((role) => (
            <Tag key={role.id} color="blue">
              {role.name}
            </Tag>
          ))}
        </>
      ),
    },
    {
      title: '状态',
      dataIndex: 'enabled',
      width: 100,
      valueType: 'select',
      valueEnum: {
        true: { text: '启用', status: 'Success' },
        false: { text: '禁用', status: 'Default' },
      },
      render: (_, record) => (
        <Tag color={record.enabled ? 'green' : 'red'}>
          {record.enabled ? '启用' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 180,
      search: false,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 180,
      fixed: 'right',
      render: (_, record) => (
        <Space>
          {access['system:admin_user:edit'] && (
            <a
              onClick={() => {
                setCurrentRow(record);
                setUpdateModalVisible(true);
              }}
            >
              编辑
            </a>
          )}
          {access['system:admin_user:delete'] && (
            <Popconfirm
              title="确定删除该用户吗？"
              onConfirm={() => handleDelete(record.id)}
              okText="确定"
              cancelText="取消"
            >
              <a style={{ color: 'red' }}>删除</a>
            </Popconfirm>
          )}
        </Space>
      ),
    },
  ];

  return (
    <PageContainer>
      <ProTable<AdminUser>
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        pagination={{
          showSizeChanger: true,
          showQuickJumper: true,
          pageSizeOptions: ['10', '20', '50', '100'],
          defaultPageSize: 10,
        }}
        toolBarRender={() => [
          access['system:admin_user:create'] && (
            <Button
              type="primary"
              key="create"
              icon={<PlusOutlined />}
              onClick={() => setCreateModalVisible(true)}
            >
              新建
            </Button>
          ),
        ]}
        request={async (params) => {
          const res = await getAdminUserList({
            page: params.current,
            pageSize: params.pageSize,
            username: params.username,
            name: params.name,
            enabled: params.enabled === 'true' ? true : params.enabled === 'false' ? false : undefined,
          });
          return {
            data: (res.data || []) as AdminUser[],
            total: res.total || 0,
            success: res.code === 0,
          };
        }}
        columns={columns}
        scroll={{ x: 1200 }}
      />

      <CreateForm
        open={createModalVisible}
        onOpenChange={setCreateModalVisible}
        onSuccess={handleSuccess}
      />

      <UpdateForm
        open={updateModalVisible}
        onOpenChange={(visible) => {
          setUpdateModalVisible(visible);
          if (!visible) setCurrentRow(undefined);
        }}
        onSuccess={handleSuccess}
        currentRow={currentRow}
      />
    </PageContainer>
  );
};

export default AdminUserList;

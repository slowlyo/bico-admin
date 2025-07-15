import {
  ActionType,
  PageContainer,
  ProColumns,
  ProTable,
} from '@ant-design/pro-components';
import { Button, Divider, message, Popconfirm, Switch } from 'antd';
import React, { useRef, useState } from 'react';
import {
  getAdminUserList,
  createAdminUser,
  updateAdminUser,
  deleteAdminUser,
  updateAdminUserStatus,
  AdminUser,
  AdminUserCreateRequest,
  AdminUserUpdateRequest,
} from '@/services/adminUser';
import AdminUserForm from './components/AdminUserForm';

/**
 * 添加管理员用户
 */
const handleAdd = async (fields: AdminUserCreateRequest) => {
  const hide = message.loading('正在添加');
  try {
    const response = await createAdminUser(fields);
    hide();
    if (response.code === 200) {
      message.success('添加成功');
      return true;
    } else {
      message.error(response.message || '添加失败');
      return false;
    }
  } catch (error) {
    hide();
    message.error('添加失败请重试！');
    return false;
  }
};

/**
 * 更新管理员用户
 */
const handleUpdate = async (id: number, fields: AdminUserUpdateRequest) => {
  const hide = message.loading('正在更新');
  try {
    const response = await updateAdminUser(id, fields);
    hide();
    if (response.code === 200) {
      message.success('更新成功');
      return true;
    } else {
      message.error(response.message || '更新失败');
      return false;
    }
  } catch (error) {
    hide();
    message.error('更新失败请重试！');
    return false;
  }
};

/**
 * 删除管理员用户
 */
const handleRemove = async (id: number) => {
  const hide = message.loading('正在删除');
  try {
    const response = await deleteAdminUser(id);
    hide();
    if (response.code === 200) {
      message.success('删除成功');
      return true;
    } else {
      message.error(response.message || '删除失败');
      return false;
    }
  } catch (error) {
    hide();
    message.error('删除失败请重试！');
    return false;
  }
};

/**
 * 更新用户状态
 */
const handleStatusChange = async (id: number, status: number) => {
  const hide = message.loading('正在更新状态');
  try {
    const response = await updateAdminUserStatus(id, status);
    hide();
    if (response.code === 200) {
      message.success('状态更新成功');
      return true;
    } else {
      message.error(response.message || '状态更新失败');
      return false;
    }
  } catch (error) {
    hide();
    message.error('状态更新失败请重试！');
    return false;
  }
};

const AdminUserList: React.FC = () => {
  const [createModalVisible, handleModalVisible] = useState<boolean>(false);
  const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>(false);
  const [stepFormValues, setStepFormValues] = useState<AdminUser | undefined>();
  const actionRef = useRef<ActionType>();

  const columns: ProColumns<AdminUser>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      hideInForm: true,
      hideInSearch: true,
      width: 80,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      formItemProps: {
        rules: [
          {
            required: true,
            message: '用户名为必填项',
          },
          {
            min: 3,
            max: 50,
            message: '用户名长度为3-50个字符',
          },
        ],
      },
    },
    {
      title: '姓名',
      dataIndex: 'name',
      formItemProps: {
        rules: [
          {
            required: true,
            message: '姓名为必填项',
          },
          {
            max: 100,
            message: '姓名长度不能超过100个字符',
          },
        ],
      },
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      hideInSearch: true,
      formItemProps: {
        rules: [
          {
            type: 'email',
            message: '请输入正确的邮箱格式',
          },
        ],
      },
    },
    {
      title: '手机号',
      dataIndex: 'phone',
      hideInSearch: true,
    },
    {
      title: '状态',
      dataIndex: 'status',
      hideInForm: true,
      valueType: 'select',
      valueEnum: {
        1: { text: '启用', status: 'Success' },
        0: { text: '禁用', status: 'Default' },
      },
      render: (_, record) => (
        <Switch
          checked={record.status === 1}
          disabled={!record.can_disable}
          onChange={async (checked) => {
            const success = await handleStatusChange(record.id, checked ? 1 : 0);
            if (success && actionRef.current) {
              actionRef.current.reload();
            }
          }}
          checkedChildren="启用"
          unCheckedChildren="禁用"
        />
      ),
    },
    {
      title: '备注',
      dataIndex: 'remark',
      hideInSearch: true,
      hideInForm: true,
      ellipsis: true,
    },

    {
      title: '最后登录时间',
      dataIndex: 'last_login_at',
      hideInForm: true,
      hideInSearch: true,
      valueType: 'dateTime',
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      hideInForm: true,
      hideInSearch: true,
      valueType: 'dateTime',
    },
    {
      title: '操作',
      dataIndex: 'option',
      valueType: 'option',
      render: (_, record) => (
        <>
          <Button
            size="small"
            type="link"
            onClick={() => {
              handleUpdateModalVisible(true);
              setStepFormValues(record);
            }}
          >
            编辑
          </Button>
          <Divider type="vertical" />
          <Popconfirm
            title="确定要删除这个管理员用户吗？"
            disabled={!record.can_delete}
            onConfirm={async () => {
              const success = await handleRemove(record.id);
              if (success && actionRef.current) {
                actionRef.current.reload();
              }
            }}
            okText="确定"
            cancelText="取消"
          >
            <Button
              type="link"
              size="small"
              danger
              disabled={!record.can_delete}
            >
              删除
            </Button>
          </Popconfirm>
        </>
      ),
    },
  ];

  return (
    <PageContainer
      header={{
        title: '管理员用户管理',
      }}
    >
      <ProTable<AdminUser>
        headerTitle="管理员用户列表"
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [
          <Button
            key="1"
            type="primary"
            onClick={() => handleModalVisible(true)}
          >
            新建管理员
          </Button>,
        ]}
        request={async (params) => {
          const response = await getAdminUserList({
            page: params.current || 1,
            page_size: params.pageSize || 10,
            username: params.username,
            name: params.name,
            status: params.status,
          });
          
          if (response.code === 200) {
            return {
              data: response.data.list || [],
              success: true,
              total: response.data.total || 0,
            };
          } else {
            message.error(response.message || '获取数据失败');
            return {
              data: [],
              success: false,
              total: 0,
            };
          }
        }}
        columns={columns}
      />
      
      <AdminUserForm
        modalVisible={createModalVisible || updateModalVisible}
        isEdit={!!stepFormValues}
        values={stepFormValues}
        onCancel={() => {
          if (createModalVisible) {
            handleModalVisible(false);
          } else {
            handleUpdateModalVisible(false);
            setStepFormValues(undefined);
          }
        }}
        onSubmit={async (value) => {
          let success = false;
          if (stepFormValues) {
            // 编辑模式
            success = await handleUpdate(stepFormValues.id, value as AdminUserUpdateRequest);
            if (success) {
              handleUpdateModalVisible(false);
              setStepFormValues(undefined);
            }
          } else {
            // 创建模式
            success = await handleAdd(value as AdminUserCreateRequest);
            if (success) {
              handleModalVisible(false);
            }
          }

          if (success && actionRef.current) {
            actionRef.current.reload();
          }
        }}
      />
    </PageContainer>
  );
};

export default AdminUserList;

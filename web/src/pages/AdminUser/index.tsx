import {
  ActionType,
  PageContainer,
  ProColumns,
  ProTable,
} from '@ant-design/pro-components';
import { Button, message, Popconfirm, Switch, Avatar } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import { Access, useAccess } from '@umijs/max';
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
import { getActiveRoles, Role } from '@/services/role';
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
  const access = useAccess();
  const [createDrawerVisible, handleDrawerVisible] = useState<boolean>(false);
  const [updateDrawerVisible, handleUpdateDrawerVisible] = useState<boolean>(false);
  const [stepFormValues, setStepFormValues] = useState<AdminUser | undefined>();
  const [roles, setRoles] = useState<Role[]>([]);
  const actionRef = useRef<ActionType>();

  // 加载角色列表用于筛选
  const loadRoles = async () => {
    try {
      const response = await getActiveRoles();
      if (response.code === 200) {
        setRoles(response.data);
      }
    } catch (error) {
      console.error('获取角色列表失败:', error);
    }
  };

  // 组件挂载时加载角色列表
  React.useEffect(() => {
    loadRoles();
  }, []);

  const columns: ProColumns<AdminUser>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      hideInForm: true,
      hideInSearch: true,
      width: 80,
    },
    {
      title: '头像',
      dataIndex: 'avatar',
      hideInForm: true,
      hideInSearch: true,
      width: 80,
      render: (_, record) => (
        <Avatar
          size={40}
          src={record.avatar}
          icon={<UserOutlined />}
        />
      ),
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
      title: '角色',
      dataIndex: 'role_id',
      hideInForm: true,
      valueType: 'select',
      fieldProps: {
        placeholder: '请选择角色',
        showSearch: true,
        filterOption: (input: string, option: any) =>
          option?.label?.toLowerCase().includes(input.toLowerCase()),
        options: roles.map(role => ({
          label: role.name,
          value: role.id,
        })),
      },
      render: (_, record) => (
        <div>
          {record.roles?.map(role => (
            <span key={role.id} style={{
              display: 'inline-block',
              background: '#f0f0f0',
              padding: '2px 8px',
              borderRadius: '4px',
              margin: '2px',
              fontSize: '12px'
            }}>
              {role.name}
            </span>
          ))}
        </div>
      ),
    },

    {
      title: '最后登录时间',
      dataIndex: 'last_login_at',
      hideInForm: true,
      hideInSearch: true,
      valueType: 'dateTime',
      sorter: true,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      hideInForm: true,
      hideInSearch: true,
      valueType: 'dateTime',
      sorter: true,
    },
    {
      title: '操作',
      dataIndex: 'option',
      valueType: 'option',
      render: (_, record) => (
        <>
          <Access accessible={access.canEditAdminUser}>
            <Button
              size="small"
              type="link"
              onClick={() => {
                handleUpdateDrawerVisible(true);
                setStepFormValues(record);
              }}
            >
              编辑
            </Button>
          </Access>
          <Access accessible={access.canDeleteAdminUser}>
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
          </Access>
        </>
      ),
    },
  ];

  return (
    <PageContainer title={false}>
      <ProTable<AdminUser>
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [
          <Access key="1" accessible={access.canCreateAdminUser}>
            <Button
              type="primary"
              onClick={() => handleDrawerVisible(true)}
            >
              新建管理员
            </Button>
          </Access>,
        ]}
        request={async (params, sort) => {
          // 处理排序参数
          let sortBy = '';
          let sortDesc = false;

          if (sort && Object.keys(sort).length > 0) {
            const sortField = Object.keys(sort)[0];
            const sortOrder = sort[sortField];
            sortBy = sortField;
            sortDesc = sortOrder === 'descend';
          }

          const response = await getAdminUserList({
            page: params.current || 1,
            page_size: params.pageSize || 10,
            username: params.username,
            name: params.name,
            status: params.status,
            role_id: params.role_id,
            sort_by: sortBy,
            sort_desc: sortDesc,
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
        pagination={{
          defaultPageSize: 20,
          showQuickJumper: true,
          showSizeChanger: true,
          pageSizeOptions: ['10', '20', '50', '100'],
        }}
      />
      
      <AdminUserForm
        drawerVisible={createDrawerVisible || updateDrawerVisible}
        isEdit={!!stepFormValues}
        values={stepFormValues}
        onCancel={() => {
          if (createDrawerVisible) {
            handleDrawerVisible(false);
          } else {
            handleUpdateDrawerVisible(false);
            setStepFormValues(undefined);
          }
        }}
        onSubmit={async (value) => {
          let success = false;
          if (stepFormValues) {
            // 编辑模式
            success = await handleUpdate(stepFormValues.id, value as AdminUserUpdateRequest);
            if (success) {
              handleUpdateDrawerVisible(false);
              setStepFormValues(undefined);
            }
          } else {
            // 创建模式
            success = await handleAdd(value as AdminUserCreateRequest);
            if (success) {
              handleDrawerVisible(false);
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

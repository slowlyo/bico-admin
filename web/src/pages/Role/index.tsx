import {
  ActionType,
  PageContainer,
  ProColumns,
  ProTable,
} from '@ant-design/pro-components';
import { Button, Divider, message, Popconfirm, Switch, Tag } from 'antd';
import React, { useRef, useState } from 'react';
import {
  getRoleList,
  createRole,
  updateRole,
  deleteRole,
  updateRoleStatus,
  Role,
  RoleCreateRequest,
  RoleUpdateRequest,
} from '@/services/role';
import RoleForm from './components/RoleForm';

/**
 * 添加角色
 */
const handleAdd = async (fields: RoleCreateRequest) => {
  const hide = message.loading('正在添加');
  try {
    const response = await createRole(fields);
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
 * 更新角色
 */
const handleUpdate = async (id: number, fields: RoleUpdateRequest) => {
  const hide = message.loading('正在更新');
  try {
    const response = await updateRole(id, fields);
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
 * 删除角色
 */
const handleRemove = async (id: number) => {
  const hide = message.loading('正在删除');
  try {
    const response = await deleteRole(id);
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
 * 更新角色状态
 */
const handleStatusChange = async (id: number, status: number) => {
  const hide = message.loading('正在更新状态');
  try {
    const response = await updateRoleStatus(id, status);
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

const RoleList: React.FC = () => {
  const [createModalVisible, handleModalVisible] = useState<boolean>(false);
  const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>(false);
  const [stepFormValues, setStepFormValues] = useState<Role | undefined>();
  const actionRef = useRef<ActionType>();

  const columns: ProColumns<Role>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      hideInForm: true,
      width: 80,
    },
    {
      title: '角色代码',
      dataIndex: 'code',
      formItemProps: {
        rules: [
          {
            required: true,
            message: '角色代码为必填项',
          },
          {
            pattern: /^[a-zA-Z][a-zA-Z0-9_]*$/,
            message: '角色代码只能包含字母、数字和下划线，且以字母开头',
          },
          {
            min: 2,
            max: 50,
            message: '角色代码长度为2-50个字符',
          },
        ],
      },
    },
    {
      title: '角色名称',
      dataIndex: 'name',
      formItemProps: {
        rules: [
          {
            required: true,
            message: '角色名称为必填项',
          },
          {
            max: 100,
            message: '角色名称长度不能超过100个字符',
          },
        ],
      },
    },
    {
      title: '描述',
      dataIndex: 'description',
      hideInSearch: true,
      ellipsis: true,
      formItemProps: {
        rules: [
          {
            max: 500,
            message: '描述长度不能超过500个字符',
          },
        ],
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      hideInForm: true,
      render: (_, record) => (
        <Switch
          checked={record.status === 1}
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
            title="确定要删除这个角色吗？"
            onConfirm={async () => {
              const success = await handleRemove(record.id);
              if (success && actionRef.current) {
                actionRef.current.reload();
              }
            }}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" size="small" danger>
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
        title: '角色管理',
      }}
    >
      <ProTable<Role>
        headerTitle="角色列表"
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
            新建角色
          </Button>,
        ]}
        request={async (params) => {
          const response = await getRoleList({
            page: params.current || 1,
            page_size: params.pageSize || 10,
            code: params.code,
            name: params.name,
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
      
      <RoleForm
        onCancel={() => handleModalVisible(false)}
        modalVisible={createModalVisible}
        isEdit={false}
        onSubmit={async (value) => {
          const success = await handleAdd(value as RoleCreateRequest);
          if (success) {
            handleModalVisible(false);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
        }}
      />

      {stepFormValues && (
        <RoleForm
          onSubmit={async (value) => {
            const success = await handleUpdate(stepFormValues.id, value as RoleUpdateRequest);
            if (success) {
              handleUpdateModalVisible(false);
              setStepFormValues(undefined);
              if (actionRef.current) {
                actionRef.current.reload();
              }
            }
          }}
          onCancel={() => {
            handleUpdateModalVisible(false);
            setStepFormValues(undefined);
          }}
          modalVisible={updateModalVisible}
          values={stepFormValues}
          isEdit={true}
        />
      )}
    </PageContainer>
  );
};

export default RoleList;

import React, { useRef, useState } from 'react';
import { PlusOutlined } from '@ant-design/icons';
import {
  ActionType,
  PageContainer,
  ProTable,
  ProColumns,
} from '@ant-design/pro-components';
import { Button, Drawer, message, Popconfirm, Switch, Tag } from 'antd';
import { getRoleList, deleteRole, updateRoleStatus, getRolePermissions } from '@/services/role';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';
import RoleDetail from './components/RoleDetail';
import PermissionAssign from './components/PermissionAssign';

export type RoleItem = {
  id: number;
  name: string;
  code: string;
  description?: string;
  status: number;
  created_at: string;
  updated_at: string;
  permissions?: any[];
};

const RoleList: React.FC = () => {
  const [createModalOpen, setCreateModalOpen] = useState<boolean>(false);
  const [updateModalOpen, setUpdateModalOpen] = useState<boolean>(false);
  const [showDetail, setShowDetail] = useState<boolean>(false);
  const [showPermissionAssign, setShowPermissionAssign] = useState<boolean>(false);
  const actionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<RoleItem>();

  const [permissionCounts, setPermissionCounts] = useState<Record<number, number>>({});

  // 获取角色权限数量
  const fetchPermissionCount = async (roleId: number): Promise<number> => {
    try {
      const response = await getRolePermissions(roleId);
      if (response.code === 200) {
        return response.data?.length || 0;
      }
      return 0;
    } catch (error) {
      console.error('获取角色权限数量失败:', error);
      return 0;
    }
  };

  // 批量获取权限数量
  const fetchAllPermissionCounts = async (roles: RoleItem[]) => {
    const counts: Record<number, number> = {};
    for (const role of roles) {
      counts[role.id] = await fetchPermissionCount(role.id);
    }
    setPermissionCounts(counts);
  };

  /**
   * 删除角色
   */
  const handleRemove = async (record: RoleItem) => {
    // 检查是否为超级管理员角色
    if (isSuperAdmin(record)) {
      message.error('不能删除超级管理员角色');
      return false;
    }

    const hide = message.loading('正在删除');
    try {
      await deleteRole(record.id);
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
   * 检查是否为超级管理员角色
   */
  const isSuperAdmin = (record: RoleItem): boolean => {
    return record.code === 'super_admin';
  };

  /**
   * 切换角色状态
   */
  const handleStatusChange = async (record: RoleItem, checked: boolean) => {
    // 超级管理员角色不允许修改状态
    if (isSuperAdmin(record)) {
      message.warning('超级管理员角色状态不可修改');
      return;
    }

    const hide = message.loading('正在更新状态');
    try {
      const newStatus = checked ? 1 : 0;
      await updateRoleStatus(record.id, newStatus);
      hide();
      message.success('状态更新成功');
      actionRef.current?.reload?.();
    } catch (error) {
      hide();
      // 错误消息由全局错误处理器显示，这里不再重复显示
    }
  };

  const columns: ProColumns<RoleItem>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      valueType: 'text',
      width: 80,
      search: false,
    },
    {
      title: '角色名称',
      dataIndex: 'name',
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
      title: '角色标识',
      dataIndex: 'code',
      valueType: 'text',
    },
    {
      title: '描述',
      dataIndex: 'description',
      valueType: 'text',
      search: false,
      ellipsis: true,
    },
    {
      title: '权限数量',
      dataIndex: 'permissions',
      search: false,
      render: (_, record) => {
        const count = permissionCounts[record.id];
        return (
          <Tag color="blue">
            {count !== undefined ? `${count} 个权限` : '加载中...'}
          </Tag>
        );
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
          checkedChildren="启用"
          unCheckedChildren="禁用"
          disabled={isSuperAdmin(record)}
        />
      ),
      valueEnum: {
        1: {
          text: '启用',
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
      search: false,
    },
    {
      title: '操作',
      dataIndex: 'option',
      valueType: 'option',
      render: (_, record) => {
        const isProtected = isSuperAdmin(record);

        if (isProtected) {
          return [
            <span key="protected" style={{ color: '#999' }}>
              受保护角色
            </span>
          ];
        }

        return [
          <a
            key="edit"
            onClick={() => {
              setCurrentRow(record);
              setUpdateModalOpen(true);
            }}
          >
            编辑
          </a>,
          <a
            key="permission"
            onClick={() => {
              setCurrentRow(record);
              setShowPermissionAssign(true);
            }}
          >
            分配权限
          </a>,
          <Popconfirm
            key="delete"
            title="确定删除这个角色吗？"
            onConfirm={async () => {
              await handleRemove(record);
            }}
          >
            <a style={{ color: 'red' }}>删除</a>
          </Popconfirm>,
        ];
      },
    },
  ];

  return (
    <PageContainer>
      <ProTable<RoleItem, API.PageParams>
        headerTitle="角色列表"
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
        request={async (params) => {
          const result = await getRoleList(params);
          // 获取权限数量
          if (result.success && result.data.length > 0) {
            fetchAllPermissionCounts(result.data);
          }
          return result;
        }}
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
        values={currentRow || {}}
        onFinish={async () => {
          setUpdateModalOpen(false);
          setCurrentRow(undefined);
          if (actionRef.current) {
            actionRef.current.reload();
          }
        }}
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
        {currentRow?.id && (
          <RoleDetail
            roleId={currentRow.id}
            onClose={() => setShowDetail(false)}
          />
        )}
      </Drawer>

      <PermissionAssign
        open={showPermissionAssign}
        onOpenChange={setShowPermissionAssign}
        roleId={currentRow?.id}
        onFinish={async () => {
          setShowPermissionAssign(false);
          setCurrentRow(undefined);
          if (actionRef.current) {
            actionRef.current.reload();
          }
        }}
      />
    </PageContainer>
  );
};

export default RoleList;

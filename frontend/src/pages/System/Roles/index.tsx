import React, { useRef, useState } from 'react';
import { PlusOutlined } from '@ant-design/icons';
import {
  ActionType,
  FooterToolbar,
  PageContainer,
  ProTable,
  ProColumns,
} from '@ant-design/pro-components';
import { Button, Drawer, message, Popconfirm, Switch, Tag } from 'antd';
import { getRoleList, deleteRole, updateRoleStatus } from '@/services/role';
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
  const [selectedRowsState, setSelectedRows] = useState<RoleItem[]>([]);

  /**
   * 删除角色
   */
  const handleRemove = async (selectedRows: RoleItem[]) => {
    const hide = message.loading('正在删除');
    if (!selectedRows) return true;
    try {
      for (const row of selectedRows) {
        await deleteRole(row.id);
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

  /**
   * 切换角色状态
   */
  const handleStatusChange = async (record: RoleItem, checked: boolean) => {
    const hide = message.loading('正在更新状态');
    try {
      const newStatus = checked ? 1 : 0;
      await updateRoleStatus(record.id, newStatus);
      hide();
      message.success('状态更新成功');
      actionRef.current?.reload?.();
    } catch (error) {
      hide();
      message.error('状态更新失败，请重试');
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
      title: '角色代码',
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
      render: (_, record) => (
        <Tag color="blue">
          {record.permissions?.length || 0} 个权限
        </Tag>
      ),
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
      render: (_, record) => [
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
        request={getRoleList}
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

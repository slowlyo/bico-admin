import { PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { Button, message, Popconfirm, Space, Tag, Drawer, Tree } from 'antd';
import React, { useRef, useState, useEffect } from 'react';
import { getAdminRoleList, deleteAdminRole, getAllPermissions, updateRolePermissions, type AdminRole, type Permission } from '@/services/admin-role';
import { useAccess } from '@umijs/max';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';

const AdminRoleList: React.FC = () => {
  const actionRef = useRef<ActionType>(null);
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [updateModalVisible, setUpdateModalVisible] = useState(false);
  const [permissionDrawerVisible, setPermissionDrawerVisible] = useState(false);
  const [currentRow, setCurrentRow] = useState<AdminRole>();
  const [allPermissions, setAllPermissions] = useState<Permission[]>([]);
  const [selectedPermissions, setSelectedPermissions] = useState<string[]>([]);
  const access = useAccess() as Record<string, boolean>;

  const handleSuccess = () => {
    setCreateModalVisible(false);
    setUpdateModalVisible(false);
    setCurrentRow(undefined);
    actionRef.current?.reload();
  };

  useEffect(() => {
    loadAllPermissions();
  }, []);

  const loadAllPermissions = async () => {
    try {
      const res = await getAllPermissions();
      if (res.code === 0) {
        setAllPermissions(res.data || []);
      }
    } catch (error) {
      console.error('加载权限失败', error);
    }
  };

  const convertPermissionsToTreeData = (permissions: Permission[]): any[] => {
    return permissions.map(perm => ({
      title: perm.label,
      key: perm.key,
      children: perm.children ? convertPermissionsToTreeData(perm.children) : undefined,
    }));
  };

  const getAllPermissionKeys = (permissions: Permission[]): string[] => {
    let keys: string[] = [];
    permissions.forEach(perm => {
      keys.push(perm.key);
      if (perm.children) {
        keys = keys.concat(getAllPermissionKeys(perm.children));
      }
    });
    return keys;
  };

  const handleDelete = async (id: number) => {
    try {
      const res = await deleteAdminRole(id);
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

  // 从权限树中查找权限节点的所有父级
  const findPermissionParents = (
    permKey: string,
    tree: Permission[],
    currentPath: string[] = []
  ): string[] | null => {
    for (const node of tree) {
      // 如果当前节点就是目标，返回当前路径（不包含自己）
      if (node.key === permKey) {
        return currentPath;
      }
      
      // 如果有子节点，递归查找
      if (node.children && node.children.length > 0) {
        const result = findPermissionParents(permKey, node.children, [...currentPath, node.key]);
        if (result !== null) {
          return result;
        }
      }
    }
    
    return null;
  };

  // 过滤权限列表，移除冗余的父级权限（保留实际勾选的权限）
  const filterRedundantPermissions = (permissions: string[]): string[] => {
    const filtered = new Set<string>();
    
    permissions.forEach(permission => {
      // 检查是否有子权限也在列表中
      const hasChildInList = permissions.some(p => {
        if (p === permission) return false;
        const parents = findPermissionParents(p, allPermissions);
        return parents && parents.includes(permission);
      });
      
      // 如果没有子权限在列表中，说明这是用户实际勾选的
      if (!hasChildInList) {
        filtered.add(permission);
      }
    });
    
    return Array.from(filtered);
  };

  const handleOpenPermissionDrawer = (record: AdminRole) => {
    setCurrentRow(record);
    // 过滤掉冗余的父级权限，只保留用户实际勾选的
    const filteredPermissions = filterRedundantPermissions(record.permissions || []);
    console.log('后端返回的权限:', record.permissions);
    console.log('过滤后显示的权限:', filteredPermissions);
    setSelectedPermissions(filteredPermissions);
    setPermissionDrawerVisible(true);
  };

  // 扩展权限列表，包含所有父级权限
  const expandPermissions = (permissions: string[]): string[] => {
    const expanded = new Set<string>();
    
    permissions.forEach(permission => {
      // 添加当前权限
      expanded.add(permission);
      // 从权限树中查找并添加所有父级权限
      const parents = findPermissionParents(permission, allPermissions);
      if (parents && parents.length > 0) {
        console.log(`权限 "${permission}" 的父级:`, parents);
        parents.forEach(parent => {
          expanded.add(parent);
        });
      } else {
        console.log(`权限 "${permission}" 没有父级或是顶级权限`);
      }
    });
    
    return Array.from(expanded);
  };

  const handleSavePermissions = async () => {
    if (!currentRow) return;
    try {
      // 在保存时扩展权限，添加所有父级权限
      const expandedPermissions = expandPermissions(selectedPermissions);
      console.log('用户选择的权限:', selectedPermissions);
      console.log('扩展后的权限（包含父级）:', expandedPermissions);
      
      const res = await updateRolePermissions(currentRow.id, {
        permissions: expandedPermissions,
      });
      if (res.code === 0) {
        message.success('权限配置成功');
        setPermissionDrawerVisible(false);
        setCurrentRow(undefined);
        actionRef.current?.reload();
      } else {
        message.error(res.msg || '权限配置失败');
      }
    } catch (error) {
      message.error('权限配置失败');
    }
  };

  const columns: ProColumns<AdminRole>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
      search: false,
    },
    {
      title: '角色名称',
      dataIndex: 'name',
      width: 150,
    },
    {
      title: '角色代码',
      dataIndex: 'code',
      width: 150,
    },
    {
      title: '描述',
      dataIndex: 'description',
      search: false,
      width: 200,
      ellipsis: true,
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
      width: 200,
      fixed: 'right',
      render: (_, record) => (
        <Space>
          {access['system:admin_role:permission'] && (
            <a onClick={() => handleOpenPermissionDrawer(record)}>
              配置权限
            </a>
          )}
          {access['system:admin_role:edit'] && (
            <a
              onClick={() => {
                setCurrentRow(record);
                setUpdateModalVisible(true);
              }}
            >
              编辑
            </a>
          )}
          {access['system:admin_role:delete'] && (
            <Popconfirm
              title="确定删除该角色吗？"
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
      <ProTable<AdminRole>
        headerTitle="角色列表"
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [
          access['system:admin_role:create'] && (
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
          const res = await getAdminRoleList({
            page: params.current,
            pageSize: params.pageSize,
            name: params.name,
            code: params.code,
            enabled: params.enabled === 'true' ? true : params.enabled === 'false' ? false : undefined,
          });
          return {
            data: (res.data || []) as AdminRole[],
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

      <Drawer
        title="配置权限"
        width={500}
        open={permissionDrawerVisible}
        onClose={() => {
          setPermissionDrawerVisible(false);
          setCurrentRow(undefined);
        }}
        extra={
          <Space>
            <Button onClick={() => setPermissionDrawerVisible(false)}>取消</Button>
            <Button type="primary" onClick={handleSavePermissions}>
              保存
            </Button>
          </Space>
        }
      >
        <Tree
          checkable
          defaultExpandAll
          checkedKeys={selectedPermissions}
          onCheck={(checkedKeys) => {
            const keys = Array.isArray(checkedKeys) ? checkedKeys : checkedKeys.checked;
            setSelectedPermissions(keys as string[]);
          }}
          treeData={convertPermissionsToTreeData(allPermissions)}
        />
      </Drawer>
    </PageContainer>
  );
};

export default AdminRoleList;

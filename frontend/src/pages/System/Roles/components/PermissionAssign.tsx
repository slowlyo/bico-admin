import React, { useState, useEffect } from 'react';
import { Modal, Tree, message, Spin, Input } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import { getPermissionTree, getRolePermissions, assignRolePermissions } from '@/services/role';
import type { PermissionItem } from '@/services/role';

export type PermissionAssignProps = {
  open: boolean;
  onOpenChange: (visible: boolean) => void;
  onFinish: () => void;
  roleId?: number;
};

const PermissionAssign: React.FC<PermissionAssignProps> = ({
  open,
  onOpenChange,
  onFinish,
  roleId,
}) => {
  const [permissionTree, setPermissionTree] = useState<PermissionItem[]>([]);
  const [filteredPermissionTree, setFilteredPermissionTree] = useState<PermissionItem[]>([]);
  const [selectedPermissions, setSelectedPermissions] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [searchValue, setSearchValue] = useState('');

  // 加载权限树
  const loadPermissionTree = async () => {
    try {
      const response = await getPermissionTree();
      if (response.code === 200) {
        setPermissionTree(response.data);
        setFilteredPermissionTree(response.data);
      }
    } catch (error) {
      console.error('加载权限树失败:', error);
      // 错误消息由全局错误处理器显示，这里不再重复显示
    }
  };

  // 加载角色权限
  const loadRolePermissions = async (id: number) => {
    try {
      const response = await getRolePermissions(id);
      if (response.code === 200) {
        // 确保data是数组，防止null或undefined导致的错误
        const data = Array.isArray(response.data) ? response.data : [];
        // 使用权限代码作为标识，与权限树的key保持一致
        const permissionCodes = data.map((p) => p.code || p.id?.toString());
        setSelectedPermissions(permissionCodes.filter(Boolean));
      }
    } catch (error) {
      console.error('加载角色权限失败:', error);
      // 错误消息由全局错误处理器显示，这里不再重复显示
    }
  };

  // 权限筛选函数
  const filterPermissions = (permissions: PermissionItem[], searchText: string): PermissionItem[] => {
    if (!searchText) return permissions;

    return permissions.filter((permission) => {
      const matchesSearch =
        permission.name.toLowerCase().includes(searchText.toLowerCase()) ||
        (permission.code && permission.code.toLowerCase().includes(searchText.toLowerCase()));

      const hasMatchingChildren = permission.children &&
        filterPermissions(permission.children, searchText).length > 0;

      if (matchesSearch || hasMatchingChildren) {
        return {
          ...permission,
          children: permission.children ? filterPermissions(permission.children, searchText) : [],
        };
      }

      return false;
    }).filter(Boolean);
  };

  // 处理搜索
  const handleSearch = (value: string) => {
    setSearchValue(value);
    if (!value) {
      setFilteredPermissionTree(permissionTree);
    } else {
      const filtered = filterPermissions(permissionTree, value);
      setFilteredPermissionTree(filtered);
    }
  };

  // 转换权限树数据格式
  const convertTreeData = (permissions: PermissionItem[]): any[] => {
    return permissions.map((permission) => ({
      title: `${permission.name} (${permission.code})`,
      key: permission.code || permission.id?.toString(), // 使用code作为key，兼容新的权限结构
      children: permission.children ? convertTreeData(permission.children) : [],
    }));
  };

  // 提交权限分配
  const handleSubmit = async () => {
    if (!roleId) {
      message.error('角色ID不存在');
      return;
    }

    setSubmitting(true);
    try {
      await assignRolePermissions(roleId, selectedPermissions);
      message.success('权限分配成功');
      onFinish();
    } catch (error) {
      // 错误消息由全局错误处理器显示，这里不再重复显示
    } finally {
      setSubmitting(false);
    }
  };

  // 加载数据
  const loadData = async () => {
    if (!roleId) return;
    
    setLoading(true);
    try {
      await Promise.all([
        loadPermissionTree(),
        loadRolePermissions(roleId),
      ]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (open && roleId) {
      loadData();
    }
  }, [open, roleId]);

  return (
    <Modal
      title="分配权限"
      open={open}
      onCancel={() => {
        onOpenChange(false);
        setSelectedPermissions([]);
      }}
      onOk={handleSubmit}
      confirmLoading={submitting}
      width={600}
      destroyOnClose
    >
      {loading ? (
        <div style={{ textAlign: 'center', padding: '50px 0' }}>
          <Spin size="large" />
        </div>
      ) : (
        <div>
          <div style={{ marginBottom: 16, color: '#666' }}>
            请选择要分配给该角色的权限：
          </div>
          <div style={{ 
            border: '1px solid #d9d9d9', 
            borderRadius: 6, 
            padding: 12, 
            maxHeight: 400, 
            overflow: 'auto' 
          }}>
            <Input
              placeholder="搜索权限名称或代码"
              prefix={<SearchOutlined />}
              value={searchValue}
              onChange={(e) => handleSearch(e.target.value)}
              style={{ marginBottom: 12 }}
              allowClear
            />
            <Tree
              checkable
              treeData={convertTreeData(filteredPermissionTree)}
              checkedKeys={selectedPermissions}
              onCheck={(checkedKeys) => {
                // 确保类型正确，支持字符串数组
                const keys = Array.isArray(checkedKeys) ? checkedKeys : checkedKeys.checked || [];
                setSelectedPermissions(keys as string[]);
              }}
              defaultExpandAll
            />
          </div>
          <div style={{ marginTop: 16, color: '#999', fontSize: '12px' }}>
            已选择 {selectedPermissions.length} 个权限
          </div>
        </div>
      )}
    </Modal>
  );
};

export default PermissionAssign;

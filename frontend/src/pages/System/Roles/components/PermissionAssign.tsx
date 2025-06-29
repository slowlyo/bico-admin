import React, { useState, useEffect } from 'react';
import { Modal, Tree, message, Spin, Input, Button, Space } from 'antd';
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
        // 后端返回的是权限代码字符串数组，直接使用
        const data = Array.isArray(response.data) ? response.data : [];
        setSelectedPermissions(data);
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

  // 获取所有权限节点的key（排除分类节点）
  const getAllPermissionKeys = (permissions: PermissionItem[]): string[] => {
    let keys: string[] = [];
    permissions.forEach(permission => {
      const code = permission.code || permission.id?.toString();
      if (code && !isCategoryNode(code)) {
        keys.push(code);
      }
      if (permission.children) {
        keys = keys.concat(getAllPermissionKeys(permission.children));
      }
    });
    return keys;
  };

  // 全选（基于当前显示的权限）
  const handleSelectAll = () => {
    const currentKeys = getAllPermissionKeys(filteredPermissionTree);
    // 合并当前选中的权限和当前显示的所有权限，去重
    const allKeys = Array.from(new Set([...selectedPermissions, ...currentKeys]));
    setSelectedPermissions(allKeys);
  };

  // 反选（反转当前显示权限的选择状态）
  const handleInvertSelection = () => {
    const currentKeys = getAllPermissionKeys(filteredPermissionTree);
    const newSelectedKeys = [...selectedPermissions];

    currentKeys.forEach(key => {
      const index = newSelectedKeys.indexOf(key);
      if (index > -1) {
        // 如果已选中，则取消选中
        newSelectedKeys.splice(index, 1);
      } else {
        // 如果未选中，则选中
        newSelectedKeys.push(key);
      }
    });

    setSelectedPermissions(newSelectedKeys);
  };

  // 检查是否为分类节点（分类节点的code通常是中文名称，不是有效的权限代码）
  const isCategoryNode = (code: string): boolean => {
    // 有效的权限代码格式为 "module:action"，如 "user:view"
    return !code.includes(':');
  };

  // 转换权限树数据格式
  const convertTreeData = (permissions: PermissionItem[]): any[] => {
    return permissions.map((permission) => {
      const code = permission.code || permission.id?.toString();
      const isCategory = isCategoryNode(code);

      return {
        title: isCategory ? permission.name : `${permission.name}`,
        key: code,
        children: permission.children ? convertTreeData(permission.children) : [],
      };
    });
  };

  // 提交权限分配
  const handleSubmit = async () => {
    if (!roleId) {
      message.error('角色ID不存在');
      return;
    }

    // 过滤掉分类节点，只提交有效的权限代码
    const validPermissions = selectedPermissions.filter(code => !isCategoryNode(code));

    setSubmitting(true);
    try {
      await assignRolePermissions(roleId, validPermissions);
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
            <div style={{ marginBottom: 12, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <Input
                placeholder="搜索权限名称"
                prefix={<SearchOutlined />}
                value={searchValue}
                onChange={(e) => handleSearch(e.target.value)}
                style={{ flex: 1, marginRight: 12 }}
                allowClear
              />
              <Space>
                <Button onClick={handleSelectAll}>
                  全选
                </Button>
                <Button onClick={handleInvertSelection}>
                  反选
                </Button>
              </Space>
            </div>
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
            已选择 {selectedPermissions.filter(code => !isCategoryNode(code)).length} 个权限
            {selectedPermissions.length > selectedPermissions.filter(code => !isCategoryNode(code)).length &&
              ` (包含 ${selectedPermissions.length - selectedPermissions.filter(code => !isCategoryNode(code)).length} 个分类)`
            }
          </div>
        </div>
      )}
    </Modal>
  );
};

export default PermissionAssign;

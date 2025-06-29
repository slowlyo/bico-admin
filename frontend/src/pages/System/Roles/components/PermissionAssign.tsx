import React, { useState, useEffect } from 'react';
import { Modal, Tree, message, Spin } from 'antd';
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
  const [selectedPermissions, setSelectedPermissions] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  // 加载权限树
  const loadPermissionTree = async () => {
    try {
      const response = await getPermissionTree();
      if (response.code === 200) {
        setPermissionTree(response.data);
      }
    } catch (error) {
      console.error('加载权限树失败:', error);
      message.error('加载权限树失败');
    }
  };

  // 加载角色权限
  const loadRolePermissions = async (id: number) => {
    try {
      const response = await getRolePermissions(id);
      if (response.code === 200) {
        // 新的权限系统使用code作为标识
        const permissionCodes = response.data.map((p) => p.code || p.id?.toString());
        setSelectedPermissions(permissionCodes.filter(Boolean));
      }
    } catch (error) {
      console.error('加载角色权限失败:', error);
      message.error('加载角色权限失败');
    }
  };

  // 转换权限树数据格式
  const convertTreeData = (permissions: PermissionItem[]): any[] => {
    return permissions.map((permission) => ({
      title: `${permission.name} (${permission.code})`,
      key: permission.code || permission.id, // 使用code作为key，兼容新的权限结构
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
      message.error('权限分配失败，请重试');
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
            <Tree
              checkable
              treeData={convertTreeData(permissionTree)}
              checkedKeys={selectedPermissions}
              onCheck={(checkedKeys) => {
                setSelectedPermissions(checkedKeys as number[]);
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

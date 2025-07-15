import React, { useEffect, useState } from 'react';
import { Drawer, Tree, Button, message, Spin, Space } from 'antd';
import type { DataNode } from 'antd/es/tree';
import {
  getPermissionTree,
  updateRolePermissions,
  type PermissionTreeNode,
  type Role
} from '@/services/role';

interface PermissionDrawerProps {
  visible: boolean;
  role: Role | null;
  onClose: () => void;
  onSuccess: () => void;
}

const PermissionDrawer: React.FC<PermissionDrawerProps> = ({
  visible,
  role,
  onClose,
  onSuccess,
}) => {
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [treeData, setTreeData] = useState<DataNode[]>([]);
  const [checkedKeys, setCheckedKeys] = useState<React.Key[]>([]);
  const [allPermissionKeys, setAllPermissionKeys] = useState<string[]>([]);
  const [expandedKeys, setExpandedKeys] = useState<React.Key[]>([]);

  // 转换权限树数据为Tree组件需要的格式
  const convertToTreeData = (permissions: PermissionTreeNode[]): DataNode[] => {
    const allKeys: string[] = [];
    const firstLevelKeys: string[] = [];

    const convertNode = (node: PermissionTreeNode): DataNode => {
      // 只有action类型的权限才加入到allKeys中（用于权限提交）
      if (node.type === 'action') {
        allKeys.push(node.key);
      }

      return {
        title: node.title,
        key: node.key,
        children: node.children.map(child => convertNode(child)),
      };
    };

    const treeData = permissions.map(node => {
      // 收集第一层的keys用于默认展开
      firstLevelKeys.push(node.key);
      return convertNode(node);
    });

    setAllPermissionKeys(allKeys);
    setExpandedKeys(firstLevelKeys);
    setTreeData(treeData);
    return treeData;
  };

  // 递归获取已选中的权限
  const getSelectedPermissions = (permissions: PermissionTreeNode[]): string[] => {
    const selected: string[] = [];

    const collectSelected = (nodes: PermissionTreeNode[]) => {
      nodes.forEach(node => {
        if (node.selected && node.type === 'action') {
          selected.push(node.key);
        }
        if (node.children && node.children.length > 0) {
          collectSelected(node.children);
        }
      });
    };

    collectSelected(permissions);
    return selected;
  };

  // 加载权限树
  const loadPermissionTree = async () => {
    if (!role) return;

    setLoading(true);
    try {
      const response = await getPermissionTree(role.id);
      if (response.code === 200) {
        convertToTreeData(response.data);
        setCheckedKeys(getSelectedPermissions(response.data));
      } else {
        message.error(response.message || '获取权限树失败');
      }
    } catch (error) {
      message.error('获取权限树失败');
    } finally {
      setLoading(false);
    }
  };

  // 处理权限选择变化
  const handleCheck = (checkedKeysValue: React.Key[] | { checked: React.Key[]; halfChecked: React.Key[] }) => {
    if (Array.isArray(checkedKeysValue)) {
      setCheckedKeys(checkedKeysValue);
    } else {
      setCheckedKeys(checkedKeysValue.checked);
    }
  };

  // 全选
  const handleSelectAll = () => {
    setCheckedKeys(allPermissionKeys);
  };

  // 反选
  const handleInvertSelection = () => {
    const invertedKeys = allPermissionKeys.filter(key => !checkedKeys.includes(key));
    setCheckedKeys(invertedKeys);
  };

  // 清空
  const handleClearAll = () => {
    setCheckedKeys([]);
  };



  // 提交权限配置
  const handleSubmit = async () => {
    if (!role) return;

    setSubmitting(true);
    try {
      // 过滤出action类型的权限代码（只有action类型的权限才需要提交）
      const permissionCodes = checkedKeys.filter(key => {
        return typeof key === 'string' && allPermissionKeys.includes(key as string);
      }) as string[];

      const response = await updateRolePermissions(role.id, {
        permissions: permissionCodes,
      });

      if (response.code === 200) {
        message.success('权限配置更新成功');
        onSuccess();
        onClose();
      } else {
        message.error(response.message || '权限配置更新失败');
      }
    } catch (error) {
      message.error('权限配置更新失败');
    } finally {
      setSubmitting(false);
    }
  };

  // 当抽屉打开时加载数据
  useEffect(() => {
    if (visible && role) {
      loadPermissionTree();
    } else if (!visible) {
      // 抽屉关闭时重置状态
      setExpandedKeys([]);
      setCheckedKeys([]);
      setTreeData([]);
    }
  }, [visible, role]);

  return (
    <Drawer
      title={`配置角色权限 - ${role?.name}`}
      width={800}
      open={visible}
      onClose={onClose}
      footer={
        <div style={{ textAlign: 'right' }}>
          <Space>
            <Button onClick={onClose}>取消</Button>
            <Button 
              type="primary" 
              loading={submitting}
              onClick={handleSubmit}
            >
              保存
            </Button>
          </Space>
        </div>
      }
    >
      <Spin spinning={loading}>
        <div style={{ marginBottom: 16 }}>
          <Space>
            <Button size="small" onClick={handleSelectAll}>
              全选
            </Button>
            <Button size="small" onClick={handleInvertSelection}>
              反选
            </Button>
            <Button size="small" onClick={handleClearAll}>
              清空
            </Button>
          </Space>
        </div>

        <div
          style={{
            border: '1px solid #d9d9d9',
            borderRadius: 6,
            padding: 8,
          }}
        >
          <Tree
            checkable
            checkedKeys={checkedKeys}
            onCheck={handleCheck}
            height={400}
            treeData={treeData}
            expandedKeys={expandedKeys}
            onExpand={setExpandedKeys}
          />
        </div>
      </Spin>
    </Drawer>
  );
};

export default PermissionDrawer;

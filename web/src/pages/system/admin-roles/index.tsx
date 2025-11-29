/**
 * 角色管理 - 使用 CrudTable 重构
 */
import type { ProColumns } from '@ant-design/pro-components';
import { ProFormText, ProFormTextArea, ProFormSwitch } from '@ant-design/pro-components';
import { Tag, Space, Drawer, Tree, Button, message } from 'antd';
import React, { useState, useEffect, useCallback, useRef } from 'react';
import type { ActionType } from '@ant-design/pro-components';
import { useAccess } from '@umijs/max';
import { CrudTable } from '@/components';
import { createCrudService } from '@/services/crud';
import { getAllPermissions, updateRolePermissions } from '@/services/system/admin-role';
import type { Permission } from '@/services/system/admin-role/types';

// 类型定义
interface AdminRole {
  id: number;
  name: string;
  code: string;
  description: string;
  enabled: boolean;
  permissions?: string[];
  created_at: string;
}

// CRUD 服务
const roleService = createCrudService<AdminRole>('/admin-roles');

// 列配置
const columns: ProColumns<AdminRole>[] = [
  { title: 'ID', dataIndex: 'id', width: 80, search: false },
  { title: '角色名称', dataIndex: 'name', width: 150 },
  { title: '角色代码', dataIndex: 'code', width: 150 },
  { title: '描述', dataIndex: 'description', search: false, width: 200, ellipsis: true },
  {
    title: '状态',
    dataIndex: 'enabled',
    width: 100,
    valueType: 'select',
    valueEnum: { true: { text: '启用', status: 'Success' }, false: { text: '禁用', status: 'Default' } },
    render: (_, r) => <Tag color={r.enabled ? 'green' : 'red'}>{r.enabled ? '启用' : '禁用'}</Tag>,
  },
  { title: '创建时间', dataIndex: 'created_at', valueType: 'dateTime', width: 180, search: false, sorter: true },
];

// 表单内容组件
const FormContent: React.FC<{ record?: AdminRole }> = ({ record }) => {
  const isEdit = !!record;
  return (
    <>
      <ProFormText name="name" label="角色名称" placeholder="请输入角色名称" rules={[{ required: true }]} />
      {!isEdit && (
        <ProFormText name="code" label="角色代码" placeholder="请输入角色代码" rules={[{ required: true }]} />
      )}
      <ProFormTextArea name="description" label="描述" placeholder="请输入描述" />
      <ProFormSwitch name="enabled" label="状态" initialValue={true} />
    </>
  );
};

// 权限树工具函数
const convertToTreeData = (permissions: Permission[]): any[] =>
  permissions.map((p) => ({ title: p.label, key: p.key, children: p.children ? convertToTreeData(p.children) : undefined }));

const findParents = (key: string, tree: Permission[], path: string[] = []): string[] | null => {
  for (const node of tree) {
    if (node.key === key) return path;
    if (node.children) {
      const result = findParents(key, node.children, [...path, node.key]);
      if (result) return result;
    }
  }
  return null;
};

const filterRedundant = (perms: string[], tree: Permission[]): string[] => {
  const filtered = new Set<string>();
  perms.forEach((p) => {
    const hasChild = perms.some((other) => {
      if (other === p) return false;
      const parents = findParents(other, tree);
      return parents?.includes(p);
    });
    if (!hasChild) filtered.add(p);
  });
  return Array.from(filtered);
};

const expandPerms = (perms: string[], tree: Permission[]): string[] => {
  const expanded = new Set<string>();
  perms.forEach((p) => {
    expanded.add(p);
    findParents(p, tree)?.forEach((parent) => expanded.add(parent));
  });
  return Array.from(expanded);
};

export default function AdminRoleList() {
  const access = useAccess() as Record<string, boolean>;
  const actionRef = useRef<ActionType>(null);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [currentRole, setCurrentRole] = useState<AdminRole>();
  const [allPermissions, setAllPermissions] = useState<Permission[]>([]);
  const [selectedKeys, setSelectedKeys] = useState<string[]>([]);

  useEffect(() => {
    getAllPermissions().then((res) => res.code === 0 && setAllPermissions(res.data || []));
  }, []);

  const handleOpenDrawer = useCallback((record: AdminRole) => {
    setCurrentRole(record);
    setSelectedKeys(filterRedundant(record.permissions || [], allPermissions));
    setDrawerOpen(true);
  }, [allPermissions]);

  const handleSavePermissions = useCallback(async () => {
    if (!currentRole) return;
    try {
      const res = await updateRolePermissions(currentRole.id, { permissions: expandPerms(selectedKeys, allPermissions) });
      if (res.code === 0) {
        message.success('权限配置成功');
        setDrawerOpen(false);
        actionRef.current?.reload();
      } else {
        message.error(res.msg || '配置失败');
      }
    } catch (e: any) {
      message.error(e.message || '配置失败');
    }
  }, [currentRole, selectedKeys, allPermissions]);

  return (
    <>
      <CrudTable<AdminRole>
        title="角色"
        permissionPrefix="system:admin_role"
        service={roleService}
        columns={columns}
        formContent={<FormContent />}
        recordToValues={(r) => ({ name: r.name, description: r.description, enabled: r.enabled })}
        transformParams={(params) => ({
          ...params,
          enabled: params.enabled === 'true' ? true : params.enabled === 'false' ? false : undefined,
        })}
        scrollX={1100}
        actionRef={actionRef}
        renderActions={(record, defaultActions) => (
          <Space>
            {access['system:admin_role:permission'] && <a onClick={() => handleOpenDrawer(record)}>配置权限</a>}
            {defaultActions}
          </Space>
        )}
      />

      <Drawer
        title="配置权限"
        width={500}
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        extra={
          <Space>
            <Button onClick={() => setDrawerOpen(false)}>取消</Button>
            <Button type="primary" onClick={handleSavePermissions}>保存</Button>
          </Space>
        }
      >
        <Tree
          checkable
          defaultExpandAll
          checkedKeys={selectedKeys}
          onCheck={(keys) => setSelectedKeys(Array.isArray(keys) ? (keys as string[]) : (keys.checked as string[]))}
          treeData={convertToTreeData(allPermissions)}
        />
      </Drawer>
    </>
  );
}

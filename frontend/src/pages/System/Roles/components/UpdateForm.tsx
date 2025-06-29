import React, { useState, useEffect } from 'react';
import {
  ModalForm,
  ProFormText,
  ProFormTextArea,
  ProFormSelect,
} from '@ant-design/pro-components';
import { message, Tree } from 'antd';
import { updateRole, getPermissionTree, getRolePermissions } from '@/services/role';
import type { RoleItem, PermissionItem } from '@/services/role';

export type UpdateFormProps = {
  open: boolean;
  onOpenChange: (visible: boolean) => void;
  onFinish: () => void;
  values: Partial<RoleItem>;
};

const UpdateForm: React.FC<UpdateFormProps> = ({
  open,
  onOpenChange,
  onFinish,
  values,
}) => {
  const [permissionTree, setPermissionTree] = useState<PermissionItem[]>([]);
  const [selectedPermissions, setSelectedPermissions] = useState<string[]>([]);

  // 加载权限树
  const loadPermissionTree = async () => {
    try {
      const response = await getPermissionTree();
      if (response.code === 200) {
        setPermissionTree(response.data);
      }
    } catch (error) {
      console.error('加载权限树失败:', error);
    }
  };

  // 加载角色权限
  const loadRolePermissions = async (roleId: number) => {
    try {
      const response = await getRolePermissions(roleId);
      if (response.code === 200) {
        // 确保data是数组，防止null或undefined导致的错误
        const data = Array.isArray(response.data) ? response.data : [];
        // 使用权限代码作为标识，与权限树的key保持一致
        const permissionCodes = data.map((p) => p.code || p.id?.toString());
        setSelectedPermissions(permissionCodes.filter(Boolean));
      }
    } catch (error) {
      console.error('加载角色权限失败:', error);
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

  const handleSubmit = async (formValues: any) => {
    if (!values.id) {
      message.error('角色ID不存在');
      return false;
    }

    try {
      await updateRole(values.id, {
        ...formValues,
        permission_ids: selectedPermissions,
      });
      message.success('更新成功');
      onFinish();
      return true;
    } catch (error) {
      // 错误消息由全局错误处理器显示，这里不再重复显示
      return false;
    }
  };

  useEffect(() => {
    if (open && values.id) {
      loadPermissionTree();
      loadRolePermissions(values.id);
    }
  }, [open, values.id]);

  return (
    <ModalForm
      title="编辑角色"
      width="600px"
      open={open}
      onOpenChange={(visible) => {
        onOpenChange(visible);
        if (!visible) {
          setSelectedPermissions([]);
        }
      }}
      onFinish={handleSubmit}
      initialValues={values}
      modalProps={{
        destroyOnClose: true,
      }}
    >
      <ProFormText
        name="name"
        label="角色名称"
        placeholder="请输入角色名称"
        rules={[
          { required: true, message: '请输入角色名称!' },
          { max: 50, message: '角色名称不能超过50个字符!' },
        ]}
      />
      
      <ProFormText
        name="code"
        label="角色代码"
        placeholder="请输入角色代码"
        rules={[
          { required: true, message: '请输入角色代码!' },
          { max: 50, message: '角色代码不能超过50个字符!' },
          { pattern: /^[a-zA-Z0-9_-]+$/, message: '角色代码只能包含字母、数字、下划线和横线!' },
        ]}
      />
      
      <ProFormTextArea
        name="description"
        label="描述"
        placeholder="请输入角色描述"
        rules={[
          { max: 255, message: '描述不能超过255个字符!' },
        ]}
      />

      <ProFormSelect
        name="status"
        label="状态"
        options={[
          { label: '启用', value: 1 },
          { label: '禁用', value: 0 },
        ]}
        placeholder="请选择状态"
        rules={[{ required: true, message: '请选择状态!' }]}
      />

      <div style={{ marginBottom: 24 }}>
        <label style={{ fontWeight: 500, marginBottom: 8, display: 'block' }}>
          权限分配
        </label>
        <div style={{ border: '1px solid #d9d9d9', borderRadius: 6, padding: 12, maxHeight: 300, overflow: 'auto' }}>
          <Tree
            checkable
            treeData={convertTreeData(permissionTree)}
            checkedKeys={selectedPermissions}
            onCheck={(checkedKeys) => {
              // 确保类型正确，支持字符串数组
              const keys = Array.isArray(checkedKeys) ? checkedKeys : checkedKeys.checked || [];
              setSelectedPermissions(keys as string[]);
            }}
            placeholder="请选择权限"
          />
        </div>
      </div>
    </ModalForm>
  );
};

export default UpdateForm;

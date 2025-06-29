import React, { useState } from 'react';
import {
  ModalForm,
  ProFormText,
  ProFormTextArea,
  ProFormSelect,
} from '@ant-design/pro-components';
import { message, Tree } from 'antd';
import { createRole, getPermissionTree } from '@/services/role';
import type { PermissionItem } from '@/services/role';

export type CreateFormProps = {
  open: boolean;
  onOpenChange: (visible: boolean) => void;
  onFinish: () => void;
};

const CreateForm: React.FC<CreateFormProps> = ({
  open,
  onOpenChange,
  onFinish,
}) => {
  const [permissionTree, setPermissionTree] = useState<PermissionItem[]>([]);
  const [selectedPermissions, setSelectedPermissions] = useState<number[]>([]);

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

  // 转换权限树数据格式
  const convertTreeData = (permissions: PermissionItem[]): any[] => {
    return permissions.map((permission) => ({
      title: permission.name,
      key: permission.id,
      children: permission.children ? convertTreeData(permission.children) : [],
    }));
  };

  const handleSubmit = async (values: any) => {
    try {
      await createRole({
        ...values,
        permission_ids: selectedPermissions,
      });
      message.success('创建成功');
      onFinish();
      setSelectedPermissions([]);
      return true;
    } catch (error) {
      // 错误消息由全局错误处理器显示，这里不再重复显示
      return false;
    }
  };

  return (
    <ModalForm
      title="新建角色"
      width="600px"
      open={open}
      onOpenChange={(visible) => {
        onOpenChange(visible);
        if (visible) {
          loadPermissionTree();
        } else {
          setSelectedPermissions([]);
        }
      }}
      onFinish={handleSubmit}
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
        initialValue={1}
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
              setSelectedPermissions(checkedKeys as number[]);
            }}
            placeholder="请选择权限"
          />
        </div>
      </div>
    </ModalForm>
  );
};

export default CreateForm;

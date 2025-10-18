import { ModalForm, ProFormText, ProFormSwitch, ProFormSelect } from '@ant-design/pro-components';
import { message } from 'antd';
import React from 'react';
import { createAdminUser, type AdminUserCreateParams } from '@/services/admin-user';
import { getAllAdminRoles } from '@/services/admin-role';

interface CreateFormProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
}

const CreateForm: React.FC<CreateFormProps> = ({ open, onOpenChange, onSuccess }) => {
  const handleCreate = async (values: AdminUserCreateParams) => {
    try {
      const res = await createAdminUser(values);
      if (res.code === 0) {
        message.success('创建成功');
        onSuccess();
        return true;
      }
      message.error(res.msg || '创建失败');
      return false;
    } catch (error) {
      message.error('创建失败');
      return false;
    }
  };

  return (
    <ModalForm
      title="新建用户"
      width={600}
      open={open}
      onOpenChange={onOpenChange}
      onFinish={handleCreate}
    >
      <ProFormText
        name="username"
        label="用户名"
        placeholder="请输入用户名"
        rules={[{ required: true, message: '请输入用户名' }]}
      />
      <ProFormText.Password
        name="password"
        label="密码"
        placeholder="请输入密码"
        rules={[{ required: true, message: '请输入密码' }]}
      />
      <ProFormText
        name="name"
        label="姓名"
        placeholder="请输入姓名"
      />
      <ProFormSelect
        name="roleIds"
        label="角色"
        mode="multiple"
        request={async () => {
          const res = await getAllAdminRoles();
          return (res.data || []).map((role) => ({
            label: role.name,
            value: role.id,
          }));
        }}
      />
      <ProFormSwitch
        name="enabled"
        label="状态"
        initialValue={true}
      />
    </ModalForm>
  );
};

export default CreateForm;

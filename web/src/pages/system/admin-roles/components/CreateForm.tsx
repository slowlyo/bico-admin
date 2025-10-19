import { ModalForm, ProFormText, ProFormTextArea, ProFormSwitch } from '@ant-design/pro-components';
import { message } from 'antd';
import React from 'react';
import { createAdminRole, type AdminRoleCreateParams } from '@/services/admin-role';

interface CreateFormProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
}

const CreateForm: React.FC<CreateFormProps> = ({ open, onOpenChange, onSuccess }) => {
  const handleCreate = async (values: AdminRoleCreateParams) => {
    try {
      const res = await createAdminRole(values);
      if (res.code === 0) {
        message.success('创建成功');
        onSuccess();
        return true;
      }
      message.error(res.msg || '创建失败');
      return false;
    } catch (error: any) {
      message.error(error.message || error.data?.msg || '创建失败');
      return false;
    }
  };

  return (
    <ModalForm
      title="新建角色"
      width={600}
      open={open}
      onOpenChange={onOpenChange}
      onFinish={handleCreate}
      modalProps={{
        destroyOnClose: true,
      }}
    >
      <ProFormText
        name="name"
        label="角色名称"
        placeholder="请输入角色名称"
        rules={[{ required: true, message: '请输入角色名称' }]}
      />
      <ProFormText
        name="code"
        label="角色代码"
        placeholder="请输入角色代码"
        rules={[{ required: true, message: '请输入角色代码' }]}
      />
      <ProFormTextArea
        name="description"
        label="描述"
        placeholder="请输入描述"
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

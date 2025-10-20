import { ModalForm, ProFormText, ProFormTextArea, ProFormSwitch } from '@ant-design/pro-components';
import { message } from 'antd';
import React from 'react';
import { updateAdminRole, type AdminRole, type AdminRoleUpdateParams } from '@/services/admin-role';

interface UpdateFormProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
  currentRow?: AdminRole;
}

const UpdateForm: React.FC<UpdateFormProps> = ({ open, onOpenChange, onSuccess, currentRow }) => {
  const handleUpdate = async (values: AdminRoleUpdateParams) => {
    if (!currentRow) return false;
    try {
      const res = await updateAdminRole(currentRow.id, values);
      if (res.code === 0) {
        message.success('更新成功');
        onSuccess();
        return true;
      }
      message.error(res.msg || '更新失败');
      return false;
    } catch (error: any) {
      message.error(error.message || error.data?.msg || '更新失败');
      return false;
    }
  };

  return (
    <ModalForm
      title="编辑角色"
      width={600}
      open={open}
      onOpenChange={onOpenChange}
      initialValues={{
        name: currentRow?.name,
        description: currentRow?.description,
        enabled: currentRow?.enabled,
      }}
      onFinish={handleUpdate}
      modalProps={{
        destroyOnHidden: true,
      }}
    >
      <ProFormText
        name="name"
        label="角色名称"
        placeholder="请输入角色名称"
        rules={[{ required: true, message: '请输入角色名称' }]}
      />
      <ProFormTextArea
        name="description"
        label="描述"
        placeholder="请输入描述"
      />
      <ProFormSwitch
        name="enabled"
        label="状态"
      />
    </ModalForm>
  );
};

export default UpdateForm;

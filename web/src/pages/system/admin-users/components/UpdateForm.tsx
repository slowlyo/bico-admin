import { ModalForm, ProFormText, ProFormSwitch, ProFormSelect } from '@ant-design/pro-components';
import { message } from 'antd';
import React from 'react';
import { updateAdminUser, type AdminUser, type AdminUserUpdateParams } from '@/services/admin-user';
import { getAllAdminRoles } from '@/services/admin-role';

interface UpdateFormProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
  currentRow?: AdminUser;
}

const UpdateForm: React.FC<UpdateFormProps> = ({ open, onOpenChange, onSuccess, currentRow }) => {
  const handleUpdate = async (values: AdminUserUpdateParams) => {
    if (!currentRow) return false;
    try {
      const res = await updateAdminUser(currentRow.id, values);
      if (res.code === 0) {
        message.success('更新成功');
        onSuccess();
        return true;
      }
      message.error(res.msg || '更新失败');
      return false;
    } catch (error) {
      message.error('更新失败');
      return false;
    }
  };

  return (
    <ModalForm
      title="编辑用户"
      width={600}
      open={open}
      onOpenChange={onOpenChange}
      initialValues={{
        name: currentRow?.name,
        enabled: currentRow?.enabled,
        roleIds: currentRow?.roles?.map((role) => role.id),
      }}
      onFinish={handleUpdate}
    >
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
      />
    </ModalForm>
  );
};

export default UpdateForm;

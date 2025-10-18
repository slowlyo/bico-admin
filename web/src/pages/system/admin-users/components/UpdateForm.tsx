import { ModalForm, ProFormText, ProFormSwitch, ProFormSelect } from '@ant-design/pro-components';
import { message, Avatar, Space, Upload, Button } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import React, { useState, useEffect } from 'react';
import { updateAdminUser, type AdminUser, type AdminUserUpdateParams } from '@/services/admin-user';
import { getAllAdminRoles } from '@/services/admin-role';

interface UpdateFormProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
  currentRow?: AdminUser;
}

const UpdateForm: React.FC<UpdateFormProps> = ({ open, onOpenChange, onSuccess, currentRow }) => {
  const [avatarUrl, setAvatarUrl] = useState<string>('');

  // 当弹窗打开或currentRow变化时，更新头像
  useEffect(() => {
    if (open && currentRow?.avatar) {
      setAvatarUrl(currentRow.avatar);
    }
  }, [open, currentRow]);

  const handleUpdate = async (values: AdminUserUpdateParams) => {
    if (!currentRow) return false;
    try {
      // 使用当前的头像URL
      const params = { ...values, avatar: avatarUrl };
      const res = await updateAdminUser(currentRow.id, params);
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

  const handleUploadChange = (info: any) => {
    if (info.file.status === 'done') {
      const url = info.file.response?.data?.url;
      if (url) {
        setAvatarUrl(url);
        message.success('头像上传成功');
      }
    } else if (info.file.status === 'error') {
      message.error('头像上传失败');
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

      <div style={{ marginBottom: 24 }}>
        <div style={{ marginBottom: 8, fontSize: 14, fontWeight: 500 }}>头像</div>
        <Space size={16}>
          <Avatar src={avatarUrl} size={80} />
          <div>
            <Upload
              name="avatar"
              showUploadList={false}
              action="/admin-api/auth/avatar"
              headers={{
                Authorization: `Bearer ${localStorage.getItem('token')}`,
              }}
              onChange={handleUploadChange}
            >
              <Button icon={<UploadOutlined />}>上传头像</Button>
            </Upload>
          </div>
        </Space>
      </div>

      <ProFormSelect
        name="roleIds"
        label="角色"
        mode="multiple"
        request={async () => {
          const res = await getAllAdminRoles();
          return (res.data || []).map((role: any) => ({
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

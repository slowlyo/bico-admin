import { ModalForm, ProFormText, ProFormSwitch, ProFormSelect } from '@ant-design/pro-components';
import { message, Avatar, Space, Upload, Button } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import React, { useState, useEffect } from 'react';
import { createAdminUser } from '@/services/system/admin-user';
import type { AdminUserCreateParams } from '@/services/system/admin-user/types';
import { getAllAdminRoles } from '@/services/system/admin-role';
import type { AdminRole } from '@/services/system/admin-role/types';
import { buildApiUrl } from '@/services/config';

interface CreateFormProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
}

// 生成随机头像URL
const generateRandomAvatar = () => {
  const randomSeed = Math.floor(Math.random() * 999999);
  return `https://api.dicebear.com/9.x/thumbs/png?seed=${randomSeed}`;
};

const CreateForm: React.FC<CreateFormProps> = ({ open, onOpenChange, onSuccess }) => {
  const [avatarUrl, setAvatarUrl] = useState<string>('');

  // 当弹窗打开时生成新的随机头像
  useEffect(() => {
    if (open) {
      setAvatarUrl(generateRandomAvatar());
    }
  }, [open]);

  const handleCreate = async (values: AdminUserCreateParams) => {
    try {
      // 使用当前的头像URL
      const params = { ...values, avatar: avatarUrl };
      const res = await createAdminUser(params);
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

  const handleUploadChange = (info: any) => {
    if (info.file.status === 'done') {
      // 假设上传接口返回的图片URL在 response.data.url
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
      title="新建用户"
      width={600}
      open={open}
      onOpenChange={onOpenChange}
      onFinish={handleCreate}
      modalProps={{
        destroyOnHidden: true,
      }}
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
      
      <div style={{ marginBottom: 24 }}>
        <div style={{ marginBottom: 8, fontSize: 14, fontWeight: 500 }}>头像</div>
        <Space size={16}>
          <Avatar src={avatarUrl} size={80} />
          <div>
            <Space direction="vertical" size={8}>
              <Upload
                name="avatar"
                showUploadList={false}
                action={buildApiUrl('/auth/avatar')}
                headers={{
                  Authorization: `Bearer ${localStorage.getItem('token')}`,
                }}
                onChange={handleUploadChange}
              >
                <Button icon={<UploadOutlined />}>上传自定义头像</Button>
              </Upload>
              <div style={{ fontSize: 12, color: '#999' }}>
                默认使用随机头像，可上传自定义图片
              </div>
            </Space>
          </div>
        </Space>
      </div>

      <ProFormSelect
        name="roleIds"
        label="角色"
        mode="multiple"
        request={async () => {
          const res = await getAllAdminRoles();
          return (res.data || []).map((role: AdminRole) => ({
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

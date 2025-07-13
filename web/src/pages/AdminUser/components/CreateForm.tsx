import { ModalForm, ProFormText, ProFormSwitch } from '@ant-design/pro-components';
import React from 'react';
import { AdminUserCreateRequest } from '@/services/adminUser';

export interface CreateFormProps {
  modalVisible: boolean;
  onCancel: () => void;
  onSubmit: (values: AdminUserCreateRequest) => Promise<void>;
}

const CreateForm: React.FC<CreateFormProps> = (props) => {
  const { modalVisible, onCancel, onSubmit } = props;

  return (
    <ModalForm
      title="新建管理员用户"
      width="400px"
      open={modalVisible}
      onOpenChange={(visible) => {
        if (!visible) {
          onCancel();
        }
      }}
      onFinish={async (value) => {
        await onSubmit(value as AdminUserCreateRequest);
      }}
    >
      <ProFormText
        name="username"
        label="用户名"
        rules={[
          {
            required: true,
            message: '用户名为必填项',
          },
          {
            min: 3,
            max: 50,
            message: '用户名长度为3-50个字符',
          },
        ]}
        placeholder="请输入用户名"
      />
      
      <ProFormText.Password
        name="password"
        label="密码"
        rules={[
          {
            required: true,
            message: '密码为必填项',
          },
          {
            min: 6,
            max: 100,
            message: '密码长度为6-100个字符',
          },
        ]}
        placeholder="请输入密码"
      />
      
      <ProFormText
        name="name"
        label="姓名"
        rules={[
          {
            required: true,
            message: '姓名为必填项',
          },
          {
            max: 100,
            message: '姓名长度不能超过100个字符',
          },
        ]}
        placeholder="请输入姓名"
      />
      
      <ProFormText
        name="email"
        label="邮箱"
        rules={[
          {
            type: 'email',
            message: '请输入正确的邮箱格式',
          },
        ]}
        placeholder="请输入邮箱（可选）"
      />
      
      <ProFormText
        name="phone"
        label="手机号"
        rules={[
          {
            pattern: /^1[3-9]\d{9}$/,
            message: '请输入正确的手机号格式',
          },
        ]}
        placeholder="请输入手机号（可选）"
      />
      
      <ProFormText
        name="avatar"
        label="头像URL"
        placeholder="请输入头像URL（可选）"
      />
      
      <ProFormSwitch
        name="enabled"
        label="启用状态"
        initialValue={true}
      />
    </ModalForm>
  );
};

export default CreateForm;

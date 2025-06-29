import React from 'react';
import {
  ModalForm,
  ProFormText,
  ProFormSelect,
  ProFormDigit,
} from '@ant-design/pro-components';
import { message } from 'antd';
import { createUser } from '@/services/user';

export interface CreateFormProps {
  open: boolean;
  onOpenChange: (visible: boolean) => void;
  onFinish: () => void;
}

const CreateForm: React.FC<CreateFormProps> = ({
  open,
  onOpenChange,
  onFinish,
}) => {
  return (
    <ModalForm
      title="新建用户"
      width="400px"
      open={open}
      onOpenChange={onOpenChange}
      onFinish={async (value) => {
        try {
          const response = await createUser(value);
          if (response.code === 200) {
            message.success('创建成功');
            onFinish();
            return true;
          } else {
            message.error(response.message || '创建失败');
            return false;
          }
        } catch (error: any) {
          message.error(error.message || '创建失败，请重试');
          return false;
        }
      }}
    >
      <ProFormText
        rules={[
          {
            required: true,
            message: '用户名为必填项',
          },
        ]}
        label="用户名"
        name="username"
        placeholder="请输入用户名"
      />
      <ProFormText
        rules={[
          {
            type: 'email',
            message: '请输入有效的邮箱地址',
          },
        ]}
        label="邮箱"
        name="email"
        placeholder="请输入邮箱"
      />
      <ProFormText.Password
        rules={[
          {
            required: true,
            message: '密码为必填项',
          },
          {
            min: 6,
            message: '密码至少6位',
          },
        ]}
        label="密码"
        name="password"
        placeholder="请输入密码"
      />
      <ProFormText
        label="昵称"
        name="nickname"
        placeholder="请输入昵称"
      />
      <ProFormText
        label="手机号"
        name="phone"
        placeholder="请输入手机号"
        rules={[
          {
            pattern: /^1[3-9]\d{9}$/,
            message: '请输入有效的手机号',
          },
        ]}
      />
      <ProFormSelect
        name="role"
        label="角色"
        valueEnum={{
          'admin': '管理员',
          'manager': '管理者',
          'user': '普通用户',
        }}
        placeholder="请选择角色"
        initialValue="user"
        rules={[{ required: true, message: '请选择角色!' }]}
      />
      <ProFormSelect
        name="status"
        label="状态"
        options={[
          { label: '正常', value: 1 },
          { label: '禁用', value: 0 },
        ]}
        placeholder="请选择状态"
        initialValue={1}
        rules={[{ required: true, message: '请选择状态!' }]}
      />
    </ModalForm>
  );
};

export default CreateForm;

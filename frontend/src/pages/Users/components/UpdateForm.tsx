import React from 'react';
import {
  ModalForm,
  ProFormText,
  ProFormSelect,
} from '@ant-design/pro-components';
import { message } from 'antd';
import { updateUser } from '@/services/user';

export interface UpdateFormProps {
  open: boolean;
  onOpenChange: (visible: boolean) => void;
  onFinish: () => void;
  values: Partial<{
    id: number;
    username: string;
    email: string;
    nickname?: string;
    phone?: string;
    status: number;
  }>;
}

const UpdateForm: React.FC<UpdateFormProps> = ({
  open,
  onOpenChange,
  onFinish,
  values,
}) => {
  return (
    <ModalForm
      title="编辑用户"
      width="400px"
      open={open}
      onOpenChange={onOpenChange}
      initialValues={values}
      onFinish={async (value) => {
        try {
          if (!values.id) {
            message.error('用户ID不能为空');
            return false;
          }
          
          const response = await updateUser(values.id, value);
          if (response.code === 200) {
            message.success('更新成功');
            onFinish();
            return true;
          } else {
            message.error(response.message || '更新失败');
            return false;
          }
        } catch (error: any) {
          message.error(error.message || '更新失败，请重试');
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
            required: true,
            message: '邮箱为必填项',
          },
          {
            type: 'email',
            message: '请输入有效的邮箱地址',
          },
        ]}
        label="邮箱"
        name="email"
        placeholder="请输入邮箱"
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
        name="status"
        label="状态"
        valueEnum={{
          1: '正常',
          0: '禁用',
        }}
        placeholder="请选择状态"
        rules={[{ required: true, message: '请选择状态!' }]}
      />
    </ModalForm>
  );
};

export default UpdateForm;

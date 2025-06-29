import React from 'react';
import {
  ModalForm,
  ProFormText,
  ProFormSelect,
} from '@ant-design/pro-components';
import { message } from 'antd';
import { updateUser } from '@/services/user';
import { getRoleList } from '@/services/role';

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
    role?: string;
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
      width="600px"
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
            // 错误消息由全局错误处理器显示，这里不再重复显示
            return false;
          }
        } catch (error: any) {
          // 错误消息由全局错误处理器显示，这里不再重复显示
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
        placeholder="请选择角色"
        showSearch
        rules={[{ required: true, message: '请选择角色!' }]}
        request={async () => {
          try {
            const response = await getRoleList({ current: 1, pageSize: 100 });
            return response.data.map((role) => ({
              label: role.name,
              value: role.code,
            }));
          } catch (error) {
            console.error('获取角色列表失败:', error);
            return [
              { label: '管理员', value: 'admin' },
              { label: '管理者', value: 'manager' },
              { label: '普通用户', value: 'user' },
            ];
          }
        }}
      />
      <ProFormSelect
        name="status"
        label="状态"
        options={[
          { label: '正常', value: 1 },
          { label: '禁用', value: 0 },
        ]}
        placeholder="请选择状态"
        rules={[{ required: true, message: '请选择状态!' }]}
      />
    </ModalForm>
  );
};

export default UpdateForm;

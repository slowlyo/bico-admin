import React from 'react';
import {
  ModalForm,
  ProFormText,
} from '@ant-design/pro-components';
import { message } from 'antd';
import { changePassword } from '@/services/auth';

export interface ChangePasswordFormProps {
  open: boolean;
  onOpenChange: (visible: boolean) => void;
  onFinish: () => void;
}

const ChangePasswordForm: React.FC<ChangePasswordFormProps> = ({
  open,
  onOpenChange,
  onFinish,
}) => {
  return (
    <ModalForm
      title="修改密码"
      width="400px"
      open={open}
      onOpenChange={onOpenChange}
      onFinish={async (values) => {
        try {
          if (values.new_password !== values.confirm_password) {
            message.error('两次输入的密码不一致');
            return false;
          }

          const response = await changePassword({
            old_password: values.old_password,
            new_password: values.new_password,
            confirm_password: values.confirm_password,
          });

          if (response.code === 200) {
            message.success('密码修改成功');
            onFinish();
            return true;
          } else {
            message.error(response.message || '密码修改失败');
            return false;
          }
        } catch (error: any) {
          message.error(error.message || '密码修改失败，请重试');
          return false;
        }
      }}
    >
      <ProFormText.Password
        name="old_password"
        label="当前密码"
        placeholder="请输入当前密码"
        rules={[
          {
            required: true,
            message: '请输入当前密码',
          },
        ]}
      />
      <ProFormText.Password
        name="new_password"
        label="新密码"
        placeholder="请输入新密码"
        rules={[
          {
            required: true,
            message: '请输入新密码',
          },
          {
            min: 6,
            message: '密码至少6位',
          },
        ]}
      />
      <ProFormText.Password
        name="confirm_password"
        label="确认密码"
        placeholder="请再次输入新密码"
        rules={[
          {
            required: true,
            message: '请确认新密码',
          },
          ({ getFieldValue }) => ({
            validator(_, value) {
              if (!value || getFieldValue('new_password') === value) {
                return Promise.resolve();
              }
              return Promise.reject(new Error('两次输入的密码不一致'));
            },
          }),
        ]}
      />
    </ModalForm>
  );
};

export default ChangePasswordForm;

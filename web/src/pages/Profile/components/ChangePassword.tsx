import React, { useState } from 'react';
import { ProForm, ProFormText } from '@ant-design/pro-components';
import { Card, message } from 'antd';
import { changePassword } from '@/services/profile';

const ChangePassword: React.FC = () => {
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (values: any) => {
    setLoading(true);
    try {
      const response = await changePassword({
        old_password: values.oldPassword,
        new_password: values.newPassword,
      });

      if (response.code === 200) {
        message.success('密码修改成功');
        // 清空表单
        return true;
      } else {
        message.error(response.message || '密码修改失败');
        return false;
      }
    } catch (error) {
      console.error('修改密码失败:', error);
      message.error('密码修改失败，请重试');
      return false;
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card title="修改密码" style={{ maxWidth: 500 }}>
      <ProForm
        layout="vertical"
        onFinish={handleSubmit}
        submitter={{
          searchConfig: {
            submitText: '修改密码',
          },
          resetButtonProps: {
            style: { display: 'none' },
          },
          submitButtonProps: {
            loading,
          },
        }}
      >
        <ProFormText.Password
          name="oldPassword"
          label="当前密码"
          rules={[
            { required: true, message: '请输入当前密码' },
          ]}
          placeholder="请输入当前密码"
        />
        
        <ProFormText.Password
          name="newPassword"
          label="新密码"
          rules={[
            { required: true, message: '请输入新密码' },
            { min: 6, message: '密码长度至少6位' },
            { max: 50, message: '密码长度不能超过50位' },
          ]}
          placeholder="请输入新密码"
        />
        
        <ProFormText.Password
          name="confirmPassword"
          label="确认新密码"
          dependencies={['newPassword']}
          rules={[
            { required: true, message: '请确认新密码' },
            ({ getFieldValue }) => ({
              validator(_, value) {
                if (!value || getFieldValue('newPassword') === value) {
                  return Promise.resolve();
                }
                return Promise.reject(new Error('两次输入的密码不一致'));
              },
            }),
          ]}
          placeholder="请再次输入新密码"
        />
      </ProForm>
    </Card>
  );
};

export default ChangePassword;

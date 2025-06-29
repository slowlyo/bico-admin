import React, { useState } from 'react';
import {
  PageContainer,
  ProCard,
  ProForm,
  ProFormText,
} from '@ant-design/pro-components';
import { Button, message, Avatar } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { useModel } from '@umijs/max';
import { updateUserProfile, changePassword } from '@/services/auth';
import ChangePasswordForm from './components/ChangePasswordForm';

const Profile: React.FC = () => {
  const { initialState, setInitialState } = useModel('@@initialState');
  const [changePasswordVisible, setChangePasswordVisible] = useState(false);

  const handleUpdateProfile = async (values: any) => {
    try {
      const response = await updateUserProfile(values);
      if (response.code === 200) {
        message.success('更新成功');
        // 更新全局状态
        await setInitialState((s) => ({
          ...s,
          currentUser: response.data,
          name: response.data.nickname || response.data.username,
        }));
        return true;
      } else {
        message.error(response.message || '更新失败');
        return false;
      }
    } catch (error: any) {
      message.error(error.message || '更新失败，请重试');
      return false;
    }
  };



  const currentUser = initialState?.currentUser;

  return (
    <PageContainer>
      <ProCard split="vertical">
        <ProCard title="个人信息" colSpan="18">
          <ProForm
            layout="horizontal"
            labelCol={{ span: 4 }}
            wrapperCol={{ span: 16 }}
            initialValues={currentUser}
            onFinish={handleUpdateProfile}
            submitter={{
              searchConfig: {
                submitText: '更新信息',
              },
              resetButtonProps: {
                style: {
                  display: 'none',
                },
              },
            }}
          >
            <ProFormText
              name="username"
              label="用户名"
              disabled
              tooltip="用户名不可修改"
            />
            <ProFormText
              name="email"
              label="邮箱"
              rules={[
                {
                  type: 'email',
                  message: '请输入有效的邮箱地址',
                },
              ]}
            />
            <ProFormText
              name="nickname"
              label="昵称"
              placeholder="请输入昵称"
            />
            <ProFormText
              name="phone"
              label="手机号"
              placeholder="请输入手机号"
              rules={[
                {
                  pattern: /^1[3-9]\d{9}$/,
                  message: '请输入有效的手机号',
                },
              ]}
            />
          </ProForm>
        </ProCard>
        
        <ProCard title="操作" colSpan="6">
          <div style={{ textAlign: 'center' }}>
            <Avatar
              size={80}
              icon={<UserOutlined />}
              src={currentUser?.avatar}
              style={{ marginBottom: 16 }}
            />
            <div style={{ marginBottom: 16 }}>
              <strong>{currentUser?.nickname || currentUser?.username}</strong>
            </div>
            <div style={{ marginBottom: 24, color: '#666' }}>
              {currentUser?.email}
            </div>
            
            <Button
              type="primary"
              icon={<LockOutlined />}
              block
              onClick={() => setChangePasswordVisible(true)}
            >
              修改密码
            </Button>
          </div>
        </ProCard>
      </ProCard>

      <ChangePasswordForm
        open={changePasswordVisible}
        onOpenChange={setChangePasswordVisible}
        onFinish={() => {
          setChangePasswordVisible(false);
          message.success('密码修改成功');
        }}
      />
    </PageContainer>
  );
};

export default Profile;

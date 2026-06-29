import { CameraOutlined, LockOutlined, UserOutlined } from '@ant-design/icons';
import { ProForm, ProFormText } from '@ant-design/pro-components';
import { useModel } from '@umijs/max';
import { Card, Upload, Avatar, App, Row, Col, Divider } from 'antd';
import { createStyles } from 'antd-style';
import React, { useState } from 'react';
import { updateProfile, changePassword, uploadAvatar } from '@/services/auth/profile';
import type { UploadFile, UploadProps } from 'antd';

const useStyles = createStyles(({ token }) => ({
  container: {
  },
  avatarSection: {
    textAlign: 'center',
    margin: '12px 0',
  },
  avatarWrapper: {
    position: 'relative',
    display: 'inline-block',
    cursor: 'pointer',
    '&:hover .avatar-mask': {
      opacity: 1,
    },
  },
  avatarMask: {
    position: 'absolute',
    top: 0,
    left: 0,
    width: '120px',
    height: '120px',
    borderRadius: '50%',
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    opacity: 0,
    transition: 'opacity 0.3s',
    color: '#fff',
    fontSize: '14px',
  },
  uploadIcon: {
    fontSize: '24px',
    marginBottom: '8px',
  },
}));

const Profile: React.FC = () => {
  const { styles } = useStyles();
  const { initialState, setInitialState } = useModel('@@initialState');
  const { message } = App.useApp();
  const [profileForm] = ProForm.useForm();
  const [passwordForm] = ProForm.useForm();
  const [uploading, setUploading] = useState(false);

  const currentUser = initialState?.currentUser;

  const handleAvatarUpload: UploadProps['customRequest'] = async (options) => {
    const { file, onSuccess, onError } = options;
    
    setUploading(true);
    try {
      const response = await uploadAvatar(file as File);
      
      if (response.code === 0 && response.data?.url) {
        const avatarUrl = response.data.url;
        
        // 更新用户信息
        const updateResponse = await updateProfile({ avatar: avatarUrl });
        
        if (updateResponse.code === 0 && updateResponse.data) {
          // 更新全局状态
          setInitialState((s) => ({
            ...s,
            currentUser: updateResponse.data,
          }));
          
          // 更新本地存储
          localStorage.setItem('currentUser', JSON.stringify(updateResponse.data));
          
          message.success('头像上传成功');
          onSuccess?.(response.data);
        }
      }
    } catch (error: any) {
      message.error(error.message || '头像上传失败');
      onError?.(error);
    } finally {
      setUploading(false);
    }
  };

  const handleProfileUpdate = async (values: API.UpdateProfileParams) => {
    try {
      const response = await updateProfile(values);
      
      if (response.code === 0 && response.data) {
        // 更新全局状态
        setInitialState((s) => ({
          ...s,
          currentUser: response.data,
        }));
        
        // 更新本地存储
        localStorage.setItem('currentUser', JSON.stringify(response.data));
        
        message.success('个人信息更新成功');
      }
    } catch (error: any) {
      message.error(error.message || '更新失败');
    }
  };

  const handlePasswordChange = async (values: API.ChangePasswordParams & { confirmPassword?: string }) => {
    if (values.newPassword !== values.confirmPassword) {
      message.error('两次输入的新密码不一致');
      return;
    }

    try {
      const response = await changePassword({
        oldPassword: values.oldPassword,
        newPassword: values.newPassword,
      });
      
      if (response.code === 0) {
        message.success('密码修改成功，请重新登录');
        passwordForm.resetFields();
        
        // 清除登录信息，跳转到登录页
        setTimeout(() => {
          localStorage.removeItem('token');
          localStorage.removeItem('currentUser');
          window.location.href = '/auth/login';
        }, 1500);
      }
    } catch (error: any) {
      message.error(error.message || '密码修改失败');
    }
  };

  return (
    <div className={styles.container}>
      <Row gutter={24}>
        <Col xs={24} lg={8}>
          <Card title="头像设置">
            <div className={styles.avatarSection}>
              <Upload
                accept="image/*"
                showUploadList={false}
                customRequest={handleAvatarUpload}
                disabled={uploading}
              >
                <div className={styles.avatarWrapper}>
                  <Avatar 
                    size={120} 
                    src={currentUser?.avatar}
                    icon={<UserOutlined />}
                  />
                  <div className={`${styles.avatarMask} avatar-mask`}>
                    <CameraOutlined className={styles.uploadIcon} />
                    <span>{uploading ? '上传中...' : '点击上传'}</span>
                  </div>
                </div>
              </Upload>
            </div>
          </Card>
        </Col>

        <Col xs={24} lg={16}>
          <Card title="基本信息" style={{ marginBottom: '24px' }}>
            <ProForm
              form={profileForm}
              layout="horizontal"
              labelCol={{ span: 4 }}
              wrapperCol={{ span: 20 }}
              initialValues={{
                name: currentUser?.name,
                username: currentUser?.username,
              }}
              onFinish={handleProfileUpdate}
              submitter={{
                searchConfig: {
                  submitText: '保存',
                },
                resetButtonProps: {
                  style: { display: 'none' },
                },
              }}
            >
              <ProFormText
                name="username"
                label="用户名"
                disabled
                fieldProps={{
                  prefix: <UserOutlined />,
                }}
              />
              <ProFormText
                name="name"
                label="姓名"
                rules={[
                  {
                    required: true,
                    message: '请输入姓名',
                  },
                ]}
                fieldProps={{
                  prefix: <UserOutlined />,
                }}
              />
            </ProForm>
          </Card>

          <Card title="修改密码">
            <ProForm
              form={passwordForm}
              layout="horizontal"
              labelCol={{ span: 4 }}
              wrapperCol={{ span: 20 }}
              onFinish={handlePasswordChange}
              submitter={{
                searchConfig: {
                  submitText: '修改密码',
                },
                resetButtonProps: {
                  style: { display: 'none' },
                },
              }}
            >
              <ProFormText.Password
                name="oldPassword"
                label="原密码"
                rules={[
                  {
                    required: true,
                    message: '请输入原密码',
                  },
                ]}
                fieldProps={{
                  prefix: <LockOutlined />,
                }}
              />
              <ProFormText.Password
                name="newPassword"
                label="新密码"
                rules={[
                  {
                    required: true,
                    message: '请输入新密码',
                  },
                  {
                    min: 6,
                    message: '密码长度至少6位',
                  },
                ]}
                fieldProps={{
                  prefix: <LockOutlined />,
                }}
              />
              <ProFormText.Password
                name="confirmPassword"
                label="确认密码"
                rules={[
                  {
                    required: true,
                    message: '请再次输入新密码',
                  },
                ]}
                fieldProps={{
                  prefix: <LockOutlined />,
                }}
              />
            </ProForm>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Profile;

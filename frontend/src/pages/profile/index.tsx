import React, { useState } from 'react';
import {
  Card,
  Form,
  Input,
  Button,
  message,
  Tabs,
  Avatar,
  Upload,
  Row,
  Col
} from 'antd';
import {
  UserOutlined,
  MailOutlined,
  PhoneOutlined,
  LockOutlined,
  UploadOutlined
} from '@ant-design/icons';
import { useGetIdentity } from '@refinedev/core';
import { authAPI } from '../../utils/request';
import './profile.css';

const { TabPane } = Tabs;

interface UserProfile {
  id: number;
  username: string;
  email: string;
  nickname: string;
  phone: string;
  avatar: string;
}

interface ChangePasswordForm {
  old_password: string;
  new_password: string;
  confirm_password: string;
}

export const Profile: React.FC = () => {
  const [profileForm] = Form.useForm();
  const [passwordForm] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [passwordLoading, setPasswordLoading] = useState(false);
  
  const { data: identity, refetch: refetchIdentity } = useGetIdentity<UserProfile>();

  // 更新用户资料
  const handleUpdateProfile = async (values: Partial<UserProfile>) => {
    setLoading(true);
    try {
      const response = await authAPI.updateProfile(values);
      if (response.code === 200) {
        message.success('个人资料更新成功');
        // 更新本地存储的用户信息
        refetchIdentity();
      } else {
        message.error(response.message || '更新失败');
      }
    } catch (error: any) {
      message.error(error.message || '更新失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  // 修改密码
  const handleChangePassword = async (values: ChangePasswordForm) => {
    if (values.new_password !== values.confirm_password) {
      message.error('新密码和确认密码不一致');
      return;
    }

    setPasswordLoading(true);
    try {
      const response = await authAPI.changePassword({
        old_password: values.old_password,
        new_password: values.new_password,
      });
      
      if (response.code === 200) {
        message.success('密码修改成功');
        passwordForm.resetFields();
      } else {
        message.error(response.message || '密码修改失败');
      }
    } catch (error: any) {
      message.error(error.message || '密码修改失败，请稍后重试');
    } finally {
      setPasswordLoading(false);
    }
  };

  // 头像上传处理
  const handleAvatarUpload = (info: any) => {
    if (info.file.status === 'done') {
      message.success('头像上传成功');
      // 这里可以处理头像上传成功后的逻辑
    } else if (info.file.status === 'error') {
      message.error('头像上传失败');
    }
  };

  return (
    <div className="profile-container">
      <Row gutter={24} justify="center">
        <Col xs={24} sm={20} md={16} lg={14} xl={12}>
          <Card className="profile-card">
            <Tabs defaultActiveKey="profile" size="large" className="profile-tabs">
              <TabPane tab="基本信息" key="profile">
                <Row gutter={[24, 24]}>
                  <Col xs={24} sm={8} md={6}>
                    <div className="profile-avatar-section">
                      <Avatar
                        size={100}
                        src={identity?.avatar}
                        icon={<UserOutlined />}
                        className="profile-avatar"
                      />
                      <div>
                        <Upload
                          name="avatar"
                          showUploadList={false}
                          action="/api/upload/avatar"
                          onChange={handleAvatarUpload}
                        >
                          <Button icon={<UploadOutlined />} className="profile-upload-btn">
                            更换头像
                          </Button>
                        </Upload>
                      </div>
                    </div>
                  </Col>

              <Col xs={24} sm={16} md={18}>
                <Form
                  form={profileForm}
                  layout="vertical"
                  initialValues={identity}
                  onFinish={handleUpdateProfile}
                  className="profile-form"
                >
                  <Row gutter={16}>
                    <Col xs={24} sm={12}>
                      <Form.Item
                        name="username"
                        label="用户名"
                        rules={[
                          { required: true, message: '请输入用户名' },
                          { min: 3, message: '用户名至少3个字符' },
                          { max: 50, message: '用户名最多50个字符' }
                        ]}
                      >
                        <Input
                          prefix={<UserOutlined />}
                          placeholder="请输入用户名"
                        />
                      </Form.Item>
                    </Col>

                    <Col xs={24} sm={12}>
                      <Form.Item
                        name="nickname"
                        label="昵称"
                        rules={[
                          { max: 50, message: '昵称最多50个字符' }
                        ]}
                      >
                        <Input
                          prefix={<UserOutlined />}
                          placeholder="请输入昵称"
                        />
                      </Form.Item>
                    </Col>
                  </Row>

                  <Row gutter={16}>
                    <Col xs={24} sm={12}>
                      <Form.Item
                        name="email"
                        label="邮箱"
                        rules={[
                          { required: true, message: '请输入邮箱' },
                          { type: 'email', message: '请输入有效的邮箱地址' }
                        ]}
                      >
                        <Input
                          prefix={<MailOutlined />}
                          placeholder="请输入邮箱"
                        />
                      </Form.Item>
                    </Col>

                    <Col xs={24} sm={12}>
                      <Form.Item
                        name="phone"
                        label="手机号"
                        rules={[
                          { max: 20, message: '手机号最多20个字符' }
                        ]}
                      >
                        <Input
                          prefix={<PhoneOutlined />}
                          placeholder="请输入手机号"
                        />
                      </Form.Item>
                    </Col>
                  </Row>

                  <Form.Item>
                    <Button
                      type="primary"
                      htmlType="submit"
                      loading={loading}
                    >
                      保存修改
                    </Button>
                  </Form.Item>
                </Form>
              </Col>
            </Row>
          </TabPane>
          
          <TabPane tab="修改密码" key="password">
            <Row justify="center">
              <Col xs={24} sm={18} md={14} lg={12}>
                <div className="password-form-container">
                  <Form
                    form={passwordForm}
                    layout="vertical"
                    onFinish={handleChangePassword}
                    className="profile-form"
                  >
                  <Form.Item
                    name="old_password"
                    label="当前密码"
                    rules={[
                      { required: true, message: '请输入当前密码' }
                    ]}
                  >
                    <Input.Password
                      prefix={<LockOutlined />}
                      placeholder="请输入当前密码"
                    />
                  </Form.Item>

                  <Form.Item
                    name="new_password"
                    label="新密码"
                    rules={[
                      { required: true, message: '请输入新密码' },
                      { min: 6, message: '密码至少6个字符' }
                    ]}
                  >
                    <Input.Password
                      prefix={<LockOutlined />}
                      placeholder="请输入新密码"
                    />
                  </Form.Item>

                  <Form.Item
                    name="confirm_password"
                    label="确认新密码"
                    dependencies={['new_password']}
                    rules={[
                      { required: true, message: '请确认新密码' },
                      ({ getFieldValue }) => ({
                        validator(_, value) {
                          if (!value || getFieldValue('new_password') === value) {
                            return Promise.resolve();
                          }
                          return Promise.reject(new Error('两次输入的密码不一致'));
                        },
                      }),
                    ]}
                  >
                    <Input.Password
                      prefix={<LockOutlined />}
                      placeholder="请再次输入新密码"
                    />
                  </Form.Item>

                  <Form.Item>
                    <Button
                      type="primary"
                      htmlType="submit"
                      loading={passwordLoading}
                      style={{ width: '100%' }}
                    >
                      修改密码
                    </Button>
                  </Form.Item>
                  </Form>
                </div>
              </Col>
            </Row>
          </TabPane>
        </Tabs>
      </Card>
        </Col>
      </Row>
    </div>
  );
};

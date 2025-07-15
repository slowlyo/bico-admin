import React, { useState, useEffect } from 'react';
import {
  ProForm,
  ProFormText,
  ProFormUploadButton
} from '@ant-design/pro-components';
import { Avatar, Row, Col, message, Form } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import { updateProfile } from '@/services/profile';
import { uploadAvatar } from '@/services/upload';
import type { UserInfo } from '@/services/auth';

interface BasicInfoProps {
  userInfo?: UserInfo;
  onUpdate: () => void;
}

const BasicInfo: React.FC<BasicInfoProps> = ({ userInfo, onUpdate }) => {
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm();

  // 当userInfo变化时，更新表单值
  useEffect(() => {
    if (userInfo) {
      form.setFieldsValue({
        username: userInfo.username,
        name: userInfo.name,
        email: userInfo.email,
        phone: userInfo.phone,
        avatar: userInfo.avatar ? [{
          uid: '-1',
          name: 'avatar',
          status: 'done',
          url: userInfo.avatar,
        }] : [],
      });
    }
  }, [userInfo, form]);

  const handleSubmit = async (values: any) => {
    setLoading(true);
    try {
      // 处理头像上传
      let avatarUrl = userInfo?.avatar;
      if (values.avatar && Array.isArray(values.avatar) && values.avatar.length > 0) {
        const avatarFile = values.avatar[0];
        if (avatarFile.originFileObj) {
          // 上传新头像
          const uploadResponse = await uploadAvatar(avatarFile.originFileObj);
          if (uploadResponse.code === 200) {
            avatarUrl = uploadResponse.data.files[0].file_path;
          } else {
            message.error('头像上传失败');
            return;
          }
        } else if (avatarFile.url) {
          // 使用现有头像URL
          avatarUrl = avatarFile.url;
        }
      }

      // 更新个人信息
      const updateData = {
        name: values.name,
        email: values.email,
        phone: values.phone,
        avatar: avatarUrl,
      };

      const response = await updateProfile(updateData);
      if (response.code === 200) {
        // 只调用onUpdate，不在这里显示成功消息
        // 成功消息由父组件统一处理
        onUpdate();
      } else {
        message.error(response.message || '更新失败');
      }
    } catch (error) {
      console.error('更新个人信息失败:', error);
      message.error('更新失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Row gutter={24}>
      <Col span={8}>
        <div style={{ textAlign: 'center', padding: '24px 0' }}>
          <Avatar
            size={120}
            src={userInfo?.avatar}
            icon={<UserOutlined />}
            style={{ marginBottom: 16 }}
          />
          <div style={{ fontSize: 16, fontWeight: 500, marginBottom: 8 }}>
            {userInfo?.nickname || userInfo?.name}
          </div>
          <div style={{ color: '#666', fontSize: 14 }}>
            {userInfo?.username}
          </div>
        </div>
      </Col>
      <Col span={16}>
        <ProForm
          layout="vertical"
          form={form}
          onFinish={handleSubmit}
          submitter={{
            searchConfig: {
              submitText: '保存修改',
            },
            resetButtonProps: {
              style: { display: 'none' },
            },
            submitButtonProps: {
              loading,
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
            name="name"
            label="姓名"
            rules={[
              { required: true, message: '请输入姓名' },
              { max: 100, message: '姓名长度不能超过100个字符' },
            ]}
            placeholder="请输入姓名"
          />
          
          <ProFormText
            name="email"
            label="邮箱"
            rules={[
              { type: 'email', message: '请输入正确的邮箱格式' },
            ]}
            placeholder="请输入邮箱"
          />
          
          <ProFormText
            name="phone"
            label="手机号"
            rules={[
              { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号格式' },
            ]}
            placeholder="请输入手机号"
          />
          
          <ProFormUploadButton
            name="avatar"
            label="头像"
            max={1}
            fieldProps={{
              name: 'files',
              listType: 'picture-card',
              showUploadList: {
                showPreviewIcon: true,
                showRemoveIcon: true,
              },
              accept: 'image/*',
              beforeUpload: (file) => {
                const isImage = file.type.startsWith('image/');
                if (!isImage) {
                  message.error('只能上传图片文件！');
                  return false;
                }
                const isLt2M = file.size / 1024 / 1024 < 2;
                if (!isLt2M) {
                  message.error('图片大小不能超过2MB！');
                  return false;
                }
                return true;
              },
            }}
            extra="支持 jpg、png 格式，文件大小不超过 2MB"
          />
        </ProForm>
      </Col>
    </Row>
  );
};

export default BasicInfo;

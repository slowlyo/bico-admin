import { LockOutlined, UserOutlined } from '@ant-design/icons';
import { ModalForm, ProForm, ProFormText } from '@ant-design/pro-components';
import { history, useModel } from '@umijs/max';
import { App } from 'antd';
import React, { useCallback } from 'react';
import AvatarUpload from '@/components/AvatarUpload';
import { changePassword, updateProfile, uploadAvatar } from '@/services/auth/profile';
import type { ChangePasswordParams, UpdateProfileParams } from '@/services/auth/types';
import { MIN_PASSWORD_LENGTH } from '@/utils/validation';

interface AccountModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

interface ChangePasswordFormValues extends ChangePasswordParams {
  confirmPassword: string;
}

// getErrorMessage 从未知异常中提取可展示的信息。
const getErrorMessage = (error: unknown, fallback: string) => {
  if (error instanceof Error && error.message) {
    // 标准异常优先展示接口返回的信息，便于用户定位失败原因。
    return error.message;
  }
  return fallback;
};

// ProfileModal 处理头像与姓名更新，并同步当前登录用户状态。
export const ProfileModal: React.FC<AccountModalProps> = ({ open, onOpenChange }) => {
  const { initialState, setInitialState } = useModel('@@initialState');
  const { message } = App.useApp();
  const currentUser = initialState?.currentUser;

  // syncCurrentUser 统一同步内存状态与本地缓存。
  const syncCurrentUser = useCallback((user: API.CurrentUser) => {
    setInitialState((state) => ({
      ...state,
      currentUser: user,
    }));
    localStorage.setItem('currentUser', JSON.stringify(user));
  }, [setInitialState]);

  // handleAvatarUpload 上传头像后立即更新当前用户资料。
  const handleAvatarUpload = useCallback(async (file: File) => {
    const response = await uploadAvatar(file);
    if (response.code !== 0 || !response.data?.url) {
      // 上传接口未返回有效地址时停止更新，避免写入空头像。
      throw new Error(response.msg || '头像上传失败');
    }

    const updateResponse = await updateProfile({ avatar: response.data.url });
    if (updateResponse.code !== 0 || !updateResponse.data) {
      // 资料更新失败时保留当前头像状态。
      throw new Error(updateResponse.msg || '头像更新失败');
    }

    syncCurrentUser(updateResponse.data);
    return updateResponse.data.avatar || response.data.url;
  }, [syncCurrentUser]);

  // handleProfileUpdate 保存姓名并在成功后关闭弹窗。
  const handleProfileUpdate = useCallback(async (values: UpdateProfileParams) => {
    try {
      const response = await updateProfile(values);
      if (response.code !== 0 || !response.data) {
        // 接口未返回用户资料时保持弹窗，允许用户继续修改。
        throw new Error(response.msg || '个人信息更新失败');
      }

      syncCurrentUser(response.data);
      message.success('个人信息更新成功');
      return true;
    } catch (error) {
      // 保存失败时保留表单内容，避免用户重复输入。
      message.error(getErrorMessage(error, '个人信息更新失败'));
      return false;
    }
  }, [message, syncCurrentUser]);

  return (
    <ModalForm<UpdateProfileParams>
      title="个人信息"
      width={520}
      open={open}
      onOpenChange={onOpenChange}
      initialValues={{ username: currentUser?.username, name: currentUser?.name }}
      preserve={false}
      onFinish={handleProfileUpdate}
      modalProps={{ destroyOnHidden: true }}
      submitter={{
        searchConfig: { submitText: '保存' },
        resetButtonProps: { style: { display: 'none' } },
      }}
    >
      <div style={{ marginBottom: 24, textAlign: 'center' }}>
        <AvatarUpload
          value={currentUser?.avatar}
          size={104}
          onUpload={handleAvatarUpload}
        />
      </div>
      <ProFormText
        name="username"
        label="用户名"
        rules={[
          { required: true, message: '请输入用户名' },
          { whitespace: true, message: '用户名不能全为空格' },
          { max: 64, message: '用户名不能超过64位' },
        ]}
        fieldProps={{ prefix: <UserOutlined /> }}
      />
      <ProFormText
        name="name"
        label="姓名"
        rules={[{ required: true, message: '请输入姓名' }]}
        fieldProps={{ prefix: <UserOutlined /> }}
      />
    </ModalForm>
  );
};

// PasswordModal 校验并修改当前用户密码。
export const PasswordModal: React.FC<AccountModalProps> = ({ open, onOpenChange }) => {
  const { message } = App.useApp();
  const [form] = ProForm.useForm<ChangePasswordFormValues>();

  // handlePasswordChange 修改成功后清除登录信息并跳转登录页。
  const handlePasswordChange = useCallback(async (values: ChangePasswordFormValues) => {
    try {
      const response = await changePassword({
        oldPassword: values.oldPassword,
        newPassword: values.newPassword,
      });
      if (response.code !== 0) {
        // 业务失败时保留弹窗与表单内容，方便用户修正。
        throw new Error(response.msg || '密码修改失败');
      }

      message.success('密码修改成功，请重新登录');
      onOpenChange(false);
      form.resetFields();
      // 短暂保留成功反馈后清理凭证，避免旧会话继续使用。
      window.setTimeout(() => {
        localStorage.removeItem('token');
        localStorage.removeItem('currentUser');
        history.replace('/auth/login');
      }, 1500);
      return true;
    } catch (error) {
      // 修改失败时展示原因，并保持当前输入。
      message.error(getErrorMessage(error, '密码修改失败'));
      return false;
    }
  }, [form, message, onOpenChange]);

  return (
    <ModalForm<ChangePasswordFormValues>
      form={form}
      title="修改密码"
      width={480}
      open={open}
      onOpenChange={onOpenChange}
      preserve={false}
      onFinish={handlePasswordChange}
      modalProps={{ destroyOnHidden: true }}
      submitter={{
        searchConfig: { submitText: '修改密码' },
        resetButtonProps: { style: { display: 'none' } },
      }}
    >
      <ProFormText.Password
        name="oldPassword"
        label="原密码"
        rules={[{ required: true, message: '请输入原密码' }]}
        fieldProps={{ prefix: <LockOutlined />, autoComplete: 'current-password' }}
      />
      <ProFormText.Password
        name="newPassword"
        label="新密码"
        rules={[
          { required: true, message: '请输入新密码' },
          { min: MIN_PASSWORD_LENGTH, message: `密码长度至少${MIN_PASSWORD_LENGTH}位` },
        ]}
        fieldProps={{ prefix: <LockOutlined />, autoComplete: 'new-password' }}
      />
      <ProFormText.Password
        name="confirmPassword"
        label="确认密码"
        dependencies={['newPassword']}
        rules={[
          { required: true, message: '请再次输入新密码' },
          ({ getFieldValue }) => ({
            // validator 保证确认密码与新密码始终一致。
            validator(_, value) {
              if (!value || getFieldValue('newPassword') === value) {
                // 空值交由必填规则处理，一致时允许提交。
                return Promise.resolve();
              }
              return Promise.reject(new Error('两次输入的新密码不一致'));
            },
          }),
        ]}
        fieldProps={{ prefix: <LockOutlined />, autoComplete: 'new-password' }}
      />
    </ModalForm>
  );
};

import { login } from '@/services/auth';
import { LockOutlined, UserOutlined } from '@ant-design/icons';
import { LoginForm, ProFormText } from '@ant-design/pro-components';
import { useModel, history } from '@umijs/max';
import { Alert, message } from 'antd';
import React, { useState } from 'react';
import { appConfig } from '@/config/app';

const LoginMessage: React.FC<{
  content: string;
}> = ({ content }) => {
  return (
    <Alert
      style={{
        marginBottom: 24,
      }}
      message={content}
      type="error"
      showIcon
    />
  );
};

const Login: React.FC = () => {
  const [userLoginState, setUserLoginState] = useState<{
    status?: 'ok' | 'error';
    message?: string;
  }>({});
  const { initialState, setInitialState } = useModel('@@initialState');

  const handleSubmit = async (values: any) => {
    try {
      // 登录
      const response = await login({
        username: values.username,
        password: values.password,
        captcha: '1234', // 临时验证码
      });

      if (response.code === 200) {
        const { token, user_info, permissions, menus } = response.data;

        // 存储token
        localStorage.setItem('token', token);
        localStorage.setItem('userInfo', JSON.stringify(user_info));
        localStorage.setItem('permissions', JSON.stringify(permissions));
        localStorage.setItem('menus', JSON.stringify(menus));

        message.success('登录成功！');

        // 更新全局状态
        await setInitialState((s) => ({
          ...s,
          currentUser: user_info,
          permissions,
          menus,
          // 添加umi layout需要的字段
          name: user_info.nickname,
          avatar: user_info.avatar,
        }));

        // 跳转到首页 - 使用history.replace避免页面刷新
        const urlParams = new URLSearchParams(window.location.search);
        const redirectUrl = urlParams.get('redirect') || '/';
        history.replace(redirectUrl);
        return;
      }

      // 登录失败
      const errorMsg = response.message || '登录失败，请重试！';
      setUserLoginState({
        status: 'error',
        message: errorMsg,
      });
      message.error(errorMsg);
    } catch (error: any) {
      console.error('登录错误:', error);

      // 处理HTTP错误响应
      let errorMessage = '登录失败，请重试！';
      if (error?.response?.data?.message) {
        errorMessage = error.response.data.message;
      } else if (error?.message) {
        errorMessage = error.message;
      }

      message.error(errorMessage);
      setUserLoginState({
        status: 'error',
        message: errorMessage,
      });
    }
  };

  const { status, message: errorMessage } = userLoginState;

  return (
    <div style={{
      minHeight: '100vh',
      background: 'linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      paddingTop: '64px',
      position: 'relative',
      overflow: 'hidden'
    }}>
      {/* 背景装饰元素 */}
      <div style={{
        position: 'absolute',
        top: '-50%',
        left: '-50%',
        width: '200%',
        height: '200%',
        background: `
          radial-gradient(circle at 25% 25%, rgba(255, 255, 255, 0.4) 0%, transparent 50%),
          radial-gradient(circle at 75% 75%, rgba(0, 0, 0, 0.05) 0%, transparent 50%)
        `,
        zIndex: 0
      }} />

      {/* 浮动圆点装饰 */}
      <div style={{
        position: 'absolute',
        top: '20%',
        right: '20%',
        width: '80px',
        height: '80px',
        background: 'rgba(255, 255, 255, 0.3)',
        borderRadius: '50%',
        filter: 'blur(1px)'
      }} />

      <div style={{
        position: 'absolute',
        bottom: '25%',
        left: '10%',
        width: '120px',
        height: '120px',
        background: 'rgba(0, 0, 0, 0.03)',
        borderRadius: '50%',
        filter: 'blur(2px)'
      }} />

      <div style={{
        position: 'absolute',
        top: '60%',
        right: '5%',
        width: '40px',
        height: '40px',
        background: 'rgba(255, 255, 255, 0.2)',
        borderRadius: '50%'
      }} />

      <div style={{
        width: '100%',
        maxWidth: '448px',
        marginTop: '-112px',
        position: 'relative',
        zIndex: 1
      }}>
        <LoginForm
          logo={<img src={appConfig.logo} alt={appConfig.title} style={{ height: '32px', width: 'auto', objectFit: 'contain' }} />}
          title={appConfig.title}
          subTitle
          onFinish={async (values) => {
            await handleSubmit(values);
          }}
        >
          {status === 'error' && (
            <LoginMessage content={errorMessage || "账户或密码错误"} />
          )}

          <ProFormText
            name="username"
            fieldProps={{
              size: 'large',
              prefix: <UserOutlined />,
            }}
            placeholder="用户名（3-50位）"
            rules={[
              {
                required: true,
                message: '用户名是必填项！',
              },
              {
                min: 3,
                max: 50,
                message: '用户名长度为3-50位！',
              },
            ]}
          />
          <ProFormText.Password
            name="password"
            fieldProps={{
              size: 'large',
              prefix: <LockOutlined />,
            }}
            placeholder="密码（最少6位）"
            rules={[
              {
                required: true,
                message: '密码是必填项！',
              },
              {
                min: 6,
                message: '密码长度至少6位！',
              },
            ]}
          />
        </LoginForm>
      </div>
    </div>
  );
};

export default Login;

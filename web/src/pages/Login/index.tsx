import { login } from '@/services/auth';
import { LockOutlined, UserOutlined } from '@ant-design/icons';
import { LoginForm, ProFormText } from '@ant-design/pro-components';
import { useModel, history } from '@umijs/max';
import { Alert, message } from 'antd';
import React, { useState } from 'react';

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
    <div className="min-h-screen bg-gray-100 flex items-center justify-center pt-16">
      <div className="w-full max-w-md -mt-28">
        <LoginForm
          logo={<div className="text-2xl font-bold text-blue-600">Bico</div>}
          title="Bico Admin"
          subTitle="管理后台"
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
              prefix: <UserOutlined className="text-blue-500" />,
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
              prefix: <LockOutlined className="text-blue-500" />,
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

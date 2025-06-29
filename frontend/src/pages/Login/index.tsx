import React, { useState } from 'react';
import { LockOutlined, UserOutlined } from '@ant-design/icons';
import { LoginForm, ProFormText, ProFormCheckbox } from '@ant-design/pro-components';
import { Alert, message } from 'antd';
import { history, useModel } from '@umijs/max';
import { flushSync } from 'react-dom';
import { login } from '@/services/auth';
import { createStyles } from 'antd-style';

const useStyles = createStyles(() => {
  return {
    container: {
      display: 'flex',
      flexDirection: 'column',
      height: '100vh',
      overflow: 'auto',
      backgroundImage: "url('https://mdn.alipayobjects.com/yuyan_qk0oxh/afts/img/V-_oS6r-i7wAAAAAAAAAAAAAFl94AQBr')",
      backgroundSize: '100% 100%',
    },
    content: {
      flex: 1,
      padding: '32px 0',
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
    },
  };
});

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
    type?: string;
  }>({});
  const { initialState, setInitialState } = useModel('@@initialState');
  const { styles } = useStyles();

  const fetchUserInfo = async () => {
    // 重新调用 getInitialState 来获取完整的用户信息和权限
    const newInitialState = await initialState?.fetchUserInfo?.();

    if (newInitialState) {
      // 手动获取权限
      try {
        const { getUserPermissions } = await import('@/services/auth');
        const permissionsInfo = await getUserPermissions();
        const userPermissions = permissionsInfo.code === 200 ? permissionsInfo.data : [];

        flushSync(() => {
          setInitialState((s) => ({
            ...s,
            currentUser: newInitialState,
            userPermissions: userPermissions,
          }));
        });
      } catch (error) {
        console.error('登录后权限获取失败:', error);
        flushSync(() => {
          setInitialState((s) => ({
            ...s,
            currentUser: newInitialState,
            userPermissions: [],
          }));
        });
      }
    }
  };

  const handleSubmit = async (values: API.LoginParams) => {
    try {
      // 登录
      const response = await login({
        username: values.username!,
        password: values.password!,
      }, {
        skipErrorHandler: true, // 跳过全局错误处理
      });

      if (response.code === 200) {
        const defaultLoginSuccessMessage = '登录成功！';
        message.success(defaultLoginSuccessMessage);

        // 保存token
        localStorage.setItem('token', response.data.token);

        await fetchUserInfo();
        const urlParams = new URL(window.location.href).searchParams;
        history.push(urlParams.get('redirect') || '/');
        return;
      }

      // 如果失败去设置用户错误信息
      setUserLoginState({
        status: 'error',
        type: 'account',
      });
    } catch (error: any) {
      const defaultLoginFailureMessage = '登录失败，请重试！';
      console.log(error);

      // 根据错误类型显示不同的错误信息
      let errorMessage = defaultLoginFailureMessage;
      if (error.response?.status === 401) {
        errorMessage = '用户名或密码错误';
      } else if (error.response?.data?.message) {
        errorMessage = error.response.data.message;
      } else if (error.message) {
        errorMessage = error.message;
      }

      message.error(errorMessage);
      setUserLoginState({
        status: 'error',
        type: 'account',
      });
    }
  };

  const { status } = userLoginState;

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        <LoginForm
          logo={<img alt="logo" src="/logo.svg" />}
          title="Bico Admin"
          subTitle="现代化管理后台系统"
          initialValues={{
            autoLogin: true,
          }}
          onFinish={async (values) => {
            await handleSubmit(values as API.LoginParams);
          }}
        >
          {status === 'error' && (
            <LoginMessage content="账户或密码错误" />
          )}

          <ProFormText
            name="username"
            fieldProps={{
              size: 'large',
              prefix: <UserOutlined />,
            }}
            placeholder="用户名"
            rules={[
              {
                required: true,
                message: '用户名是必填项！',
              },
            ]}
          />
          <ProFormText.Password
            name="password"
            fieldProps={{
              size: 'large',
              prefix: <LockOutlined />,
            }}
            placeholder="密码"
            rules={[
              {
                required: true,
                message: '密码是必填项！',
              },
            ]}
          />

          <div style={{ marginBottom: 24 }}>
            <ProFormCheckbox noStyle name="autoLogin">
              自动登录
            </ProFormCheckbox>
          </div>
        </LoginForm>
      </div>
    </div>
  );
};

export default Login;

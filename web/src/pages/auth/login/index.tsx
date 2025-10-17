import { LockOutlined, UserOutlined } from "@ant-design/icons";
import {
    LoginForm,
    ProFormCheckbox,
    ProFormText,
    ProForm,
} from "@ant-design/pro-components";
import {
    FormattedMessage,
    Helmet,
    SelectLang,
    useIntl,
    useModel,
} from "@umijs/max";
import { Alert, App } from "antd";
import { createStyles } from "antd-style";
import React, { useState } from "react";
import { flushSync } from "react-dom";
import { Footer } from "@/components";
import { login } from '@/services/admin';
import { saveCredentials, getCredentials, clearCredentials } from '@/utils/crypto';

const useStyles = createStyles(({ token }) => {
    return {
        lang: {
            width: 42,
            height: 42,
            lineHeight: "42px",
            position: "fixed",
            right: 16,
            borderRadius: token.borderRadius,
            ":hover": {
                backgroundColor: token.colorBgTextHover,
            },
        },
        container: {
            display: "flex",
            flexDirection: "column",
            height: "100vh",
            overflow: "auto",
            backgroundImage:
                "url('https://mdn.alipayobjects.com/yuyan_qk0oxh/afts/img/V-_oS6r-i7wAAAAAAAAAAAAAFl94AQBr')",
            backgroundSize: "100% 100%",
        },
    };
});

const Lang = () => {
    const { styles } = useStyles();

    return (
        <div className={styles.lang} data-lang>
            {SelectLang && <SelectLang />}
        </div>
    );
};

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
    const [userLoginState, setUserLoginState] = useState<{status?: string; message?: string}>({});
    const [form] = ProForm.useForm();
    const { initialState, setInitialState } = useModel("@@initialState");
    const { styles } = useStyles();
    const { message } = App.useApp();
    const intl = useIntl();
    
    const appName = initialState?.appConfig?.name || 'Bico Admin';
    const appLogo = initialState?.appConfig?.logo || '/logo.png';
    
    // 检查是否已登录，已登录则自动跳转
    React.useEffect(() => {
        const token = localStorage.getItem('token');
        
        if (token) {
            // 已登录，获取重定向地址或跳转首页
            const urlParams = new URL(window.location.href).searchParams;
            const redirect = urlParams.get('redirect') || '/';
            
            // 使用 history.push 而不是 window.location.href，避免整页刷新
            window.location.href = redirect;
        }
    }, []);
    
    // 组件加载时，检查是否有记住的密码
    React.useEffect(() => {
        const credentials = getCredentials();
        if (credentials) {
            form.setFieldsValue({
                username: credentials.username,
                password: credentials.password,
                rememberPassword: true,
            });
        }
    }, [form]);

    const fetchUserInfo = async () => {
        const userInfo = await initialState?.fetchUserInfo?.();
        if (userInfo) {
            flushSync(() => {
                setInitialState((s) => ({
                    ...s,
                    currentUser: userInfo,
                }));
            });
        }
    };

    const handleSubmit = async (values: API.LoginParams & { rememberPassword?: boolean }) => {
        try {
            // 调用后端登录接口
            const response = await login(values);
            
            // 登录成功（响应拦截器已处理错误情况，这里只会收到成功的响应）
            if (!response.data) {
                throw new Error('登录响应数据为空');
            }
            const { token, user } = response.data;
            
            // 保存 token 和用户信息
            localStorage.setItem('token', token);
            localStorage.setItem('currentUser', JSON.stringify(user));
            
            // 处理记住密码
            if (values.rememberPassword) {
                saveCredentials(values.username || '', values.password || '');
            } else {
                clearCredentials();
            }
            
            const defaultLoginSuccessMessage = intl.formatMessage({
                id: "pages.login.success",
                defaultMessage: "登录成功！",
            });
            message.success(defaultLoginSuccessMessage);
            
            // 更新全局状态
            await fetchUserInfo();
            
            // 跳转到重定向页面或首页
            const urlParams = new URL(window.location.href).searchParams;
            window.location.href = urlParams.get("redirect") || "/";
        } catch (error: any) {
            // 响应拦截器已将后端错误信息放到 error.message 中
            setUserLoginState({ 
                status: 'error',
                message: error.message || '登录失败，请重试！'
            });
        }
    };
    const { status } = userLoginState;

    return (
        <div className={styles.container}>
            <Helmet>
                <title>
                    {intl.formatMessage({
                        id: "menu.login",
                        defaultMessage: "登录页",
                    })}
                    {appName && ` - ${appName}`}
                </title>
            </Helmet>
            <Lang />
            <div
                style={{
                    flex: "1",
                    marginTop: "50px",
                    padding: "32px 0",
                }}
            >
                <LoginForm
                    form={form}
                    contentStyle={{
                        minWidth: 280,
                        marginTop: "70px",
                        maxWidth: "75vw",
                    }}
                    logo={<img alt="logo" src={appLogo} />}
                    title={appName}
                    initialValues={{
                        rememberPassword: false,
                    }}
                    onFinish={async (values) => {
                        await handleSubmit(values as API.LoginParams);
                    }}
                >
                    {status === "error" && userLoginState.message && (
                        <LoginMessage content={userLoginState.message} />
                    )}
                    <ProFormText
                        name="username"
                        fieldProps={{
                            size: "large",
                            prefix: <UserOutlined />,
                            autoFocus: false,
                        }}
                        placeholder={intl.formatMessage({
                            id: "pages.login.username.placeholder",
                            defaultMessage: "请输入用户名",
                        })}
                        rules={[
                            {
                                required: true,
                                message: (
                                    <FormattedMessage
                                        id="pages.login.username.required"
                                        defaultMessage="请输入用户名!"
                                    />
                                ),
                            },
                        ]}
                    />
                    <ProFormText.Password
                        name="password"
                        fieldProps={{
                            size: "large",
                            prefix: <LockOutlined />,
                            autoFocus: false,
                        }}
                        placeholder={intl.formatMessage({
                            id: "pages.login.password.placeholder",
                            defaultMessage: "请输入密码",
                        })}
                        rules={[
                            {
                                required: true,
                                message: (
                                    <FormattedMessage
                                        id="pages.login.password.required"
                                        defaultMessage="请输入密码！"
                                    />
                                ),
                            },
                        ]}
                    />
                    <div
                        style={{
                            marginBottom: 24,
                        }}
                    >
                        <ProFormCheckbox noStyle name="rememberPassword">
                            <FormattedMessage
                                id="pages.login.rememberPassword"
                                defaultMessage="记住密码"
                            />
                        </ProFormCheckbox>
                    </div>
                </LoginForm>
            </div>
            <Footer />
        </div>
    );
};

export default Login;

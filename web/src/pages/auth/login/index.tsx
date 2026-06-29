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
    // SelectLang,
    history,
    useIntl,
    useModel,
} from "@umijs/max";
import { Alert, App } from "antd";
import { createStyles } from "antd-style";
import React, { useState } from "react";
import { flushSync } from "react-dom";
import { Footer } from "@/components";
import { login, getCaptcha } from '@/services/auth';
import { saveCredentials, getCredentials, clearCredentials } from '@/utils/crypto';
import { getSafeRedirectFromSearch } from '@/utils/redirect';

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
            backgroundImage: "url('/login-bg.png')",
            backgroundSize: "100% 100%",
        },
    };
});

// 语言切换组件（已注释）
// const Lang = () => {
//     const { styles } = useStyles();

//     return (
//         <div className={styles.lang} data-lang>
//             {SelectLang && <SelectLang />}
//         </div>
//     );
// };

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
    const [captchaData, setCaptchaData] = useState<{id: string; image: string}>({id: '', image: 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=='});
    const [form] = ProForm.useForm();
    const { initialState, setInitialState } = useModel("@@initialState");
    const { styles } = useStyles();
    const { message } = App.useApp();
    const intl = useIntl();
    
    const appName = initialState?.appConfig?.name || 'Bico Admin';
    const appLogo = initialState?.appConfig?.logo || '/logo.png';

    const getRedirectPath = () => getSafeRedirectFromSearch(history.location.search);
    
    // 检查是否已登录，已登录则自动跳转
    React.useEffect(() => {
        const token = localStorage.getItem('token');
        
        if (token) {
            // 已登录，获取重定向地址或跳转首页
            history.replace(getRedirectPath());
        }
    }, []);
    
    // 获取验证码
    const fetchCaptcha = async () => {
        try {
            const response = await getCaptcha();
            if (response.data) {
                setCaptchaData({
                    id: response.data.id,
                    image: response.data.image,
                });
            }
        } catch (error) {
            console.error('获取验证码失败', error);
        }
    };
    
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
        fetchCaptcha();
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

    const handleSubmit = async (values: API.LoginParams & { rememberPassword?: boolean; captchaCode?: string }) => {
        try {
            // 调用后端登录接口
            const response = await login({
                ...values,
                captchaId: captchaData.id,
                captchaCode: values.captchaCode || '',
            });
            
            // 登录成功（响应拦截器已处理错误情况，这里只会收到成功的响应）
            if (!response.data?.token) {
                throw new Error('登录响应数据为空');
            }
            
            // 保存 token
            localStorage.setItem('token', response.data.token);
            
            // 处理记住密码
            if (values.rememberPassword) {
                saveCredentials(values.username || '', values.password || '');
            } else {
                clearCredentials();
            }
            
            // 更新全局状态（获取用户信息）
            await fetchUserInfo();
            
            const defaultLoginSuccessMessage = intl.formatMessage({
                id: "pages.login.success",
                defaultMessage: "登录成功！",
            });
            message.success(defaultLoginSuccessMessage);
            
            // 跳转到重定向页面或首页
            history.replace(getRedirectPath());
        } catch (error: any) {
            // 响应拦截器已将后端错误信息放到 error.message 中
            setUserLoginState({ 
                status: 'error',
                message: error.message || '登录失败，请重试！'
            });
            // 登录失败后刷新验证码
            fetchCaptcha();
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
            {/* <Lang /> */}
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
                    <div style={{ display: 'flex', gap: 8 }}>
                        <ProFormText
                            name="captchaCode"
                            fieldProps={{
                                size: "large",
                                autoFocus: false,
                                maxLength: 4,
                            }}
                            placeholder="请输入验证码"
                            rules={[
                                {
                                    required: true,
                                    message: "请输入验证码！",
                                },
                            ]}
                            style={{ flex: 1 }}
                        />
                        <img
                            src={captchaData.image}
                            alt="验证码"
                            onClick={fetchCaptcha}
                            style={{
                                height: 40,
                                cursor: 'pointer',
                                borderRadius: 8,
                                border: '1px solid #d9d9d9',
                                backgroundColor: '#fff',
                            }}
                            title="点击刷新验证码"
                        />
                    </div>
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

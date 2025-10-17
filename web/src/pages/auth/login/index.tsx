import { LockOutlined, UserOutlined } from "@ant-design/icons";
import {
    LoginForm,
    ProFormCheckbox,
    ProFormText,
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
import Settings from "../../../../config/defaultSettings";

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
    const [userLoginState, setUserLoginState] = useState<API.LoginResult>({});
    const { initialState, setInitialState } = useModel("@@initialState");
    const { styles } = useStyles();
    const { message } = App.useApp();
    const intl = useIntl();

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

    const handleSubmit = async (values: API.LoginParams) => {
        try {
            // 使用静态数据直接登录成功
            const msg: API.LoginResult = {
                status: "ok",
                type: "account",
                currentAuthority: "admin",
            };

            const defaultLoginSuccessMessage = intl.formatMessage({
                id: "pages.login.success",
                defaultMessage: "登录成功！",
            });
            message.success(defaultLoginSuccessMessage);
            await fetchUserInfo();
            const urlParams = new URL(window.location.href).searchParams;
            window.location.href = urlParams.get("redirect") || "/";
        } catch (error) {
            const defaultLoginFailureMessage = intl.formatMessage({
                id: "pages.login.failure",
                defaultMessage: "登录失败，请重试！",
            });
            console.log(error);
            message.error(defaultLoginFailureMessage);
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
                    {Settings.title && ` - ${Settings.title}`}
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
                    contentStyle={{
                        minWidth: 280,
                        marginTop: "70px",
                        maxWidth: "75vw",
                    }}
                    logo={<img alt="logo" src="/logo.png" />}
                    title="Bico Admin"
                    initialValues={{
                        autoLogin: true,
                    }}
                    onFinish={async (values) => {
                        await handleSubmit(values as API.LoginParams);
                    }}
                >
                    {status === "error" && (
                        <LoginMessage
                            content={intl.formatMessage({
                                id: "pages.login.accountLogin.errorMessage",
                                defaultMessage:
                                    "账户或密码错误(admin/ant.design)",
                            })}
                        />
                    )}
                    <ProFormText
                        name="username"
                        fieldProps={{
                            size: "large",
                            prefix: <UserOutlined />,
                        }}
                        placeholder={intl.formatMessage({
                            id: "pages.login.username.placeholder",
                            defaultMessage: "用户名: admin or user",
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
                        }}
                        placeholder={intl.formatMessage({
                            id: "pages.login.password.placeholder",
                            defaultMessage: "密码: ant.design",
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
                        <ProFormCheckbox noStyle name="autoLogin">
                            <FormattedMessage
                                id="pages.login.rememberMe"
                                defaultMessage="自动登录"
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

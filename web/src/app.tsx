import type { Settings as LayoutSettings } from "@ant-design/pro-components";
import { SettingDrawer } from "@ant-design/pro-components";
import type { RequestConfig, RunTimeLayoutConfig } from "@umijs/max";
import { history } from "@umijs/max";
import React from "react";
import { AvatarDropdown, AvatarName, Footer } from "@/components";
import { getCurrentUser as fetchCurrentUser } from "@/services/auth";
import { getAppConfig } from "@/services/common";
import { buildLoginUrl, getCurrentPathWithSearch, LOGIN_PATH } from "@/utils/redirect";
import defaultSettings from "../config/defaultSettings";
import { errorConfig } from "./requestErrorConfig";
const loginPath = LOGIN_PATH;

/**
 * @see https://umijs.org/docs/api/runtime-config#getinitialstate
 * */
export async function getInitialState(): Promise<{
    settings?: Partial<LayoutSettings>;
    currentUser?: API.CurrentUser;
    appConfig?: API.AppConfig;
    loading?: boolean;
    fetchUserInfo?: () => Promise<API.CurrentUser | undefined>;
}> {
    const fetchUserInfo = async () => {
        const token = localStorage.getItem("token");

        if (!token) {
            return undefined;
        }

        try {
            // 调用后端接口获取当前用户信息
            const response = await fetchCurrentUser();

            if (response.code === 0 && response.data) {
                // 更新 localStorage 中的用户信息
                localStorage.setItem(
                    "currentUser",
                    JSON.stringify(response.data)
                );
                return response.data;
            }

            return undefined;
        } catch (e) {
            console.error("获取用户信息失败:", e);
            // 清除无效的 token
            localStorage.removeItem("token");
            localStorage.removeItem("currentUser");
            return undefined;
        }
    };

    // 获取应用配置
    let appConfig: API.AppConfig | undefined;
    try {
        const configResponse = await getAppConfig();
        if (configResponse.code === 0 && configResponse.data) {
            appConfig = configResponse.data;
        }
    } catch (e) {
        console.error("获取应用配置失败:", e);
    }

    // 如果不是登录页面，执行
    const { pathname } = window.location;
    if (pathname !== loginPath) {
        const currentUser = await fetchUserInfo();
        return {
            fetchUserInfo,
            currentUser,
            appConfig,
            settings: {
                ...defaultSettings,
                title: appConfig?.name || defaultSettings.title,
                logo: appConfig?.logo || defaultSettings.logo,
            } as Partial<LayoutSettings>,
        };
    }
    return {
        fetchUserInfo,
        appConfig,
        settings: {
            ...defaultSettings,
            title: appConfig?.name || defaultSettings.title,
            logo: appConfig?.logo || defaultSettings.logo,
        } as Partial<LayoutSettings>,
    };
}

// ProLayout 支持的api https://procomponents.ant.design/components/layout
export const layout: RunTimeLayoutConfig = ({
    initialState,
    setInitialState,
}) => {
    return {
        // actionsRender: () => [<SelectLang key="SelectLang" />],
        avatarProps: {
            src: initialState?.currentUser?.avatar,
            title: <AvatarName />,
            render: (_, avatarChildren) => {
                return <AvatarDropdown menu>{avatarChildren}</AvatarDropdown>;
            },
        },
        // waterMarkProps: {
        //     content: initialState?.currentUser?.name,
        // },
        footerRender: () => <Footer />,
        onPageChange: () => {
            const { pathname } = history.location;
            // 如果没有登录，重定向到 login
            if (!initialState?.currentUser && pathname !== loginPath) {
                history.replace(buildLoginUrl(getCurrentPathWithSearch()));
            }
        },
        bgLayoutImgList: [
            {
                src: "/bg1.png",
                left: 85,
                bottom: 100,
                height: "303px",
            },
            {
                src: "/bg2.png",
                bottom: -68,
                right: -45,
                height: "303px",
            },
            {
                src: "/bg3.png",
                bottom: 0,
                left: 0,
                width: "331px",
            },
        ],
        links: [],
        menuHeaderRender: undefined,
        // 自定义 403 页面
        // unAccessible: <div>unAccessible</div>,
        // 增加一个 loading 的状态
        childrenRender: (children) => {
            // if (initialState?.loading) return <PageLoading />;
            return (
                <>
                    {children}
                    {initialState?.appConfig?.debug && (
                        <SettingDrawer
                            disableUrlParams
                            enableDarkTheme
                            settings={initialState?.settings}
                            onSettingChange={(settings) => {
                                setInitialState((preInitialState) => ({
                                    ...preInitialState,
                                    settings: {
                                        ...settings,
                                        title:
                                            preInitialState?.appConfig?.name ||
                                            defaultSettings.title,
                                        logo:
                                            preInitialState?.appConfig?.logo ||
                                            defaultSettings.logo,
                                    },
                                }));
                            }}
                        />
                    )}
                </>
            );
        },
        ...initialState?.settings,
    };
};

/**
 * @name request 配置，可以配置错误处理
 * 它基于 axios 和 ahooks 的 useRequest 提供了一套统一的网络请求和错误处理方案。
 * @doc https://umijs.org/docs/max/request#配置
 */
export const request: RequestConfig = {
    baseURL: "",
    ...errorConfig,
};

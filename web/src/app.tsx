import type {
    MenuDataItem,
    Settings as LayoutSettings,
} from "@ant-design/pro-components";
import { SettingDrawer } from "@ant-design/pro-components";
import type { RequestConfig, RunTimeLayoutConfig } from "@umijs/max";
import { history } from "@umijs/max";
import React from "react";
import { AvatarDropdown, AvatarName } from "@/components";
import { getCurrentUser as fetchCurrentUser } from "@/services/auth";
import { getAppConfig } from "@/services/common";
import { buildLoginUrl, getCurrentPathWithSearch, LOGIN_PATH } from "@/utils/redirect";
import defaultSettings from "../config/defaultSettings";
import { errorConfig } from "./requestErrorConfig";
const loginPath = LOGIN_PATH;

/**
 * MobileDrawerScrollLocker 监听 ProLayout 移动端侧边栏，打开时锁住页面滚动。
 */
const MobileDrawerScrollLocker: React.FC = () => {
    React.useEffect(() => {
        // 运行时布局只在浏览器里同步滚动状态，避免非浏览器环境访问 DOM。
        if (typeof window === "undefined") {
            return;
        }

        const mobileQuery = window.matchMedia("(max-width: 768px)");
        const body = document.body;
        const html = document.documentElement;
        const initialBodyOverflow = body.style.overflow;
        const initialHtmlOverflow = html.style.overflow;

        /**
         * syncScrollLock 根据当前视口和 Drawer 打开状态同步滚动锁。
         */
        const syncScrollLock = () => {
            const hasOpenDrawer = !!document.querySelector(".ant-drawer-open");
            // 只有移动端打开抽屉时才锁滚动，避免影响桌面端弹层和普通页面滚动。
            if (mobileQuery.matches && hasOpenDrawer) {
                body.style.overflow = "hidden";
                html.style.overflow = "hidden";
                return;
            }

            // 抽屉关闭或离开移动端时恢复进入组件前的滚动样式。
            body.style.overflow = initialBodyOverflow;
            html.style.overflow = initialHtmlOverflow;
        };

        const observer = new MutationObserver(syncScrollLock);
        observer.observe(document.body, {
            childList: true,
            subtree: true,
            attributes: true,
            attributeFilter: ["class", "style"],
        });
        mobileQuery.addEventListener("change", syncScrollLock);
        syncScrollLock();

        return () => {
            observer.disconnect();
            mobileQuery.removeEventListener("change", syncScrollLock);
            body.style.overflow = initialBodyOverflow;
            html.style.overflow = initialHtmlOverflow;
        };
    }, []);

    return null;
};

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

    // 获取当前页面路径，用于判断初始化时是否需要拉取用户信息
    const { pathname } = history.location;
    let appConfig: API.AppConfig | undefined;
    let currentUser: API.CurrentUser | undefined;

    try {
        if (pathname !== loginPath) {
            // 在非登录页面，并行获取应用配置和用户信息，大幅减少串行请求带来的首屏白屏时间
            const [configResponse, userInfo] = await Promise.all([
                getAppConfig().catch((e) => {
                    console.error("获取应用配置失败:", e);
                    return null;
                }),
                fetchUserInfo(),
            ]);

            // 接口成功返回时解析并填充应用配置
            if (configResponse && configResponse.code === 0 && configResponse.data) {
                appConfig = configResponse.data;
            }
            currentUser = userInfo;
        } else {
            // 处于登录页面时，仅需要获取应用全局配置，无需调用用户信息接口
            const configResponse = await getAppConfig().catch((e) => {
                console.error("获取应用配置失败:", e);
                return null;
            });
            if (configResponse && configResponse.code === 0 && configResponse.data) {
                appConfig = configResponse.data;
            }
        }
    } catch (e) {
        console.error("系统初始化数据失败:", e);
    }

    // 根据路由分发不同的初始全局状态
    if (pathname !== loginPath) {
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
    /**
     * 过滤菜单数据
     */
    const filterMenuData = (menuData: MenuDataItem[]): MenuDataItem[] => {
        const removeDeveloper = !initialState?.appConfig?.debug;

        return (menuData || [])
            .filter((item) => {
                // 调试模式关闭时隐藏开发者工具，避免生产环境暴露调试入口。
                if (removeDeveloper && item?.path === "/developer") {
                    return false;
                }
                return true;
            })
            .map((item) => {
                if (!item?.children || item.children.length === 0) {
                    return item;
                }
                return {
                    ...item,
                    children: filterMenuData(item.children),
                };
            });
    };

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
        onPageChange: () => {
            const { pathname } = history.location;
            // 如果没有登录，重定向到 login
            if (!initialState?.currentUser && pathname !== loginPath) {
                history.replace(buildLoginUrl(getCurrentPathWithSearch()));
            }
        },
        bgLayoutImgList: [
            {
                src: "bg1.png",
                left: 85,
                bottom: 100,
                height: "303px",
            },
            {
                src: "bg2.png",
                bottom: -68,
                right: -45,
                height: "303px",
            },
            {
                src: "bg3.png",
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
                    <MobileDrawerScrollLocker />
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
        menuDataRender: (menuData) => {
            return filterMenuData(menuData);
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

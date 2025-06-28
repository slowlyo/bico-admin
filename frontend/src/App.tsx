import {
  Refine,
  Authenticated,
} from "@refinedev/core";
import { DashboardOutlined } from "@ant-design/icons";

import {
  ErrorComponent,
  useNotificationProvider,
  ThemedLayoutV2,
  ThemedSiderV2,
} from "@refinedev/antd";
import "@refinedev/antd/dist/reset.css";

import dataProvider from "@refinedev/simple-rest";
import { App as AntdApp } from "antd";
import { BrowserRouter, Route, Routes, Outlet } from "react-router";
import { useTranslation } from "react-i18next";
import { ConfigProvider } from "antd";
import zhCN from "antd/locale/zh_CN";
import routerBindings, {
  CatchAllNavigate,
  UnsavedChangesNotifier,
  DocumentTitleHandler,
} from "@refinedev/react-router";

import { AppIcon } from "./components/app-icon";
import { ColorModeContextProvider } from "./contexts/color-mode";
import { Header } from "./components/header";
import { Login } from "./pages/login";
import { Dashboard } from "./pages/dashboard";
import { Profile } from "./pages/profile";
import { authProvider } from "./authProvider";
import { axiosInstance } from "./utils/request";
import config from "./config";
import "./i18n";

function App() {
  const { t, i18n } = useTranslation();

  const i18nProvider = {
    translate: (key: string, options?: any) => t(key, options) as string,
    changeLocale: (lang: string) => i18n.changeLanguage(lang),
    getLocale: () => i18n.language,
  };

  return (
    <BrowserRouter>
      <ColorModeContextProvider>
        <ConfigProvider locale={zhCN}>
          <AntdApp>
          <Refine
            dataProvider={dataProvider(config.adminApiUrl, axiosInstance)}
            notificationProvider={useNotificationProvider}
            authProvider={authProvider}
            routerProvider={routerBindings}
            i18nProvider={i18nProvider}
            resources={[
              {
                name: "dashboard",
                list: "/",
                meta: {
                  label: "控制台",
                  icon: <DashboardOutlined />,
                },
              },
            ]}
            options={{
              syncWithLocation: true,
              warnWhenUnsavedChanges: true,
              useNewQueryKeys: true,
              title: { text: config.appName, icon: <AppIcon /> },
            }}
          >
                <Routes>
                  <Route
                    element={
                      <Authenticated
                        key="authenticated-inner"
                        fallback={<CatchAllNavigate to="/login" />}
                      >
                        <ThemedLayoutV2
                          Header={Header}
                          Sider={(props) => (
                            <ThemedSiderV2
                              {...props}
                              fixed
                              render={({ items, dashboard }) => {
                                return (
                                  <>
                                    {dashboard}
                                    {items}
                                  </>
                                );
                              }}
                            />
                          )}
                        >
                          <Outlet />
                        </ThemedLayoutV2>
                      </Authenticated>
                    }
                  >
                    <Route index element={<Dashboard />} />
                    <Route path="/profile" element={<Profile />} />
                    <Route path="*" element={<ErrorComponent />} />
                  </Route>
                  <Route
                    element={
                      <Authenticated
                        key="authenticated-outer"
                        fallback={<Outlet />}
                      >
                        <CatchAllNavigate to="/" />
                      </Authenticated>
                    }
                  >
                    <Route path="/login" element={<Login />} />
                  </Route>
                </Routes>

                <UnsavedChangesNotifier />
                <DocumentTitleHandler />
              </Refine>
          </AntdApp>
        </ConfigProvider>
      </ColorModeContextProvider>
    </BrowserRouter>
  );
}

export default App;

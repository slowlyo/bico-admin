import { defineOverridesPreferences } from '@vben/preferences';

/**
 * @description 项目配置文件
 * 只需要覆盖项目中的一部分配置，不需要的配置不用覆盖，会自动使用默认配置
 * !!! 更改配置后请清空缓存，否则可能不生效
 */
export const overridesPreferences = defineOverridesPreferences({
  app: {
    // 菜单来自本地路由，权限码只控制按钮和带 authority 的路由。
    accessMode: 'frontend',
    defaultHomePath: '/dashboard',
    name: import.meta.env.VITE_APP_TITLE,
  },
  // 登录页和侧边栏共用这一张图。
  logo: {
    source: '/logo.png',
  },
  copyright: {
    companyName: 'Bico Admin',
    companySiteLink: '',
    date: '2026',
    enable: false,
  },
  widget: {
    // 顶栏不再放语言、时区、通知。登录页的语言按钮也读 languageToggle。
    languageToggle: false,
    notification: false,
    timezone: false,
  },
  // 主题模式，light 亮色，dark 暗色
  theme: {
    // 默认亮色模式
    mode: 'light',
  },
});

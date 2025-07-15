import { defineConfig } from "@umijs/max";

export default defineConfig({
  history:{
    type: 'hash'
  },
  antd: {},
  access: {},
  model: {},
  initialState: {},
  request: {},
  layout: {
    title: "@umijs/max",
  },
  proxy: {
    "/admin": {
      target: "http://localhost:8899",
      changeOrigin: true,
      ws: true,
    },
  },
  routes: [
    {
      path: "/login",
      component: "./Login",
      layout: false,
    },
    {
      path: "/",
      redirect: "/home",
    },
    {
      name: "首页",
      path: "/home",
      component: "./Home",
      icon: "HomeOutlined",
    },
    {
      name: "个人信息",
      path: "/profile",
      component: "./Profile",
      hideInMenu: true,
    },
    {
      name: "系统管理",
      path: "/system",
      icon: "SettingOutlined",
      access: "canSeeSystemManagement",
      routes: [
        {
          path: "/system",
          redirect: "/system/admin-users",
        },
        {
          name: "管理员",
          path: "/system/admin-users",
          component: "./AdminUser",
          access: "canViewAdminUsers",
        },
        {
          name: "角色",
          path: "/system/roles",
          component: "./Role",
          access: "canViewRoles",
        },
      ],
    },
    {
      name: "权限演示",
      path: "/access",
      component: "./Access",
      hideInMenu: true,
    },
    {
      path: "*",
      component: "./404",
      layout: false,
    },
  ],

  npmClient: "pnpm",
});

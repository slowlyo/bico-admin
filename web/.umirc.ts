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
    },
    {
      name: "权限演示",
      path: "/access",
      component: "./Access",
    },
    {
      name: " CRUD 示例",
      path: "/table",
      component: "./Table",
    },
    {
      name: "管理员用户管理",
      path: "/admin-users",
      component: "./AdminUser",
    },
  ],

  npmClient: "pnpm",
  tailwindcss: {},
});

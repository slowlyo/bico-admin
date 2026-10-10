# 前端接口层

> 接口调用统一放在 `web/apps/admin/src/api/`，基于 Vben 的 `RequestClient` 封装。

## 目录结构

```
src/api/
├── request.ts        # 请求客户端 + 拦截器
├── index.ts          # 统一导出
├── common.ts         # 通用接口（应用配置等）
├── dashboard.ts      # 工作台
├── core/             # 登录、用户信息、菜单
│   ├── auth.ts
│   ├── user.ts
│   └── menu.ts
├── system/           # 系统管理（用户、角色）
│   ├── user.ts
│   └── role.ts
└── demo/             # 示例接口
```

按业务域分目录，一个资源一个文件，文件内同时放该资源的类型定义和接口函数。

## 请求客户端（request.ts）

`requestClient` 是页面使用的客户端，`responseReturn: 'data'` 表示调用方直接拿到 `data` 字段。

拦截器处理：

- 请求头：自动带上 `Authorization: Bearer {token}` 和 `Accept-Language`。
- 响应格式：后端统一为 `{ code, msg, data }`，`code === 0` 视为成功。
- 登录失效：401 时清除 token 并按偏好设置跳转登录或弹出登录框。
- 错误提示：优先使用后端 `msg`，没有时按状态码提示。

`baseRequestClient` 不带业务拦截器，仅用于刷新 token 等特殊场景。

## 编写接口文件

标准 CRUD 模块用 `createCrudApi`（`api/crud.ts`）一行生成 `list/create/update/remove/removeBatch`：

```ts
import type { Defined, PageQuery } from '#/api/crud';
import type { definitions, paths } from '#/api/generated/schema';

import { createCrudApi } from '#/api/crud';

export type AdminUser = Defined<definitions['model.AdminUser']>;
export type AdminUserQuery = NonNullable<paths['/admin-users']['get']['parameters']['query']> &
  PageQuery;

// 泛型依次为：行类型、列表查询参数、创建载荷、更新载荷（不传则与创建相同）。
export const adminUserApi = createCrudApi<
  AdminUser,
  AdminUserQuery,
  definitions['handler.createUserReq'],
  definitions['handler.updateUserReq']
>('/admin-users');

// 使用
adminUserApi.list({ page: 1, pageSize: 20 });
adminUserApi.update(1, payload);
adminUserApi.removeBatch([1, 2]); // DELETE /admin-users/batch，ids 放在请求体里
```

非标准接口（如上传头像、权限树）仍写成独立函数，放在同一个文件里。

实体、查询和载荷类型来自后端 Swagger 生成的 `api/generated/schema.d.ts`（禁止手改）。后端改了请求或响应结构后，执行 `make swagger web-types`，前端用到旧字段的地方会在 `vue-tsc` 中报错。

约定：

- 标准 CRUD 用 `{资源}Api` 对象，其他接口用 `{动作}{资源}Api` 函数，如 `listAllAdminRolesApi`。
- 路径不含 `/admin-api` 前缀，由 `apiURL` 统一提供。
- 后端标准 CRUD 路由：`GET ""`、`GET /:id`、`POST ""`、`PUT /:id`、`DELETE /:id`、`DELETE /batch`。
- 文件上传使用 `requestClient.upload(url, { file })`。

## 后端接口约定

| 类型 | 路径 | 说明 |
| --- | --- | --- |
| 公开 | `/admin-api/auth/login`、`/captcha`、`/app-config` | 无需登录 |
| 登录后 | `/admin-api/auth/*`、`/upload`、`/dashboard/overview` | 需要 token |
| CRUD 模块 | `/admin-api/admin-users`、`/admin-roles` | 需要 token 和对应权限 |

完整接口文档见 Swagger：`/swagger/admin/index.html`。

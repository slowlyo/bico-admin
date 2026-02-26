# AI Agent 项目说明（同步版）

本文档提供当前仓库的快速事实索引，内容以代码为准。

## 项目概览

- 后端：Go + Gin + GORM + Cobra + Viper + Zap
- 前端：Vue 3 + Vite + Element Plus + TypeScript
- 默认数据库：SQLite（`storage/data.db`）
- 缓存：Memory / Redis
- 服务默认端口：`8080`
- 后台入口：`/admin/`
- API 前缀：`/admin-api`

## 目录结构（当前）

```text
internal/
├── admin/      # 后台管理模块（认证、用户、角色、上传等）
├── api/        # 预留业务模块（当前未注册业务路由）
├── core/       # 基础设施（app/config/db/cache/logger/server/upload/middleware）
├── job/        # 定时任务模块
├── migrate/    # 数据迁移
└── pkg/        # 工具包（response/jwt/password/crud/excel/pagination 等）

web/            # 前端工程
config/         # 配置文件
docs/           # 项目文档
```

## 启动流程

入口：[cmd/main.go](../cmd/main.go)

1. `app.BuildContext(configPath)` 构建基础设施上下文。
2. `server.RegisterCoreRoutes(...)` 注册健康检查、Swagger、静态资源。
3. `app.RegisterModules(...)` 注册 `admin`、`api`、`job` 模块。
4. `app.Run(ctx)` 启动 HTTP 服务与调度器。

## 关键路由

- `GET /health`：健康检查
- `GET /swagger/*any`：Swagger 页面
- `GET /admin-api/captcha`：验证码
- `POST /admin-api/auth/login`：登录
- `POST /admin-api/auth/logout`：退出
- `GET /admin-api/auth/current-user`：当前用户
- `GET /admin-api/app-config`：前端应用配置（动态读取 ConfigManager）

## 认证与权限

- JWT 中间件：`internal/core/middleware/jwt.go`
- 权限中间件：`internal/admin/middleware/permission.go`
- 用户状态中间件：`internal/admin/middleware/user_status.go`
- 权限键规范：`模块:资源:动作`，例如 `system:admin_user:list`

## 响应约定

统一响应结构在 [internal/pkg/response/response.go](../internal/pkg/response/response.go)：

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

说明：

- 业务错误多数为 HTTP 200 + `code != 0`
- `BadRequest` 为 HTTP 400
- `TooManyRequests` 为 HTTP 429
- `NotFound` 当前实现为 HTTP 200 + `code=404`

## 配置与热更新

配置文件查找顺序：

1. `-c/--config` 指定路径
2. `./config.yaml`
3. `./config/config.yaml`

`ConfigManager` 会监听文件变化并更新内存配置对象；但大多数基础设施（端口、数据库、限流器、日志器、上传器等）不会自动重建，修改后通常仍需重启。

## CRUD 扩展入口

- CRUD 框架：`internal/pkg/crud`
- 模块装配点：`internal/admin/module.go`
- 新增后台 CRUD 模块时，将 `handler.NewXXXHandler(ctx.DB)` 加入 `modules` 列表。

## 常用命令

```bash
make help
make install
make migrate
make serve
make web
make package
```

## 相关文档

- [项目结构](./structure.md)
- [认证 API](./auth-api.md)
- [配置说明](./config.md)
- [配置热更新](./config-hot-reload.md)
- [CRUD 框架](./crud-pkg.md)
- [Docker 部署](./docker-deploy.md)

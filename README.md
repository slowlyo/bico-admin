# Bico Admin

<img src="./web/public/logo.png" width="200" />

基于 Go + Vue 3 的后台管理系统，支持前后端分离开发与一体化部署。

## 当前状态

- 后端模块：`admin`、`api(预留)`、`job`
- 已实现能力：登录鉴权、RBAC 权限、用户/角色管理、验证码、上传、限流、定时任务、Swagger
- 默认 API 前缀：`/admin-api`
- 默认后台入口：`/admin/`

## 技术栈

### 后端

- [Gin](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [Viper](https://github.com/spf13/viper)
- [Cobra](https://github.com/spf13/cobra)
- [Zap](https://github.com/uber-go/zap)
- [Swaggo](https://github.com/swaggo/swag)

### 前端

- [Vue 3](https://vuejs.org/)
- [TypeScript](https://www.typescriptlang.org/)
- [Vite](https://vitejs.dev/)
- [Element Plus](https://element-plus.org/)
- [TailwindCSS](https://tailwindcss.com/)
- [pnpm](https://pnpm.io/)

## 环境要求

- Go `1.25+`
- Node.js `>=20.19.0`
- pnpm `>=8.8.0`

## 快速开始

### 1. 配置文件

配置文件查找优先级：

1. 命令行参数 `-c/--config`
2. 项目根目录 `config.yaml`
3. 项目目录 `config/config.yaml`

建议先修改 `config/config.yaml`：

```bash
vim config/config.yaml
```

### 2. 安装依赖

```bash
make install
```

### 3. 执行迁移

```bash
make migrate
```

首次迁移会初始化默认管理员账户：

- 用户名：`admin`
- 密码：`admin`

### 4. 启动（前后端分离开发）

本地分离开发建议将 `server.embed_static` 设为 `false`。

```bash
# 推荐：一条命令同时启动前后端
make dev

# 或分开启动
make air
make web
```

开发访问地址：

- 前端：`http://localhost:3006/admin/`
- 后端：`http://localhost:8080`
- Swagger：`http://localhost:8080/swagger/index.html`
- 健康检查：`http://localhost:8080/health`

### 5. 一体化发布（嵌入前端）

```bash
make package
./bin/bico-admin serve -c config/config.yaml
```

访问：`http://localhost:8080/admin/`

## 开发指南

### 新增后台 CRUD 模块

使用 `internal/pkg/crud` 声明式注册模块，详见 [CRUD 包文档](./docs/crud-pkg.md)。

在 `internal/admin/module.go` 的 `modules` 中注册新模块：

```go
modules := []crud.Module{
    handler.NewAdminUserHandler(ctx.DB),
    handler.NewAdminRoleHandler(ctx.DB),
    handler.NewArticleHandler(ctx.DB),
}
```

### 前端路由

在 `web/src/router/modules/` 对应模块路由中配置页面与权限键。

## 常用命令

```bash
make help        # 查看所有可用命令
make serve       # 启动后端服务
make air         # 后端热重载
make dev         # 同时启动前后端开发服务
make web         # 启动前端开发服务器
make migrate     # 数据库迁移
make swagger     # 生成 Swagger 文档
make build       # 编译后端
make build-web   # 编译前端
make package     # 构建嵌入前端的生产版本
make package-win # 构建 Windows 版本
make clean       # 清理构建产物
make tidy        # 整理后端依赖
```

## 项目文档

- [项目结构说明](./docs/structure.md)
- [后端 CRUD 框架](./docs/crud-pkg.md)
- [前端 CRUD 组件](./docs/frontend-crud.md)
- [前端服务封装](./docs/frontend-services.md)
- [认证 API](./docs/auth-api.md)
- [缓存机制](./docs/cache.md)
- [限流中间件](./docs/rate-limit.md)
- [配置说明](./docs/config.md)
- [配置热更新](./docs/config-hot-reload.md)
- [Docker 部署](./docs/docker-deploy.md)

## License

MIT

# Bico Admin

<img width=200 src="./web/src/assets/img/common/logo.png"/>
<br>

基于 Go 语言的多端管理系统，支持 master、admin、api 三个端点的统一管理。

## 快速开始

```bash
# 克隆项目并进入目录
git clone --depth=1 https://github.com/slowlyo/bico-admin.git && cd bico-admin

# 安装依赖
make deps

# 启动开发服务
make dev

# 生成 Swagger 文档
make swagger

# 运行测试
make test
```

启动后访问：
- 管理后台：http://localhost:8899/admin
- API 文档：http://localhost:8899/swagger/index.html

## 技术栈

**后端**
- Go 1.23 + Gin + GORM
- MySQL/PostgreSQL/SQLite + Redis
- Google Wire (依赖注入)
- Zap (日志) + Viper (配置)

**前端**
- [Art Design Pro](https://github.com/Daymychen/art-design-pro) 一款专注于用户体验和视觉设计的后台管理系统模版
- Vue 3 + TypeScript + Element Plus
- Vite + Pinia + Vue Router
- VueUse + Auto Import + Components

## 项目特点

- 🎯 **多端分离**: master(主控) / admin(管理) / api(接口)
- 🏗️ **分层架构**: Handler → Service → Repository → Database
- 🔐 **权限控制**: 基于角色的访问控制 (RBAC)
- 📝 **API 文档**: 自动生成 Swagger 文档
- 🗄️ **多数据库**: 支持 MySQL/PostgreSQL/SQLite
- 🤖 **代码生成**: 基于 MCP 的智能代码生成器，快速生成 CRUD 模块

## 文档导航

- [开发流程指南](docs/development-guide.md) - AI Agent 完整 CRUD 开发流程
- [项目结构说明](docs/project-structure.md) - 详细的目录结构和架构设计
- [数据库配置](docs/database.md) - 多数据库支持和配置说明
- [文件上传](docs/file-upload.md) - 文件上传功能说明
- [MCP 开发工具](docs/mcp-devtools.md) - 开发工具服务使用指南
- [API 文档](docs/swagger.json) - Swagger API 文档

## 开发命令

```bash
make help     # 查看所有可用命令
make run      # 启动服务
make build    # 构建二进制文件
make clean    # 清理文件
make wire     # 生成依赖注入代码
make swagger  # 生成 API 文档

# MCP 开发工具
make devtools         # 启动 MCP 服务 (HTTP模式)
make devtools-help    # 查看 MCP 工具帮助
```

> **注意**: MCP 服务使用 HTTP 传输方式，通过 URL 连接简单可靠。详见 [MCP 开发工具文档](docs/mcp-devtools.md)。
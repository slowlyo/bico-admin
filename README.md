# Bico Admin - AI友好的管理后台框架

🤖 **专为AI开发者设计的现代化管理后台框架**

Bico Admin是一个专门为AI辅助开发优化的全栈管理后台框架，提供清晰的代码结构、完善的文档和智能化的开发体验。

## ✨ 特性

- 🤖 **AI友好设计** - 清晰的代码结构和完善的注释，便于AI理解和协助开发
- ⚡ **高性能后端** - 基于Go Fiber的高性能Web框架
- 🗄️ **现代化ORM** - 使用GORM进行数据库操作，支持MySQL
- 🎨 **优雅前端** - 基于Refine的现代化React管理界面
- 🔧 **开箱即用** - 预配置的开发环境和常用功能模块
- 📚 **完善文档** - 详细的API文档和开发指南
- 🔒 **安全可靠** - 内置认证授权和安全防护机制

## 🏗️ 技术栈

### 后端
- **框架**: [Go Fiber](https://gofiber.io/) - 高性能Web框架
- **ORM**: [GORM](https://gorm.io/) - Go语言ORM库
- **数据库**: MySQL 8.0+
- **认证**: JWT + RBAC权限控制
- **文档**: Swagger/OpenAPI 3.0

### 前端
- **框架**: [Refine](https://refine.dev/) - 企业级React框架
- **UI库**: Ant Design
- **状态管理**: React Query
- **路由**: React Router
- **构建工具**: Vite

## 📁 项目结构

```
bico-admin/
├── backend/                    # 后端Go服务
│   ├── cmd/                   # 应用程序统一入口
│   │   └── server/           # 服务器启动文件
│   │       └── main.go
│   ├── core/                 # 框架核心（系统默认功能，可覆盖更新）
│   │   ├── config/          # 配置管理
│   │   │   ├── config.go
│   │   │   ├── database.go
│   │   │   └── redis.go
│   │   ├── middleware/      # 核心中间件
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   ├── logger.go
│   │   │   └── rate_limit.go
│   │   ├── model/          # 系统基础模型
│   │   │   ├── base.go
│   │   │   ├── user.go
│   │   │   ├── role.go
│   │   │   └── permission.go
│   │   ├── repository/     # 核心数据访问层
│   │   │   ├── base.go
│   │   │   ├── user.go
│   │   │   └── auth.go
│   │   ├── service/        # 核心业务服务
│   │   │   ├── auth.go
│   │   │   ├── user.go
│   │   │   └── rbac.go
│   │   ├── handler/        # 核心处理器
│   │   │   ├── auth.go
│   │   │   ├── user.go
│   │   │   └── system.go
│   │   └── router/         # 核心路由
│   │       ├── auth.go
│   │       ├── system.go
│   │       └── middleware.go
│   ├── admin/              # 后台管理模块（用户自定义业务）
│   │   ├── handler/        # 后台管理处理器
│   │   │   ├── dashboard.go
│   │   │   ├── content.go
│   │   │   └── settings.go
│   │   ├── service/        # 后台管理业务服务
│   │   │   ├── dashboard.go
│   │   │   ├── content.go
│   │   │   └── analytics.go
│   │   ├── model/          # 后台管理数据模型
│   │   │   ├── content.go
│   │   │   ├── category.go
│   │   │   └── settings.go
│   │   ├── repository/     # 后台管理数据访问
│   │   │   ├── content.go
│   │   │   └── analytics.go
│   │   └── router/         # 后台管理路由
│   │       └── admin.go
│   ├── api/                # API模块（对外API服务）
│   │   ├── handler/        # API处理器
│   │   │   ├── app.go
│   │   │   ├── user.go
│   │   │   └── content.go
│   │   ├── service/        # API业务服务
│   │   │   ├── app.go
│   │   │   └── content.go
│   │   ├── model/          # API数据模型
│   │   │   ├── app_user.go
│   │   │   └── api_response.go
│   │   ├── repository/     # API数据访问
│   │   │   └── app.go
│   │   └── router/         # API路由
│   │       └── api.go
│   ├── pkg/                # 公共包（可对外暴露）
│   │   ├── utils/          # 工具函数
│   │   ├── validator/      # 数据验证
│   │   ├── response/       # 响应格式
│   │   └── constants/      # 常量定义
│   ├── business/           # 业务方法封装（常用CRUD操作）
│   │   ├── base.go         # 基础业务方法
│   │   ├── crud.go         # CRUD操作封装
│   │   ├── list.go         # 列表查询封装
│   │   ├── pagination.go   # 分页封装
│   │   └── validation.go   # 业务验证封装
│   ├── docs/               # API文档
│   ├── migrations/         # 数据库迁移文件
│   ├── storage/            # 文件存储目录
│   ├── go.mod
│   ├── go.sum
│   └── .env.example
├── frontend/               # 后台管理前端（遵循Refine标准目录结构）
│   ├── public/            # 静态资源
│   ├── src/               # 源代码（具体结构遵循Refine规范）
│   ├── package.json
│   ├── vite.config.ts
│   └── tsconfig.json
├── docker/                # Docker配置
│   ├── Dockerfile.backend
│   ├── Dockerfile.frontend
│   └── docker-compose.yml
├── scripts/               # 脚本文件
│   ├── build.sh          # 构建脚本
│   ├── deploy.sh         # 部署脚本
│   ├── dev.sh            # 开发环境启动
│   └── update-core.sh    # 框架核心更新脚本
├── docs/                  # 项目文档
│   ├── api/              # API文档
│   ├── development/      # 开发指南
│   ├── deployment/       # 部署指南
│   └── framework/        # 框架更新指南
├── .gitignore
├── README.md
└── LICENSE
```

### 🏗️ 架构设计说明

#### 后端模块化设计

1. **core/** - 框架核心目录
   - 包含系统默认功能和基础组件
   - 可以通过脚本整体覆盖更新
   - 提供用户管理、权限控制、认证等基础功能

2. **admin/** - 后台管理模块
   - 后台管理相关的业务逻辑
   - 包含内容管理、数据分析、系统设置等功能
   - 独立的业务代码，不受框架更新影响

3. **api/** - 对外API模块
   - 为移动端、PC端、第三方提供API服务
   - 统一的API接口和业务逻辑
   - 支持多端数据访问

4. **business/** - 业务方法封装
   - 封装常用的CRUD操作：`CreateOne`, `UpdateById`, `DeleteById`, `GetById`, `List`
   - 提供分页、排序、筛选等通用功能
   - 统一的业务验证和数据处理逻辑
   - AI友好的标准化业务方法

#### 前端结构

- **frontend/** - 后台管理界面，基于Refine框架
- 遵循Refine标准目录结构，便于开发和维护

#### 框架更新机制

- `core/` 目录可以通过 `update-core.sh` 脚本整体更新
- 业务模块（admin、api）和业务封装（business）不受框架更新影响
- 配置文件和自定义代码得到保护

#### 业务方法封装示例

```go
// business/crud.go
type CRUDService[T any] struct {
    db *gorm.DB
}

func (s *CRUDService[T]) CreateOne(data *T) error
func (s *CRUDService[T]) UpdateById(id uint, data *T) error
func (s *CRUDService[T]) DeleteById(id uint) error
func (s *CRUDService[T]) GetById(id uint) (*T, error)
func (s *CRUDService[T]) List(params ListParams) (*ListResult[T], error)
```

## 🚀 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Redis 6.0+ (可选)
- Docker (可选)

### 安装步骤

1. **克隆项目**
   ```bash
   git clone https://github.com/your-username/bico-admin.git
   cd bico-admin
   ```

2. **后端设置**
   ```bash
   cd backend
   cp .env.example .env
   # 编辑 .env 文件配置数据库连接
   go mod tidy

   # 启动统一服务（包含admin和api模块）
   go run cmd/server/main.go
   ```

3. **前端设置**
   ```bash
   # 启动后台管理界面
   cd frontend
   npm install
   npm run dev
   ```

4. **数据库初始化**
   ```bash
   # 运行数据库迁移
   cd backend
   go run migrations/migrate.go
   ```

### 使用Docker

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f backend
docker-compose logs -f frontend
```

### 开发模式

```bash
# 使用开发脚本快速启动
./scripts/dev.sh

# 或者分别启动
./scripts/dev.sh backend   # 启动后端服务
./scripts/dev.sh frontend  # 启动前端服务
```

### 访问地址

- 后台管理界面: http://localhost:3000
- 后台管理API: http://localhost:8080/admin/api
- 对外API接口: http://localhost:8080/api
- API文档: http://localhost:8080/docs

## 📖 文档中心

详细的项目文档请查看 [docs](./docs/) 目录：

- 📚 [文档索引](./docs/README.md) - 所有文档的导航入口
- 🚀 [开发指南](./docs/DEVELOPMENT.md) - 开发环境搭建和日常开发流程
- 🏗️ [项目架构](./docs/ARCHITECTURE.md) - 系统整体架构设计
- 🌐 [API设计](./docs/API_DESIGN.md) - RESTful API设计规范
- 🚀 [部署指南](./docs/DEPLOYMENT.md) - 生产环境部署说明
- 📝 [Git指南](./docs/GIT_GUIDE.md) - Git工作流程和最佳实践

### AI开发友好特性

1. **清晰的代码结构** - 遵循Go和React最佳实践
2. **完善的注释** - 每个函数和组件都有详细说明
3. **标准化命名** - 使用一致的命名规范
4. **模块化设计** - 高内聚低耦合的模块划分
5. **类型安全** - 完整的TypeScript类型定义

### 代码规范

- 后端遵循Go官方代码规范
- 前端使用ESLint + Prettier
- 提交信息遵循Conventional Commits

### API设计

- RESTful API设计
- 统一的响应格式
- 完整的错误处理
- Swagger文档自动生成

## 🔧 配置说明

### 后端配置 (.env)

```env
# 服务器配置
PORT=8080           # 统一服务端口
ENV=development

# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=bico_admin

# Redis配置（可选）
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT配置
JWT_SECRET=your-secret-key
JWT_EXPIRE=24h

# 文件上传配置
UPLOAD_PATH=./storage/uploads
MAX_UPLOAD_SIZE=10MB

# 日志配置
LOG_LEVEL=info
LOG_PATH=./storage/logs

# 路由前缀配置
ADMIN_PREFIX=/admin    # 后台管理路由前缀
API_PREFIX=/api        # 对外API路由前缀
```

### 前端配置

```typescript
// frontend/src/config/index.ts
export const config = {
  adminApiUrl: process.env.VITE_ADMIN_API_URL || 'http://localhost:8080/admin/api',
  apiUrl: process.env.VITE_API_URL || 'http://localhost:8080/api',
  appName: 'Bico Admin',
  version: '1.0.0'
}
```

### 业务方法使用示例

```go
// 在service中使用业务封装
import "your-project/backend/business"

type ContentService struct {
    crud *business.CRUDService[model.Content]
}

func (s *ContentService) CreateContent(data *model.Content) error {
    return s.crud.CreateOne(data)
}

func (s *ContentService) UpdateContent(id uint, data *model.Content) error {
    return s.crud.UpdateById(id, data)
}

func (s *ContentService) GetContentList(params business.ListParams) (*business.ListResult[model.Content], error) {
    return s.crud.List(params)
}
```

### 框架更新

```bash
# 更新框架核心（保留业务代码）
./scripts/update-core.sh

# 备份当前业务代码
./scripts/backup-business.sh

# 恢复业务代码
./scripts/restore-business.sh
```

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 🙏 致谢

- [Go Fiber](https://gofiber.io/) - 高性能Web框架
- [GORM](https://gorm.io/) - 优秀的Go ORM
- [Refine](https://refine.dev/) - 强大的React框架
- [Ant Design](https://ant.design/) - 企业级UI设计语言

## 📞 联系方式

- 项目主页: [GitHub](https://github.com/your-username/bico-admin)
- 问题反馈: [Issues](https://github.com/your-username/bico-admin/issues)
- 邮箱: your-email@example.com

---

⭐ 如果这个项目对您有帮助，请给我们一个星标！
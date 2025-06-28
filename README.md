# Bico Admin - AI友好的管理后台框架

🤖 **专为AI开发者设计的现代化管理后台框架**

Bico Admin是一个专门为AI辅助开发优化的全栈管理后台框架，提供清晰的代码结构、完善的文档和智能化的开发体验。

---

## 🚨 AI开发者必读

> **⚠️ 重要提示**：AI开发者在协助开发前，请务必先查看 `docs/` 目录下的文档！
>
> 🎯 **快速导航**：
> - 🏗️ **简化架构**：[docs/SIMPLIFIED_ARCHITECTURE.md](docs/SIMPLIFIED_ARCHITECTURE.md)
> - 🛠️ **数据库工具**：[docs/DATABASE_TOOLS.md](docs/DATABASE_TOOLS.md)
> - 📋 **API规范**：[docs/API_DESIGN.md](docs/API_DESIGN.md)
> - 🔧 **开发指南**：[docs/DEVELOPMENT.md](docs/DEVELOPMENT.md)
> - ❓ **常见问题**：[docs/FAQ.md](docs/FAQ.md)

---

## 📚 AI开发者指南

> **重要提示：AI开发者请优先查看 `docs/` 目录下的详细文档**

为了更好地理解和使用本项目，AI开发者应该首先查看以下文档：

### 🎯 核心文档
- **[简化架构设计](docs/SIMPLIFIED_ARCHITECTURE.md)** - 了解简化后的架构设计和模块划分
- **[数据库工具指南](docs/DATABASE_TOOLS.md)** - 掌握通用数据库操作工具的使用方法
- **[API设计规范](docs/API_DESIGN.md)** - 掌握API接口设计标准和响应格式
- **[开发指南](docs/DEVELOPMENT.md)** - 快速上手开发流程和最佳实践

### 📖 专项文档
- **[后端开发](docs/backend/README.md)** - Go Fiber + GORM 后端开发详解
- **[前端开发](docs/frontend/README.md)** - Refine + React 前端开发指南
- **[部署指南](docs/DEPLOYMENT.md)** - 生产环境部署配置
- **[常见问题](docs/FAQ.md)** - 开发过程中的常见问题解答

### 🔧 工具文档
- **[Git工作流](docs/GIT_GUIDE.md)** - 代码提交和分支管理规范
- **[缓存工具](docs/CACHE_USAGE.md)** - Redis缓存操作工具使用指南
- **[认证增强](docs/AUTH_IMPROVEMENTS.md)** - 认证逻辑完善和安全功能说明
- **[用户信息编辑](docs/USER_PROFILE_FEATURE.md)** - 用户个人资料编辑和密码修改功能

**建议AI开发者的阅读顺序：**
1. 先阅读 `docs/SIMPLIFIED_ARCHITECTURE.md` 了解简化后的架构设计
2. 查看 `docs/DATABASE_TOOLS.md` 学习通用数据库操作工具
3. 查看 `docs/API_DESIGN.md` 理解API设计规范
4. 根据开发需求查看对应的后端或前端文档
5. 遇到问题时参考 `docs/FAQ.md`

---

## 🤖 AI助手快速入门

如果你是AI助手，在协助开发时请遵循以下步骤：

1. **📖 首先查看文档**：使用 `codebase-retrieval` 工具查询 `docs/` 目录下的相关文档
2. **🏗️ 理解架构**：重点关注项目的模块化设计和目录结构
3. **📋 遵循规范**：严格按照 `docs/API_DESIGN.md` 中的API设计规范
4. **🔍 代码分析**：在修改代码前，先使用检索工具了解相关代码结构
5. **✅ 测试验证**：修改后建议编写或运行相关测试

**重要提醒**：本项目采用模块化设计，`core/` 目录为框架核心，`modules/` 目录为业务模块，请在正确的目录中进行开发。

---

## ✨ 特性

- 🤖 **AI友好设计** - 清晰的代码结构和完善的注释，便于AI理解和协助开发
- ⚡ **高性能后端** - 基于Go Fiber的高性能Web框架
- 🗄️ **现代化ORM** - 使用GORM进行数据库操作，支持MySQL
- 🎨 **优雅前端** - 基于Refine的现代化React管理界面
- 🔧 **开箱即用** - 预配置的开发环境和常用功能模块
- 📚 **完善文档** - 详细的API文档和开发指南
- 🔒 **安全可靠** - 内置认证授权和安全防护机制
- 👤 **用户管理** - 完整的用户信息编辑和密码修改功能

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
│   ├── core/                 # 框架核心 🔥 核心层
│   │   ├── config/          # 配置管理
│   │   │   ├── config.go
│   │   │   ├── database.go
│   │   │   └── redis.go
│   │   ├── database/        # 通用数据库操作工具 🛠️ AI友好
│   │   │   ├── base.go      # 基础CRUD操作
│   │   │   ├── pagination.go # 分页查询工具
│   │   │   └── operations.go # 完整操作工具
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
│   │   ├── handler/        # 核心处理器（业务逻辑直接实现）
│   │   │   ├── auth.go
│   │   │   ├── user.go
│   │   │   └── system.go
│   │   └── router/         # 核心路由
│   │       ├── auth.go
│   │       ├── system.go
│   │       └── middleware.go
│   ├── modules/            # 业务模块目录 📦 业务实现
│   │   ├── admin/          # 后台管理模块（业务逻辑直接实现）
│   │   │   ├── handler/    # 后台管理处理器（业务逻辑直接实现）
│   │   │   │   └── dashboard.go
│   │   │   ├── model/      # 后台管理数据模型
│   │   │   │   ├── content.go
│   │   │   │   ├── category.go
│   │   │   │   └── settings.go
│   │   │   └── router/     # 后台管理路由
│   │   │       └── admin.go
│   │   └── api/            # API模块（对外API服务）
│   │       ├── handler/    # API处理器（业务逻辑直接实现）
│   │       │   └── app.go
│   │       ├── model/      # API数据模型
│   │       │   ├── app_user.go
│   │       │   └── api_response.go
│   │       └── router/     # API路由
│   │           └── api.go
│   ├── pkg/                # 公共包（可对外暴露）
│   │   ├── utils/          # 工具函数
│   │   ├── validator/      # 数据验证
│   │   ├── response/       # 响应格式
│   │   └── constants/      # 常量定义

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
├── docs/                  # 项目文档 📚 AI开发者必读
│   ├── SIMPLIFIED_ARCHITECTURE.md # 简化架构设计 🏗️ 优先阅读
│   ├── DATABASE_TOOLS.md # 数据库工具指南 🛠️ 工具使用必读
│   ├── API_DESIGN.md     # API设计规范 📋 开发必遵循
│   ├── DEVELOPMENT.md    # 开发指南 🔧 快速上手
│   ├── DEPLOYMENT.md     # 部署指南
│   ├── FAQ.md           # 常见问题 ❓ 遇到问题先查看
│   ├── backend/         # 后端开发文档
│   └── frontend/        # 前端开发文档
├── .gitignore
├── README.md
└── LICENSE
```

### 🏗️ 简化架构设计说明

#### 后端简化架构

1. **core/** - 框架核心目录 🔥
   - **database/** - 通用数据库操作工具，提供标准化的CRUD、分页、批量操作
   - **handler/** - 核心处理器，业务逻辑直接在handler中实现
   - **model/** - 数据模型和验证
   - **middleware/** - 中间件和路由配置

2. **modules/admin/** - 后台管理模块 📦
   - **handler/** - 管理处理器，业务逻辑直接实现，调用core/database工具
   - **model/** - 业务数据模型
   - **router/** - 路由配置

3. **modules/api/** - 对外API模块 📦
   - **handler/** - API处理器，业务逻辑直接实现，调用core/database工具
   - **model/** - API数据模型
   - **router/** - API路由配置

#### 🎯 架构特点

- **业务逻辑归位**：所有业务逻辑直接在对应的handler中实现
- **工具层分离**：core/database只提供纯粹的数据库操作工具
- **减少抽象**：移除service和repository层，代码更直接
- **AI友好**：结构清晰，便于AI理解和生成代码

#### 前端结构

- **frontend/** - 后台管理界面，基于Refine框架
- 遵循Refine标准目录结构，便于开发和维护

#### 框架更新机制

- `core/` 目录可以通过 `update-core.sh` 脚本整体更新
- 业务模块（admin、api）和业务封装（business）不受框架更新影响
- 配置文件和自定义代码得到保护

#### 🔥 AI友好的数据库工具示例

```go
// core/database/operations.go
type Operations[T any] struct {
    *PaginationOperations[T]
}

func (ops *Operations[T]) Create(data *T) error
func (ops *Operations[T]) UpdateFields(id uint, updates map[string]interface{}) error
func (ops *Operations[T]) Delete(id uint) error
func (ops *Operations[T]) Get(id uint) (*T, error)
func (ops *Operations[T]) List(params PaginationParams) (*PaginationResult[T], error)
func (ops *Operations[T]) GetByCondition(condition string, args ...interface{}) (*T, error)
func (ops *Operations[T]) BatchCreate(params BatchCreateParams[T]) error
func (ops *Operations[T]) BatchUpdate(params BatchUpdateParams) error
```

#### Handler层业务逻辑直接实现示例

```go
// modules/admin/handler/dashboard.go
type DashboardHandler struct {
    db      *gorm.DB
    userOps *database.Operations[model.User]
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
    return &DashboardHandler{
        db:      db,
        userOps: database.NewOperations[model.User](db),
    }
}

func (h *DashboardHandler) CreateUser(c *fiber.Ctx) error {
    var req model.UserCreateRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request")
    }

    // 业务逻辑直接在handler中实现
    user, err := h.createUser(req)
    if err != nil {
        return response.BadRequest(c, err.Error())
    }

    return response.Success(c, user)
}

// 业务逻辑方法
func (h *DashboardHandler) createUser(req model.UserCreateRequest) (*model.UserResponse, error) {
    // 检查用户名是否已存在
    if existingUser, _ := h.userOps.GetByCondition("username = ?", req.Username); existingUser != nil {
        return nil, errors.New("username already exists")
    }

    // 创建用户
    user := model.User{
        Username: req.Username,
        Email:    req.Email,
        Status:   model.UserStatusActive,
    }

    // 使用通用工具保存
    if err := h.userOps.Create(&user); err != nil {
        return nil, err
    }

    userResponse := user.ToResponse()
    return &userResponse, nil
}
```

## 🚀 快速开始

> **AI开发者提示**：在开始开发前，请先阅读 `docs/DEVELOPMENT.md` 了解详细的开发流程和规范。

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
- 🗄️ [缓存工具](./docs/CACHE_USAGE.md) - Redis缓存操作工具使用指南
- 🔐 [认证增强](./docs/AUTH_IMPROVEMENTS.md) - 认证逻辑完善和安全功能说明

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

## 🤖 AI开发者注意事项

### 重要提醒
1. **优先查看文档**：在进行任何开发工作前，请先查看 `docs/` 目录下的相关文档
2. **遵循简化架构**：严格按照简化架构进行开发，业务逻辑直接在handler中实现
3. **使用通用工具**：充分利用 `core/database/` 下的通用数据库操作工具
4. **API规范**：所有API开发必须遵循 `docs/API_DESIGN.md` 中的设计规范
5. **代码检索**：使用 `codebase-retrieval` 工具深入了解现有代码结构后再进行修改

### 常用开发模式
- **新增功能**：在 `modules/admin/handler/` 或 `modules/api/handler/` 中直接实现业务逻辑
- **数据库操作**：使用 `core/database/` 下的通用工具进行数据操作
- **业务逻辑**：直接在handler中实现，保持代码简洁明了
- **API接口**：遵循RESTful设计，使用统一的响应格式
- **前端组件**：基于Refine框架和Ant Design组件库开发

### 文档更新
当添加新功能时，请同步更新相关文档：
- API变更 → 更新 `docs/API_DESIGN.md`
- 架构变更 → 更新 `docs/SIMPLIFIED_ARCHITECTURE.md`
- 工具使用 → 更新 `docs/DATABASE_TOOLS.md`
- 新增配置 → 更新 `docs/DEVELOPMENT.md`

## 📞 联系方式

- 项目主页: [GitHub](https://github.com/your-username/bico-admin)
- 问题反馈: [Issues](https://github.com/your-username/bico-admin/issues)
- 邮箱: your-email@example.com

---

⭐ 如果这个项目对您有帮助，请给我们一个星标！
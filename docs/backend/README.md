# 后端开发指南

## 🎯 技术栈

- **语言**: Go 1.21+
- **框架**: [Fiber](https://gofiber.io/) - 高性能Web框架
- **ORM**: [GORM](https://gorm.io/) - Go语言ORM库
- **数据库**: MySQL 8.0+
- **缓存**: Redis 6.0+ (可选)
- **认证**: JWT
- **验证**: [go-playground/validator](https://github.com/go-playground/validator)

## 🏗️ 项目结构

```
backend/
├── cmd/server/           # 应用程序入口
├── core/                 # 框架核心功能
│   ├── config/          # 配置管理
│   ├── middleware/      # 中间件
│   ├── model/          # 数据模型
│   ├── repository/     # 数据访问层
│   ├── service/        # 业务逻辑层
│   ├── handler/        # HTTP处理器
│   └── router/         # 路由配置
├── admin/              # 后台管理模块
├── api/                # 对外API模块
├── business/           # 业务方法封装
├── pkg/                # 公共包
├── docs/               # API文档
├── migrations/         # 数据库迁移
└── storage/            # 文件存储
```

## 🚀 快速开始

### 1. 环境准备
```bash
# 安装Go 1.21+
go version

# 安装依赖
cd backend
go mod tidy

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件
```

### 2. 数据库设置
```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE bico_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 运行迁移
go run cmd/server/main.go
```

### 3. 启动服务
```bash
# 开发模式
go run cmd/server/main.go

# 或使用Makefile
make dev-backend
```

## 📋 开发规范

### 代码结构规范

#### 1. 分层架构
```
Handler (HTTP层) → Service (业务层) → Repository (数据层) → Model (数据模型)
```

#### 2. 命名规范
- **文件名**: 使用小写字母和下划线，如 `user_service.go`
- **包名**: 使用小写字母，如 `package service`
- **结构体**: 使用大驼峰命名，如 `UserService`
- **方法**: 使用大驼峰命名，如 `CreateUser`
- **变量**: 使用小驼峰命名，如 `userID`

#### 3. 目录组织
```go
// 每个模块包含完整的分层结构
module/
├── handler/     # HTTP处理器
├── service/     # 业务逻辑
├── repository/  # 数据访问
├── model/       # 数据模型
└── router/      # 路由配置
```

### 代码示例

#### 1. 模型定义
```go
// core/model/user.go
package model

type User struct {
    BaseModel
    Username string     `json:"username" gorm:"uniqueIndex;size:50;not null" validate:"required,min=3,max=50"`
    Email    string     `json:"email" gorm:"uniqueIndex;size:100;not null" validate:"required,email"`
    Password string     `json:"-" gorm:"size:255;not null" validate:"required,min=6"`
    Status   UserStatus `json:"status" gorm:"default:1"`
}

type UserStatus int

const (
    UserStatusInactive UserStatus = 0
    UserStatusActive   UserStatus = 1
)
```

#### 2. 仓储层
```go
// core/repository/user.go
package repository

type UserRepository interface {
    Create(user *model.User) error
    GetByID(id uint) (*model.User, error)
    GetByUsername(username string) (*model.User, error)
    Update(id uint, user *model.User) error
    Delete(id uint) error
}

type UserRepositoryImpl struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(user *model.User) error {
    return r.db.Create(user).Error
}
```

#### 3. 服务层
```go
// core/service/user.go
package service

type UserService interface {
    CreateUser(req model.UserCreateRequest) (*model.UserResponse, error)
    GetUser(id uint) (*model.UserResponse, error)
    UpdateUser(id uint, req model.UserUpdateRequest) (*model.UserResponse, error)
    DeleteUser(id uint) error
}

type UserServiceImpl struct {
    userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
    return &UserServiceImpl{userRepo: userRepo}
}

func (s *UserServiceImpl) CreateUser(req model.UserCreateRequest) (*model.UserResponse, error) {
    // 业务逻辑处理
    user := model.User{
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
    }
    
    if err := user.HashPassword(); err != nil {
        return nil, err
    }
    
    if err := s.userRepo.Create(&user); err != nil {
        return nil, err
    }
    
    response := user.ToResponse()
    return &response, nil
}
```

#### 4. 处理器层
```go
// core/handler/user.go
package handler

type UserHandler struct {
    userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var req model.UserCreateRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body")
    }
    
    if errors := validator.Validate(req); len(errors) > 0 {
        return response.ValidationError(c, errors)
    }
    
    user, err := h.userService.CreateUser(req)
    if err != nil {
        return response.BadRequest(c, err.Error())
    }
    
    return response.SuccessWithMessage(c, "User created successfully", user)
}
```

## 🔧 业务封装使用

### CRUD操作封装
```go
// 使用通用CRUD服务
import "bico-admin/business"

type ContentService struct {
    crud *business.CRUDService[model.Content]
}

func NewContentService(db *gorm.DB) *ContentService {
    return &ContentService{
        crud: business.NewCRUDService[model.Content](db),
    }
}

// 使用封装的方法
func (s *ContentService) CreateContent(data *model.Content) error {
    return s.crud.CreateOne(data)
}

func (s *ContentService) GetContentList(params business.ListParams) (*business.ListResult[model.Content], error) {
    return s.crud.List(params)
}
```

## 🔐 认证和授权

### JWT中间件使用
```go
// 在路由中使用认证中间件
func SetupProtectedRoutes(app fiber.Router, db *gorm.DB) {
    jwtSecret := config.GetJWTSecret()
    
    // 需要认证的路由组
    protected := app.Group("/")
    protected.Use(middleware.AuthMiddleware(jwtSecret))
    
    userHandler := handler.NewUserHandler(db)
    protected.Get("/profile", userHandler.GetProfile)
    protected.Put("/profile", userHandler.UpdateProfile)
}
```

### 权限验证
```go
// 在处理器中获取当前用户
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
    userID := middleware.GetUserID(c)
    if userID == 0 {
        return response.Unauthorized(c, "User not authenticated")
    }
    
    // 处理业务逻辑
    user, err := h.userService.GetUser(userID)
    if err != nil {
        return response.NotFound(c, "User not found")
    }
    
    return response.Success(c, user)
}
```

## 📊 数据验证

### 使用验证器
```go
// 在处理器中验证请求数据
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var req model.UserCreateRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body")
    }
    
    // 验证请求参数
    if errors := validator.Validate(req); len(errors) > 0 {
        return response.ValidationError(c, errors)
    }
    
    // 处理业务逻辑
    // ...
}
```

### 自定义验证规则
```go
// 在validator包中注册自定义验证器
func init() {
    validator.Validator.RegisterValidation("username", validateUsername)
}

func validateUsername(fl validator.FieldLevel) bool {
    username := fl.Field().String()
    // 自定义验证逻辑
    return len(username) >= 3 && len(username) <= 50
}
```

## 🗄️ 数据库操作

### GORM最佳实践
```go
// 1. 预加载关联数据
user, err := db.Preload("Roles").First(&user, id).Error

// 2. 批量操作
var users []model.User
db.Where("status = ?", model.UserStatusActive).Find(&users)

// 3. 事务处理
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err
    }
    if err := tx.Create(&profile).Error; err != nil {
        return err
    }
    return nil
})

// 4. 软删除
db.Delete(&user, id) // 软删除
db.Unscoped().Delete(&user, id) // 硬删除
```

## 🧪 测试

### 单元测试
```go
// user_service_test.go
package service

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestUserService_CreateUser(t *testing.T) {
    // 创建mock仓储
    mockRepo := &MockUserRepository{}
    userService := NewUserService(mockRepo)
    
    // 设置mock期望
    mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)
    
    // 执行测试
    req := model.UserCreateRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    user, err := userService.CreateUser(req)
    
    // 断言结果
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "testuser", user.Username)
    
    // 验证mock调用
    mockRepo.AssertExpectations(t)
}
```

### 集成测试
```go
// integration_test.go
func TestUserAPI(t *testing.T) {
    // 设置测试数据库
    db := setupTestDB()
    defer cleanupTestDB(db)
    
    // 创建测试应用
    app := setupTestApp(db)
    
    // 测试创建用户
    req := httptest.NewRequest("POST", "/api/users", strings.NewReader(`{
        "username": "testuser",
        "email": "test@example.com",
        "password": "password123"
    }`))
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, 201, resp.StatusCode)
}
```

## 📝 最佳实践

1. **错误处理**
   - 使用自定义错误类型
   - 提供清晰的错误信息
   - 记录详细的错误日志

2. **性能优化**
   - 使用数据库索引
   - 实现查询缓存
   - 避免N+1查询问题

3. **安全性**
   - 验证所有输入数据
   - 使用参数化查询
   - 实现适当的权限控制

4. **代码质量**
   - 编写单元测试
   - 使用代码检查工具
   - 遵循Go语言规范

5. **文档**
   - 为公共API编写文档
   - 使用有意义的注释
   - 保持文档更新

# Core Database 通用工具使用指南

## 概述

Core Database 层提供了完整的通用数据库操作工具，涵盖了常见的CRUD操作、分页查询、批量操作等功能。这些工具经过精心设计，具有良好的错误处理、类型安全和性能优化。

## 核心组件

### 1. BaseOperations[T] - 基础数据库操作

提供最基本的数据库操作方法。

#### 创建操作
```go
// 创建用户操作实例
userOps := database.NewOperations[model.User](db)

// 创建单个记录
user := &model.User{
    Username: "john",
    Email:    "john@example.com",
}
err := userOps.CreateOne(user)
// 或使用别名
err := userOps.Create(user)
```

#### 查询操作
```go
// 根据ID获取记录
user, err := userOps.GetById(1)
// 或使用别名
user, err := userOps.Get(1)

// 根据ID列表批量获取
users, err := userOps.GetByIds([]uint{1, 2, 3})

// 根据条件获取单个记录
user, err := userOps.GetByCondition("username = ?", "admin")

// 根据多个条件获取记录
user, err := userOps.GetByConditions(map[string]interface{}{
    "username": "admin",
    "status":   1,
})

// 根据条件获取多个记录
users, err := userOps.FindByCondition("status = ?", 1)
users, err := userOps.FindByConditions(map[string]interface{}{
    "status": 1,
    "role":   "admin",
})

// 检查记录是否存在
exists, err := userOps.Exists(1)

// 统计记录总数
count, err := userOps.Count()

// 根据条件统计
count, err := userOps.CountWithCondition("status = ?", 1)
count, err := userOps.CountWithConditions(map[string]interface{}{
    "status": 1,
})
```

#### 更新操作
```go
// 根据ID更新记录
user.Username = "john_updated"
err := userOps.UpdateById(1, user)
// 或使用别名
err := userOps.Update(1, user)

// 使用map更新指定字段
err := userOps.UpdateByIdWithMap(1, map[string]interface{}{
    "username": "john_updated",
    "status":   1,
})
// 或使用别名
err := userOps.UpdateFields(1, map[string]interface{}{
    "username": "john_updated",
    "status":   1,
})
```

#### 删除操作
```go
// 软删除
err := userOps.DeleteById(1)
// 或使用别名
err := userOps.Delete(1)

// 硬删除
err := userOps.HardDeleteById(1)
```

### 2. PaginationOperations[T] - 分页查询工具

继承 BaseOperations，提供高级查询和批量操作。

#### 分页列表查询
```go
params := database.PaginationParams{
    Page:         1,
    PageSize:     10,
    Sort:         "created_at",
    Order:        "desc",
    Search:       "john",
    SearchFields: []string{"username", "email", "nickname"},
    Filters: map[string]interface{}{
        "status": 1,
        "role":   "admin",
    },
    Preloads: []string{"Roles", "Profile"},
}

result, err := userOps.Paginate(params)
// 或使用别名
result, err := userOps.List(params)

// result.Data - 数据列表
// result.Total - 总记录数
// result.Page - 当前页
// result.PageSize - 每页大小
// result.TotalPages - 总页数
```

#### 简化的分页查询
```go
// 基本分页查询
result, err := userOps.GetWithPagination(1, 10, map[string]interface{}{
    "status": 1,
})

// 搜索分页查询
result, err := userOps.SearchWithPagination(1, 10, "john", []string{"username", "email"})
```

#### 批量操作
```go
// 批量创建
users := []model.User{
    {Username: "user1", Email: "user1@example.com"},
    {Username: "user2", Email: "user2@example.com"},
}

params := database.BatchCreateParams[model.User]{
    Data:      users,
    BatchSize: 100, // 批次大小，默认100
}
err := userOps.BatchCreate(params)

// 批量更新
updateParams := database.BatchUpdateParams{
    IDs: []uint{1, 2, 3},
    Updates: map[string]interface{}{
        "status":     1,
        "updated_at": time.Now(),
    },
}
err := userOps.BatchUpdate(updateParams)

// 批量删除
deleteParams := database.BatchDeleteParams{
    IDs:        []uint{1, 2, 3},
    HardDelete: false, // false=软删除, true=硬删除
}
err := userOps.BatchDelete(deleteParams)
```

### 3. Operations[T] - 完整操作工具

包含所有功能的完整工具集。

#### 事务操作
```go
// 执行事务
err := userOps.Transaction(func(tx *gorm.DB) error {
    // 在事务中使用新的操作实例
    txUserOps := userOps.WithDB(tx)
    
    // 创建用户
    if err := txUserOps.Create(&user); err != nil {
        return err
    }
    
    // 分配角色
    if err := txUserOps.UpdateFields(user.ID, map[string]interface{}{
        "role_id": roleID,
    }); err != nil {
        return err
    }
    
    return nil
})
```

#### 原生SQL操作
```go
// 执行原生查询
rows := userOps.Raw("SELECT * FROM users WHERE status = ?", 1)

// 执行原生命令
err := userOps.Exec("UPDATE users SET last_login_at = NOW() WHERE id = ?", userID)
```

## Handler 层使用示例

### 标准CRUD Handler
```go
type UserHandler struct {
    db      *gorm.DB
    userOps *database.Operations[model.User]
}

func NewUserHandler(db *gorm.DB) *UserHandler {
    return &UserHandler{
        db:      db,
        userOps: database.NewOperations[model.User](db),
    }
}

// 获取用户列表
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
    
    params := database.PaginationParams{
        Page:         page,
        PageSize:     pageSize,
        Search:       c.Query("search"),
        SearchFields: []string{"username", "email", "nickname"},
    }
    
    result, err := h.userOps.List(params)
    if err != nil {
        return response.InternalServerError(c, "Failed to get users")
    }
    
    return response.Success(c, result)
}

// 创建用户
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var user model.User
    if err := c.BodyParser(&user); err != nil {
        return response.BadRequest(c, "Invalid request body")
    }
    
    if err := h.userOps.Create(&user); err != nil {
        return response.InternalServerError(c, "Failed to create user")
    }
    
    return response.Success(c, user)
}

// 更新用户
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
    id, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return response.BadRequest(c, "Invalid user ID")
    }
    
    var updates map[string]interface{}
    if err := c.BodyParser(&updates); err != nil {
        return response.BadRequest(c, "Invalid request body")
    }
    
    if err := h.userOps.UpdateFields(uint(id), updates); err != nil {
        return response.InternalServerError(c, "Failed to update user")
    }
    
    return response.SuccessWithMessage(c, "User updated successfully", nil)
}

// 删除用户
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
    id, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return response.BadRequest(c, "Invalid user ID")
    }
    
    if err := h.userOps.Delete(uint(id)); err != nil {
        return response.InternalServerError(c, "Failed to delete user")
    }
    
    return response.SuccessWithMessage(c, "User deleted successfully", nil)
}
```

### 业务逻辑实现示例
```go
// 用户登录业务逻辑
func (h *AuthHandler) performLogin(req model.UserLoginRequest) (*model.UserResponse, string, error) {
    // 根据用户名或邮箱查找用户
    user, err := h.userOps.GetByCondition("username = ? OR email = ?", req.Username, req.Username)
    if err != nil {
        return nil, "", errors.New("invalid credentials")
    }

    // 验证密码
    if !user.CheckPassword(req.Password) {
        return nil, "", errors.New("invalid credentials")
    }

    // 检查用户状态
    if user.Status != model.UserStatusActive {
        return nil, "", errors.New("user account is not active")
    }

    // 更新最后登录时间
    now := time.Now()
    h.userOps.UpdateFields(user.ID, map[string]interface{}{
        "last_login_at": &now,
    })

    // 生成JWT token
    token, err := middleware.GenerateToken(user.ID, user.Username, h.jwtSecret, h.jwtExpire)
    if err != nil {
        return nil, "", err
    }

    userResponse := user.ToResponse()
    return &userResponse, token, nil
}
```

## 工具函数

### 参数验证和构建
```go
// 验证分页参数
params := database.PaginationParams{Page: 0, PageSize: -1}
database.ValidatePaginationParams(&params)
// params.Page = 1, params.PageSize = 10

// 构建搜索参数
params := database.BuildSearchParams(1, 10, "john", []string{"username", "email"})

// 构建过滤参数
params := database.BuildFilterParams(1, 10, map[string]interface{}{
    "status": 1,
})

// 构建完整参数
params := database.BuildCompleteParams(1, 10, "created_at", "desc", "john", 
    []string{"username", "email"}, 
    map[string]interface{}{"status": 1}, 
    []string{"Roles"})
```

## 最佳实践

### 1. 错误处理
```go
// ✅ 推荐：详细的错误处理
user, err := userOps.Get(id)
if err != nil {
    if strings.Contains(err.Error(), "记录不存在") {
        return response.NotFound(c, "User not found")
    }
    return response.InternalServerError(c, "Failed to get user")
}
```

### 2. 参数验证
```go
// ✅ 推荐：验证必要参数
params := database.PaginationParams{
    Page:     page,
    PageSize: pageSize,
    Search:   search,
}
database.ValidatePaginationParams(&params)
```

### 3. 性能优化
```go
// ✅ 推荐：使用预加载
params := database.PaginationParams{
    Preloads: []string{"Roles", "Profile"},
}

// ✅ 推荐：批量操作
err := userOps.BatchUpdate(database.BatchUpdateParams{
    IDs: []uint{1, 2, 3},
    Updates: map[string]interface{}{"status": 1},
})
```

### 4. 事务处理
```go
// 复杂业务逻辑使用事务
err := userOps.Transaction(func(tx *gorm.DB) error {
    txUserOps := userOps.WithDB(tx)
    
    // 在事务中执行多个操作
    if err := txUserOps.Create(&user); err != nil {
        return err
    }
    
    if err := txUserOps.UpdateFields(user.ID, updates); err != nil {
        return err
    }
    
    return nil
})
```

通过这些通用工具，你可以快速实现大部分常见的数据库操作，同时保持代码的一致性和可维护性。所有的业务逻辑都直接在handler中实现，保持了架构的简洁性。

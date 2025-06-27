# API 设计规范

## 🎯 设计原则

### 1. RESTful 设计
- 使用标准的HTTP方法 (GET, POST, PUT, DELETE)
- 资源导向的URL设计
- 无状态的请求处理
- 统一的响应格式

### 2. 版本控制
- URL路径版本控制: `/api/v1/users`
- 向后兼容性保证
- 废弃API的优雅处理

### 3. 安全性
- JWT Token 认证
- 权限验证
- 请求限流
- 数据验证

## 🌐 API 路由结构

### 路由前缀
```
/api          # 对外API接口
/admin/api    # 后台管理API接口
/auth         # 认证相关接口
/system       # 系统信息接口
```

### 完整路由示例
```
# 认证接口
POST   /auth/login           # 用户登录
POST   /auth/register        # 用户注册
POST   /auth/logout          # 用户登出
GET    /auth/profile         # 获取用户资料
PUT    /auth/profile         # 更新用户资料
POST   /auth/change-password # 修改密码

# 对外API接口
GET    /api/app/info         # 获取应用信息
GET    /api/app/config       # 获取应用配置
GET    /api/user/profile     # 获取用户资料
PUT    /api/user/profile     # 更新用户资料
GET    /api/content          # 获取内容列表
GET    /api/content/:id      # 获取单个内容

# 后台管理API接口
GET    /admin/api/dashboard       # 获取仪表板数据
GET    /admin/api/dashboard/stats # 获取统计数据
GET    /admin/api/users           # 获取用户列表
POST   /admin/api/users           # 创建用户
GET    /admin/api/users/:id       # 获取单个用户
PUT    /admin/api/users/:id       # 更新用户
DELETE /admin/api/users/:id       # 删除用户

# 系统接口
GET    /health               # 健康检查
GET    /system/info          # 系统信息
```

## 📋 HTTP 方法规范

| 方法   | 用途           | 示例                    |
|--------|----------------|-------------------------|
| GET    | 获取资源       | GET /api/users          |
| POST   | 创建资源       | POST /api/users         |
| PUT    | 更新整个资源   | PUT /api/users/1        |
| PATCH  | 部分更新资源   | PATCH /api/users/1      |
| DELETE | 删除资源       | DELETE /api/users/1     |

## 📊 响应格式规范

### 统一响应结构
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    // 响应数据
  }
}
```

### 分页响应结构
```json
{
  "code": 200,
  "message": "Success",
  "data": [
    // 数据列表
  ],
  "total": 100,
  "page": 1,
  "size": 10
}
```

### 错误响应结构
```json
{
  "code": 400,
  "message": "Validation failed",
  "data": {
    "errors": [
      {
        "field": "email",
        "tag": "email",
        "value": "invalid-email",
        "message": "email must be a valid email"
      }
    ]
  }
}
```

## 🔢 状态码规范

### 成功状态码
- `200 OK` - 请求成功
- `201 Created` - 资源创建成功
- `204 No Content` - 请求成功但无返回内容

### 客户端错误
- `400 Bad Request` - 请求参数错误
- `401 Unauthorized` - 未认证
- `403 Forbidden` - 权限不足
- `404 Not Found` - 资源不存在
- `409 Conflict` - 资源冲突
- `422 Unprocessable Entity` - 验证失败
- `429 Too Many Requests` - 请求过于频繁

### 服务器错误
- `500 Internal Server Error` - 服务器内部错误
- `502 Bad Gateway` - 网关错误
- `503 Service Unavailable` - 服务不可用

## 🔐 认证和授权

### JWT Token 格式
```
Authorization: Bearer <token>
```

### Token 响应
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 权限验证
- 所有需要认证的接口都需要携带有效的JWT Token
- 权限不足时返回403状态码
- Token过期时返回401状态码

## 📝 请求参数规范

### 查询参数 (Query Parameters)
```
GET /api/users?page=1&size=10&sort=id&order=desc&search=admin
```

### 路径参数 (Path Parameters)
```
GET /api/users/1
PUT /api/users/1
DELETE /api/users/1
```

### 请求体 (Request Body)
```json
{
  "username": "admin",
  "email": "admin@example.com",
  "password": "password123"
}
```

## 📋 分页参数规范

### 请求参数
```json
{
  "page": 1,        // 页码，从1开始
  "size": 10,       // 每页数量，默认10，最大100
  "sort": "id",     // 排序字段
  "order": "desc",  // 排序方向：asc/desc
  "search": "关键词" // 搜索关键词
}
```

### 响应格式
```json
{
  "code": 200,
  "message": "Success",
  "data": [
    // 数据列表
  ],
  "total": 100,      // 总记录数
  "page": 1,         // 当前页码
  "size": 10,        // 每页数量
  "total_pages": 10  // 总页数
}
```

## 🔍 筛选和搜索

### 筛选参数
```
GET /api/users?status=active&role=admin&created_at_gte=2024-01-01
```

### 搜索参数
```
GET /api/users?search=admin&search_fields=username,email
```

## 📊 批量操作

### 批量创建
```json
POST /api/users/batch
{
  "data": [
    {"username": "user1", "email": "user1@example.com"},
    {"username": "user2", "email": "user2@example.com"}
  ]
}
```

### 批量更新
```json
PUT /api/users/batch
{
  "ids": [1, 2, 3],
  "data": {
    "status": "active"
  }
}
```

### 批量删除
```json
DELETE /api/users/batch
{
  "ids": [1, 2, 3]
}
```

## 🔄 API 版本控制

### URL版本控制
```
/api/v1/users    # 版本1
/api/v2/users    # 版本2
```

### 版本兼容性
- 新版本保持向后兼容
- 废弃的API提供迁移指南
- 版本生命周期管理

## 📖 API 文档

### Swagger/OpenAPI
- 自动生成API文档
- 交互式API测试
- 代码示例生成

### 文档访问
```
GET /docs         # API文档首页
GET /docs/swagger # Swagger UI
GET /docs/openapi # OpenAPI规范
```

## 🧪 API 测试

### 测试用例
- 正常流程测试
- 异常情况测试
- 边界值测试
- 性能测试

### 测试工具
- Postman 集合
- 自动化测试脚本
- 压力测试工具

## 📈 性能优化

### 缓存策略
- 响应缓存
- 数据库查询缓存
- 静态资源缓存

### 限流策略
- 基于IP的限流
- 基于用户的限流
- 基于接口的限流

## 🔍 监控和日志

### 请求日志
```json
{
  "timestamp": "2024-01-01T00:00:00Z",
  "method": "GET",
  "path": "/api/users",
  "status": 200,
  "duration": "10ms",
  "ip": "127.0.0.1",
  "user_id": 1
}
```

### 错误日志
```json
{
  "timestamp": "2024-01-01T00:00:00Z",
  "level": "error",
  "message": "Database connection failed",
  "error": "connection refused",
  "path": "/api/users",
  "user_id": 1
}
```

## 📝 最佳实践

1. **URL设计**
   - 使用名词而不是动词
   - 使用复数形式
   - 保持URL简洁明了

2. **错误处理**
   - 提供清晰的错误信息
   - 使用标准的HTTP状态码
   - 包含错误代码和描述

3. **数据验证**
   - 服务端验证所有输入
   - 提供详细的验证错误信息
   - 使用标准的验证规则

4. **安全性**
   - 始终验证用户权限
   - 对敏感数据进行加密
   - 防止SQL注入和XSS攻击

5. **性能**
   - 合理使用缓存
   - 优化数据库查询
   - 实现请求限流

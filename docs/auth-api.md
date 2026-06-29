# 后台用户认证 API 文档

## 概述

后台用户认证系统已实现，包括登录、退出登录功能，使用 JWT 进行身份验证。

## 数据库表结构

### admin_users 表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | uint | 主键，自增 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| username | string(64) | 用户名，唯一索引，非空 |
| password | string(255) | 密码，非空（明文存储，生产环境应使用 bcrypt） |
| name | string(64) | 姓名 |
| avatar | string(255) | 头像 URL |
| enabled | bool | 启用状态，默认 true |

## API 接口

### 1. 用户登录

**接口地址：** `POST /admin-api/login`

**请求参数：**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**成功响应：**
```json
{
  "code": 0,
  "msg": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "name": "管理员",
      "avatar": "",
      "enabled": true
    }
  }
}
```

**错误响应：**

- 用户不存在：
```json
{
  "code": 401,
  "msg": "用户不存在"
}
```

- 密码错误：
```json
{
  "code": 401,
  "msg": "密码错误"
}
```

- 用户已禁用：
```json
{
  "code": 401,
  "msg": "用户已被禁用"
}
```

### 2. 用户退出登录

**接口地址：** `POST /admin-api/logout`

**请求头：**
```
Authorization: Bearer {token}
```

**成功响应：**
```json
{
  "code": 0,
  "msg": "退出成功"
}
```

**说明：** 退出登录会将 token 加入黑名单，token 将在 7 天后自动过期清除

## JWT 配置

JWT 相关配置位于 `config/config.yaml`：

```yaml
jwt:
  secret: "bico-admin-secret-key-change-in-production"  # 密钥，生产环境必须修改
  expire_hours: 168  # 过期时间（小时），默认 7 天
```

**注意：** 生产环境务必修改 secret 为强随机字符串！

## 使用方法

### 1. 执行数据库迁移

```bash
go run ./cmd/main.go migrate
```

### 2. 创建管理员用户

```bash
go run ./scripts/create_admin.go
```

默认创建的管理员账号：
- 用户名：admin
- 密码：admin123

### 3. 启动服务

```bash
go run ./cmd/main.go serve
```

或编译后运行：

```bash
go build -o bico-admin ./cmd/main.go
./bico-admin serve
```

### 4. 测试接口

使用提供的测试脚本：

```bash
./scripts/test_auth.sh
```

或手动测试：

```bash
# 登录
curl -X POST http://localhost:8080/admin-api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 退出登录
curl -X POST http://localhost:8080/admin-api/logout
```

## 代码结构

```
internal/admin/
├── model/
│   └── admin_user.go          # AdminUser 模型
├── service/
│   └── auth_service.go        # 认证服务（登录逻辑）
├── handler/
│   └── auth_handler.go        # 认证处理器（HTTP 处理）
└── router.go                  # 路由注册

internal/shared/jwt/
├── jwt.go                     # JWT 管理器
└── token.go                   # JWT 令牌生成和解析

config/
└── config.yaml                # 配置文件（包含 JWT 配置）
```

## 注意事项

1. ✅ **密码存储：** 已使用 bcrypt 加密存储
2. **JWT Secret：** 默认密钥仅供开发使用，生产环境必须修改
3. ✅ **Token 黑名单：** 已实现退出登录 token 黑名单功能（基于缓存）
4. ✅ **Token 验证：** 已实现 JWT 认证中间件（`internal/core/middleware/jwt.go`）
5. ✅ **权限验证：** 已实现权限验证中间件（`internal/admin/middleware/permission.go`）
6. ✅ **用户状态检查：** 已实现用户状态中间件，自动拦截已禁用用户

## 主要改进

### 1. SQLite 数据库存储优化
- 数据库文件统一存放在 `./storage/data.db`
- 便于备份和管理

### 2. 查询优化
- 使用 GORM Model 执行查询，而非表名硬编码
- 代码更清晰，便于维护

### 3. 密码加密
- 使用 bcrypt 对密码进行加密存储
- 登录时自动验证加密密码

### 4. 自动初始化管理员
- 执行迁移时自动检查用户表
- 如果为空，自动创建默认管理员账户（admin/admin）
- 避免手动创建用户的繁琐步骤

### 5. Token 黑名单
- 退出登录时将 token 加入黑名单
- 基于缓存系统实现（支持 memory/redis）
- Token 在黑名单中保留 7 天后自动清除

## 后续优化建议

1. ✅ ~~实现 JWT 验证中间件~~ （已实现）
2. 添加 token 刷新机制
3. 添加登录日志记录
4. ✅ ~~实现用户权限管理~~ （已实现完整的 RBAC 权限系统）
5. 添加 API 访问频率限制（可使用 rate limiting 中间件）
6. 实现登录验证码功能
7. 支持多因素认证（MFA）

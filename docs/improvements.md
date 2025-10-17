# 认证系统优化总结

## 优化内容

根据您的需求，对后台用户认证系统进行了以下 5 项重要优化：

### 1. ✅ SQLite 数据库路径调整

**问题：** 数据库文件分散在根目录，不便于管理

**优化：**
- 将 SQLite 数据库文件统一存放到 `./storage/data.db`
- 配置文件路径：`config/config.yaml` → `sqlite.path: storage/data.db`
- 便于统一备份和管理数据

**代码位置：**
- `config/config.yaml` - 配置路径修改
- `.gitignore` - 添加 storage 目录到忽略列表

---

### 2. ✅ 使用 Model 执行查询

**问题：** 使用表名字符串查询，容易出错且不便维护

**优化前：**
```go
err := s.db.Table("admin_users").Where("username = ?", loginReq.Username).First(&user).Error
```

**优化后：**
```go
var user model.AdminUser
err := s.db.Where("username = ?", loginReq.Username).First(&user).Error
```

**优点：**
- 代码更清晰，避免硬编码表名
- 利用 GORM Model 的自动表名映射
- 类型安全，IDE 可以自动补全
- 便于重构和维护

**代码位置：**
- `internal/admin/service/auth_service.go` - `Login()` 方法

---

### 3. ✅ 自动初始化管理员账户

**问题：** 每次迁移后需要手动创建管理员账户，操作繁琐

**优化：**
- 在数据库迁移时自动检查 `admin_users` 表
- 如果表为空，自动创建默认管理员账户
- 默认账户：`admin/admin`（密码已加密）

**执行效果：**
```bash
$ go run ./cmd/main.go migrate
📦 开始数据库迁移...
✅ 初始化管理员账户成功 (用户名: admin, 密码: admin)
✅ 数据库迁移完成
```

**代码位置：**
- `internal/migrate/migrate.go` - `initAdminUser()` 函数

---

### 4. ✅ 完善密码 bcrypt 加密

**问题：** 之前密码为明文存储，存在安全隐患

**优化：**
- 创建独立的密码加密工具包 `internal/shared/password`
- 使用 bcrypt 算法加密密码（DefaultCost = 10）
- 登录时自动验证加密密码
- 初始化管理员时自动加密密码

**实现：**
```go
// 加密密码
hashedPassword, _ := password.Hash("admin")

// 验证密码
isValid := password.Verify(hashedPassword, plainPassword)
```

**代码位置：**
- `internal/shared/password/password.go` - 密码加密工具
- `internal/admin/service/auth_service.go` - 使用密码验证
- `internal/migrate/migrate.go` - 初始化时加密密码

---

### 5. ✅ 实现 Token 黑名单（退出登录）

**问题：** JWT 是无状态的，退出登录无法让 token 失效

**优化：**
- 退出登录时将 token 加入黑名单
- 基于缓存系统实现（支持 memory/redis）
- Token 在黑名单中保留 7 天，过期自动清除
- 提供 `IsTokenBlacklisted()` 方法供中间件使用

**使用方式：**
```bash
# 退出登录需要携带 token
curl -X POST http://localhost:8080/admin-api/logout \
  -H "Authorization: Bearer {your_token}"
```

**实现细节：**
- Token 存储键：`token:blacklist:{token}`
- 有效期：7 天（与 JWT 过期时间一致）
- 缓存驱动：可配置 memory 或 redis

**代码位置：**
- `internal/admin/service/auth_service.go` - `Logout()` 和 `IsTokenBlacklisted()`
- `internal/admin/handler/auth_handler.go` - 从请求头获取 token
- `internal/core/cache/` - 缓存系统

---

## 技术栈

- **密码加密：** `golang.org/x/crypto/bcrypt`
- **JWT：** 自实现（HMAC-SHA256）
- **缓存：** 内存缓存 / Redis（可配置）
- **数据库：** GORM + SQLite

---

## 配置文件

`config/config.yaml`:
```yaml
database:
  driver: sqlite
  sqlite:
    path: storage/data.db  # 数据库路径

cache:
  driver: memory  # 缓存驱动: memory / redis

jwt:
  secret: "bico-admin-secret-key-change-in-production"
  expire_hours: 168  # 7天
```

---

## 测试验证

### 1. 数据库迁移
```bash
$ go run ./cmd/main.go migrate
📦 开始数据库迁移...
✅ 初始化管理员账户成功 (用户名: admin, 密码: admin)
✅ 数据库迁移完成
```

### 2. 登录测试
```bash
$ curl -X POST http://localhost:8080/admin-api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

# 返回
{
  "code": 0,
  "msg": "登录成功",
  "data": {
    "token": "eyJhbG...",
    "user": {
      "id": 1,
      "username": "admin",
      "name": "系统管理员",
      "enabled": true
    }
  }
}
```

### 3. 退出登录测试
```bash
$ curl -X POST http://localhost:8080/admin-api/logout \
  -H "Authorization: Bearer eyJhbG..."

# 返回
{
  "code": 0,
  "msg": "退出成功"
}
```

---

## 文件变更清单

### 新增文件
- `internal/admin/model/admin_user.go` - 管理员用户模型
- `internal/admin/service/auth_service.go` - 认证服务
- `internal/admin/handler/auth_handler.go` - 认证处理器
- `internal/shared/password/password.go` - 密码加密工具
- `internal/shared/jwt/jwt.go` - JWT 管理器
- `internal/shared/jwt/token.go` - JWT 实现
- `docs/improvements.md` - 本文档

### 修改文件
- `config/config.yaml` - 添加 JWT 配置，修改数据库路径
- `internal/core/config/config.go` - 添加 JWTConfig 结构
- `internal/migrate/migrate.go` - 添加自动初始化管理员功能
- `internal/admin/router.go` - 注册登录/退出路由
- `internal/core/app/container.go` - 注册依赖注入
- `internal/core/cache/factory.go` - 优化接口定义
- `internal/core/cache/redis.go` - 优化接口实现
- `.gitignore` - 添加 storage 目录和编译文件

---

## 后续建议

1. **JWT 验证中间件**
   - 实现统一的 token 验证中间件
   - 自动检查 token 黑名单
   - 保护需要认证的接口

2. **Token 刷新机制**
   - 实现 refresh token
   - 延长用户会话时间

3. **登录日志**
   - 记录登录时间、IP、设备信息
   - 便于审计和安全分析

4. **用户权限管理**
   - 实现 RBAC 角色权限控制
   - 细粒度的接口访问控制

5. **API 频率限制**
   - 防止暴力破解
   - 保护系统资源

---

## 总结

所有需求已全部实现：
- ✅ SQLite 数据库存放到 ./storage 目录
- ✅ 使用 Model 执行查询而非表名
- ✅ 迁移时自动初始化 admin/admin 账户
- ✅ 完善 bcrypt 密码加密逻辑
- ✅ 退出登录时 token 加入黑名单

系统现在更加安全、易用、易维护！

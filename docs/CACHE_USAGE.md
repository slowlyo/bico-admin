# 缓存工具使用指南

## 概述

本项目封装了常用的Redis缓存操作方法，提供了简单易用的缓存接口。

## 基本使用

### 获取缓存实例

```go
import "bico-admin/core/cache"

// 获取全局缓存实例
c := cache.GetCache()

// 或者创建新实例
c := cache.NewCache()
```

### 基本操作

#### 设置缓存

```go
// 设置对象缓存
user := &model.User{ID: 1, Username: "admin"}
err := c.Set("user:1", user, 30*time.Minute)

// 设置字符串缓存
err := c.SetString("key", "value", time.Hour)
```

#### 获取缓存

```go
// 获取对象缓存
var user model.User
err := c.Get("user:1", &user)

// 获取字符串缓存
value, err := c.GetString("key")
```

#### 删除缓存

```go
// 删除单个缓存
err := c.Delete("key")

// 删除匹配模式的缓存
err := c.DeletePattern("user:*")
```

#### 检查缓存是否存在

```go
exists, err := c.Exists("key")
```

### 高级操作

#### 计数器操作

```go
// 递增计数器
count, err := c.Increment("counter")

// 按指定值递增
count, err := c.IncrementBy("counter", 5)
```

#### 条件设置

```go
// 仅当key不存在时设置
success, err := c.SetNX("key", "value", time.Hour)
```

#### 批量操作

```go
// 批量获取
values, err := c.MGet("key1", "key2", "key3")

// 批量设置
err := c.MSet("key1", "value1", "key2", "value2")
```

## 认证相关缓存

### Token黑名单

```go
// 添加token到黑名单
err := c.AddTokenToBlacklist(token, time.Hour)

// 检查token是否在黑名单中
isBlacklisted, err := c.IsTokenBlacklisted(token)
```

### 用户会话

```go
// 设置用户会话
sessionData := map[string]interface{}{
    "user_id": 1,
    "username": "admin",
    "login_time": time.Now(),
}
err := c.SetUserSession(1, sessionData, time.Hour)

// 获取用户会话
var session map[string]interface{}
err := c.GetUserSession(1, &session)

// 删除用户会话
err := c.DeleteUserSession(1)
```

### 登录尝试限制

```go
// 增加登录尝试次数
count, err := c.IncrementLoginAttempt("user:192.168.1.1", 15*time.Minute)

// 获取登录尝试次数
count, err := c.GetLoginAttemptCount("user:192.168.1.1")

// 重置登录尝试次数
err := c.ResetLoginAttempt("user:192.168.1.1")
```

## 缓存键前缀常量

```go
const (
    TokenBlacklistPrefix = "token_blacklist:"
    UserSessionPrefix    = "user_session:"
    LoginAttemptPrefix   = "login_attempt:"
    CaptchaPrefix        = "captcha:"
)
```

## 错误处理

```go
value, err := c.GetString("key")
if err != nil {
    if err.Error() == "key not found" {
        // 键不存在
        return nil
    }
    // 其他错误
    return err
}
```

## 最佳实践

1. **合理设置过期时间**：避免缓存永不过期导致内存泄漏
2. **使用有意义的键名**：建议使用前缀和分隔符，如 `user:1:profile`
3. **错误处理**：始终检查缓存操作的错误返回值
4. **批量操作**：对于多个相关操作，优先使用批量方法
5. **避免大对象**：缓存大对象可能影响性能，考虑分片存储

## 配置

缓存依赖Redis配置，确保在 `config.yaml` 中正确配置Redis连接信息：

```yaml
redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
```

## 注意事项

1. 缓存实例依赖Redis连接，确保Redis服务正常运行
2. 序列化使用JSON格式，确保缓存的对象可以被JSON序列化
3. 缓存操作是异步的，不保证强一致性
4. 在高并发场景下，注意缓存雪崩和缓存穿透问题

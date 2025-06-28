# 认证逻辑完善说明

## 概述

本次更新完善了认证逻辑，实现了token黑名单机制、登录尝试限制、用户会话管理等安全功能，并封装了常用的缓存操作方法。

## 主要改进

### 1. 缓存工具封装 (`core/cache/cache.go`)

#### 基本功能
- **通用缓存操作**：Set、Get、Delete、Exists等
- **字符串缓存**：SetString、GetString
- **模式删除**：DeletePattern支持通配符删除
- **计数器操作**：Increment、IncrementBy
- **条件设置**：SetNX仅在key不存在时设置
- **批量操作**：MGet、MSet
- **过期时间管理**：Expire、TTL

#### 认证专用方法
- **Token黑名单**：AddTokenToBlacklist、IsTokenBlacklisted
- **用户会话**：SetUserSession、GetUserSession、DeleteUserSession
- **登录尝试限制**：IncrementLoginAttempt、GetLoginAttemptCount、ResetLoginAttempt

### 2. 登出逻辑完善

#### Token黑名单机制
```go
// 将token添加到黑名单
err := h.cache.AddTokenToBlacklist(tokenString, h.jwtExpire)

// 在JWT中间件中检查黑名单
isBlacklisted, err := cacheInstance.IsTokenBlacklisted(tokenString)
if err == nil && isBlacklisted {
    return response.Unauthorized(c, "Token has been revoked")
}
```

#### 会话清理
- 登出时自动删除用户会话缓存
- 清理相关的用户状态信息

### 3. 登录安全增强

#### 登录尝试限制
- **限制规则**：同一用户名+IP组合，15分钟内最多尝试5次
- **自动重置**：登录成功后自动重置尝试次数
- **渐进式惩罚**：可扩展为更复杂的限制策略

```go
// 检查登录尝试次数
attemptCount, err := h.cache.GetLoginAttemptCount(loginIdentifier)
if err == nil && attemptCount >= 5 {
    return response.Unauthorized(c, "Too many login attempts, please try again later")
}

// 登录失败时增加次数
h.cache.IncrementLoginAttempt(loginIdentifier, 15*time.Minute)

// 登录成功时重置次数
h.cache.ResetLoginAttempt(loginIdentifier)
```

#### 用户会话管理
- **会话数据**：存储用户ID、用户名、登录时间、登录IP
- **自动过期**：会话过期时间与JWT一致
- **状态同步**：登录时创建，登出时删除

### 4. JWT中间件增强

#### 黑名单检查
- 在token解析前检查是否在黑名单中
- 被撤销的token立即失效
- 提供明确的错误信息

#### 性能优化
- 缓存检查在token解析前进行，减少不必要的计算
- 使用Redis的高性能特性

## 安全特性

### 1. Token安全
- **即时撤销**：登出后token立即失效
- **防重放攻击**：黑名单机制防止token重复使用
- **自动清理**：黑名单条目自动过期，避免内存泄漏

### 2. 登录安全
- **暴力破解防护**：限制登录尝试次数
- **IP级别限制**：结合用户名和IP进行限制
- **时间窗口**：15分钟的限制窗口，平衡安全性和用户体验

### 3. 会话安全
- **状态一致性**：登录状态在缓存和JWT中保持一致
- **自动清理**：会话数据自动过期
- **多设备支持**：支持同一用户多设备登录管理

## 配置要求

### Redis配置
```yaml
redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
```

### 环境变量
- 确保Redis服务正常运行
- 配置正确的连接参数

## 使用示例

### 基本认证流程
1. **登录**：验证用户凭据，创建JWT和会话
2. **请求验证**：检查JWT有效性和黑名单状态
3. **登出**：撤销JWT，清理会话

### 缓存操作
```go
// 获取缓存实例
cache := cache.GetCache()

// 设置缓存
err := cache.Set("key", data, time.Hour)

// 获取缓存
var data MyStruct
err := cache.Get("key", &data)
```

## 错误处理

### 缓存错误
- Redis连接失败时，认证功能降级但不中断服务
- 缓存操作失败时记录日志但不影响主要流程

### 认证错误
- 提供明确的错误信息
- 区分不同类型的认证失败

## 性能考虑

### 缓存策略
- 合理设置过期时间
- 使用批量操作减少网络开销
- 避免缓存大对象

### 内存管理
- 自动清理过期数据
- 使用有意义的键名便于管理

## 扩展性

### 功能扩展
- 支持更复杂的限制策略
- 可添加验证码机制
- 支持多级缓存

### 监控扩展
- 登录尝试统计
- 缓存命中率监控
- 安全事件记录

## 注意事项

1. **Redis依赖**：缓存功能依赖Redis，确保服务可用性
2. **数据一致性**：缓存和数据库之间的数据一致性
3. **安全配置**：根据实际需求调整限制参数
4. **监控告警**：建议添加相关监控和告警机制

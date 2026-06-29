# 限流中间件

## 功能介绍

基于令牌桶算法的限流中间件，用于防止 API 被恶意刷新和保护服务器资源。

## 配置说明

在 `config/config.yaml` 中配置限流参数：

```yaml
rate_limit:
  enabled: true  # 是否启用限流
  rps: 100       # 每秒请求数（Requests Per Second）
  burst: 200     # 突发流量桶容量
```

### 参数详解

- **enabled**: 是否启用限流功能
  - `true`: 启用全局限流
  - `false`: 禁用限流

- **rps**: 每秒允许的请求数量
  - 示例：设置为 100，表示每秒最多允许 100 个请求
  - 建议根据服务器性能调整

- **burst**: 突发流量桶容量
  - 允许短时间内的请求突发
  - 建议设置为 rps 的 2-3 倍
  - 示例：rps=100, burst=200，表示瞬间可处理 200 个请求

## 使用场景

### 1. 全局限流（已启用）

框架默认在所有路由上启用基于 IP 的全局限流：

```go
// 已自动集成，无需手动配置
engine.Use(rateLimiter.RateLimit())
```

**特点**：
- 按客户端 IP 地址限流
- 每个 IP 独立计数
- 自动清理过期限流器

### 2. 用户维度限流

对已认证用户使用用户 ID 限流：

```go
// 在路由中使用
authorized := engine.Group("/api", jwtAuth, rateLimiter.RateLimitByUser())
```

**特点**：
- 未认证用户按 IP 限流
- 已认证用户按 user_id 限流
- 更精确的流量控制

### 3. 自定义限流

根据业务需求自定义限流 key：

```go
// 按 API 路径限流
rateLimiter.RateLimitByKey(func(c *gin.Context) string {
    return c.Request.URL.Path
})

// 按用户角色限流
rateLimiter.RateLimitByKey(func(c *gin.Context) string {
    role := c.GetString("user_role")
    return "role:" + role
})

// 按组合条件限流
rateLimiter.RateLimitByKey(func(c *gin.Context) string {
    userID := c.GetString("user_id")
    path := c.Request.URL.Path
    return fmt.Sprintf("%s:%s", userID, path)
})
```

## 响应格式

当触发限流时，返回 HTTP 429 状态码：

```json
{
  "code": 429,
  "msg": "请求过于频繁，请稍后再试"
}
```

## 性能优化

### 内存管理

限流器会自动清理过期的限流记录：
- 清理周期：5 分钟
- 触发条件：限流器数量 > 10000

### 配置建议

**低流量服务**（< 10 QPS）：
```yaml
rate_limit:
  enabled: true
  rps: 10
  burst: 20
```

**中流量服务**（100 QPS）：
```yaml
rate_limit:
  enabled: true
  rps: 100
  burst: 200
```

**高流量服务**（1000+ QPS）：
```yaml
rate_limit:
  enabled: true
  rps: 1000
  burst: 2000
```

**开发环境**（关闭限流）：
```yaml
rate_limit:
  enabled: false
  rps: 0
  burst: 0
```

## 监控和调试

### 观察限流效果

使用压测工具测试限流：

```bash
# 使用 ab (Apache Bench)
ab -n 1000 -c 100 http://localhost:8080/admin-api/health

# 使用 wrk
wrk -t 10 -c 100 -d 30s http://localhost:8080/admin-api/health
```

### 日志记录

限流触发时，客户端会收到 429 响应，可通过访问日志统计：

```bash
# 统计 429 响应数量
grep "429" logs/access.log | wc -l
```

## 高级用法

### 不同路由不同限流策略

```go
// 登录接口：严格限流（防止暴力破解）
loginLimiter := middleware.NewRateLimiter(5, 10)  // 每秒5次，突发10次
engine.POST("/login", loginLimiter.RateLimit(), loginHandler)

// 普通 API：正常限流
normalLimiter := middleware.NewRateLimiter(100, 200)
engine.GET("/api/users", normalLimiter.RateLimit(), usersHandler)

// 文件上传：宽松限流
uploadLimiter := middleware.NewRateLimiter(10, 20)
engine.POST("/upload", uploadLimiter.RateLimit(), uploadHandler)
```

### 白名单配置

```go
func (rl *RateLimiter) RateLimitWithWhitelist(whitelist []string) gin.HandlerFunc {
    whitelistMap := make(map[string]bool)
    for _, ip := range whitelist {
        whitelistMap[ip] = true
    }
    
    return func(c *gin.Context) {
        ip := c.ClientIP()
        if whitelistMap[ip] {
            c.Next()
            return
        }
        
        limiter := rl.getLimiter(ip)
        if !limiter.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "code": 429,
                "msg":  "请求过于频繁，请稍后再试",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}

// 使用
whitelist := []string{"127.0.0.1", "192.168.1.100"}
engine.Use(rateLimiter.RateLimitWithWhitelist(whitelist))
```

## 注意事项

1. **反向代理配置**
   - 如果使用 Nginx 等反向代理，确保正确传递真实 IP
   - Gin 会自动从 `X-Real-IP` 或 `X-Forwarded-For` 获取 IP

2. **分布式部署**
   - 当前实现为单机内存限流
   - 多机部署需使用 Redis 实现分布式限流

3. **配置调整**
   - 限流配置支持热更新
   - 修改配置文件后自动生效

4. **测试环境**
   - 建议在测试环境关闭限流或设置较大值
   - 避免影响自动化测试

## 故障排查

### 问题：所有请求都被限流

**原因**：`rps` 设置过低或 `burst` 不足

**解决**：调整配置参数
```yaml
rate_limit:
  rps: 200      # 提高限流阈值
  burst: 400    # 增加突发容量
```

### 问题：限流不生效

**原因**：限流未启用或配置错误

**解决**：
1. 检查 `enabled: true`
2. 确认 `rps > 0` 和 `burst > 0`
3. 重启服务生效

### 问题：内存占用过高

**原因**：限流器数量过多未清理

**解决**：
- 当前会自动清理（每 5 分钟）
- 可调整清理策略或降低限流器数量阈值

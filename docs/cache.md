# Cache 缓存模块

## 能力概览

统一缓存接口 `Cache`，支持两种实现：

1. `memory`：进程内缓存
2. `redis`：Redis 缓存

接口定义位置：`internal/core/cache/cache.go`

## 配置示例

```yaml
cache:
  driver: memory  # memory / redis
  redis:
    host: localhost
    port: 6379
    password: ""
    db: 0
```

## 行为差异

### Memory 缓存

- 存储值为原始 Go 对象
- 每分钟清理一次过期键
- 进程重启后缓存丢失

实现：`internal/core/cache/memory.go`

### Redis 缓存

- 写入时做 JSON 序列化
- 读取时做 JSON 反序列化
- 支持跨实例共享缓存

实现：`internal/core/cache/redis.go`

## 使用示例

```go
c, err := cache.NewCache(&cfg.Cache)
if err != nil {
    return err
}
defer c.Close()

_ = c.Set("user:1", map[string]any{"name": "admin"}, 10*time.Minute)
v, err := c.Get("user:1")
exists := c.Exists("user:1")
_ = c.Delete("user:1")
_ = c.Clear()
```

## 注意事项

1. 仅 Redis 实现使用 JSON 序列化；Memory 不做 JSON 转换。
2. `expiration=0` 时，Memory 视为不过期；Redis 由 Redis 语义处理。
3. 切换 `memory -> redis` 后，建议重启服务并检查 Redis 连通性。

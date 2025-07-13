# 缓存包使用说明

这是一个通用的缓存包，支持内存和Redis两种驱动，默认使用内存缓存。

## 特性

- 🚀 **多驱动支持**: 支持内存缓存和Redis缓存
- 🔧 **配置灵活**: 通过配置文件轻松切换驱动
- 💾 **内存缓存**: 高性能的内存缓存，支持过期时间和自动清理
- 🔗 **Redis缓存**: 支持Redis集群，连接池管理
- 📦 **高级功能**: 提供JSON序列化、GetOrSet、Remember等高级功能
- 🛡️ **类型安全**: 完整的错误处理和类型安全
- ⚡ **高性能**: 优化的内存使用和并发安全

## 配置

### 默认配置（内存缓存）

```yaml
cache:
  driver: "memory"          # 缓存驱动: memory(内存) 或 redis
  memory:
    max_size: 10000         # 最大缓存条目数，0表示无限制
    default_expiration: "30m"  # 默认过期时间
    cleanup_interval: "10m"    # 清理过期条目的间隔时间
```

### Redis缓存配置

```yaml
cache:
  driver: "redis"
  redis:
    host: "localhost"
    port: 6379
    password: ""
    database: 1
    pool_size: 10
    min_idle_conns: 5
    dial_timeout: "5s"
    read_timeout: "3s"
    write_timeout: "3s"
    key_prefix: "bico:cache:"
```

## 基本使用

### 1. 在项目中使用（推荐）

```go
// 通过依赖注入获取缓存管理器
func NewUserService(cacheManager *cache.Manager) *UserService {
    return &UserService{
        cache: cacheManager,
    }
}

// 在服务中使用缓存
func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
    // 尝试从缓存获取
    var user User
    err := s.cache.GetOrSetJSON(ctx, "user:"+userID, &user, func() (interface{}, error) {
        // 缓存未命中，从数据库获取
        return s.repo.GetUser(ctx, userID)
    }, 10*time.Minute)
    
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}
```

### 2. 直接使用缓存包

```go
import "bico-admin/pkg/cache"

// 创建缓存实例
config := cache.Config{
    Driver: "memory", // 或 "redis"
    Memory: cache.MemoryConfig{
        MaxSize:           10000,
        DefaultExpiration: 30 * time.Minute,
        CleanupInterval:   10 * time.Minute,
    },
}

cacheInstance, err := cache.NewCache(config)
if err != nil {
    panic(err)
}
defer cacheInstance.Close()

// 创建缓存管理器
manager := cache.NewManager(cacheInstance)
```

## API 参考

### 基本操作

```go
// 设置字符串值
err := manager.SetString(ctx, "key", "value", 5*time.Minute)

// 获取字符串值
value, err := manager.GetString(ctx, "key")

// 检查键是否存在
exists, err := manager.Exists(ctx, "key")

// 删除缓存
err := manager.Delete(ctx, "key")

// 设置过期时间
err := manager.Expire(ctx, "key", 10*time.Minute)
```

### JSON操作

```go
// 设置JSON对象
user := &User{ID: 1, Name: "John"}
err := manager.SetJSON(ctx, "user:1", user, 10*time.Minute)

// 获取JSON对象
var user User
err := manager.GetJSON(ctx, "user:1", &user)
```

### 高级功能

```go
// GetOrSet 模式
value, err := manager.GetOrSet(ctx, "expensive:key", func() (string, error) {
    // 执行耗时操作
    return computeExpensiveValue(), nil
}, 1*time.Hour)

// Remember 模式（记忆化）
result, err := manager.Remember(ctx, "api:data", 30*time.Minute, func() (interface{}, error) {
    // 调用外部API
    return callExternalAPI()
})

// GetOrSetJSON 模式
var data APIResponse
err := manager.GetOrSetJSON(ctx, "api:response", &data, func() (interface{}, error) {
    return fetchFromAPI()
}, 15*time.Minute)
```

### 统计信息

```go
// 获取缓存统计
stats, err := manager.GetStats(ctx)
fmt.Printf("驱动: %s, 键数量: %d\n", stats.Driver, stats.KeyCount)
```

## 常见模式

### 1. 缓存穿透保护

```go
func (s *Service) GetUser(userID string) (*User, error) {
    var user User
    err := s.cache.GetOrSetJSON(ctx, "user:"+userID, &user, func() (interface{}, error) {
        u, err := s.db.GetUser(userID)
        if err != nil {
            if errors.Is(err, ErrUserNotFound) {
                // 缓存空值防止缓存穿透
                return &User{ID: userID, NotFound: true}, nil
            }
            return nil, err
        }
        return u, nil
    }, 5*time.Minute)
    
    if err != nil {
        return nil, err
    }
    
    if user.NotFound {
        return nil, ErrUserNotFound
    }
    
    return &user, nil
}
```

### 2. 多级缓存

```go
func (s *Service) GetData(key string) (string, error) {
    // L1缓存（内存，短期）
    if value, err := s.l1Cache.GetString(ctx, key); err == nil {
        return value, nil
    }
    
    // L2缓存（Redis，长期）
    value, err := s.l2Cache.GetOrSet(ctx, key, func() (string, error) {
        // 从数据源获取
        return s.dataSource.Get(key)
    }, 1*time.Hour)
    
    if err != nil {
        return "", err
    }
    
    // 回填L1缓存
    s.l1Cache.SetString(ctx, key, value, 5*time.Minute)
    
    return value, nil
}
```

## 性能建议

1. **内存缓存**: 适合单机部署，高频访问的小数据
2. **Redis缓存**: 适合分布式部署，大数据量，需要持久化
3. **合理设置过期时间**: 避免内存泄漏
4. **使用键前缀**: 避免键冲突
5. **监控缓存命中率**: 优化缓存策略

## 错误处理

```go
// 检查特定错误类型
if errors.Is(err, cache.ErrCacheNotFound) {
    // 处理缓存未找到
}

// 检查缓存错误
var cacheErr *cache.CacheError
if errors.As(err, &cacheErr) {
    log.Printf("缓存操作失败: %s, 键: %s", cacheErr.Op, cacheErr.Key)
}
```

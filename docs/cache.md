# Cache 缓存模块

## 功能特性

提供统一的缓存接口，支持 **Memory** 和 **Redis** 两种实现：

- **Memory缓存**：基于内存的缓存，带自动过期清理机制
- **Redis缓存**：基于Redis的分布式缓存

## 快速开始

### 1. 配置文件

在 `config/config.yaml` 中配置缓存：

```yaml
cache:
  driver: memory  # memory / redis
  
  # Redis 配置（仅当 driver 为 redis 时需要）
  redis:
    host: localhost
    port: 6379
    password: ""
    db: 0
```

### 2. 使用示例

```go
package main

import (
    "time"
    "bico-admin/internal/core/cache"
    "bico-admin/internal/core/config"
)

func main() {
    // 加载配置
    cfg, _ := config.LoadConfig("config/config.yaml")
    
    // 创建缓存实例
    c, err := cache.NewCache(&cfg.Cache)
    if err != nil {
        panic(err)
    }
    defer c.Close()
    
    // 设置缓存（10分钟过期）
    c.Set("user:1", map[string]string{
        "name": "张三",
        "age":  "25",
    }, 10*time.Minute)
    
    // 获取缓存
    user, err := c.Get("user:1")
    if err != nil {
        // 处理缓存miss
    }
    
    // 检查是否存在
    if c.Exists("user:1") {
        // ...
    }
    
    // 删除缓存
    c.Delete("user:1")
    
    // 清空所有缓存
    c.Clear()
}
```

### 3. 直接使用Memory缓存

如果不需要配置，可以直接创建Memory缓存：

```go
memCache := cache.NewMemoryCache()
defer memCache.Close()

memCache.Set("key", "value", 5*time.Minute)
val, _ := memCache.Get("key")
```

### 4. 直接使用Redis缓存

```go
import "bico-admin/internal/core/config"

// 需要先定义配置
redisConfig := &config.RedisConfig{
    Host:     "localhost",
    Port:     6379,
    Password: "",
    DB:       0,
}

redisCache, err := cache.NewRedisCache(redisConfig)
if err != nil {
    panic(err)
}
defer redisCache.Close()

redisCache.Set("key", "value", 5*time.Minute)
```

## API 接口

```go
type Cache interface {
    // Set 设置缓存
    Set(key string, value interface{}, expiration time.Duration) error
    
    // Get 获取缓存
    Get(key string) (interface{}, error)
    
    // Delete 删除缓存
    Delete(key string) error
    
    // Exists 检查key是否存在
    Exists(key string) bool
    
    // Clear 清空所有缓存
    Clear() error
    
    // Close 关闭缓存连接
    Close() error
}
```

## 注意事项

1. **Memory缓存**自动启动后台协程定时清理过期数据（每分钟一次）
2. **Redis缓存**在创建时会自动测试连接，连接失败会返回错误
3. 缓存的值会通过JSON序列化存储，因此需要确保数据可序列化
4. `expiration` 为 0 表示永不过期（Memory缓存），Redis默认行为由Redis配置决定

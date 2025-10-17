package cache

import (
	"fmt"
)

// Config 缓存配置接口
type Config interface {
	GetDriver() string
	GetRedisConfig() RedisConfig
}

// RedisConfig Redis配置接口
type RedisConfig interface {
	GetHost() string
	GetPort() int
	GetPassword() string
	GetDB() int
}

// NewCache 创建缓存实例工厂方法
func NewCache(config Config) (Cache, error) {
	driver := config.GetDriver()
	
	switch driver {
	case "memory":
		return NewMemoryCache(), nil
	case "redis":
		redisConfig := config.GetRedisConfig()
		return NewRedisCache(redisConfig)
	default:
		return nil, fmt.Errorf("不支持的缓存驱动: %s", driver)
	}
}

package shared

import (
	"bico-admin/pkg/cache"
	"bico-admin/pkg/config"
)

// NewCacheFromConfig 从应用配置创建缓存实例
func NewCacheFromConfig(cfg *config.Config) (cache.Cache, error) {
	cacheConfig := cache.Config{
		Driver: cfg.Cache.Driver,
		Memory: cache.MemoryConfig{
			MaxSize:           cfg.Cache.Memory.MaxSize,
			DefaultExpiration: cfg.Cache.Memory.DefaultExpiration,
			CleanupInterval:   cfg.Cache.Memory.CleanupInterval,
		},
		Redis: cache.RedisConfig{
			Host:         cfg.Cache.Redis.Host,
			Port:         cfg.Cache.Redis.Port,
			Password:     cfg.Cache.Redis.Password,
			Database:     cfg.Cache.Redis.Database,
			PoolSize:     cfg.Cache.Redis.PoolSize,
			MinIdleConns: cfg.Cache.Redis.MinIdleConns,
			DialTimeout:  cfg.Cache.Redis.DialTimeout,
			ReadTimeout:  cfg.Cache.Redis.ReadTimeout,
			WriteTimeout: cfg.Cache.Redis.WriteTimeout,
			KeyPrefix:    cfg.Cache.Redis.KeyPrefix,
		},
	}

	return cache.NewCache(cacheConfig)
}

// NewCacheManagerFromConfig 从应用配置创建缓存管理器
func NewCacheManagerFromConfig(cfg *config.Config) (*cache.Manager, error) {
	cacheInstance, err := NewCacheFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cache.NewManager(cacheInstance), nil
}

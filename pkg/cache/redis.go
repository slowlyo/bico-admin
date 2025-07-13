package cache

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisCache Redis缓存实现
type redisCache struct {
	client    *redis.Client
	keyPrefix string
}

// NewRedisCache 创建Redis缓存
func NewRedisCache(config RedisConfig) (Cache, error) {
	// 设置默认值
	if config.Host == "" {
		config.Host = "localhost"
	}
	if config.Port == 0 {
		config.Port = 6379
	}
	if config.PoolSize == 0 {
		config.PoolSize = 10
	}
	if config.MinIdleConns == 0 {
		config.MinIdleConns = 5
	}
	if config.DialTimeout == 0 {
		config.DialTimeout = 5 * time.Second
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = 3 * time.Second
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 3 * time.Second
	}

	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.Database,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, &CacheError{
			Op:  "connect",
			Err: fmt.Errorf("Redis连接失败: %w", err),
		}
	}

	return &redisCache{
		client:    client,
		keyPrefix: config.KeyPrefix,
	}, nil
}

// buildKey 构建带前缀的键
func (c *redisCache) buildKey(key string) string {
	if c.keyPrefix == "" {
		return key
	}
	return c.keyPrefix + key
}

// stripPrefix 移除键前缀
func (c *redisCache) stripPrefix(key string) string {
	if c.keyPrefix == "" {
		return key
	}
	return strings.TrimPrefix(key, c.keyPrefix)
}

// Get 获取缓存值
func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := c.client.Get(ctx, c.buildKey(key)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrCacheNotFound
		}
		return "", &CacheError{
			Op:  "get",
			Key: key,
			Err: err,
		}
	}
	return result, nil
}

// Set 设置缓存值
func (c *redisCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := c.client.Set(ctx, c.buildKey(key), value, expiration).Err()
	if err != nil {
		return &CacheError{
			Op:  "set",
			Key: key,
			Err: err,
		}
	}
	return nil
}

// Delete 删除缓存
func (c *redisCache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, c.buildKey(key)).Err()
	if err != nil {
		return &CacheError{
			Op:  "delete",
			Key: key,
			Err: err,
		}
	}
	return nil
}

// Exists 检查键是否存在
func (c *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, c.buildKey(key)).Result()
	if err != nil {
		return false, &CacheError{
			Op:  "exists",
			Key: key,
			Err: err,
		}
	}
	return result > 0, nil
}

// Clear 清空所有缓存
func (c *redisCache) Clear(ctx context.Context) error {
	var err error
	if c.keyPrefix == "" {
		// 如果没有前缀，清空整个数据库
		err = c.client.FlushDB(ctx).Err()
	} else {
		// 如果有前缀，只删除匹配前缀的键
		keys, scanErr := c.client.Keys(ctx, c.keyPrefix+"*").Result()
		if scanErr != nil {
			return &CacheError{
				Op:  "clear",
				Err: scanErr,
			}
		}
		if len(keys) > 0 {
			err = c.client.Del(ctx, keys...).Err()
		}
	}

	if err != nil {
		return &CacheError{
			Op:  "clear",
			Err: err,
		}
	}
	return nil
}

// Keys 获取匹配模式的所有键
func (c *redisCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	// 构建带前缀的模式
	searchPattern := c.buildKey(pattern)
	
	keys, err := c.client.Keys(ctx, searchPattern).Result()
	if err != nil {
		return nil, &CacheError{
			Op:  "keys",
			Err: err,
		}
	}

	// 移除前缀
	result := make([]string, len(keys))
	for i, key := range keys {
		result[i] = c.stripPrefix(key)
	}

	return result, nil
}

// TTL 获取键的剩余生存时间
func (c *redisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	result, err := c.client.TTL(ctx, c.buildKey(key)).Result()
	if err != nil {
		return 0, &CacheError{
			Op:  "ttl",
			Key: key,
			Err: err,
		}
	}
	return result, nil
}

// Expire 设置键的过期时间
func (c *redisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	result, err := c.client.Expire(ctx, c.buildKey(key), expiration).Result()
	if err != nil {
		return &CacheError{
			Op:  "expire",
			Key: key,
			Err: err,
		}
	}

	if !result {
		return &CacheError{
			Op:  "expire",
			Key: key,
			Err: fmt.Errorf("key not found"),
		}
	}

	return nil
}

// Close 关闭缓存连接
func (c *redisCache) Close() error {
	return c.client.Close()
}

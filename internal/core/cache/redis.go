package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache Redis缓存实现
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache 创建 Redis 缓存实例
func NewRedisCache(cfg interface {
	GetHost() string
	GetPort() int
	GetPassword() string
	GetDB() int
}) (*RedisCache, error) {
	addr := fmt.Sprintf("%s:%d", cfg.GetHost(), cfg.GetPort())
	
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.GetPassword(),
		DB:       cfg.GetDB(),
	})
	
	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis 连接失败: %w", err)
	}
	
	return &RedisCache{client: client, ctx: ctx}, nil
}

// Set 设置缓存
func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	return r.client.Set(r.ctx, key, string(data), expiration).Err()
}

// Get 获取缓存
func (r *RedisCache) Get(key string) (interface{}, error) {
	result, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}
	
	var value interface{}
	if err := json.Unmarshal([]byte(result), &value); err != nil {
		return nil, err
	}
	
	return value, nil
}

// Delete 删除缓存
func (r *RedisCache) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Exists 检查key是否存在
func (r *RedisCache) Exists(key string) bool {
	count, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false
	}
	return count > 0
}

// Clear 清空所有缓存
func (r *RedisCache) Clear() error {
	return r.client.FlushDB(r.ctx).Err()
}

// Close 关闭Redis连接
func (r *RedisCache) Close() error {
	return r.client.Close()
}

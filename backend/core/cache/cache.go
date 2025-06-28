package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"bico-admin/core/config"
)

// Cache 缓存操作封装
type Cache struct {
	client *redis.Client
	ctx    context.Context
}

// NewCache 创建缓存实例
func NewCache() *Cache {
	return &Cache{
		client: config.GetRedis(),
		ctx:    context.Background(),
	}
}

// Set 设置缓存
func (c *Cache) Set(key string, value interface{}, expiration time.Duration) error {
	if c.client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	// 序列化值
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.client.Set(c.ctx, key, data, expiration).Err()
}

// Get 获取缓存
func (c *Cache) Get(key string, dest interface{}) error {
	if c.client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	data, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found")
		}
		return fmt.Errorf("failed to get value: %w", err)
	}

	// 反序列化值
	return json.Unmarshal([]byte(data), dest)
}

// GetString 获取字符串缓存
func (c *Cache) GetString(key string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("redis client not initialized")
	}

	return c.client.Get(c.ctx, key).Result()
}

// SetString 设置字符串缓存
func (c *Cache) SetString(key, value string, expiration time.Duration) error {
	if c.client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	return c.client.Set(c.ctx, key, value, expiration).Err()
}

// Delete 删除缓存
func (c *Cache) Delete(key string) error {
	if c.client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	return c.client.Del(c.ctx, key).Err()
}

// DeletePattern 删除匹配模式的缓存
func (c *Cache) DeletePattern(pattern string) error {
	if c.client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	keys, err := c.client.Keys(c.ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys: %w", err)
	}

	if len(keys) > 0 {
		return c.client.Del(c.ctx, keys...).Err()
	}

	return nil
}

// Exists 检查缓存是否存在
func (c *Cache) Exists(key string) (bool, error) {
	if c.client == nil {
		return false, fmt.Errorf("redis client not initialized")
	}

	count, err := c.client.Exists(c.ctx, key).Result()
	return count > 0, err
}

// Expire 设置缓存过期时间
func (c *Cache) Expire(key string, expiration time.Duration) error {
	if c.client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	return c.client.Expire(c.ctx, key, expiration).Err()
}

// TTL 获取缓存剩余时间
func (c *Cache) TTL(key string) (time.Duration, error) {
	if c.client == nil {
		return 0, fmt.Errorf("redis client not initialized")
	}

	return c.client.TTL(c.ctx, key).Result()
}

// Increment 递增计数器
func (c *Cache) Increment(key string) (int64, error) {
	if c.client == nil {
		return 0, fmt.Errorf("redis client not initialized")
	}

	return c.client.Incr(c.ctx, key).Result()
}

// IncrementBy 按指定值递增计数器
func (c *Cache) IncrementBy(key string, value int64) (int64, error) {
	if c.client == nil {
		return 0, fmt.Errorf("redis client not initialized")
	}

	return c.client.IncrBy(c.ctx, key, value).Result()
}

// SetNX 仅当key不存在时设置
func (c *Cache) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	if c.client == nil {
		return false, fmt.Errorf("redis client not initialized")
	}

	// 序列化值
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.client.SetNX(c.ctx, key, data, expiration).Result()
}

// SetEX 设置带过期时间的缓存
func (c *Cache) SetEX(key string, value interface{}, expiration time.Duration) error {
	return c.Set(key, value, expiration)
}

// MGet 批量获取缓存
func (c *Cache) MGet(keys ...string) ([]interface{}, error) {
	if c.client == nil {
		return nil, fmt.Errorf("redis client not initialized")
	}

	return c.client.MGet(c.ctx, keys...).Result()
}

// MSet 批量设置缓存
func (c *Cache) MSet(pairs ...interface{}) error {
	if c.client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	return c.client.MSet(c.ctx, pairs...).Err()
}

// FlushAll 清空所有缓存
func (c *Cache) FlushAll() error {
	if c.client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	return c.client.FlushAll(c.ctx).Err()
}

// Close 关闭连接
func (c *Cache) Close() error {
	if c.client == nil {
		return nil
	}

	return c.client.Close()
}

// 全局缓存实例
var globalCache *Cache

// GetCache 获取全局缓存实例
func GetCache() *Cache {
	if globalCache == nil {
		globalCache = NewCache()
	}
	return globalCache
}

// 常用的缓存键前缀
const (
	TokenBlacklistPrefix = "token_blacklist:"
	UserSessionPrefix    = "user_session:"
	LoginAttemptPrefix   = "login_attempt:"
	CaptchaPrefix        = "captcha:"
)

// Token黑名单相关方法
func (c *Cache) AddTokenToBlacklist(token string, expiration time.Duration) error {
	key := TokenBlacklistPrefix + token
	return c.SetString(key, "1", expiration)
}

func (c *Cache) IsTokenBlacklisted(token string) (bool, error) {
	key := TokenBlacklistPrefix + token
	return c.Exists(key)
}

// 用户会话相关方法
func (c *Cache) SetUserSession(userID uint, sessionData interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("%s%d", UserSessionPrefix, userID)
	return c.Set(key, sessionData, expiration)
}

func (c *Cache) GetUserSession(userID uint, dest interface{}) error {
	key := fmt.Sprintf("%s%d", UserSessionPrefix, userID)
	return c.Get(key, dest)
}

func (c *Cache) DeleteUserSession(userID uint) error {
	key := fmt.Sprintf("%s%d", UserSessionPrefix, userID)
	return c.Delete(key)
}

// 登录尝试次数相关方法
func (c *Cache) IncrementLoginAttempt(identifier string, expiration time.Duration) (int64, error) {
	key := LoginAttemptPrefix + identifier
	count, err := c.Increment(key)
	if err != nil {
		return 0, err
	}
	
	// 如果是第一次尝试，设置过期时间
	if count == 1 {
		err = c.Expire(key, expiration)
	}
	
	return count, err
}

func (c *Cache) GetLoginAttemptCount(identifier string) (int64, error) {
	key := LoginAttemptPrefix + identifier
	countStr, err := c.GetString(key)
	if err != nil {
		if err.Error() == "key not found" {
			return 0, nil
		}
		return 0, err
	}
	
	var count int64
	if err := json.Unmarshal([]byte(countStr), &count); err != nil {
		return 0, err
	}
	
	return count, nil
}

func (c *Cache) ResetLoginAttempt(identifier string) error {
	key := LoginAttemptPrefix + identifier
	return c.Delete(key)
}

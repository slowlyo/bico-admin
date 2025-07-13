package cache

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"
)

// memoryItem 内存缓存项
type memoryItem struct {
	value      string
	expiration int64 // 过期时间戳，0表示永不过期
}

// isExpired 检查是否过期
func (item *memoryItem) isExpired() bool {
	if item.expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.expiration
}

// memoryCache 内存缓存实现
type memoryCache struct {
	mu              sync.RWMutex
	items           map[string]*memoryItem
	maxSize         int
	defaultExp      time.Duration
	cleanupInterval time.Duration
	stopCleanup     chan bool
}

// NewMemoryCache 创建内存缓存
func NewMemoryCache(config MemoryConfig) (Cache, error) {
	// 设置默认值
	if config.MaxSize == 0 {
		config.MaxSize = 10000 // 默认最大10000个条目
	}
	if config.DefaultExpiration == 0 {
		config.DefaultExpiration = 30 * time.Minute // 默认30分钟过期
	}
	if config.CleanupInterval == 0 {
		config.CleanupInterval = 10 * time.Minute // 默认10分钟清理一次
	}

	cache := &memoryCache{
		items:           make(map[string]*memoryItem),
		maxSize:         config.MaxSize,
		defaultExp:      config.DefaultExpiration,
		cleanupInterval: config.CleanupInterval,
		stopCleanup:     make(chan bool),
	}

	// 启动清理协程
	go cache.startCleanup()

	return cache, nil
}

// Get 获取缓存值
func (c *memoryCache) Get(ctx context.Context, key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return "", ErrCacheNotFound
	}

	if item.isExpired() {
		// 延迟删除过期项
		go func() {
			c.mu.Lock()
			delete(c.items, key)
			c.mu.Unlock()
		}()
		return "", ErrCacheNotFound
	}

	return item.value, nil
}

// Set 设置缓存值
func (c *memoryCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查容量限制
	if c.maxSize > 0 && len(c.items) >= c.maxSize {
		// 如果键已存在，允许更新
		if _, exists := c.items[key]; !exists {
			return &CacheError{
				Op:  "set",
				Key: key,
				Err: fmt.Errorf("cache is full (max size: %d)", c.maxSize),
			}
		}
	}

	var exp int64
	if expiration > 0 {
		exp = time.Now().Add(expiration).UnixNano()
	} else if expiration == 0 && c.defaultExp > 0 {
		exp = time.Now().Add(c.defaultExp).UnixNano()
	}
	// expiration < 0 表示永不过期，exp保持为0

	c.items[key] = &memoryItem{
		value:      value,
		expiration: exp,
	}

	return nil
}

// Delete 删除缓存
func (c *memoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
	return nil
}

// Exists 检查键是否存在
func (c *memoryCache) Exists(ctx context.Context, key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return false, nil
	}

	if item.isExpired() {
		// 延迟删除过期项
		go func() {
			c.mu.Lock()
			delete(c.items, key)
			c.mu.Unlock()
		}()
		return false, nil
	}

	return true, nil
}

// Clear 清空所有缓存
func (c *memoryCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*memoryItem)
	return nil
}

// Keys 获取匹配模式的所有键
func (c *memoryCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var keys []string
	for key, item := range c.items {
		// 检查是否过期
		if item.isExpired() {
			continue
		}

		// 简单的通配符匹配
		matched, err := filepath.Match(pattern, key)
		if err != nil {
			return nil, &CacheError{
				Op:  "keys",
				Err: fmt.Errorf("invalid pattern %s: %w", pattern, err),
			}
		}

		if matched {
			keys = append(keys, key)
		}
	}

	return keys, nil
}

// TTL 获取键的剩余生存时间
func (c *memoryCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return -2 * time.Second, nil // Redis约定：-2表示键不存在
	}

	if item.expiration == 0 {
		return -1 * time.Second, nil // Redis约定：-1表示永不过期
	}

	if item.isExpired() {
		return -2 * time.Second, nil
	}

	remaining := time.Duration(item.expiration - time.Now().UnixNano())
	return remaining, nil
}

// Expire 设置键的过期时间
func (c *memoryCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.items[key]
	if !exists {
		return &CacheError{
			Op:  "expire",
			Key: key,
			Err: fmt.Errorf("key not found"),
		}
	}

	if expiration > 0 {
		item.expiration = time.Now().Add(expiration).UnixNano()
	} else {
		item.expiration = 0 // 永不过期
	}

	return nil
}

// Close 关闭缓存连接
func (c *memoryCache) Close() error {
	close(c.stopCleanup)
	return nil
}

// startCleanup 启动清理协程
func (c *memoryCache) startCleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopCleanup:
			return
		}
	}
}

// cleanup 清理过期项
func (c *memoryCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixNano()
	for key, item := range c.items {
		if item.expiration > 0 && item.expiration < now {
			delete(c.items, key)
		}
	}
}

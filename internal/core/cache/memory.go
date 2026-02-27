package cache

import (
	"errors"
	"sync"
	"time"
)

// memoryItem 内存缓存项
type memoryItem struct {
	value      interface{}
	expiration int64
}

// MemoryCache 内存缓存实现
type MemoryCache struct {
	items     map[string]memoryItem
	mu        sync.RWMutex
	stopCh    chan struct{}
	closeOnce sync.Once
}

// NewMemoryCache 创建内存缓存实例
func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		items:  make(map[string]memoryItem),
		stopCh: make(chan struct{}),
	}

	// 启动过期清理协程
	go cache.cleanupExpired()

	return cache
}

// Set 设置缓存
func (m *MemoryCache) Set(key string, value interface{}, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var exp int64
	if expiration > 0 {
		exp = time.Now().Add(expiration).UnixNano()
	}

	m.items[key] = memoryItem{
		value:      value,
		expiration: exp,
	}

	return nil
}

// Get 获取缓存
func (m *MemoryCache) Get(key string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.items[key]
	if !exists {
		return nil, errors.New("key not found")
	}

	// 检查是否过期
	if item.expiration > 0 && time.Now().UnixNano() > item.expiration {
		return nil, errors.New("key expired")
	}

	return item.value, nil
}

// Delete 删除缓存
func (m *MemoryCache) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, key)
	return nil
}

// Exists 检查key是否存在
func (m *MemoryCache) Exists(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.items[key]
	if !exists {
		return false
	}

	// 检查是否过期
	if item.expiration > 0 && time.Now().UnixNano() > item.expiration {
		return false
	}

	return true
}

// Clear 清空所有缓存
func (m *MemoryCache) Clear() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items = make(map[string]memoryItem)
	return nil
}

// Close 关闭缓存连接（内存缓存无需关闭）
func (m *MemoryCache) Close() error {
	m.closeOnce.Do(func() {
		close(m.stopCh)
	})
	return nil
}

// cleanupExpired 清理过期缓存
func (m *MemoryCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.mu.Lock()
			now := time.Now().UnixNano()
			for key, item := range m.items {
				// 仅删除已过期条目，避免影响未过期缓存。
				if item.expiration > 0 && now > item.expiration {
					delete(m.items, key)
				}
			}
			m.mu.Unlock()
		case <-m.stopCh:
			// 收到停止信号后退出清理协程，防止 goroutine 泄漏。
			return
		}
	}
}

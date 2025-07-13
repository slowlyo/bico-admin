package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// Manager 缓存管理器，提供更高级的缓存操作
type Manager struct {
	cache Cache
}

// NewManager 创建缓存管理器
func NewManager(cache Cache) *Manager {
	return &Manager{cache: cache}
}

// GetString 获取字符串值
func (m *Manager) GetString(ctx context.Context, key string) (string, error) {
	return m.cache.Get(ctx, key)
}

// SetString 设置字符串值
func (m *Manager) SetString(ctx context.Context, key string, value string, expiration time.Duration) error {
	return m.cache.Set(ctx, key, value, expiration)
}

// GetJSON 获取JSON对象
func (m *Manager) GetJSON(ctx context.Context, key string, dest interface{}) error {
	value, err := m.cache.Get(ctx, key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(value), dest); err != nil {
		return &CacheError{
			Op:  "get_json",
			Key: key,
			Err: fmt.Errorf("JSON反序列化失败: %w", err),
		}
	}

	return nil
}

// SetJSON 设置JSON对象
func (m *Manager) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return &CacheError{
			Op:  "set_json",
			Key: key,
			Err: fmt.Errorf("JSON序列化失败: %w", err),
		}
	}

	return m.cache.Set(ctx, key, string(data), expiration)
}

// GetOrSet 获取缓存，如果不存在则设置
func (m *Manager) GetOrSet(ctx context.Context, key string, setter func() (string, error), expiration time.Duration) (string, error) {
	// 先尝试获取
	value, err := m.cache.Get(ctx, key)
	if err == nil {
		return value, nil
	}

	// 如果不是未找到错误，直接返回
	if err != ErrCacheNotFound {
		return "", err
	}

	// 调用setter函数获取值
	value, err = setter()
	if err != nil {
		return "", err
	}

	// 设置缓存
	if setErr := m.cache.Set(ctx, key, value, expiration); setErr != nil {
		// 设置失败不影响返回值，只记录错误
		return value, setErr
	}

	return value, nil
}

// GetOrSetJSON 获取JSON缓存，如果不存在则设置
func (m *Manager) GetOrSetJSON(ctx context.Context, key string, dest interface{}, setter func() (interface{}, error), expiration time.Duration) error {
	// 先尝试获取
	err := m.GetJSON(ctx, key, dest)
	if err == nil {
		return nil
	}

	// 如果不是未找到错误，直接返回
	if err != ErrCacheNotFound {
		return err
	}

	// 调用setter函数获取值
	value, err := setter()
	if err != nil {
		return err
	}

	// 设置缓存
	if setErr := m.SetJSON(ctx, key, value, expiration); setErr != nil {
		// 设置失败不影响返回值，只记录错误
		// 但仍需要将值复制到dest
		if err := copyValue(value, dest); err != nil {
			return err
		}
		return setErr
	}

	// 将值复制到dest
	return copyValue(value, dest)
}

// Remember 记忆化缓存，如果缓存不存在则执行函数并缓存结果
func (m *Manager) Remember(ctx context.Context, key string, expiration time.Duration, fn func() (interface{}, error)) (interface{}, error) {
	// 先尝试获取缓存
	var result interface{}
	err := m.GetJSON(ctx, key, &result)
	if err == nil {
		return result, nil
	}

	// 如果不是未找到错误，直接返回
	if err != ErrCacheNotFound {
		return nil, err
	}

	// 执行函数获取值
	value, err := fn()
	if err != nil {
		return nil, err
	}

	// 设置缓存
	if setErr := m.SetJSON(ctx, key, value, expiration); setErr != nil {
		// 设置失败不影响返回值
		return value, setErr
	}

	return value, nil
}

// Delete 删除缓存
func (m *Manager) Delete(ctx context.Context, key string) error {
	return m.cache.Delete(ctx, key)
}

// Exists 检查键是否存在
func (m *Manager) Exists(ctx context.Context, key string) (bool, error) {
	return m.cache.Exists(ctx, key)
}

// Clear 清空所有缓存
func (m *Manager) Clear(ctx context.Context) error {
	return m.cache.Clear(ctx)
}

// Keys 获取匹配模式的所有键
func (m *Manager) Keys(ctx context.Context, pattern string) ([]string, error) {
	return m.cache.Keys(ctx, pattern)
}

// TTL 获取键的剩余生存时间
func (m *Manager) TTL(ctx context.Context, key string) (time.Duration, error) {
	return m.cache.TTL(ctx, key)
}

// Expire 设置键的过期时间
func (m *Manager) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return m.cache.Expire(ctx, key, expiration)
}

// Close 关闭缓存连接
func (m *Manager) Close() error {
	return m.cache.Close()
}

// GetCache 获取底层缓存实例
func (m *Manager) GetCache() Cache {
	return m.cache
}

// copyValue 复制值到目标变量
func copyValue(src, dest interface{}) error {
	srcValue := reflect.ValueOf(src)
	destValue := reflect.ValueOf(dest)

	if destValue.Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer")
	}

	destElem := destValue.Elem()
	if !destElem.CanSet() {
		return fmt.Errorf("dest cannot be set")
	}

	// 如果类型相同，直接设置
	if srcValue.Type().AssignableTo(destElem.Type()) {
		destElem.Set(srcValue)
		return nil
	}

	// 尝试JSON转换
	data, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("JSON反序列化失败: %w", err)
	}

	return nil
}

// Stats 缓存统计信息
type Stats struct {
	Driver    string            `json:"driver"`
	KeyCount  int               `json:"key_count"`
	Memory    string            `json:"memory,omitempty"`
	Info      map[string]string `json:"info,omitempty"`
}

// GetStats 获取缓存统计信息
func (m *Manager) GetStats(ctx context.Context) (*Stats, error) {
	stats := &Stats{}

	// 获取所有键的数量
	keys, err := m.cache.Keys(ctx, "*")
	if err != nil {
		return nil, err
	}
	stats.KeyCount = len(keys)

	// 根据缓存类型获取不同的统计信息
	switch cache := m.cache.(type) {
	case *redisCache:
		stats.Driver = "redis"
		// 可以添加Redis特定的统计信息
		if info, err := cache.client.Info(ctx).Result(); err == nil {
			stats.Info = map[string]string{"redis_info": info}
		}
	case *memoryCache:
		stats.Driver = "memory"
		// 可以添加内存缓存特定的统计信息
		stats.Memory = fmt.Sprintf("%d items", len(cache.items))
	default:
		stats.Driver = "unknown"
	}

	return stats, nil
}

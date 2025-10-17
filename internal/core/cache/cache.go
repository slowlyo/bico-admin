package cache

import "time"

// Cache 缓存接口
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

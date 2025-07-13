package cache

import (
	"context"
	"errors"
	"time"
)

// Cache 缓存接口
type Cache interface {
	// Get 获取缓存值
	Get(ctx context.Context, key string) (string, error)

	// Set 设置缓存值
	Set(ctx context.Context, key string, value string, expiration time.Duration) error

	// Delete 删除缓存
	Delete(ctx context.Context, key string) error

	// Exists 检查键是否存在
	Exists(ctx context.Context, key string) (bool, error)

	// Clear 清空所有缓存
	Clear(ctx context.Context) error

	// Keys 获取匹配模式的所有键
	Keys(ctx context.Context, pattern string) ([]string, error)

	// TTL 获取键的剩余生存时间
	TTL(ctx context.Context, key string) (time.Duration, error)

	// Expire 设置键的过期时间
	Expire(ctx context.Context, key string, expiration time.Duration) error

	// Close 关闭缓存连接
	Close() error
}

// Config 缓存配置
type Config struct {
	// Driver 缓存驱动类型: "memory" 或 "redis"，默认为 "memory"
	Driver string `mapstructure:"driver" yaml:"driver"`

	// Memory 内存缓存配置
	Memory MemoryConfig `mapstructure:"memory" yaml:"memory"`

	// Redis Redis缓存配置
	Redis RedisConfig `mapstructure:"redis" yaml:"redis"`
}

// MemoryConfig 内存缓存配置
type MemoryConfig struct {
	// MaxSize 最大缓存条目数，0表示无限制，默认为10000
	MaxSize int `mapstructure:"max_size" yaml:"max_size"`

	// DefaultExpiration 默认过期时间，默认为30分钟
	DefaultExpiration time.Duration `mapstructure:"default_expiration" yaml:"default_expiration"`

	// CleanupInterval 清理过期条目的间隔时间，默认为10分钟
	CleanupInterval time.Duration `mapstructure:"cleanup_interval" yaml:"cleanup_interval"`
}

// RedisConfig Redis缓存配置
type RedisConfig struct {
	// Host Redis服务器地址，默认为localhost
	Host string `mapstructure:"host" yaml:"host"`

	// Port Redis服务器端口，默认为6379
	Port int `mapstructure:"port" yaml:"port"`

	// Password Redis密码，默认为空
	Password string `mapstructure:"password" yaml:"password"`

	// Database Redis数据库编号，默认为0
	Database int `mapstructure:"database" yaml:"database"`

	// PoolSize 连接池大小，默认为10
	PoolSize int `mapstructure:"pool_size" yaml:"pool_size"`

	// MinIdleConns 最小空闲连接数，默认为5
	MinIdleConns int `mapstructure:"min_idle_conns" yaml:"min_idle_conns"`

	// DialTimeout 连接超时时间，默认为5秒
	DialTimeout time.Duration `mapstructure:"dial_timeout" yaml:"dial_timeout"`

	// ReadTimeout 读取超时时间，默认为3秒
	ReadTimeout time.Duration `mapstructure:"read_timeout" yaml:"read_timeout"`

	// WriteTimeout 写入超时时间，默认为3秒
	WriteTimeout time.Duration `mapstructure:"write_timeout" yaml:"write_timeout"`

	// KeyPrefix 键前缀，默认为空
	KeyPrefix string `mapstructure:"key_prefix" yaml:"key_prefix"`
}

// CacheError 缓存错误类型
type CacheError struct {
	Op  string // 操作名称
	Key string // 缓存键
	Err error  // 原始错误
}

func (e *CacheError) Error() string {
	if e.Key != "" {
		return "cache " + e.Op + " " + e.Key + ": " + e.Err.Error()
	}
	return "cache " + e.Op + ": " + e.Err.Error()
}

func (e *CacheError) Unwrap() error {
	return e.Err
}

// ErrCacheNotFound 缓存未找到错误
var ErrCacheNotFound = &CacheError{Op: "get", Err: errors.New("cache not found")}

// NewCache 创建缓存实例
func NewCache(config Config) (Cache, error) {
	switch config.Driver {
	case "redis":
		return NewRedisCache(config.Redis)
	case "memory", "":
		return NewMemoryCache(config.Memory)
	default:
		return nil, &CacheError{
			Op:  "new",
			Err: errors.New("unsupported cache driver: " + config.Driver),
		}
	}
}

// NewCacheFromAppConfig 从应用配置创建缓存实例
func NewCacheFromAppConfig(appConfig interface{}) (Cache, error) {
	// 这里使用interface{}是为了避免循环导入
	// 实际使用时会传入 *config.Config 类型

	// 通过反射或类型断言获取配置
	// 这里简化处理，实际项目中可以在shared包中提供转换函数
	config := Config{
		Driver: "memory", // 默认使用内存缓存
		Memory: MemoryConfig{
			MaxSize:           10000,
			DefaultExpiration: 30 * time.Minute,
			CleanupInterval:   10 * time.Minute,
		},
	}

	return NewCache(config)
}

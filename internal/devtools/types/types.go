package types

import "time"

// ConfigInfo 配置信息结构
type ConfigInfo struct {
	App      AppInfo      `json:"app"`
	Server   ServerInfo   `json:"server"`
	Database DatabaseInfo `json:"database"`
	Redis    RedisInfo    `json:"redis"`
	Log      LogInfo      `json:"log"`
	JWT      JWTInfo      `json:"jwt"`
	Cache    CacheInfo    `json:"cache"`
}

// AppInfo 应用配置信息
type AppInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	Debug       bool   `json:"debug"`
}

// ServerInfo 服务器配置信息
type ServerInfo struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}

// DatabaseInfo 数据库配置信息
type DatabaseInfo struct {
	Driver          string        `json:"driver"`
	Host            string        `json:"host,omitempty"`
	Port            int           `json:"port,omitempty"`
	Username        string        `json:"username,omitempty"`
	Password        string        `json:"password,omitempty"`
	Database        string        `json:"database"`
	Charset         string        `json:"charset,omitempty"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	MaxOpenConns    int           `json:"max_open_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
}

// RedisInfo Redis配置信息
type RedisInfo struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	Password     string        `json:"password,omitempty"`
	Database     int           `json:"database"`
	PoolSize     int           `json:"pool_size"`
	MinIdleConns int           `json:"min_idle_conns"`
	DialTimeout  time.Duration `json:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

// LogInfo 日志配置信息
type LogInfo struct {
	Level      string `json:"level"`
	Format     string `json:"format"`
	Output     string `json:"output"`
	Filename   string `json:"filename,omitempty"`
	MaxSize    int    `json:"max_size"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
	Compress   bool   `json:"compress"`
}

// JWTInfo JWT配置信息
type JWTInfo struct {
	Secret     string        `json:"secret"`
	Issuer     string        `json:"issuer"`
	ExpireTime time.Duration `json:"expire_time"`
}

// CacheInfo 缓存配置信息
type CacheInfo struct {
	Driver string                 `json:"driver"`
	Memory *MemoryCacheInfo       `json:"memory,omitempty"`
	Redis  *RedisCacheInfo        `json:"redis,omitempty"`
}

// MemoryCacheInfo 内存缓存配置
type MemoryCacheInfo struct {
	MaxSize           int           `json:"max_size"`
	DefaultExpiration time.Duration `json:"default_expiration"`
	CleanupInterval   time.Duration `json:"cleanup_interval"`
}

// RedisCacheInfo Redis缓存配置
type RedisCacheInfo struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	Password     string        `json:"password,omitempty"`
	Database     int           `json:"database"`
	PoolSize     int           `json:"pool_size"`
	MinIdleConns int           `json:"min_idle_conns"`
	DialTimeout  time.Duration `json:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	KeyPrefix    string        `json:"key_prefix"`
}

// ToolResponse MCP工具响应的通用结构
type ToolResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

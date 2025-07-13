package config

import (
	"fmt"
	"strings"
	"time"
)

// Config 应用配置结构
type Config struct {
	App      AppConfig      `mapstructure:"app" yaml:"app"`
	Server   ServerConfig   `mapstructure:"server" yaml:"server"`
	Database DatabaseConfig `mapstructure:"database" yaml:"database"`
	Redis    RedisConfig    `mapstructure:"redis" yaml:"redis"`
	Log      LogConfig      `mapstructure:"log" yaml:"log"`
	JWT      JWTConfig      `mapstructure:"jwt" yaml:"jwt"`
	Cache    CacheConfig    `mapstructure:"cache" yaml:"cache"`
}

// AppConfig 应用基础配置
type AppConfig struct {
	Name        string `mapstructure:"name" yaml:"name"`
	Version     string `mapstructure:"version" yaml:"version"`
	Environment string `mapstructure:"environment" yaml:"environment"`
	Debug       bool   `mapstructure:"debug" yaml:"debug"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string        `mapstructure:"host" yaml:"host"`
	Port         int           `mapstructure:"port" yaml:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" yaml:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout" yaml:"idle_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver" yaml:"driver"`
	Host            string        `mapstructure:"host" yaml:"host"`
	Port            int           `mapstructure:"port" yaml:"port"`
	Username        string        `mapstructure:"username" yaml:"username"`
	Password        string        `mapstructure:"password" yaml:"password"`
	Database        string        `mapstructure:"database" yaml:"database"`
	Charset         string        `mapstructure:"charset" yaml:"charset"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string        `mapstructure:"host" yaml:"host"`
	Port         int           `mapstructure:"port" yaml:"port"`
	Password     string        `mapstructure:"password" yaml:"password"`
	Database     int           `mapstructure:"database" yaml:"database"`
	PoolSize     int           `mapstructure:"pool_size" yaml:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns" yaml:"min_idle_conns"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" yaml:"write_timeout"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level" yaml:"level"`
	Format     string `mapstructure:"format" yaml:"format"`
	Output     string `mapstructure:"output" yaml:"output"`
	Filename   string `mapstructure:"filename" yaml:"filename"`
	MaxSize    int    `mapstructure:"max_size" yaml:"max_size"`
	MaxAge     int    `mapstructure:"max_age" yaml:"max_age"`
	MaxBackups int    `mapstructure:"max_backups" yaml:"max_backups"`
	Compress   bool   `mapstructure:"compress" yaml:"compress"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string        `mapstructure:"secret" yaml:"secret"`
	Issuer     string        `mapstructure:"issuer" yaml:"issuer"`
	ExpireTime time.Duration `mapstructure:"expire_time" yaml:"expire_time"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Driver string            `mapstructure:"driver" yaml:"driver"`
	Memory CacheMemoryConfig `mapstructure:"memory" yaml:"memory"`
	Redis  CacheRedisConfig  `mapstructure:"redis" yaml:"redis"`
}

// CacheMemoryConfig 内存缓存配置
type CacheMemoryConfig struct {
	MaxSize           int           `mapstructure:"max_size" yaml:"max_size"`
	DefaultExpiration time.Duration `mapstructure:"default_expiration" yaml:"default_expiration"`
	CleanupInterval   time.Duration `mapstructure:"cleanup_interval" yaml:"cleanup_interval"`
}

// CacheRedisConfig Redis缓存配置
type CacheRedisConfig struct {
	Host         string        `mapstructure:"host" yaml:"host"`
	Port         int           `mapstructure:"port" yaml:"port"`
	Password     string        `mapstructure:"password" yaml:"password"`
	Database     int           `mapstructure:"database" yaml:"database"`
	PoolSize     int           `mapstructure:"pool_size" yaml:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns" yaml:"min_idle_conns"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" yaml:"write_timeout"`
	KeyPrefix    string        `mapstructure:"key_prefix" yaml:"key_prefix"`
}

// GetDSN 获取数据库连接字符串
func (d *DatabaseConfig) GetDSN() string {
	switch strings.ToLower(d.Driver) {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			d.Username, d.Password, d.Host, d.Port, d.Database, d.Charset)
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
			d.Host, d.Port, d.Username, d.Password, d.Database)
	case "sqlite":
		return d.Database // SQLite 只需要文件路径
	default:
		return ""
	}
}

// GetRedisAddr 获取Redis地址
func (r *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// IsProduction 判断是否为生产环境
func (a *AppConfig) IsProduction() bool {
	return strings.ToLower(a.Environment) == "production"
}

// IsDevelopment 判断是否为开发环境
func (a *AppConfig) IsDevelopment() bool {
	return strings.ToLower(a.Environment) == "development"
}

package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config 配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
	Cache    CacheConfig    `mapstructure:"cache"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Upload   UploadConfig   `mapstructure:"upload"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port        int    `mapstructure:"port"`
	Mode        string `mapstructure:"mode"`
	EmbedStatic bool   `mapstructure:"embed_static"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name string `mapstructure:"name"`
	Logo string `mapstructure:"logo"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver       string           `mapstructure:"driver"`
	SQLite       SQLiteConfig     `mapstructure:"sqlite"`
	MySQL        MySQLConfig      `mapstructure:"mysql"`
	MaxIdleConns int              `mapstructure:"max_idle_conns"`
	MaxOpenConns int              `mapstructure:"max_open_conns"`
}

// SQLiteConfig SQLite 配置
type SQLiteConfig struct {
	Path string `mapstructure:"path"`
}

// MySQLConfig MySQL 配置
type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Charset  string `mapstructure:"charset"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Driver string      `mapstructure:"driver"` // memory / redis
	Redis  RedisConfig `mapstructure:"redis"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	Driver       string         `mapstructure:"driver"` // local / qiniu / aliyun / tencent
	MaxSize      int64          `mapstructure:"max_size"`
	AllowedTypes []string       `mapstructure:"allowed_types"`
	Local        LocalUploadConfig `mapstructure:"local"`
}

// LocalUploadConfig 本地存储配置
type LocalUploadConfig struct {
	BasePath  string `mapstructure:"base_path"`
	ServePath string `mapstructure:"serve_path"`
	URLPrefix string `mapstructure:"url_prefix"`
}

// GetDriver 获取缓存驱动
func (c *CacheConfig) GetDriver() string {
	return c.Driver
}

// GetRedisConfig 获取Redis配置
func (c *CacheConfig) GetRedisConfig() interface {
	GetHost() string
	GetPort() int
	GetPassword() string
	GetDB() int
} {
	return &c.Redis
}

// GetHost 获取Redis主机
func (r *RedisConfig) GetHost() string {
	return r.Host
}

// GetPort 获取Redis端口
func (r *RedisConfig) GetPort() int {
	return r.Port
}

// GetPassword 获取Redis密码
func (r *RedisConfig) GetPassword() string {
	return r.Password
}

// GetDB 获取Redis数据库
func (r *RedisConfig) GetDB() int {
	return r.DB
}

// LoadConfig 加载配置文件
// 支持多路径自动查找，优先级：
// 1. 指定的路径（如果存在）
// 2. ./config.yaml（项目根目录，Docker 友好）
// 3. ./config/config.yaml（传统位置）
func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}
	
	// 查找配置文件
	actualPath, err := findConfigFile(configPath)
	if err != nil {
		return nil, err
	}
	
	v := viper.New()
	v.SetConfigFile(actualPath)
	
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	
	return cfg, nil
}

// findConfigFile 查找配置文件
func findConfigFile(configPath string) (string, error) {
	// 如果指定了路径且文件存在，直接使用
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}
	}
	
	// 尝试默认路径列表
	defaultPaths := []string{
		"config.yaml",         // 项目根目录（Docker 友好）
		"config/config.yaml",  // 传统位置
	}
	
	for _, path := range defaultPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	
	// 所有路径都找不到，返回详细错误
	return "", fmt.Errorf("配置文件未找到，已尝试的路径: %s, %v", configPath, defaultPaths)
}

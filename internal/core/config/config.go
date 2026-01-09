package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Config 配置结构体
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	App       AppConfig       `mapstructure:"app"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Log       LogConfig       `mapstructure:"log"`
	Cache     CacheConfig     `mapstructure:"cache"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Upload    UploadConfig    `mapstructure:"upload"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port        int    `mapstructure:"port"`
	Mode        string `mapstructure:"mode"`
	EmbedStatic bool   `mapstructure:"embed_static"`
	AdminPath   string `mapstructure:"admin_path"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name string `mapstructure:"name"`
	Logo string `mapstructure:"logo"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver       string         `mapstructure:"driver"`
	SQLite       SQLiteConfig   `mapstructure:"sqlite"`
	MySQL        MySQLConfig    `mapstructure:"mysql"`
	Postgres     PostgresConfig `mapstructure:"postgres"`
	MaxIdleConns int            `mapstructure:"max_idle_conns"`
	MaxOpenConns int            `mapstructure:"max_open_conns"`
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

// PostgresConfig PostgreSQL 配置
type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"sslmode"`
	TimeZone string `mapstructure:"timezone"`
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
	Driver       string             `mapstructure:"driver"` // local / qiniu / aliyun
	MaxSize      int64              `mapstructure:"max_size"`
	AllowedTypes []string           `mapstructure:"allowed_types"`
	Local        LocalUploadConfig  `mapstructure:"local"`
	Qiniu        QiniuUploadConfig  `mapstructure:"qiniu"`
	Aliyun       AliyunUploadConfig `mapstructure:"aliyun"`
}

// LocalUploadConfig 本地存储配置
type LocalUploadConfig struct {
	BasePath  string `mapstructure:"base_path"`
	ServePath string `mapstructure:"serve_path"`
	URLPrefix string `mapstructure:"url_prefix"`
}

// QiniuUploadConfig 七牛云存储配置
type QiniuUploadConfig struct {
	AccessKey    string `mapstructure:"access_key"`
	SecretKey    string `mapstructure:"secret_key"`
	Bucket       string `mapstructure:"bucket"`
	Domain       string `mapstructure:"domain"`
	Zone         string `mapstructure:"zone"`
	UseHTTPS     bool   `mapstructure:"use_https"`
	UseCDNDomain bool   `mapstructure:"use_cdn_domain"`
}

// AliyunUploadConfig 阿里云OSS存储配置
type AliyunUploadConfig struct {
	AccessKeyId     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	Bucket          string `mapstructure:"bucket"`
	Endpoint        string `mapstructure:"endpoint"`
	Domain          string `mapstructure:"domain"`
	UseHTTPS        bool   `mapstructure:"use_https"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled bool `mapstructure:"enabled"` // 是否启用限流
	RPS     int  `mapstructure:"rps"`     // 每秒请求数
	Burst   int  `mapstructure:"burst"`   // 突发流量桶容量
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

// ConfigManager 配置管理器（支持热更新）
type ConfigManager struct {
	config *Config
	mu     sync.RWMutex
	viper  *viper.Viper
	logger *zap.Logger
}

// NewConfigManager 创建配置管理器
func NewConfigManager(configPath string, logger *zap.Logger) (*ConfigManager, error) {
	actualPath, err := findConfigFile(configPath)
	if err != nil {
		return nil, err
	}

	v := viper.New()
	v.SetConfigFile(actualPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	cm := &ConfigManager{
		config: cfg,
		viper:  v,
		logger: logger,
	}

	// 监听配置文件变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		cm.onConfigChange(e)
	})

	return cm, nil
}

// onConfigChange 配置变更回调
func (cm *ConfigManager) onConfigChange(e fsnotify.Event) {
	cm.logger.Info("检测到配置文件变化", zap.String("file", e.Name))

	newConfig := &Config{}
	if err := cm.viper.Unmarshal(newConfig); err != nil {
		cm.logger.Error("重新加载配置失败", zap.Error(err))
		return
	}

	cm.mu.Lock()
	cm.config = newConfig
	cm.mu.Unlock()

	cm.logger.Info("配置已热更新")
}

// GetConfig 获取当前配置（线程安全）
func (cm *ConfigManager) GetConfig() *Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

// GetRateLimitConfig 获取限流配置
func (cm *ConfigManager) GetRateLimitConfig() RateLimitConfig {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.RateLimit
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
		"config.yaml",        // 项目根目录（Docker 友好）
		"config/config.yaml", // 传统位置
	}

	for _, path := range defaultPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// 所有路径都找不到，返回详细错误
	return "", fmt.Errorf("配置文件未找到，已尝试的路径: %s, %v", configPath, defaultPaths)
}

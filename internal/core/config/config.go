package config

import (
	"fmt"
	"os"
	"strings"
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
	Port           int      `mapstructure:"port"`
	Mode           string   `mapstructure:"mode"`
	EmbedStatic    bool     `mapstructure:"embed_static"`
	AdminPath      string   `mapstructure:"admin_path"`
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name string `mapstructure:"name"`
	Logo string `mapstructure:"logo"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver       string         `mapstructure:"driver"`
	AutoMigrate  bool           `mapstructure:"auto_migrate"`
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
	config   *Config
	mu       sync.RWMutex
	viper    *viper.Viper
	logger   *zap.Logger
	onChange func(*Config)
}

// NewConfigManager 创建配置管理器
func NewConfigManager(configPath string, logger *zap.Logger) (*ConfigManager, error) {
	actualPath, err := findConfigFile(configPath)
	if err != nil {
		return nil, err
	}

	v := newViper(actualPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	cfg, err := decodeConfig(v)
	if err != nil {
		return nil, err
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

	newConfig, err := decodeConfig(cm.viper)
	if err != nil {
		cm.logger.Error("重新加载配置失败", zap.Error(err))
		return
	}

	cm.mu.Lock()
	cm.config = newConfig
	hook := cm.onChange
	cm.mu.Unlock()

	// 先更新内存配置，再通知外部。通知里可以读到新配置。
	if hook != nil {
		hook(newConfig)
	}

	cm.logger.Info("配置已热更新")
}

// SetOnChange 注册配置重载后的回调。
//
// 说明：只保留最后一个回调。用于把日志级别这类不能靠再次读取配置生效的项接进来。
func (cm *ConfigManager) SetOnChange(fn func(*Config)) {
	cm.mu.Lock()
	cm.onChange = fn
	cm.mu.Unlock()
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
	// 查找配置文件
	actualPath, err := findConfigFile(configPath)
	if err != nil {
		return nil, err
	}

	v := newViper(actualPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	return decodeConfig(v)
}

// newViper 创建支持环境变量覆盖的配置读取器。
// 敏感配置使用 BICO_ 前缀，避免生产密钥写入配置文件。
func newViper(configPath string) *viper.Viper {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetEnvPrefix("BICO")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	_ = v.BindEnv("jwt.secret")
	_ = v.BindEnv("database.mysql.password")
	_ = v.BindEnv("database.postgres.password")
	_ = v.BindEnv("cache.redis.password")
	_ = v.BindEnv("upload.qiniu.access_key")
	_ = v.BindEnv("upload.qiniu.secret_key")
	_ = v.BindEnv("upload.aliyun.access_key_id")
	_ = v.BindEnv("upload.aliyun.access_key_secret")
	return v
}

// decodeConfig 解析并校验配置。
// 热更新和首次加载共用此路径，防止非法配置进入运行时。
func decodeConfig(v *viper.Viper) (*Config, error) {
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate 校验启动所需的安全配置。
// release 模式必须使用独立强密钥，避免示例配置被直接用于生产。
func (c *Config) Validate() error {
	if c.Server.Mode != "debug" && c.Server.Mode != "release" && c.Server.Mode != "test" {
		return fmt.Errorf("server.mode 仅支持 debug、release 或 test")
	}
	if c.JWT.ExpireHours <= 0 {
		return fmt.Errorf("jwt.expire_hours 必须大于 0")
	}
	if c.Server.Mode == "release" {
		secret := strings.TrimSpace(c.JWT.Secret)
		if len(secret) < 32 || secret == "bico-admin-secret-key-change-in-production" {
			return fmt.Errorf("release 模式必须通过 BICO_JWT_SECRET 配置至少 32 位的 JWT 密钥")
		}
		for _, origin := range c.Server.AllowedOrigins {
			if origin == "*" {
				return fmt.Errorf("release 模式不允许 CORS 通配来源")
			}
		}
	}
	return nil
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

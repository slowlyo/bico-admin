package config

import "github.com/spf13/viper"

// Config 配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
	Cache    CacheConfig    `mapstructure:"cache"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
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

// GetDriver 获取缓存驱动
func (c *CacheConfig) GetDriver() string {
	return c.Driver
}

// GetRedisConfig 获取Redis配置
func (c *CacheConfig) GetRedisConfig() *RedisConfig {
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
func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}
	
	v := viper.New()
	v.SetConfigFile(configPath)
	
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}
	
	return cfg, nil
}

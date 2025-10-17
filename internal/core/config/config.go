package config

import "github.com/spf13/viper"

// Config 配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
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

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var globalConfig *Config

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// 设置配置文件路径
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// 默认配置文件路径
		v.SetConfigName("app")
		v.SetConfigType("yaml")
		v.AddConfigPath("./config")
		v.AddConfigPath("../config")
		v.AddConfigPath("../../config")
	}

	// 环境变量配置
	v.SetEnvPrefix("BICO")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 尝试读取环境特定的配置文件
	env := v.GetString("app.environment")
	if env != "" {
		envConfigFile := fmt.Sprintf("app.%s.yaml", env)
		envConfigPath := filepath.Join(filepath.Dir(v.ConfigFileUsed()), envConfigFile)

		if _, err := os.Stat(envConfigPath); err == nil {
			envViper := viper.New()
			envViper.SetConfigFile(envConfigPath)
			if err := envViper.ReadInConfig(); err == nil {
				// 合并环境特定配置
				if err := v.MergeConfigMap(envViper.AllSettings()); err != nil {
					return nil, fmt.Errorf("合并环境配置失败: %w", err)
				}
			}
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	globalConfig = &config
	return &config, nil
}

// Get 获取全局配置
func Get() *Config {
	return globalConfig
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	if config.App.Name == "" {
		return fmt.Errorf("应用名称不能为空")
	}

	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("服务器端口必须在1-65535之间")
	}

	if config.Database.Driver == "" {
		return fmt.Errorf("数据库驱动不能为空")
	}

	// 根据数据库类型验证不同的配置
	switch strings.ToLower(config.Database.Driver) {
	case "mysql", "postgres":
		if config.Database.Host == "" {
			return fmt.Errorf("数据库主机不能为空")
		}
		if config.Database.Username == "" {
			return fmt.Errorf("数据库用户名不能为空")
		}
	case "sqlite":
		// SQLite 只需要验证数据库文件路径
		if config.Database.Database == "" {
			return fmt.Errorf("SQLite数据库文件路径不能为空")
		}
	default:
		return fmt.Errorf("不支持的数据库驱动: %s", config.Database.Driver)
	}

	if config.Database.Database == "" {
		return fmt.Errorf("数据库名称不能为空")
	}

	if config.JWT.Secret == "" {
		return fmt.Errorf("JWT密钥不能为空")
	}

	return nil
}

// LoadFromEnv 从环境变量加载配置
func LoadFromEnv() (*Config, error) {
	v := viper.New()

	// 设置环境变量前缀
	v.SetEnvPrefix("BICO")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 设置默认值
	setDefaults(v)

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("从环境变量解析配置失败: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	globalConfig = &config
	return &config, nil
}

// setDefaults 设置默认配置值
func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "bico-admin")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.debug", true)

	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", "30s")
	v.SetDefault("server.write_timeout", "30s")
	v.SetDefault("server.idle_timeout", "60s")

	// Database defaults
	v.SetDefault("database.driver", "sqlite")
	v.SetDefault("database.database", "data/bico_admin.db") // SQLite 默认文件名
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.charset", "utf8mb4")
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.max_open_conns", 10)
	v.SetDefault("database.conn_max_lifetime", "1h")

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.database", 0)
	v.SetDefault("redis.pool_size", 10)
	v.SetDefault("redis.min_idle_conns", 5)
	v.SetDefault("redis.dial_timeout", "5s")
	v.SetDefault("redis.read_timeout", "3s")
	v.SetDefault("redis.write_timeout", "3s")

	// Log defaults
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.output", "stdout")
	v.SetDefault("log.max_size", 100)
	v.SetDefault("log.max_age", 7)
	v.SetDefault("log.max_backups", 10)
	v.SetDefault("log.compress", true)

	// JWT defaults
	v.SetDefault("jwt.issuer", "bico-admin")
	v.SetDefault("jwt.expire_time", "24h")
}

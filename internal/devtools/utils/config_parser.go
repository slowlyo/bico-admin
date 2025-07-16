package utils

import (
	"encoding/json"
	"fmt"

	"bico-admin/internal/devtools/types"
	"bico-admin/pkg/config"
)

// ConfigParser 配置解析器
type ConfigParser struct{}

// NewConfigParser 创建配置解析器
func NewConfigParser() *ConfigParser {
	return &ConfigParser{}
}

// ParseConfig 解析配置为MCP友好的格式
func (p *ConfigParser) ParseConfig(cfg *config.Config) (*types.ConfigInfo, error) {
	if cfg == nil {
		return nil, fmt.Errorf("配置对象为空")
	}

	configInfo := &types.ConfigInfo{
		App: types.AppInfo{
			Name:        cfg.App.Name,
			Version:     cfg.App.Version,
			Environment: cfg.App.Environment,
			Debug:       cfg.App.Debug,
		},
		Server: types.ServerInfo{
			Host:         cfg.Server.Host,
			Port:         cfg.Server.Port,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
		Database: types.DatabaseInfo{
			Driver:          cfg.Database.Driver,
			Host:            cfg.Database.Host,
			Port:            cfg.Database.Port,
			Username:        cfg.Database.Username,
			Password:        maskPassword(cfg.Database.Password),
			Database:        cfg.Database.Database,
			Charset:         cfg.Database.Charset,
			MaxIdleConns:    cfg.Database.MaxIdleConns,
			MaxOpenConns:    cfg.Database.MaxOpenConns,
			ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
		},
		Redis: types.RedisInfo{
			Host:         cfg.Redis.Host,
			Port:         cfg.Redis.Port,
			Password:     maskPassword(cfg.Redis.Password),
			Database:     cfg.Redis.Database,
			PoolSize:     cfg.Redis.PoolSize,
			MinIdleConns: cfg.Redis.MinIdleConns,
			DialTimeout:  cfg.Redis.DialTimeout,
			ReadTimeout:  cfg.Redis.ReadTimeout,
			WriteTimeout: cfg.Redis.WriteTimeout,
		},
		Log: types.LogInfo{
			Level:      cfg.Log.Level,
			Format:     cfg.Log.Format,
			Output:     cfg.Log.Output,
			Filename:   cfg.Log.Filename,
			MaxSize:    cfg.Log.MaxSize,
			MaxAge:     cfg.Log.MaxAge,
			MaxBackups: cfg.Log.MaxBackups,
			Compress:   cfg.Log.Compress,
		},
		JWT: types.JWTInfo{
			Secret:     maskPassword(cfg.JWT.Secret),
			Issuer:     cfg.JWT.Issuer,
			ExpireTime: cfg.JWT.ExpireTime,
		},
		Cache: types.CacheInfo{
			Driver: cfg.Cache.Driver,
		},
	}

	// 处理缓存配置
	if cfg.Cache.Driver == "memory" {
		configInfo.Cache.Memory = &types.MemoryCacheInfo{
			MaxSize:           cfg.Cache.Memory.MaxSize,
			DefaultExpiration: cfg.Cache.Memory.DefaultExpiration,
			CleanupInterval:   cfg.Cache.Memory.CleanupInterval,
		}
	}

	if cfg.Cache.Driver == "redis" {
		configInfo.Cache.Redis = &types.RedisCacheInfo{
			Host:         cfg.Cache.Redis.Host,
			Port:         cfg.Cache.Redis.Port,
			Password:     maskPassword(cfg.Cache.Redis.Password),
			Database:     cfg.Cache.Redis.Database,
			PoolSize:     cfg.Cache.Redis.PoolSize,
			MinIdleConns: cfg.Cache.Redis.MinIdleConns,
			DialTimeout:  cfg.Cache.Redis.DialTimeout,
			ReadTimeout:  cfg.Cache.Redis.ReadTimeout,
			WriteTimeout: cfg.Cache.Redis.WriteTimeout,
			KeyPrefix:    cfg.Cache.Redis.KeyPrefix,
		}
	}

	return configInfo, nil
}

// ToJSON 将配置转换为JSON字符串
func (p *ConfigParser) ToJSON(configInfo *types.ConfigInfo) (string, error) {
	data, err := json.MarshalIndent(configInfo, "", "  ")
	if err != nil {
		return "", fmt.Errorf("序列化配置失败: %w", err)
	}
	return string(data), nil
}

// maskPassword 遮蔽密码信息
func maskPassword(password string) string {
	if password == "" {
		return ""
	}
	if len(password) <= 4 {
		return "****"
	}
	return password[:2] + "****" + password[len(password)-2:]
}

// ValidateConfig 验证配置的完整性
func (p *ConfigParser) ValidateConfig(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("配置对象为空")
	}

	// 验证应用配置
	if cfg.App.Name == "" {
		return fmt.Errorf("应用名称不能为空")
	}

	// 验证服务器配置
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("服务器端口无效: %d", cfg.Server.Port)
	}

	// 验证数据库配置
	if cfg.Database.Driver == "" {
		return fmt.Errorf("数据库驱动不能为空")
	}

	return nil
}

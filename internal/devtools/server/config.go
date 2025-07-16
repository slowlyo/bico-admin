package server

import (
	"time"
)

// Config MCP服务器配置
type Config struct {
	// 服务器基本信息
	Name        string `yaml:"name" json:"name"`
	Version     string `yaml:"version" json:"version"`
	Description string `yaml:"description" json:"description"`

	// 传输配置
	Transport TransportConfig `yaml:"transport" json:"transport"`

	// 日志配置
	Log LogConfig `yaml:"log" json:"log"`

	// 工具配置
	Tools ToolsConfig `yaml:"tools" json:"tools"`
}

// TransportConfig 传输层配置
type TransportConfig struct {
	// HTTP配置
	HTTP HTTPConfig `yaml:"http" json:"http"`
}

// HTTPConfig HTTP传输配置
type HTTPConfig struct {
	Host         string        `yaml:"host" json:"host"`
	Port         int           `yaml:"port" json:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" json:"idle_timeout"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `yaml:"level" json:"level"`
	Format string `yaml:"format" json:"format"`
}

// ToolsConfig 工具配置
type ToolsConfig struct {
	// 是否启用配置读取工具
	ConfigReader bool `yaml:"config_reader" json:"config_reader"`

	// 是否启用数据库操作工具
	DatabaseTool bool `yaml:"database_tool" json:"database_tool"`

	// 是否启用CRUD生成工具（后续实现）
	CRUDGenerator bool `yaml:"crud_generator" json:"crud_generator"`

	// 是否启用数据库辅助工具（后续实现）
	DBHelper bool `yaml:"db_helper" json:"db_helper"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Name:        "Bico Admin DevTools",
		Version:     "1.0.0",
		Description: "Bico Admin 开发工具 MCP 服务",
		Transport: TransportConfig{
			HTTP: HTTPConfig{
				Host:         "127.0.0.1",
				Port:         18901, // 使用不常用的端口
				ReadTimeout:  30 * time.Second,
				WriteTimeout: 30 * time.Second,
				IdleTimeout:  60 * time.Second,
			},
		},
		Log: LogConfig{
			Level:  "info",
			Format: "console",
		},
		Tools: ToolsConfig{
			ConfigReader:  true,
			DatabaseTool:  true,  // 数据库操作工具
			CRUDGenerator: false, // 后续实现
			DBHelper:      false, // 后续实现
		},
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Name == "" {
		c.Name = "Bico Admin DevTools"
	}

	if c.Version == "" {
		c.Version = "1.0.0"
	}

	// 验证HTTP配置
	if c.Transport.HTTP.Host == "" {
		c.Transport.HTTP.Host = "127.0.0.1"
	}
	if c.Transport.HTTP.Port <= 0 {
		c.Transport.HTTP.Port = 18901 // 使用不常用的端口
	}
	if c.Transport.HTTP.ReadTimeout <= 0 {
		c.Transport.HTTP.ReadTimeout = 30 * time.Second
	}
	if c.Transport.HTTP.WriteTimeout <= 0 {
		c.Transport.HTTP.WriteTimeout = 30 * time.Second
	}
	if c.Transport.HTTP.IdleTimeout <= 0 {
		c.Transport.HTTP.IdleTimeout = 60 * time.Second
	}

	if c.Log.Level == "" {
		c.Log.Level = "info"
	}

	if c.Log.Format == "" {
		c.Log.Format = "console"
	}

	return nil
}

package server

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
	// 传输模式：http, stdio
	Mode string `yaml:"mode" json:"mode"`
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

	// 是否启用代码生成工具
	CodeGenerator bool `yaml:"code_generator" json:"code_generator"`

	// 是否启用目录树工具
	DirectoryTree bool `yaml:"directory_tree" json:"directory_tree"`

	// 是否启用表结构查看工具
	TableSchema bool `yaml:"table_schema" json:"table_schema"`

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
			Mode: "stdio", // 默认使用 stdio 模式
		},
		Log: LogConfig{
			Level:  "info",
			Format: "console",
		},
		Tools: ToolsConfig{
			ConfigReader:  true,
			DatabaseTool:  true,  // 数据库操作工具
			CodeGenerator: true,  // 代码生成工具
			DirectoryTree: true,  // 目录树工具
			TableSchema:   true,  // 表结构查看工具
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

	// 验证传输模式
	if c.Transport.Mode == "" {
		c.Transport.Mode = "stdio"
	}

	if c.Log.Level == "" {
		c.Log.Level = "info"
	}

	if c.Log.Format == "" {
		c.Log.Format = "console"
	}

	return nil
}

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
	var configFileFound bool

	// 设置配置文件路径
	if configPath != "" {
		// 检查指定的配置文件是否存在
		if _, err := os.Stat(configPath); err != nil {
			return nil, fmt.Errorf("指定的配置文件不存在: %s", configPath)
		}
		v.SetConfigFile(configPath)
		configFileFound = true
	} else {
		// 智能查找配置文件
		configFile, err := findConfigFile()
		if err != nil {
			return nil, fmt.Errorf("查找配置文件失败: %w", err)
		}
		if configFile != "" {
			v.SetConfigFile(configFile)
			configFileFound = true
		} else {
			// 如果没有找到配置文件，设置搜索路径但不强制要求文件存在
			v.SetConfigName("app")
			v.SetConfigType("yml")

			// 添加搜索路径
			configPaths := findConfigPaths()
			for _, path := range configPaths {
				v.AddConfigPath(path)
			}
			configFileFound = false
		}
	}

	// 环境变量配置
	v.SetEnvPrefix("BICO")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 设置默认值（无论是否找到配置文件都设置）
	setDefaults(v)

	// 尝试读取配置文件
	if configFileFound {
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	} else {
		// 没有找到配置文件，尝试读取但不报错
		_ = v.ReadInConfig() // 忽略错误，使用默认值
	}

	// 尝试读取环境特定的配置文件（仅当主配置文件存在时）
	if configFileFound && v.ConfigFileUsed() != "" {
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
	if strings.TrimSpace(config.App.Name) == "" {
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

	if strings.TrimSpace(config.JWT.Secret) == "" {
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
	v.SetDefault("jwt.secret", "bico-admin-default-jwt-secret-key-for-development-only")
	v.SetDefault("jwt.issuer", "bico-admin")
	v.SetDefault("jwt.expire_time", "24h")

	// Cache defaults
	v.SetDefault("cache.driver", "memory")
	v.SetDefault("cache.memory.max_size", 10000)
	v.SetDefault("cache.memory.default_expiration", "30m")
	v.SetDefault("cache.memory.cleanup_interval", "10m")
	v.SetDefault("cache.redis.host", "localhost")
	v.SetDefault("cache.redis.port", 6379)
	v.SetDefault("cache.redis.database", 1)
	v.SetDefault("cache.redis.pool_size", 10)
	v.SetDefault("cache.redis.min_idle_conns", 5)
	v.SetDefault("cache.redis.dial_timeout", "5s")
	v.SetDefault("cache.redis.read_timeout", "3s")
	v.SetDefault("cache.redis.write_timeout", "3s")
	v.SetDefault("cache.redis.key_prefix", "bico:cache:")

	// Upload defaults
	v.SetDefault("upload.max_file_size", "10MB")
	v.SetDefault("upload.max_files", 10)
	v.SetDefault("upload.allowed_types", []string{
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp",
		".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
		".txt", ".md", ".csv",
		".zip", ".rar", ".7z",
	})
	v.SetDefault("upload.upload_dir", "data/uploads")
	v.SetDefault("upload.base_url", "")

	// Frontend defaults
	v.SetDefault("frontend.mode", "embed")
	v.SetDefault("frontend.static_dir", "web/dist")
	v.SetDefault("frontend.index_file", "web/dist/index.html")
	v.SetDefault("frontend.assets_dir", "web/dist/assets")
}

// findConfigPaths 智能查找配置文件路径
func findConfigPaths() []string {
	var paths []string

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		// 如果获取失败，使用默认路径
		return []string{".", "./config", "../config", "../../config"}
	}

	// 查找项目根目录（包含 go.mod 文件的目录）
	projectRoot := findProjectRoot(currentDir)
	if projectRoot != "" {
		// 项目根目录本身（支持配置文件直接放在根目录）
		paths = append(paths, projectRoot)
		// 项目根目录的 config 子目录
		paths = append(paths, filepath.Join(projectRoot, "config"))
		// 项目根目录的 configs 子目录
		paths = append(paths, filepath.Join(projectRoot, "configs"))
	}

	// 尝试从可执行文件路径推断项目根目录
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)

		// 可执行文件所在目录
		paths = append(paths, execDir)

		// 如果可执行文件在 bin 目录下，项目根目录可能是上一级
		if filepath.Base(execDir) == "bin" {
			parentDir := filepath.Dir(execDir)
			if isProjectRoot(parentDir) {
				paths = append(paths, parentDir)
				paths = append(paths, filepath.Join(parentDir, "config"))
				paths = append(paths, filepath.Join(parentDir, "configs"))
			}
		}

		// 也尝试从可执行文件目录向上查找
		if execProjectRoot := findProjectRoot(execDir); execProjectRoot != "" && execProjectRoot != projectRoot {
			paths = append(paths, execProjectRoot)
			paths = append(paths, filepath.Join(execProjectRoot, "config"))
			paths = append(paths, filepath.Join(execProjectRoot, "configs"))
		}
	}

	// 添加相对路径作为备选
	paths = append(paths, ".", "./config", "./configs", "../config", "../configs", "../../config", "../../configs")

	// 去重
	seen := make(map[string]bool)
	var uniquePaths []string
	for _, path := range paths {
		if !seen[path] {
			seen[path] = true
			uniquePaths = append(uniquePaths, path)
		}
	}

	return uniquePaths
}

// findProjectRoot 查找项目根目录（包含 go.mod 的目录）
func findProjectRoot(startDir string) string {
	dir := startDir

	// 向上查找，最多查找10层
	for range 10 {
		// 检查是否存在 go.mod 文件
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		}

		// 向上一级目录
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			// 已经到达根目录
			break
		}
		dir = parentDir
	}

	return ""
}

// isProjectRoot 检查目录是否是项目根目录
func isProjectRoot(dir string) bool {
	goModPath := filepath.Join(dir, "go.mod")
	_, err := os.Stat(goModPath)
	return err == nil
}

// findConfigFile 智能查找配置文件
func findConfigFile() (string, error) {
	// 支持的配置文件名称（按优先级排序）
	configNames := []string{
		"app.yml", "app.yaml",
		"config.yml", "config.yaml",
	}

	// 搜索路径（按优先级排序）
	searchPaths := findConfigPaths()

	// 添加当前目录作为最高优先级
	currentDir, err := os.Getwd()
	if err == nil {
		searchPaths = append([]string{currentDir}, searchPaths...)
	}

	// 在每个路径中查找配置文件
	for _, searchPath := range searchPaths {
		for _, configName := range configNames {
			configFile := filepath.Join(searchPath, configName)
			if _, err := os.Stat(configFile); err == nil {
				return configFile, nil
			}
		}
	}

	return "", nil // 没有找到配置文件，返回空字符串（不是错误）
}

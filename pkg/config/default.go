package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

// GetDefaultConfig 获取默认配置
// 当没有找到配置文件时，使用此默认配置
func GetDefaultConfig() *Config {
	return &Config{
		App:      getDefaultAppConfig(),
		Server:   getDefaultServerConfig(),
		Database: getDefaultDatabaseConfig(),
		Redis:    getDefaultRedisConfig(),
		Log:      getDefaultLogConfig(),
		JWT:      getDefaultJWTConfig(),
		Cache:    getDefaultCacheConfig(),
		Upload:   getDefaultUploadConfig(),
		Frontend: getDefaultFrontendConfig(),
	}
}

// getDefaultAppConfig 获取默认应用配置
func getDefaultAppConfig() AppConfig {
	return AppConfig{
		Name:        "Bico Admin",
		Version:     "1.0.0",
		Environment: "development",
		Debug:       true,
	}
}

// getDefaultServerConfig 获取默认服务器配置
func getDefaultServerConfig() ServerConfig {
	return ServerConfig{
		Host:         "0.0.0.0",
		Port:         8080,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// getDefaultDatabaseConfig 获取默认数据库配置
func getDefaultDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Driver:          "sqlite",
		Host:            "localhost",
		Port:            3306,
		Username:        "",
		Password:        "",
		Database:        "data/bico_admin.db", // SQLite 默认文件路径
		Charset:         "utf8mb4",
		MaxIdleConns:    5,
		MaxOpenConns:    10,
		ConnMaxLifetime: 1 * time.Hour,
	}
}

// getDefaultRedisConfig 获取默认Redis配置
func getDefaultRedisConfig() RedisConfig {
	return RedisConfig{
		Host:         "localhost",
		Port:         6379,
		Password:     "",
		Database:     0,
		PoolSize:     10,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
}

// getDefaultLogConfig 获取默认日志配置
func getDefaultLogConfig() LogConfig {
	return LogConfig{
		Level:      "info",
		Format:     "json",
		Output:     "stdout",
		Filename:   "logs/app.log",
		MaxSize:    100, // MB
		MaxAge:     7,   // days
		MaxBackups: 10,
		Compress:   true,
	}
}

// getDefaultJWTConfig 获取默认JWT配置
func getDefaultJWTConfig() JWTConfig {
	// 生成随机密钥
	secret := generateRandomSecret()

	return JWTConfig{
		Secret:     secret,
		Issuer:     "bico-admin",
		ExpireTime: 24 * time.Hour,
	}
}

// getDefaultCacheConfig 获取默认缓存配置
func getDefaultCacheConfig() CacheConfig {
	return CacheConfig{
		Driver: "memory",
		Memory: CacheMemoryConfig{
			MaxSize:           10000,
			DefaultExpiration: 30 * time.Minute,
			CleanupInterval:   10 * time.Minute,
		},
		Redis: CacheRedisConfig{
			Host:         "localhost",
			Port:         6379,
			Password:     "",
			Database:     1,
			PoolSize:     10,
			MinIdleConns: 5,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			KeyPrefix:    "bico:cache:",
		},
	}
}

// getDefaultUploadConfig 获取默认上传配置
func getDefaultUploadConfig() UploadConfig {
	return UploadConfig{
		MaxFileSize: "10MB",
		MaxFiles:    10,
		AllowedTypes: []string{
			".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", // 图片
			".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", // 文档
			".txt", ".md", ".csv", // 文本
			".zip", ".rar", ".7z", // 压缩包
		},
		UploadDir: "data/uploads",
		BaseURL:   "",
	}
}

// getDefaultFrontendConfig 获取默认前端配置
func getDefaultFrontendConfig() FrontendConfig {
	return FrontendConfig{
		Mode:      "embed",
		StaticDir: "web/dist",
		IndexFile: "web/dist/index.html",
		AssetsDir: "web/dist/assets",
	}
}

// generateRandomSecret 生成随机密钥
func generateRandomSecret() string {
	// 生成32字节的随机数据
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// 如果生成随机数失败，使用固定的开发密钥（仅用于开发环境）
		log.Printf("警告: 生成随机JWT密钥失败，使用默认密钥: %v", err)
		return "bico-admin-default-jwt-secret-key-for-development-only"
	}
	return hex.EncodeToString(bytes)
}

// LoadWithDefaults 加载配置，如果失败则使用默认配置
func LoadWithDefaults(configPath string) (*Config, error) {
	// 首先尝试加载配置文件
	config, err := Load(configPath)
	if err != nil {
		log.Printf("加载配置文件失败，使用默认配置: %v", err)

		// 使用默认配置
		defaultConfig := GetDefaultConfig()

		// 验证默认配置
		if validateErr := validateConfig(defaultConfig); validateErr != nil {
			return nil, fmt.Errorf("默认配置验证失败: %w", validateErr)
		}

		// 设置全局配置
		globalConfig = defaultConfig

		log.Printf("已加载默认配置，应用将以开发模式运行")
		log.Printf("数据库: %s", defaultConfig.Database.GetDSN())
		log.Printf("服务器: %s:%d", defaultConfig.Server.Host, defaultConfig.Server.Port)
		log.Printf("JWT密钥: %s", maskSecret(defaultConfig.JWT.Secret))

		return defaultConfig, nil
	}

	return config, nil
}

// maskSecret 遮蔽密钥显示
func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "****"
	}
	return secret[:4] + "****" + secret[len(secret)-4:]
}

// LoadFromEnvWithDefaults 从环境变量加载配置，如果失败则使用默认配置
func LoadFromEnvWithDefaults() (*Config, error) {
	// 首先尝试从环境变量加载
	config, err := LoadFromEnv()
	if err != nil {
		log.Printf("从环境变量加载配置失败，使用默认配置: %v", err)

		// 使用默认配置
		defaultConfig := GetDefaultConfig()

		// 验证默认配置
		if validateErr := validateConfig(defaultConfig); validateErr != nil {
			return nil, fmt.Errorf("默认配置验证失败: %w", validateErr)
		}

		// 设置全局配置
		globalConfig = defaultConfig

		log.Printf("已加载默认配置，应用将以开发模式运行")

		return defaultConfig, nil
	}

	return config, nil
}

// MustLoad 必须加载配置，如果失败则使用默认配置
// 这个函数永远不会返回错误，适用于应用启动时使用
func MustLoad(configPath string) *Config {
	config, err := LoadWithDefaults(configPath)
	if err != nil {
		// 如果连默认配置都失败了，说明代码有问题，直接panic
		panic(fmt.Sprintf("无法加载配置: %v", err))
	}
	return config
}

// MustLoadFromEnv 必须从环境变量加载配置，如果失败则使用默认配置
func MustLoadFromEnv() *Config {
	config, err := LoadFromEnvWithDefaults()
	if err != nil {
		// 如果连默认配置都失败了，说明代码有问题，直接panic
		panic(fmt.Sprintf("无法加载配置: %v", err))
	}
	return config
}

// IsUsingDefaults 检查当前是否使用的是默认配置
func IsUsingDefaults() bool {
	if globalConfig == nil {
		return false
	}

	// 通过检查JWT密钥是否是随机生成的来判断
	// 如果是从配置文件加载的，通常不会是随机生成的64位十六进制字符串
	return len(globalConfig.JWT.Secret) == 64 && isHexString(globalConfig.JWT.Secret)
}

// isHexString 检查字符串是否为十六进制
func isHexString(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil
}

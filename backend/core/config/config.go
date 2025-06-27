package config

import (
	"os"
	"strconv"
)

// Config 应用配置结构
type Config struct {
	// 服务器配置
	Port string
	Env  string

	// 数据库配置
	Database DatabaseConfig

	// Redis配置
	Redis RedisConfig

	// JWT配置
	JWT JWTConfig

	// 文件上传配置
	Upload UploadConfig

	// 日志配置
	Log LogConfig

	// 路由前缀配置
	RoutePrefix RoutePrefixConfig
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string
	Expire string
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	Path    string
	MaxSize string
}

// LogConfig 日志配置
type LogConfig struct {
	Level string
	Path  string
}

// RoutePrefixConfig 路由前缀配置
type RoutePrefixConfig struct {
	Admin string
	API   string
}

// New 创建新的配置实例
func New() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 3306),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "bico_admin"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key"),
			Expire: getEnv("JWT_EXPIRE", "24h"),
		},
		Upload: UploadConfig{
			Path:    getEnv("UPLOAD_PATH", "./storage/uploads"),
			MaxSize: getEnv("MAX_UPLOAD_SIZE", "10MB"),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			Path:  getEnv("LOG_PATH", "./storage/logs"),
		},
		RoutePrefix: RoutePrefixConfig{
			Admin: getEnv("ADMIN_PREFIX", "/admin"),
			API:   getEnv("API_PREFIX", "/api"),
		},
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为整数
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

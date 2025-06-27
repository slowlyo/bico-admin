package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// LoggerMiddleware 日志中间件配置
func LoggerMiddleware() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} - ${ip} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Output:     os.Stdout,
	})
}

// FileLoggerMiddleware 文件日志中间件
func FileLoggerMiddleware(logPath string) fiber.Handler {
	// 确保日志目录存在
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return LoggerMiddleware() // 如果创建目录失败，使用标准输出
	}

	// 创建日志文件
	logFile := logPath + "/access.log"
	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return LoggerMiddleware() // 如果创建文件失败，使用标准输出
	}

	return logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} - ${ip} - ${latency} - ${ua}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Output:     file,
	})
}

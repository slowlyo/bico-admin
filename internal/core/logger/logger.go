package logger

import (
	"os"

	"bico-admin/internal/core/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 全局日志实例
var Logger *zap.Logger

// InitLogger 初始化日志
func InitLogger(cfg *config.LogConfig) (*zap.Logger, error) {
	// 解析日志级别
	level := parseLevel(cfg.Level)
	
	// 解析日志格式
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		// console 格式，更易读
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	
	// 解析输出位置
	writer := parseOutput(cfg.Output)
	
	core := zapcore.NewCore(encoder, writer, level)
	
	// 添加调用者信息
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	
	// 设置全局实例
	Logger = logger
	
	return logger, nil
}

// parseLevel 解析日志级别
func parseLevel(levelStr string) zapcore.Level {
	switch levelStr {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// parseOutput 解析输出位置
func parseOutput(output string) zapcore.WriteSyncer {
	switch output {
	case "stdout":
		return zapcore.AddSync(os.Stdout)
	case "stderr":
		return zapcore.AddSync(os.Stderr)
	default:
		// 如果是文件路径，写入文件
		file, err := os.OpenFile(output, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			// 如果打开文件失败，回退到 stdout
			return zapcore.AddSync(os.Stdout)
		}
		return zapcore.AddSync(file)
	}
}

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(msg, fields...)
	}
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(msg, fields...)
	}
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(msg, fields...)
	}
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Fatal(msg, fields...)
	}
}

// Sync 同步日志
func Sync() error {
	if Logger != nil {
		return Logger.Sync()
	}
	return nil
}

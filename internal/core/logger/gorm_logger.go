package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger 自定义 GORM 日志实现
type GormLogger struct {
	zapLogger                 *zap.Logger
	logLevel                  gormlogger.LogLevel
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
}

// NewGormLogger 创建 GORM 日志实例
func NewGormLogger(zapLogger *zap.Logger, logLevel gormlogger.LogLevel) *GormLogger {
	return &GormLogger{
		zapLogger:                 zapLogger,
		logLevel:                  logLevel,
		slowThreshold:             200 * time.Millisecond, // 慢查询阈值
		ignoreRecordNotFoundError: true,                   // 忽略记录不存在错误
	}
}

// LogMode 设置日志级别
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// Info 信息日志
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Info {
		l.zapLogger.Info(fmt.Sprintf(msg, data...))
	}
}

// Warn 警告日志
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Warn {
		l.zapLogger.Warn(fmt.Sprintf(msg, data...))
	}
}

// Error 错误日志
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Error {
		l.zapLogger.Error(fmt.Sprintf(msg, data...))
	}
}

// Trace SQL 执行日志
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 错误日志
	if err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.ignoreRecordNotFoundError) {
		l.zapLogger.Error("SQL 执行错误",
			zap.Error(err),
			zap.Duration("耗时", elapsed),
			zap.Int64("影响行数", rows),
			zap.String("SQL", sql),
		)
		return
	}

	// 慢查询日志
	if l.slowThreshold > 0 && elapsed > l.slowThreshold {
		l.zapLogger.Warn("慢查询检测",
			zap.Duration("耗时", elapsed),
			zap.Duration("阈值", l.slowThreshold),
			zap.Int64("影响行数", rows),
			zap.String("SQL", sql),
		)
		return
	}

	// 正常查询日志（仅在 Info 级别及以上输出）
	if l.logLevel >= gormlogger.Info {
		l.zapLogger.Debug("SQL 执行",
			zap.Duration("耗时", elapsed),
			zap.Int64("影响行数", rows),
			zap.String("SQL", sql),
		)
	}
}

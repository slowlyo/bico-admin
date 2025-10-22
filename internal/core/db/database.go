package db

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"bico-admin/internal/core/config"
	"bico-admin/internal/core/logger"
	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// InitDB 初始化数据库连接
func InitDB(cfg *config.DatabaseConfig, zapLogger *zap.Logger, isDebug bool) (*gorm.DB, error) {
	var dialector gorm.Dialector
	
	switch cfg.Driver {
	case "sqlite":
		dialector = buildSQLiteDialector(cfg)
	case "mysql":
		dialector = buildMySQLDialector(cfg)
	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", cfg.Driver)
	}

	// 根据 debug 模式设置日志级别
	var logLevel gormlogger.LogLevel
	if isDebug {
		logLevel = gormlogger.Info // debug 模式下输出所有 SQL
	} else {
		logLevel = gormlogger.Warn // 生产模式仅输出警告和错误
	}

	// 使用自定义日志
	gormLog := logger.NewGormLogger(zapLogger, logLevel)

	db, err := gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: gormLog,
	})
	
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// buildSQLiteDialector 构建 SQLite 驱动配置
func buildSQLiteDialector(cfg *config.DatabaseConfig) gorm.Dialector {
	dbPath := cfg.SQLite.Path
	if dbPath == "" {
		dbPath = "data.db"
	}
	
	// 确保数据库文件的父目录存在
	dir := filepath.Dir(dbPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("创建数据库目录失败: %v\n", err)
		}
	}
	
	return sqlite.Open(dbPath)
}

// buildMySQLDialector 构建 MySQL 驱动配置
func buildMySQLDialector(cfg *config.DatabaseConfig) gorm.Dialector {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.MySQL.Username,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.Database,
		cfg.MySQL.Charset,
	)
	return mysql.Open(dsn)
}

package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"bico-admin/core/model"
)

// InitDatabase 初始化数据库连接
func InitDatabase(cfg *Config) (*gorm.DB, error) {
	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	// 配置GORM日志级别
	var logLevel logger.LogLevel
	switch cfg.Log.Level {
	case "debug":
		logLevel = logger.Info
	case "info":
		logLevel = logger.Warn
	default:
		logLevel = logger.Error
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 自动迁移
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database connected successfully")
	return db, nil
}

// autoMigrate 自动迁移数据库表
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		// 在这里添加其他需要迁移的模型
	)
}

// GetDB 获取数据库实例（单例模式）
var dbInstance *gorm.DB

func GetDB() *gorm.DB {
	return dbInstance
}

func SetDB(db *gorm.DB) {
	dbInstance = db
}

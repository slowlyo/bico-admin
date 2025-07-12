package shared

import (
	"fmt"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"bico-admin/internal/shared/model"
	"bico-admin/pkg/config"
	"bico-admin/pkg/database"
	pkgLogger "bico-admin/pkg/logger"
)

// ProviderSet 共享组件Provider集合
var ProviderSet = wire.NewSet(
	ProvideDatabase,
	ProvideRedis,
)

// ProvideDatabase 提供数据库连接
func ProvideDatabase(cfg *config.Config) (*gorm.DB, error) {
	// 设置日志级别
	logLevel := logger.Silent
	if cfg.App.Debug {
		logLevel = logger.Info
	}

	var db *gorm.DB
	var err error

	switch cfg.Database.Driver {
	case "mysql":
		dbConfig := database.MySQLConfig{
			Host:            cfg.Database.Host,
			Port:            cfg.Database.Port,
			Username:        cfg.Database.Username,
			Password:        cfg.Database.Password,
			Database:        cfg.Database.Database,
			Charset:         cfg.Database.Charset,
			MaxIdleConns:    cfg.Database.MaxIdleConns,
			MaxOpenConns:    cfg.Database.MaxOpenConns,
			ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
			LogLevel:        logLevel,
		}
		db, err = database.NewMySQL(dbConfig)

	case "postgres":
		dbConfig := database.PostgresConfig{
			Host:            cfg.Database.Host,
			Port:            cfg.Database.Port,
			Username:        cfg.Database.Username,
			Password:        cfg.Database.Password,
			Database:        cfg.Database.Database,
			MaxIdleConns:    cfg.Database.MaxIdleConns,
			MaxOpenConns:    cfg.Database.MaxOpenConns,
			ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
			LogLevel:        logLevel,
		}
		db, err = database.NewPostgres(dbConfig)

	case "sqlite":
		dbConfig := database.SQLiteConfig{
			Database:        cfg.Database.Database,
			MaxIdleConns:    cfg.Database.MaxIdleConns,
			MaxOpenConns:    cfg.Database.MaxOpenConns,
			ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
			LogLevel:        logLevel,
		}
		db, err = database.NewSQLite(dbConfig)

	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", cfg.Database.Driver)
	}

	if err != nil {
		return nil, err
	}

	// 自动迁移数据库表
	if err := autoMigrate(db); err != nil {
		pkgLogger.Error("数据库迁移失败", zap.Error(err))
		return nil, err
	}

	return db, nil
}

// ProvideRedis 提供Redis连接
func ProvideRedis(cfg *config.Config) (*redis.Client, error) {
	redisConfig := database.RedisConfig{
		Host:         cfg.Redis.Host,
		Port:         cfg.Redis.Port,
		Password:     cfg.Redis.Password,
		Database:     cfg.Redis.Database,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
	}

	return database.NewRedis(redisConfig)
}

// autoMigrate 自动迁移数据库表
func autoMigrate(db *gorm.DB) error {
	// 导入模型包
	models := []interface{}{
		&model.User{},
		// 其他模型...
	}

	// 执行自动迁移
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	pkgLogger.Info("数据库迁移完成")
	return nil
}

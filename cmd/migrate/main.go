package main

import (
	"flag"
	"fmt"
	"log"

	"bico-admin/internal/shared/model"
	"bico-admin/pkg/config"
	"bico-admin/pkg/database"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	var (
		configPath = flag.String("config", "config/app.yml", "配置文件路径")
		action     = flag.String("action", "migrate", "操作类型: migrate(迁移), rollback(回滚), fresh(重新创建)")
	)
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 连接数据库
	db, err := connectDatabase(cfg)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 执行操作
	switch *action {
	case "migrate":
		if err := migrate(db); err != nil {
			log.Fatalf("迁移失败: %v", err)
		}
		fmt.Println("✅ 数据库迁移完成")
	case "rollback":
		if err := rollback(db); err != nil {
			log.Fatalf("回滚失败: %v", err)
		}
		fmt.Println("✅ 数据库回滚完成")
	case "fresh":
		if err := fresh(db); err != nil {
			log.Fatalf("重新创建失败: %v", err)
		}
		fmt.Println("✅ 数据库重新创建完成")
	default:
		log.Fatalf("不支持的操作类型: %s", *action)
	}
}

// connectDatabase 连接数据库
func connectDatabase(cfg *config.Config) (*gorm.DB, error) {
	logLevel := logger.Silent
	if cfg.App.Debug {
		logLevel = logger.Info
	}

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
		return database.NewMySQL(dbConfig)

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
		return database.NewPostgres(dbConfig)

	case "sqlite":
		dbConfig := database.SQLiteConfig{
			Database:        cfg.Database.Database,
			MaxIdleConns:    cfg.Database.MaxIdleConns,
			MaxOpenConns:    cfg.Database.MaxOpenConns,
			ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
			LogLevel:        logLevel,
		}
		return database.NewSQLite(dbConfig)

	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", cfg.Database.Driver)
	}
}

// migrate 执行迁移
func migrate(db *gorm.DB) error {
	models := []interface{}{
		&model.User{},
		&model.AdminUser{},
		// 其他模型...
	}

	return db.AutoMigrate(models...)
}

// rollback 回滚迁移 (简单实现，删除表)
func rollback(db *gorm.DB) error {
	tables := []string{
		"admin_users",
		"users",
		// 其他表...
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			return fmt.Errorf("删除表 %s 失败: %w", table, err)
		}
		fmt.Printf("已删除表: %s\n", table)
	}

	return nil
}

// fresh 重新创建数据库
func fresh(db *gorm.DB) error {
	// 先回滚
	if err := rollback(db); err != nil {
		return fmt.Errorf("回滚失败: %w", err)
	}

	// 再迁移
	if err := migrate(db); err != nil {
		return fmt.Errorf("迁移失败: %w", err)
	}

	return nil
}

package db

import (
	"path/filepath"
	"testing"

	"bico-admin/internal/core/config"

	"go.uber.org/zap"
)

// TestSQLiteUsesWAL 验证默认连接开启 WAL，并在锁等待后返回。
func TestSQLiteUsesWAL(t *testing.T) {
	database, err := InitDB(&config.DatabaseConfig{
		Driver:       "sqlite",
		SQLite:       config.SQLiteConfig{Path: filepath.Join(t.TempDir(), "data.db")},
		MaxIdleConns: 1,
		MaxOpenConns: 1,
	}, zap.NewNop(), false)
	if err != nil {
		t.Fatalf("打开 SQLite 失败: %v", err)
	}
	sqlDB, err := database.DB()
	if err != nil {
		t.Fatalf("获取连接失败: %v", err)
	}
	defer sqlDB.Close()

	var mode string
	if err := database.Raw("PRAGMA journal_mode").Scan(&mode).Error; err != nil {
		t.Fatalf("读取 journal_mode 失败: %v", err)
	}
	if mode != "wal" {
		t.Fatalf("journal_mode = %s，期望 wal", mode)
	}

	var timeout int
	if err := database.Raw("PRAGMA busy_timeout").Scan(&timeout).Error; err != nil {
		t.Fatalf("读取 busy_timeout 失败: %v", err)
	}
	if timeout != 5000 {
		t.Fatalf("busy_timeout = %d，期望 5000", timeout)
	}
}

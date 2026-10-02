package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bico-admin/internal/core/config"
)

// TestSetLevelChangesExistingLogger 验证调高日志级别后，已创建的 logger 立即输出更低级别。
func TestSetLevelChangesExistingLogger(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.log")
	lg, err := InitLogger(&config.LogConfig{Level: "info", Format: "json", Output: path})
	if err != nil {
		t.Fatalf("初始化日志失败: %v", err)
	}
	defer lg.Sync()

	lg.Info("keep")
	lg.Debug("drop")
	SetLevel("debug")
	lg.Debug("show")
	if err := lg.Sync(); err != nil {
		t.Fatalf("同步日志失败: %v", err)
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("读取日志失败: %v", err)
	}
	text := string(body)
	if !strings.Contains(text, "keep") || strings.Contains(text, "drop") || !strings.Contains(text, "show") {
		t.Fatalf("日志级别未按预期过滤: %s", text)
	}
}

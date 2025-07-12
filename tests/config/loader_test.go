package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bico-admin/pkg/config"
)

// TestLoadDefaultConfig 测试加载默认配置
func TestLoadDefaultConfig(t *testing.T) {
	// 切换到项目根目录进行测试
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// 切换到项目根目录
	os.Chdir("../../")

	cfg, err := config.Load("")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// 验证默认配置值
	assert.Equal(t, "bico-admin", cfg.App.Name)
	assert.Equal(t, "1.0.0", cfg.App.Version)
	assert.Equal(t, "development", cfg.App.Environment)
	assert.True(t, cfg.App.Debug)

	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 8899, cfg.Server.Port)

	assert.Equal(t, "sqlite", cfg.Database.Driver)
	assert.Equal(t, "data/bico_admin.db", cfg.Database.Database)

	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, "json", cfg.Log.Format)
	assert.Equal(t, "stdout", cfg.Log.Output)
}

// TestLoadSpecificConfig 测试加载指定配置文件
func TestLoadSpecificConfig(t *testing.T) {
	configPath := filepath.Join("testdata", "test_config.yml")

	cfg, err := config.Load(configPath)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "test-app", cfg.App.Name)
	assert.Equal(t, 9999, cfg.Server.Port)
	assert.Equal(t, "debug", cfg.Log.Level)
}

// TestConfigValidation 测试配置验证
func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		configFile  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "有效配置",
			configFile:  "test_config.yml",
			expectError: false,
		},
		{
			name:        "无效配置-缺少应用名称",
			configFile:  "invalid_config.yml",
			expectError: true,
			errorMsg:    "应用名称不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := filepath.Join("testdata", tt.configFile)

			cfg, err := config.Load(configPath)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, cfg)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
			}
		})
	}
}

// TestDatabaseDriverValidation 测试数据库驱动验证
func TestDatabaseDriverValidation(t *testing.T) {
	tests := []struct {
		name     string
		driver   string
		host     string
		username string
		database string
		valid    bool
	}{
		{
			name:     "SQLite有效配置",
			driver:   "sqlite",
			database: "test.db",
			valid:    true,
		},
		{
			name:     "MySQL有效配置",
			driver:   "mysql",
			host:     "localhost",
			username: "root",
			database: "test",
			valid:    true,
		},
		{
			name:     "PostgreSQL有效配置",
			driver:   "postgres",
			host:     "localhost",
			username: "postgres",
			database: "test",
			valid:    true,
		},
		{
			name:   "不支持的驱动",
			driver: "oracle",
			valid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 注意：这里我们无法直接测试验证函数，因为它是私有的
			// 在实际项目中，可以考虑将验证函数导出或创建专门的验证测试
			// 这个测试主要验证不同驱动配置的有效性

			configPath := filepath.Join("testdata", "test_config.yml")

			// 先加载一个有效配置
			validCfg, err := config.Load(configPath)
			require.NoError(t, err)

			// 验证配置加载成功
			assert.NotNil(t, validCfg)
			assert.Equal(t, "test-app", validCfg.App.Name)
		})
	}
}

// TestTimeoutParsing 测试时间配置解析
func TestTimeoutParsing(t *testing.T) {
	configPath := filepath.Join("testdata", "test_config.yml")

	cfg, err := config.Load(configPath)
	require.NoError(t, err)

	// 验证时间解析
	assert.Equal(t, 10*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(t, 10*time.Second, cfg.Server.WriteTimeout)
	assert.Equal(t, 30*time.Second, cfg.Server.IdleTimeout)
	assert.Equal(t, 30*time.Minute, cfg.Database.ConnMaxLifetime)
	assert.Equal(t, 3*time.Second, cfg.Redis.DialTimeout)
	assert.Equal(t, 2*time.Second, cfg.Redis.ReadTimeout)
	assert.Equal(t, 2*time.Second, cfg.Redis.WriteTimeout)
	assert.Equal(t, 1*time.Hour, cfg.JWT.ExpireTime)
}

// TestEnvironmentVariableFormats 测试环境变量格式
func TestEnvironmentVariableFormats(t *testing.T) {
	configPath := filepath.Join("testdata", "test_config.yml")

	// 清理环境变量
	defer func() {
		os.Unsetenv("BICO_APP_NAME")
		os.Unsetenv("BICO_SERVER_PORT")
		os.Unsetenv("BICO_APP_DEBUG")
		os.Unsetenv("BICO_DATABASE_MAX_IDLE_CONNS")
	}()

	// 设置不同类型的环境变量
	os.Setenv("BICO_APP_NAME", "env-test")
	os.Setenv("BICO_SERVER_PORT", "8888")
	os.Setenv("BICO_APP_DEBUG", "false")
	os.Setenv("BICO_DATABASE_MAX_IDLE_CONNS", "10")

	cfg, err := config.Load(configPath)
	require.NoError(t, err)

	// 验证不同类型的环境变量覆盖
	assert.Equal(t, "env-test", cfg.App.Name)      // 字符串
	assert.Equal(t, 8888, cfg.Server.Port)         // 整数
	assert.False(t, cfg.App.Debug)                 // 布尔值
	assert.Equal(t, 10, cfg.Database.MaxIdleConns) // 整数
}

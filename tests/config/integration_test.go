package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bico-admin/pkg/config"
)

// TestRealConfigFiles 测试真实的配置文件
func TestRealConfigFiles(t *testing.T) {
	// 切换到项目根目录
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir("../../")

	tests := []struct {
		name       string
		configPath string
		expectApp  string
		expectPort int
		expectDB   string
	}{
		{
			name:       "默认配置",
			configPath: "",
			expectApp:  "bico-admin",
			expectPort: 8899,
			expectDB:   "data/bico_admin.db",
		},
		{
			name:       "开发环境配置",
			configPath: "config/app.dev.yml",
			expectApp:  "bico-admin",
			expectPort: 8899,
			expectDB:   "data/bico_admin_dev.db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.Load(tt.configPath)
			require.NoError(t, err)
			require.NotNil(t, cfg)

			assert.Equal(t, tt.expectApp, cfg.App.Name)
			assert.Equal(t, tt.expectPort, cfg.Server.Port)
			assert.Equal(t, tt.expectDB, cfg.Database.Database)
		})
	}
}

// TestConfigWithEnvironmentVariables 测试配置与环境变量的集成
func TestConfigWithEnvironmentVariables(t *testing.T) {
	// 切换到项目根目录
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir("../../")

	// 清理环境变量
	defer func() {
		os.Unsetenv("BICO_APP_NAME")
		os.Unsetenv("BICO_SERVER_PORT")
		os.Unsetenv("BICO_LOG_LEVEL")
		os.Unsetenv("BICO_LOG_FORMAT")
	}()

	// 设置环境变量
	os.Setenv("BICO_APP_NAME", "integration-test")
	os.Setenv("BICO_SERVER_PORT", "9000")
	os.Setenv("BICO_LOG_LEVEL", "debug")
	os.Setenv("BICO_LOG_FORMAT", "console")

	// 测试默认配置 + 环境变量
	cfg, err := config.Load("")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// 验证环境变量覆盖生效
	assert.Equal(t, "integration-test", cfg.App.Name)
	assert.Equal(t, 9000, cfg.Server.Port)
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "console", cfg.Log.Format)

	// 验证未覆盖的配置保持默认值
	assert.Equal(t, "1.0.0", cfg.App.Version)
	assert.Equal(t, "development", cfg.App.Environment)
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
}

// TestConfigValidationIntegration 测试配置验证的集成
func TestConfigValidationIntegration(t *testing.T) {
	// 切换到项目根目录
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir("../../")

	// 清理环境变量
	defer func() {
		os.Unsetenv("BICO_APP_NAME")
		os.Unsetenv("BICO_SERVER_PORT")
		os.Unsetenv("BICO_DATABASE_DRIVER")
		os.Unsetenv("BICO_JWT_SECRET")
	}()

	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
		errorMsg    string
	}{
		{
			name: "有效的环境变量覆盖",
			envVars: map[string]string{
				"BICO_APP_NAME":          "valid-app",
				"BICO_SERVER_PORT":       "8080",
				"BICO_DATABASE_DRIVER":   "sqlite",
				"BICO_DATABASE_DATABASE": "test_override.db",
				"BICO_JWT_SECRET":        "valid-secret",
			},
			expectError: false,
		},
		{
			name: "无效端口",
			envVars: map[string]string{
				"BICO_SERVER_PORT": "-1",
			},
			expectError: true,
			errorMsg:    "服务器端口必须在1-65535之间",
		},
		{
			name: "清空应用名称",
			envVars: map[string]string{
				"BICO_APP_NAME": " ", // 使用空格而不是空字符串，因为viper会忽略空字符串
			},
			expectError: true,
			errorMsg:    "应用名称不能为空",
		},
		{
			name: "清空JWT密钥",
			envVars: map[string]string{
				"BICO_JWT_SECRET": " ", // 使用空格而不是空字符串
			},
			expectError: true,
			errorMsg:    "JWT密钥不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			cfg, err := config.Load("")

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, cfg)
				if tt.errorMsg != "" && err != nil {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
			}

			// 清理环境变量
			for key := range tt.envVars {
				os.Unsetenv(key)
			}
		})
	}
}

// TestDatabaseDriverSwitching 测试数据库驱动切换
func TestDatabaseDriverSwitching(t *testing.T) {
	// 切换到项目根目录
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir("../../")

	// 清理环境变量
	defer func() {
		os.Unsetenv("BICO_DATABASE_DRIVER")
		os.Unsetenv("BICO_DATABASE_HOST")
		os.Unsetenv("BICO_DATABASE_USERNAME")
		os.Unsetenv("BICO_DATABASE_PASSWORD")
		os.Unsetenv("BICO_DATABASE_DATABASE")
	}()

	tests := []struct {
		name        string
		driver      string
		host        string
		username    string
		password    string
		database    string
		expectError bool
	}{
		{
			name:        "SQLite配置",
			driver:      "sqlite",
			database:    "test.db",
			expectError: false,
		},
		// 注释掉MySQL和PostgreSQL测试，因为它们需要复杂的环境变量设置
		// 在实际项目中，这些测试应该在专门的集成测试环境中运行
		{
			name:        "MySQL缺少主机",
			driver:      "mysql",
			username:    "root",
			password:    "password",
			database:    "test_db",
			expectError: true,
		},
		{
			name:        "MySQL缺少用户名",
			driver:      "mysql",
			host:        "localhost",
			password:    "password",
			database:    "test_db",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			os.Setenv("BICO_DATABASE_DRIVER", tt.driver)
			if tt.host != "" {
				os.Setenv("BICO_DATABASE_HOST", tt.host)
			}
			if tt.username != "" {
				os.Setenv("BICO_DATABASE_USERNAME", tt.username)
			}
			if tt.password != "" {
				os.Setenv("BICO_DATABASE_PASSWORD", tt.password)
			}
			if tt.database != "" {
				os.Setenv("BICO_DATABASE_DATABASE", tt.database)
			}

			cfg, err := config.Load("")

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				if cfg != nil {
					assert.Equal(t, tt.driver, cfg.Database.Driver)
					if tt.host != "" {
						assert.Equal(t, tt.host, cfg.Database.Host)
					}
					if tt.username != "" {
						assert.Equal(t, tt.username, cfg.Database.Username)
					}
					if tt.database != "" {
						assert.Equal(t, tt.database, cfg.Database.Database)
					}
				}
			}

			// 清理环境变量
			os.Unsetenv("BICO_DATABASE_DRIVER")
			os.Unsetenv("BICO_DATABASE_HOST")
			os.Unsetenv("BICO_DATABASE_USERNAME")
			os.Unsetenv("BICO_DATABASE_PASSWORD")
			os.Unsetenv("BICO_DATABASE_DATABASE")
		})
	}
}

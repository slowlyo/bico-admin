package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"bico-admin/pkg/config"
)

// ConfigTestSuite 配置测试套件
type ConfigTestSuite struct {
	suite.Suite
	testDataDir string
}

// SetupSuite 测试套件初始化
func (suite *ConfigTestSuite) SetupSuite() {
	// 获取测试数据目录
	suite.testDataDir = filepath.Join("testdata")
}

// TearDownTest 每个测试后清理
func (suite *ConfigTestSuite) TearDownTest() {
	// 清理环境变量
	os.Unsetenv("BICO_APP_NAME")
	os.Unsetenv("BICO_SERVER_PORT")
	os.Unsetenv("BICO_LOG_LEVEL")
	os.Unsetenv("BICO_DATABASE_DRIVER")
}

// TestLoadValidConfig 测试加载有效配置
func (suite *ConfigTestSuite) TestLoadValidConfig() {
	configPath := filepath.Join(suite.testDataDir, "test_config.yml")

	cfg, err := config.Load(configPath)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), cfg)

	// 验证应用配置
	assert.Equal(suite.T(), "test-app", cfg.App.Name)
	assert.Equal(suite.T(), "1.0.0", cfg.App.Version)
	assert.Equal(suite.T(), "test", cfg.App.Environment)
	assert.True(suite.T(), cfg.App.Debug)

	// 验证服务器配置
	assert.Equal(suite.T(), "127.0.0.1", cfg.Server.Host)
	assert.Equal(suite.T(), 9999, cfg.Server.Port)
	assert.Equal(suite.T(), 10*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(suite.T(), 10*time.Second, cfg.Server.WriteTimeout)
	assert.Equal(suite.T(), 30*time.Second, cfg.Server.IdleTimeout)

	// 验证数据库配置
	assert.Equal(suite.T(), "sqlite", cfg.Database.Driver)
	assert.Equal(suite.T(), "test.db", cfg.Database.Database)
	assert.Equal(suite.T(), 2, cfg.Database.MaxIdleConns)
	assert.Equal(suite.T(), 5, cfg.Database.MaxOpenConns)
	assert.Equal(suite.T(), 30*time.Minute, cfg.Database.ConnMaxLifetime)

	// 验证Redis配置
	assert.Equal(suite.T(), "localhost", cfg.Redis.Host)
	assert.Equal(suite.T(), 6379, cfg.Redis.Port)
	assert.Equal(suite.T(), "", cfg.Redis.Password)
	assert.Equal(suite.T(), 15, cfg.Redis.Database)
	assert.Equal(suite.T(), 5, cfg.Redis.PoolSize)
	assert.Equal(suite.T(), 2, cfg.Redis.MinIdleConns)
	assert.Equal(suite.T(), 3*time.Second, cfg.Redis.DialTimeout)
	assert.Equal(suite.T(), 2*time.Second, cfg.Redis.ReadTimeout)
	assert.Equal(suite.T(), 2*time.Second, cfg.Redis.WriteTimeout)

	// 验证日志配置
	assert.Equal(suite.T(), "debug", cfg.Log.Level)
	assert.Equal(suite.T(), "console", cfg.Log.Format)
	assert.Equal(suite.T(), "stdout", cfg.Log.Output)
	assert.Equal(suite.T(), "test.log", cfg.Log.Filename)
	assert.Equal(suite.T(), 50, cfg.Log.MaxSize)
	assert.Equal(suite.T(), 3, cfg.Log.MaxAge)
	assert.Equal(suite.T(), 5, cfg.Log.MaxBackups)
	assert.False(suite.T(), cfg.Log.Compress)

	// 验证JWT配置
	assert.Equal(suite.T(), "test-secret-key", cfg.JWT.Secret)
	assert.Equal(suite.T(), "test-issuer", cfg.JWT.Issuer)
	assert.Equal(suite.T(), 1*time.Hour, cfg.JWT.ExpireTime)
}

// TestLoadInvalidConfig 测试加载无效配置
func (suite *ConfigTestSuite) TestLoadInvalidConfig() {
	configPath := filepath.Join(suite.testDataDir, "invalid_config.yml")

	cfg, err := config.Load(configPath)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), cfg)
	assert.Contains(suite.T(), err.Error(), "配置验证失败")
}

// TestLoadNonExistentConfig 测试加载不存在的配置文件
func (suite *ConfigTestSuite) TestLoadNonExistentConfig() {
	configPath := filepath.Join(suite.testDataDir, "non_existent.yml")

	cfg, err := config.Load(configPath)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), cfg)
	assert.Contains(suite.T(), err.Error(), "读取配置文件失败")
}

// TestEnvironmentVariableOverride 测试环境变量覆盖
func (suite *ConfigTestSuite) TestEnvironmentVariableOverride() {
	configPath := filepath.Join(suite.testDataDir, "test_config.yml")

	// 设置环境变量
	os.Setenv("BICO_APP_NAME", "env-override-app")
	os.Setenv("BICO_SERVER_PORT", "7777")
	os.Setenv("BICO_LOG_LEVEL", "error")
	os.Setenv("BICO_DATABASE_DRIVER", "sqlite") // 使用sqlite避免需要额外的host和username

	cfg, err := config.Load(configPath)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), cfg)

	// 验证环境变量覆盖生效
	assert.Equal(suite.T(), "env-override-app", cfg.App.Name)
	assert.Equal(suite.T(), 7777, cfg.Server.Port)
	assert.Equal(suite.T(), "error", cfg.Log.Level)
	assert.Equal(suite.T(), "sqlite", cfg.Database.Driver)

	// 验证未覆盖的配置保持原值
	assert.Equal(suite.T(), "1.0.0", cfg.App.Version)
	assert.Equal(suite.T(), "127.0.0.1", cfg.Server.Host)
}

// TestRunConfigTestSuite 运行配置测试套件
func TestRunConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

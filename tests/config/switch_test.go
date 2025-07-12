package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"bico-admin/pkg/config"
)

// ConfigSwitchTestSuite 配置切换测试套件
type ConfigSwitchTestSuite struct {
	suite.Suite
	testDataDir string
}

// SetupSuite 测试套件初始化
func (suite *ConfigSwitchTestSuite) SetupSuite() {
	suite.testDataDir = filepath.Join("testdata")
}

// TearDownTest 每个测试后清理环境变量
func (suite *ConfigSwitchTestSuite) TearDownTest() {
	// 清理所有可能的环境变量
	envVars := []string{
		"BICO_APP_NAME", "BICO_APP_VERSION", "BICO_APP_ENVIRONMENT", "BICO_APP_DEBUG",
		"BICO_SERVER_HOST", "BICO_SERVER_PORT",
		"BICO_DATABASE_DRIVER", "BICO_DATABASE_HOST", "BICO_DATABASE_USERNAME", "BICO_DATABASE_PASSWORD", "BICO_DATABASE_DATABASE",
		"BICO_LOG_LEVEL", "BICO_LOG_FORMAT", "BICO_LOG_OUTPUT", "BICO_LOG_FILENAME",
		"BICO_JWT_SECRET", "BICO_JWT_ISSUER", "BICO_JWT_EXPIRE_TIME",
	}

	for _, env := range envVars {
		os.Unsetenv(env)
	}
}

// TestSwitchBetweenConfigs 测试在不同配置文件之间切换
func (suite *ConfigSwitchTestSuite) TestSwitchBetweenConfigs() {
	// 加载第一个配置
	config1Path := filepath.Join(suite.testDataDir, "test_config.yml")
	cfg1, err := config.Load(config1Path)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), cfg1)

	// 验证第一个配置的值
	assert.Equal(suite.T(), "test-app", cfg1.App.Name)
	assert.Equal(suite.T(), 9999, cfg1.Server.Port)
	assert.Equal(suite.T(), "debug", cfg1.Log.Level)
	assert.Equal(suite.T(), "sqlite", cfg1.Database.Driver)

	// 加载第二个配置（覆盖配置）
	config2Path := filepath.Join(suite.testDataDir, "test_config_override.yml")
	cfg2, err := config.Load(config2Path)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), cfg2)

	// 验证第二个配置的值
	assert.Equal(suite.T(), "override-app", cfg2.App.Name)
	assert.Equal(suite.T(), 8888, cfg2.Server.Port)
	assert.Equal(suite.T(), "error", cfg2.Log.Level)
	assert.Equal(suite.T(), "mysql", cfg2.Database.Driver)

	// 验证配置确实不同
	assert.NotEqual(suite.T(), cfg1.App.Name, cfg2.App.Name)
	assert.NotEqual(suite.T(), cfg1.Server.Port, cfg2.Server.Port)
	assert.NotEqual(suite.T(), cfg1.Log.Level, cfg2.Log.Level)
	assert.NotEqual(suite.T(), cfg1.Database.Driver, cfg2.Database.Driver)
}

// TestEnvironmentOverrideWithDifferentConfigs 测试环境变量在不同配置文件中的覆盖
func (suite *ConfigSwitchTestSuite) TestEnvironmentOverrideWithDifferentConfigs() {
	// 设置环境变量
	os.Setenv("BICO_APP_NAME", "env-override")
	os.Setenv("BICO_SERVER_PORT", "7777")
	os.Setenv("BICO_LOG_LEVEL", "warn")

	// 加载第一个配置文件
	config1Path := filepath.Join(suite.testDataDir, "test_config.yml")
	cfg1, err := config.Load(config1Path)
	require.NoError(suite.T(), err)

	// 验证环境变量覆盖生效
	assert.Equal(suite.T(), "env-override", cfg1.App.Name)
	assert.Equal(suite.T(), 7777, cfg1.Server.Port)
	assert.Equal(suite.T(), "warn", cfg1.Log.Level)
	// 验证未覆盖的值保持配置文件中的值
	assert.Equal(suite.T(), "test", cfg1.App.Environment)
	assert.Equal(suite.T(), "sqlite", cfg1.Database.Driver)

	// 加载第二个配置文件
	config2Path := filepath.Join(suite.testDataDir, "test_config_override.yml")
	cfg2, err := config.Load(config2Path)
	require.NoError(suite.T(), err)

	// 验证环境变量覆盖在第二个配置中也生效
	assert.Equal(suite.T(), "env-override", cfg2.App.Name)
	assert.Equal(suite.T(), 7777, cfg2.Server.Port)
	assert.Equal(suite.T(), "warn", cfg2.Log.Level)
	// 验证未覆盖的值使用第二个配置文件中的值
	assert.Equal(suite.T(), "override", cfg2.App.Environment)
	assert.Equal(suite.T(), "mysql", cfg2.Database.Driver)
}

// TestPartialEnvironmentOverride 测试部分环境变量覆盖
func (suite *ConfigSwitchTestSuite) TestPartialEnvironmentOverride() {
	// 只设置部分环境变量
	os.Setenv("BICO_APP_NAME", "partial-override")
	os.Setenv("BICO_DATABASE_DRIVER", "sqlite") // 使用sqlite避免验证问题
	os.Setenv("BICO_DATABASE_DATABASE", "partial_test.db")

	configPath := filepath.Join(suite.testDataDir, "test_config.yml")
	cfg, err := config.Load(configPath)
	require.NoError(suite.T(), err)

	// 验证被覆盖的值
	assert.Equal(suite.T(), "partial-override", cfg.App.Name)
	assert.Equal(suite.T(), "sqlite", cfg.Database.Driver)
	assert.Equal(suite.T(), "partial_test.db", cfg.Database.Database)

	// 验证未被覆盖的值保持原配置文件中的值
	assert.Equal(suite.T(), "1.0.0", cfg.App.Version)
	assert.Equal(suite.T(), "test", cfg.App.Environment)
	assert.Equal(suite.T(), 9999, cfg.Server.Port)
	assert.Equal(suite.T(), "debug", cfg.Log.Level)
}

// TestComplexEnvironmentOverride 测试复杂环境变量覆盖
func (suite *ConfigSwitchTestSuite) TestComplexEnvironmentOverride() {
	// 设置多个层级的环境变量
	os.Setenv("BICO_APP_NAME", "complex-app")
	os.Setenv("BICO_APP_ENVIRONMENT", "production")
	os.Setenv("BICO_APP_DEBUG", "false")
	os.Setenv("BICO_SERVER_HOST", "0.0.0.0")
	os.Setenv("BICO_SERVER_PORT", "8080")
	os.Setenv("BICO_DATABASE_DRIVER", "sqlite")
	os.Setenv("BICO_DATABASE_DATABASE", "production.db")
	os.Setenv("BICO_LOG_LEVEL", "info")
	os.Setenv("BICO_LOG_FORMAT", "json")
	os.Setenv("BICO_LOG_OUTPUT", "file")
	os.Setenv("BICO_LOG_FILENAME", "production.log")
	os.Setenv("BICO_JWT_SECRET", "production-secret")
	os.Setenv("BICO_JWT_EXPIRE_TIME", "24h")

	configPath := filepath.Join(suite.testDataDir, "test_config.yml")
	cfg, err := config.Load(configPath)
	require.NoError(suite.T(), err)

	// 验证所有环境变量覆盖都生效
	assert.Equal(suite.T(), "complex-app", cfg.App.Name)
	assert.Equal(suite.T(), "production", cfg.App.Environment)
	assert.False(suite.T(), cfg.App.Debug)
	assert.Equal(suite.T(), "0.0.0.0", cfg.Server.Host)
	assert.Equal(suite.T(), 8080, cfg.Server.Port)
	assert.Equal(suite.T(), "sqlite", cfg.Database.Driver)
	assert.Equal(suite.T(), "production.db", cfg.Database.Database)
	assert.Equal(suite.T(), "info", cfg.Log.Level)
	assert.Equal(suite.T(), "json", cfg.Log.Format)
	assert.Equal(suite.T(), "file", cfg.Log.Output)
	assert.Equal(suite.T(), "production.log", cfg.Log.Filename)
	assert.Equal(suite.T(), "production-secret", cfg.JWT.Secret)
}

// TestConfigSwitchWithValidation 测试配置切换时的验证
func (suite *ConfigSwitchTestSuite) TestConfigSwitchWithValidation() {
	// 测试从有效配置切换到无效配置
	validConfigPath := filepath.Join(suite.testDataDir, "test_config.yml")
	invalidConfigPath := filepath.Join(suite.testDataDir, "invalid_config.yml")

	// 加载有效配置
	validCfg, err := config.Load(validConfigPath)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), validCfg)

	// 尝试加载无效配置
	invalidCfg, err := config.Load(invalidConfigPath)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), invalidCfg)
	assert.Contains(suite.T(), err.Error(), "配置验证失败")

	// 验证有效配置仍然可用（如果有全局配置的话）
	assert.Equal(suite.T(), "test-app", validCfg.App.Name)
}

// TestRunConfigSwitchTestSuite 运行配置切换测试套件
func TestRunConfigSwitchTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigSwitchTestSuite))
}

# 配置热更新

## 概述

配置热更新允许在不重启应用的情况下，自动检测配置文件变化并重新加载配置。这对于调整运行时参数（如日志级别、限流参数等）非常有用。

## ⚠️ 重要说明

**并非所有配置都支持热更新**。配置分为两类：

### 1. 静态配置（不支持热更新）
这些配置在应用启动时固定，修改后需要重启应用：
- 数据库连接配置 (`database`)
- 服务器端口配置 (`server.port`)
- JWT 密钥 (`jwt.secret`)
- 上传驱动配置 (`upload.driver`)

### 2. 动态配置（支持热更新）
这些配置可以通过 `ConfigManager` 动态获取，支持热更新：
- 限流参数 (`rate_limit.rps`, `rate_limit.burst`)
- 日志级别 (`log.level`) - 需要应用层实现
- 业务相关配置

## 使用方法

### 1. 在服务中使用 ConfigManager

如果你的服务需要访问可能会变化的配置，应该依赖 `ConfigManager` 而不是 `Config`：

```go
type MyService struct {
    configManager *config.ConfigManager
}

func NewMyService(cm *config.ConfigManager) *MyService {
    return &MyService{
        configManager: cm,
    }
}

func (s *MyService) DoSomething() {
    // 每次调用时获取最新配置
    cfg := s.configManager.GetConfig()
    
    // 使用配置
    if cfg.RateLimit.Enabled {
        // ...
    }
}
```

### 2. 在 Handler 中使用

```go
type MyHandler struct {
    configManager *config.ConfigManager
}

func (h *MyHandler) HandleRequest(c *gin.Context) {
    cfg := h.configManager.GetConfig()
    // 使用最新配置处理请求
}
```

## 支持的配置项

当前支持热更新的配置：
- ✅ 限流配置（`rate_limit`）
- ✅ 日志级别（`log.level`）
- ✅ 缓存配置（`cache`）
- ⚠️ 数据库连接（需重启）
- ⚠️ 服务器端口（需重启）

## 使用方法

### 1. 基础使用（自动热更新）

框架已内置配置热更新，修改 `config/config.yaml` 后自动生效：

```bash
# 1. 修改配置文件
vim config/config.yaml

# 2. 保存文件
# 配置会自动重新加载，无需重启服务
```

### 2. 使用 ConfigManager

如果需要在代码中获取实时配置：

```go
package main

import (
    "bico-admin/internal/core/config"
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewProduction()
    
    // 创建配置管理器
    cm, err := config.NewConfigManager("config/config.yaml", logger)
    if err != nil {
        logger.Fatal("加载配置失败", zap.Error(err))
    }
    
    // 获取当前配置
    cfg := cm.GetConfig()
    
    // 配置会自动更新，无需手动重新加载
    rateLimitCfg := cm.GetRateLimitConfig()
    logger.Info("限流配置", 
        zap.Bool("enabled", rateLimitCfg.Enabled),
        zap.Int("rps", rateLimitCfg.RPS),
    )
}
```

## 配置更新示例

### 动态调整限流

```yaml
# 修改前
rate_limit:
  enabled: true
  rps: 100
  burst: 200

# 修改后（立即生效）
rate_limit:
  enabled: true
  rps: 500     # 提升限流阈值
  burst: 1000  # 增加突发容量
```

### 动态调整日志级别

```yaml
# 线上排查问题时，临时调整日志级别
log:
  level: debug   # 从 info 改为 debug，立即生效
  format: json
  output: stdout

# 问题排查完成后，恢复
log:
  level: info
  format: json
  output: stdout
```

## 监控配置变更

### 查看日志

配置更新时会输出日志：

```json
{
  "level": "info",
  "time": "2025-11-15T14:30:00+08:00",
  "msg": "检测到配置文件变化",
  "file": "config/config.yaml"
}

{
  "level": "info",
  "time": "2025-11-15T14:30:00+08:00",
  "msg": "配置已热更新"
}
```

### 错误处理

如果配置文件格式错误，会保留旧配置：

```json
{
  "level": "error",
  "time": "2025-11-15T14:30:00+08:00",
  "msg": "重新加载配置失败",
  "error": "yaml: line 10: mapping values are not allowed in this context"
}
```

## 实现原理

### 文件监听

使用 `fsnotify` 库监听配置文件变化：

```go
viper.WatchConfig()
viper.OnConfigChange(func(e fsnotify.Event) {
    // 配置文件发生变化时触发
    // 自动重新加载配置
})
```

### 线程安全

使用读写锁保证并发安全：

```go
type ConfigManager struct {
    config *Config
    mu     sync.RWMutex  // 读写锁
}

// 读取配置
func (cm *ConfigManager) GetConfig() *Config {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    return cm.config
}

// 更新配置
func (cm *ConfigManager) onConfigChange(e fsnotify.Event) {
    cm.mu.Lock()
    cm.config = newConfig
    cm.mu.Unlock()
}
```

## 最佳实践

### 1. 配置备份

修改前先备份：

```bash
cp config/config.yaml config/config.yaml.backup
```

### 2. 验证配置格式

使用 YAML 验证工具：

```bash
# 安装 yamllint
pip install yamllint

# 验证配置
yamllint config/config.yaml
```

### 3. 逐步调整

生产环境建议小步快跑：

```yaml
# 第一步：小幅调整
rate_limit:
  rps: 120  # 从 100 增加到 120

# 观察效果后，再继续调整
rate_limit:
  rps: 150
```

### 4. 监控影响

配置更新后，监控关键指标：
- 请求响应时间
- 错误率
- 系统资源使用

## 限制和注意事项

### 不支持热更新的配置

以下配置修改后需要重启服务：

1. **服务器端口**
   ```yaml
   server:
     port: 8080  # 修改需重启
   ```

2. **数据库连接**
   ```yaml
   database:
     driver: mysql  # 修改需重启
     host: localhost
   ```

3. **嵌入静态文件**
   ```yaml
   server:
     embed_static: true  # 修改需重启
   ```

### 多实例部署

如果部署多个实例，需要同步更新所有实例的配置文件：

```bash
# 使用 ansible 批量更新
ansible all -m copy -a "src=config/config.yaml dest=/app/config/"

# 或使用配置中心（如 etcd、consul）
```

### 配置回滚

如果发现配置错误：

```bash
# 方案1: 恢复备份
cp config/config.yaml.backup config/config.yaml

# 方案2: Git 回滚
git checkout config/config.yaml

# 方案3: 直接编辑修正
vim config/config.yaml
```

## 高级用法

### 配置变更回调

```go
// 自定义配置变更处理
func onRateLimitConfigChange(oldCfg, newCfg RateLimitConfig) {
    if oldCfg.RPS != newCfg.RPS {
        logger.Info("限流阈值变更",
            zap.Int("old_rps", oldCfg.RPS),
            zap.Int("new_rps", newCfg.RPS),
        )
        
        // 更新限流器
        rateLimiter.UpdateConfig(newCfg.RPS, newCfg.Burst)
    }
}
```

### 配置验证

```go
// 验证配置合法性
func validateConfig(cfg *Config) error {
    if cfg.RateLimit.Enabled && cfg.RateLimit.RPS <= 0 {
        return errors.New("限流已启用，但 RPS 配置无效")
    }
    if cfg.RateLimit.Burst < cfg.RateLimit.RPS {
        return errors.New("突发容量应大于等于 RPS")
    }
    return nil
}
```

## 故障排查

### 问题：配置修改不生效

**排查步骤**：
1. 检查配置文件路径是否正确
2. 查看日志是否有 "配置已热更新" 消息
3. 确认修改的配置项支持热更新
4. 检查 YAML 格式是否正确

### 问题：频繁触发配置重载

**原因**：某些编辑器保存时会触发多次文件变更事件

**解决**：
- 正常现象，不影响使用
- 如需优化，可添加防抖逻辑

### 问题：配置回滚失败

**原因**：配置文件权限或格式错误

**解决**：
```bash
# 检查文件权限
ls -l config/config.yaml

# 验证 YAML 格式
yamllint config/config.yaml

# 强制恢复默认配置
git checkout config/config.yaml
```

## 性能影响

- **CPU 开销**：监听文件变化几乎无开销
- **内存开销**：单个配置结构体 < 1KB
- **更新延迟**：< 100ms（文件系统事件响应时间）

## 安全考虑

1. **配置文件权限**
   ```bash
   chmod 600 config/config.yaml  # 仅所有者可读写
   ```

2. **敏感信息管理**
   - JWT Secret
   - 数据库密码
   - API 密钥
   
   建议使用环境变量或密钥管理系统。

3. **审计日志**
   - 记录配置变更时间
   - 记录变更内容
   - 便于追溯问题

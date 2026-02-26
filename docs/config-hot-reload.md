# 配置热更新说明

## 当前实现

项目通过 `ConfigManager` 监听配置文件变化并更新内存配置对象：

- `viper.WatchConfig()`
- `viper.OnConfigChange(...)`

对应代码：`internal/core/config/config.go`。

## 重要边界

“配置对象已更新”不等于“运行中的基础设施已自动重建”。

当前内置模块中：

1. **可动态读取到新值**

- `GET /admin-api/app-config`
- 原因：`ConfigService.GetAppConfig()` 每次请求都从 `ConfigManager` 取值。

2. **修改后通常需要重启服务**

- `server.*`（端口、静态资源路由前缀等）
- `database.*`
- `cache.*`
- `jwt.*`
- `upload.*`
- `log.*`
- `rate_limit.*`（限流器在启动时创建）

## 推荐操作

### 需要实时观察配置值时

在业务代码中依赖 `ConfigManager`，每次执行时读取最新配置：

```go
cfg := cm.GetConfig()
```

### 修改基础设施相关配置时

直接重启服务：

```bash
# 本地
make serve

# Docker
docker-compose restart app
```

## 运行期日志

配置文件变化时会看到：

- `检测到配置文件变化`
- `配置已热更新`

如果 YAML 非法，会记录 `重新加载配置失败`，并继续保留旧配置。

# 日志功能文档

## 📖 概述

项目已集成 **Uber Zap** 日志库，提供高性能、结构化的日志功能，并在 debug 模式下自动输出 GORM 的 SQL 执行日志。

## 🔧 配置说明

在 `config/config.yaml` 中配置日志参数：

```yaml
log:
  level: debug        # 日志级别: debug / info / warn / error / fatal
  format: console     # 日志格式: console(彩色易读) / json(适合收集)
  output: stdout      # 输出位置: stdout / stderr / 文件路径(如: logs/app.log)

server:
  mode: debug         # debug 模式下会输出 SQL 语句
```

## ✨ 核心功能

### 1. 日志级别
- **debug**: 调试信息（开发环境推荐）
- **info**: 常规信息
- **warn**: 警告信息
- **error**: 错误信息
- **fatal**: 致命错误（会终止程序）

### 2. 日志格式
- **console**: 彩色输出，方便开发时阅读
- **json**: 结构化输出，便于日志收集系统解析

### 3. SQL 输出
- **debug 模式**: 自动输出所有 SQL 执行语句
- **release 模式**: 仅输出警告和错误
- **慢查询检测**: 超过 200ms 的查询会被标记为慢查询

## 📝 使用示例

### 在 Handler 或 Service 中使用

```go
package handler

import (
    "bico-admin/internal/core/logger"
    "go.uber.org/zap"
)

func (h *Handler) SomeMethod() {
    // 简单日志
    logger.Info("用户登录成功")
    
    // 带字段的结构化日志
    logger.Info("用户登录", 
        zap.String("username", "admin"),
        zap.Int("user_id", 123),
        zap.String("ip", "192.168.1.1"),
    )
    
    // 错误日志
    logger.Error("操作失败", 
        zap.Error(err),
        zap.String("operation", "update_user"),
    )
    
    // 调试日志
    logger.Debug("调试信息", 
        zap.Any("data", someData),
    )
}
```

### 在 DI 中注入使用

如果需要在特定组件中使用 logger，可以通过 DI 注入：

```go
type SomeService struct {
    logger *zap.Logger
    db     *gorm.DB
}

func NewSomeService(logger *zap.Logger, db *gorm.DB) *SomeService {
    return &SomeService{
        logger: logger,
        db:     db,
    }
}

func (s *SomeService) DoSomething() {
    s.logger.Info("执行某操作")
}
```

## 🎯 SQL 日志示例

### Debug 模式输出示例

```
2024-10-18T22:58:30.123+0800    DEBUG   SQL 执行
    {"耗时": "2.3ms", "影响行数": 1, "SQL": "SELECT * FROM `admin_user` WHERE `username` = 'admin' LIMIT 1"}

2024-10-18T22:58:30.456+0800    WARN    慢查询检测
    {"耗时": "350ms", "阈值": "200ms", "影响行数": 100, "SQL": "SELECT * FROM `admin_user`"}
```

### 生产模式
- 不输出常规 SQL
- 仅输出慢查询警告和错误日志

## 🔍 日志字段说明

常用的 zap 字段类型：

```go
zap.String("key", "value")      // 字符串
zap.Int("key", 123)             // 整数
zap.Int64("key", 123456789)     // 长整数
zap.Bool("key", true)           // 布尔值
zap.Float64("key", 3.14)        // 浮点数
zap.Duration("key", duration)   // 时间段
zap.Time("key", time.Now())     // 时间
zap.Error(err)                  // 错误（自动使用 "error" 作为 key）
zap.Any("key", interface{})     // 任意类型（会自动序列化）
```

## 🚀 最佳实践

1. **开发环境**: 使用 `console` 格式 + `debug` 级别
2. **生产环境**: 使用 `json` 格式 + `info` 级别 + 日志文件
3. **结构化日志**: 优先使用 zap 字段而非格式化字符串
4. **避免敏感信息**: 不要记录密码、token 等敏感数据
5. **适度记录**: 避免在高频循环中大量打日志

## 📦 项目结构

```
internal/core/logger/
├── logger.go         # 主日志实现
└── gorm_logger.go    # GORM 日志适配器
```

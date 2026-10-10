# 定时任务文档

## 📖 概述

项目使用 **robfig/cron** 库实现定时任务调度，支持秒级精度的 cron 表达式，所有任务在应用启动时自动注册并启动。

## 🏗️ 架构设计

```
internal/core/scheduler/
└── scheduler.go      # 调度器封装（封装 cron 库）

internal/job/
├── module.go         # 模块入口（把任务注册到 ctx.Scheduler）
├── register.go       # 任务注册器（注册所有定时任务）
└── task/            # 任务实现目录
    ├── clean.go      # 清理任务（示例）
    └── sync.go       # 同步任务（示例，默认不注册）
```

### 核心组件

#### 1. **Scheduler（调度器）**
- 封装 `robfig/cron` 库
- 提供任务注册、启动、停止功能
- 集成日志记录（任务开始、成功、失败）

#### 2. **Task（任务接口）**
```go
type Task interface {
    Run() error
}
```
所有任务必须实现此接口。

#### 3. **RegisterJobs（任务注册器）**
- 统一管理所有任务的注册
- 在应用启动时自动调用

## 📝 Cron 表达式说明

项目使用**6位 cron 表达式**（支持秒级调度）：

```
格式: 秒 分 时 日 月 周

字段说明:
┌──────────── 秒 (0-59)
│ ┌────────── 分钟 (0-59)
│ │ ┌──────── 小时 (0-23)
│ │ │ ┌────── 日 (1-31)
│ │ │ │ ┌──── 月 (1-12)
│ │ │ │ │ ┌── 周 (0-6, 0=周日)
│ │ │ │ │ │
* * * * * *
```

### 常用表达式示例

| 表达式 | 说明 |
|--------|------|
| `0 0 3 * * *` | 每天凌晨 3:00:00 执行 |
| `0 */30 * * * *` | 每 30 分钟执行一次 |
| `0 0 * * * *` | 每小时整点执行 |
| `0 0 0 * * 1` | 每周一午夜执行 |
| `0 0 2 1 * *` | 每月 1 号凌晨 2 点执行 |
| `*/10 * * * * *` | 每 10 秒执行一次 |
| `0 0 9-17 * * 1-5` | 工作日 9:00-17:00 每小时执行 |

### 特殊字符

- `*` - 匹配任意值
- `,` - 列举多个值，如 `0,15,30,45`
- `-` - 范围，如 `9-17`
- `/` - 步长，如 `*/5` 表示每 5 个单位

## 🚀 使用指南

### 1. 创建新任务

#### 步骤 1：创建任务文件

在 `internal/job/task/` 下创建新任务：

```go
// internal/job/task/report.go
package task

import (
    "go.uber.org/zap"
    "gorm.io/gorm"
)

// ReportTask 报表生成任务
type ReportTask struct {
    db     *gorm.DB
    logger *zap.Logger
}

// NewReportTask 创建报表任务
func NewReportTask(db *gorm.DB, logger *zap.Logger) *ReportTask {
    return &ReportTask{
        db:     db,
        logger: logger,
    }
}

// Run 执行报表生成
func (t *ReportTask) Run() error {
    t.logger.Info("开始生成日报")
    
    // 实现报表生成逻辑
    // ...
    
    t.logger.Info("日报生成完成")
    return nil
}
```

#### 步骤 2：注册任务

在 `internal/job/register.go` 中注册：

```go
func RegisterJobs(scheduler *scheduler.Scheduler, db *gorm.DB, cache cache.Cache, logger *zap.Logger) error {
    // ... 现有任务 ...
    
    // 注册报表任务（每天早上 8 点执行）
    reportTask := task.NewReportTask(db, logger)
    if err := scheduler.AddTask("0 0 8 * * *", reportTask, "ReportTask"); err != nil {
        return err
    }
    
    return nil
}
```

### 2. 任务依赖注入

任务可以注入所需的依赖：

```go
type MyTask struct {
    db     *gorm.DB        // 数据库
    cache  cache.Cache     // 缓存
    logger *zap.Logger     // 日志
    // 可以注入任何你需要的依赖
}

func NewMyTask(db *gorm.DB, cache cache.Cache, logger *zap.Logger) *MyTask {
    return &MyTask{
        db:     db,
        cache:  cache,
        logger: logger,
    }
}
```

### 3. 错误处理

任务执行失败时会自动记录错误日志：

```go
func (t *MyTask) Run() error {
    // 返回错误会被调度器捕获并记录
    if err := t.doSomething(); err != nil {
        return fmt.Errorf("执行失败: %w", err)
    }
    return nil
}
```

## 📊 现有任务说明

### CleanTask（清理任务）

- **执行时间**: 每天凌晨 3:00
- **Cron 表达式**: `0 0 3 * * *`
- **功能**: 清理过期数据（如缓存、临时文件等）
- **文件**: `internal/job/task/clean.go`

### SyncTask（同步任务，默认禁用）

- **状态**: 仅作为示例，`internal/job/register.go` 中的注册代码默认被注释
- **启用方式**: 取消注释后即按每小时执行（Cron 表达式 `0 0 * * * *`）
- **功能**: 示例逻辑为统计用户数量，实际使用时替换为同步数据、更新统计信息
- **文件**: `internal/job/task/sync.go`

## 🔧 高级用法

### 1. 动态任务

根据配置动态启用/禁用任务：

```go
func RegisterJobs(...) error {
    // 根据配置决定是否启用
    if cfg.Job.EnableCleanTask {
        cleanTask := task.NewCleanTask(cache, logger)
        scheduler.AddTask("0 0 3 * * *", cleanTask, "CleanTask")
    }
    
    return nil
}
```

### 2. 任务链

通过组合多个小任务实现复杂流程：

```go
type CompositeTask struct {
    tasks []Task
}

func (t *CompositeTask) Run() error {
    for _, task := range t.tasks {
        if err := task.Run(); err != nil {
            return err
        }
    }
    return nil
}
```

### 3. 并发控制

对于耗时任务，可以添加互斥锁防止重复执行：

```go
type LongRunningTask struct {
    mu     sync.Mutex
    logger *zap.Logger
}

func (t *LongRunningTask) Run() error {
    if !t.mu.TryLock() {
        t.logger.Warn("任务正在执行中，跳过本次调度")
        return nil
    }
    defer t.mu.Unlock()
    
    // 执行耗时操作
    return nil
}
```

## 📈 监控与日志

### 日志输出

所有任务执行都会产生日志：

```json
{
  "level": "info",
  "time": "2025-10-18T23:00:00+08:00",
  "msg": "定时任务开始执行",
  "task": "CleanTask"
}

{
  "level": "info",
  "time": "2025-10-18T23:00:01+08:00",
  "msg": "定时任务执行成功",
  "task": "CleanTask"
}
```

### 错误日志

任务失败时：

```json
{
  "level": "error",
  "time": "2025-10-18T23:00:00+08:00",
  "msg": "定时任务执行失败",
  "task": "SyncTask",
  "error": "database connection lost"
}
```

## 🎯 最佳实践

### 1. 任务命名

- ✅ 使用清晰的任务名：`CleanExpiredCacheTask`
- ❌ 避免模糊命名：`Task1`, `MyTask`

### 2. 执行时间

- ✅ 避开业务高峰期
- ✅ 资源密集型任务在凌晨执行
- ❌ 避免多个大任务同时执行

### 3. 幂等性

确保任务可以安全地重复执行：

```go
func (t *ImportTask) Run() error {
    // ✅ 使用事务确保原子性
    return t.db.Transaction(func(tx *gorm.DB) error {
        // 导入逻辑
        return nil
    })
}
```

### 4. 超时控制

对于可能超时的任务，添加 context 控制：

```go
func (t *LongTask) Run() error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
    defer cancel()
    
    return t.doWorkWithContext(ctx)
}
```

### 5. 错误恢复

对于关键任务，可以添加重试机制：

```go
func (t *ImportantTask) Run() error {
    var err error
    for i := 0; i < 3; i++ {
        if err = t.execute(); err == nil {
            return nil
        }
        time.Sleep(time.Second * time.Duration(i+1))
    }
    return err
}
```

## 🐛 调试技巧

### 1. 测试单个任务

创建测试函数直接运行任务：

```go
func TestCleanTask(t *testing.T) {
    task := task.NewCleanTask(cache, logger)
    err := task.Run()
    assert.NoError(t, err)
}
```

### 2. 临时调整执行时间

开发时可以调整为每分钟执行：

```go
// 开发环境
scheduler.AddTask("0 * * * * *", task, "TestTask")

// 生产环境
scheduler.AddTask("0 0 3 * * *", task, "TestTask")
```

### 3. 手动触发任务

可以创建 API 接口手动触发任务：

```go
func (h *JobHandler) TriggerTask(c *gin.Context) {
    taskName := c.Param("name")
    // 查找并执行任务
    // ...
}
```

## 🔄 生命周期

1. **应用启动** → 创建 Scheduler
2. **注册任务** → RegisterJobs 被调用
3. **启动调度** → scheduler.Start()
4. **定时执行** → 按 cron 表达式触发
5. **应用关闭** → scheduler.Stop()

## 📚 参考资源

- [robfig/cron 文档](https://github.com/robfig/cron)
- [Cron 表达式在线生成器](https://crontab.guru/)
- [Go Zap 日志库](https://github.com/uber-go/zap)

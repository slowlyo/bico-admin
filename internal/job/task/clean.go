package task

import (
	"time"

	"bico-admin/internal/core/cache"
	"go.uber.org/zap"
)

// CleanTask 清理任务（清理过期的 token 黑名单）
type CleanTask struct {
	cache  cache.Cache
	logger *zap.Logger
}

// NewCleanTask 创建清理任务
func NewCleanTask(cache cache.Cache, logger *zap.Logger) *CleanTask {
	return &CleanTask{
		cache:  cache,
		logger: logger,
	}
}

// Run 执行清理任务
func (t *CleanTask) Run() error {
	t.logger.Info("开始清理过期数据")
	
	// 清理逻辑示例：清理过期的缓存数据
	// 注意：实际的清理逻辑取决于你的业务需求
	// 这里只是一个示例，展示如何访问缓存
	
	// 示例：删除特定前缀的过期键
	// 注意：Memory 缓存会自动清理过期数据，这里只是演示
	
	cleanedCount := 0
	// 实际应用中可以遍历并清理特定模式的键
	// 例如：blacklist:* 等
	
	t.logger.Info("清理任务完成",
		zap.Int("cleaned_count", cleanedCount),
		zap.Time("executed_at", time.Now()),
	)
	
	return nil
}

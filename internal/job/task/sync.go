package task

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SyncTask 同步任务（更新统计数据）
type SyncTask struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewSyncTask 创建同步任务
func NewSyncTask(db *gorm.DB, logger *zap.Logger) *SyncTask {
	return &SyncTask{
		db:     db,
		logger: logger,
	}
}

// Run 执行同步任务
func (t *SyncTask) Run() error {
	t.logger.Info("开始同步数据")
	
	// 同步逻辑示例：更新统计数据
	// 实际应用中可以：
	// 1. 同步外部 API 数据
	// 2. 更新统计信息
	// 3. 生成报表数据
	// 4. 同步用户状态等
	
	// 示例：统计用户数量
	var totalUsers int64
	if err := t.db.Table("admin_users").Count(&totalUsers).Error; err != nil {
		t.logger.Error("统计用户数量失败", zap.Error(err))
		return err
	}
	
	t.logger.Info("同步任务完成",
		zap.Int64("total_users", totalUsers),
		zap.Time("executed_at", time.Now()),
	)
	
	return nil
}

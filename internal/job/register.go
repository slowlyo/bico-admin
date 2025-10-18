package job

import (
	"bico-admin/internal/core/cache"
	"bico-admin/internal/job/task"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RegisterJobs 注册所有定时任务
func RegisterJobs(scheduler *Scheduler, db *gorm.DB, cache cache.Cache, logger *zap.Logger) error {
	// 注意：cron 表达式格式：秒 分 时 日 月 周
	// 例如：
	// "0 0 2 * * *"  - 每天凌晘02:00:00执行
	// "0 */30 * * * *" - 每30分钟执行一次
	// "0 0 0 * * 0"  - 每周日午夜执行
	
	// 注册清理任务（每天凌晨 3 点执行）
	cleanTask := task.NewCleanTask(cache, logger)
	if err := scheduler.AddTask("0 0 3 * * *", cleanTask, "CleanTask"); err != nil {
		return err
	}
	
	// 注册同步任务（每小时执行一次）
	syncTask := task.NewSyncTask(db, logger)
	if err := scheduler.AddTask("0 0 * * * *", syncTask, "SyncTask"); err != nil {
		return err
	}
	
	// 可以继续注册更多任务...
	
	return nil
}

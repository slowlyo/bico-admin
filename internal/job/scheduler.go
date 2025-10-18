package job

import (
	"fmt"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Task 任务接口
type Task interface {
	Run() error
}

// Scheduler 定时任务调度器
type Scheduler struct {
	cron   *cron.Cron
	logger *zap.Logger
}

// NewScheduler 创建调度器
func NewScheduler(logger *zap.Logger) *Scheduler {
	return &Scheduler{
		cron: cron.New(cron.WithSeconds()), // 支持秒级调度
		logger: logger,
	}
}

// AddTask 添加定时任务
func (s *Scheduler) AddTask(spec string, task Task, name string) error {
	_, err := s.cron.AddFunc(spec, func() {
		s.logger.Info("定时任务开始执行", zap.String("task", name))
		
		if err := task.Run(); err != nil {
			s.logger.Error("定时任务执行失败",
				zap.String("task", name),
				zap.Error(err),
			)
		} else {
			s.logger.Info("定时任务执行成功", zap.String("task", name))
		}
	})
	
	if err != nil {
		return fmt.Errorf("添加任务失败: %w", err)
	}
	
	s.logger.Info("定时任务已注册",
		zap.String("task", name),
		zap.String("schedule", spec),
	)
	
	return nil
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.cron.Start()
	s.logger.Info("定时任务调度器已启动")
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.cron.Stop()
	s.logger.Info("定时任务调度器已停止")
}

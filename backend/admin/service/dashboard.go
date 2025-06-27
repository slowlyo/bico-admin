package service

import (
	"gorm.io/gorm"

	"bico-admin/core/model"
)

// DashboardService 仪表板服务接口
type DashboardService interface {
	GetDashboardData() (map[string]interface{}, error)
	GetStats() (map[string]interface{}, error)
}

// DashboardServiceImpl 仪表板服务实现
type DashboardServiceImpl struct {
	db *gorm.DB
}

// NewDashboardService 创建仪表板服务实例
func NewDashboardService(db *gorm.DB) DashboardService {
	return &DashboardServiceImpl{
		db: db,
	}
}

// GetDashboardData 获取仪表板数据
func (s *DashboardServiceImpl) GetDashboardData() (map[string]interface{}, error) {
	// 获取基本统计数据
	var userCount, roleCount, permissionCount int64
	
	s.db.Model(&model.User{}).Count(&userCount)
	s.db.Model(&model.Role{}).Count(&roleCount)
	s.db.Model(&model.Permission{}).Count(&permissionCount)

	// 获取最近用户
	var recentUsers []model.User
	s.db.Order("created_at desc").Limit(5).Find(&recentUsers)

	data := map[string]interface{}{
		"stats": map[string]interface{}{
			"total_users":       userCount,
			"total_roles":       roleCount,
			"total_permissions": permissionCount,
		},
		"recent_users": recentUsers,
		"system_info": map[string]interface{}{
			"version": "1.0.0",
			"uptime":  "24h",
		},
	}

	return data, nil
}

// GetStats 获取统计数据
func (s *DashboardServiceImpl) GetStats() (map[string]interface{}, error) {
	var userCount, activeUserCount, roleCount, permissionCount int64
	
	// 总用户数
	s.db.Model(&model.User{}).Count(&userCount)
	
	// 活跃用户数
	s.db.Model(&model.User{}).Where("status = ?", model.UserStatusActive).Count(&activeUserCount)
	
	// 角色数
	s.db.Model(&model.Role{}).Count(&roleCount)
	
	// 权限数
	s.db.Model(&model.Permission{}).Count(&permissionCount)

	stats := map[string]interface{}{
		"users": map[string]interface{}{
			"total":  userCount,
			"active": activeUserCount,
		},
		"roles":       roleCount,
		"permissions": permissionCount,
		"growth": map[string]interface{}{
			"users_this_month": 0, // TODO: 实现月度增长统计
			"users_last_month": 0,
		},
	}

	return stats, nil
}

package service

import (
	"gorm.io/gorm"

	"bico-admin/core/model"
	"bico-admin/core/repository"
)

// AppService 应用服务接口
type AppService interface {
	GetAppInfo() (map[string]interface{}, error)
	GetAppConfig() (map[string]interface{}, error)
	GetUserProfile(userID uint) (*model.UserResponse, error)
	UpdateUserProfile(userID uint, data map[string]interface{}) (*model.UserResponse, error)
	GetContentList(page, pageSize int, category string) (map[string]interface{}, error)
	GetContent(id uint) (map[string]interface{}, error)
}

// AppServiceImpl 应用服务实现
type AppServiceImpl struct {
	db       *gorm.DB
	userRepo repository.UserRepository
}

// NewAppService 创建应用服务实例
func NewAppService(db *gorm.DB) AppService {
	return &AppServiceImpl{
		db:       db,
		userRepo: repository.NewUserRepository(db),
	}
}

// GetAppInfo 获取应用信息
func (s *AppServiceImpl) GetAppInfo() (map[string]interface{}, error) {
	info := map[string]interface{}{
		"name":        "Bico Admin",
		"version":     "1.0.0",
		"description": "AI友好的管理后台框架",
		"author":      "Bico Team",
		"website":     "https://github.com/your-username/bico-admin",
	}

	return info, nil
}

// GetAppConfig 获取应用配置
func (s *AppServiceImpl) GetAppConfig() (map[string]interface{}, error) {
	config := map[string]interface{}{
		"upload": map[string]interface{}{
			"max_size":     "10MB",
			"allowed_types": []string{"jpg", "jpeg", "png", "gif", "pdf", "doc", "docx"},
		},
		"pagination": map[string]interface{}{
			"default_page_size": 10,
			"max_page_size":     100,
		},
		"features": map[string]interface{}{
			"user_registration": true,
			"email_verification": false,
			"two_factor_auth":   false,
		},
	}

	return config, nil
}

// GetUserProfile 获取用户资料
func (s *AppServiceImpl) GetUserProfile(userID uint) (*model.UserResponse, error) {
	user, err := s.userRepo.GetWithRoles(userID)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// UpdateUserProfile 更新用户资料
func (s *AppServiceImpl) UpdateUserProfile(userID uint, data map[string]interface{}) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// 更新允许的字段
	if nickname, ok := data["nickname"].(string); ok {
		user.Nickname = nickname
	}
	if phone, ok := data["phone"].(string); ok {
		user.Phone = phone
	}
	if avatar, ok := data["avatar"].(string); ok {
		user.Avatar = avatar
	}

	// 保存更新
	if err := s.userRepo.Update(userID, user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// GetContentList 获取内容列表
func (s *AppServiceImpl) GetContentList(page, pageSize int, category string) (map[string]interface{}, error) {
	// TODO: 实现内容列表获取逻辑
	// 这里返回模拟数据
	content := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"id":       1,
				"title":    "示例内容1",
				"content":  "这是示例内容的描述",
				"category": "技术",
				"status":   "published",
			},
			{
				"id":       2,
				"title":    "示例内容2",
				"content":  "这是另一个示例内容的描述",
				"category": "生活",
				"status":   "published",
			},
		},
		"pagination": map[string]interface{}{
			"page":        page,
			"page_size":   pageSize,
			"total":       2,
			"total_pages": 1,
		},
	}

	return content, nil
}

// GetContent 获取单个内容
func (s *AppServiceImpl) GetContent(id uint) (map[string]interface{}, error) {
	// TODO: 实现单个内容获取逻辑
	// 这里返回模拟数据
	content := map[string]interface{}{
		"id":         id,
		"title":      "示例内容",
		"content":    "这是详细的内容描述...",
		"category":   "技术",
		"status":     "published",
		"created_at": "2024-01-01T00:00:00Z",
		"updated_at": "2024-01-01T00:00:00Z",
	}

	return content, nil
}

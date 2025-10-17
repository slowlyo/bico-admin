package model

import "bico-admin/internal/shared/model"

// AdminUser 后台用户模型
type AdminUser struct {
	model.BaseModel
	Username string `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`
	Name     string `gorm:"size:64" json:"name"`
	Avatar   string `gorm:"size:255" json:"avatar"`
	Enabled  bool   `gorm:"default:true" json:"enabled"`
}

// TableName 指定表名
func (AdminUser) TableName() string {
	return "admin_users"
}

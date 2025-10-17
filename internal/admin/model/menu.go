package model

import "bico-admin/internal/shared/model"

// Menu 菜单模型
type Menu struct {
	model.BaseModel
	ParentID uint   `gorm:"default:0;index" json:"parent_id"`
	Name     string `gorm:"size:64;not null" json:"name"`
	Path     string `gorm:"size:255" json:"path"`
	Icon     string `gorm:"size:128" json:"icon"`
	Sort     int    `gorm:"default:0" json:"sort"`
	Type     int8   `gorm:"default:1;comment:类型 1菜单 2按钮" json:"type"`
	Status   int8   `gorm:"default:1;comment:状态 1正常 0禁用" json:"status"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}

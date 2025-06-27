package model

import (
	"bico-admin/core/model"
)

// Category 分类模型
type Category struct {
	model.BaseModel
	Name        string          `json:"name" gorm:"uniqueIndex;size:100;not null" validate:"required,max=100"`
	Code        string          `json:"code" gorm:"uniqueIndex;size:50;not null" validate:"required,max=50"`
	Description string          `json:"description" gorm:"size:255"`
	ParentID    *uint           `json:"parent_id" gorm:"index"`
	Sort        int             `json:"sort" gorm:"default:0"`
	Status      CategoryStatus  `json:"status" gorm:"default:1"`
	
	// 关联关系
	Parent   *Category  `json:"parent" gorm:"foreignKey:ParentID"`
	Children []Category `json:"children" gorm:"foreignKey:ParentID"`
	Contents []Content  `json:"contents" gorm:"foreignKey:CategoryID"`
}

// CategoryStatus 分类状态
type CategoryStatus int

const (
	CategoryStatusInactive CategoryStatus = 0 // 未激活
	CategoryStatusActive   CategoryStatus = 1 // 激活
)

// CategoryCreateRequest 创建分类请求
type CategoryCreateRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Code        string `json:"code" validate:"required,max=50"`
	Description string `json:"description" validate:"max=255"`
	ParentID    *uint  `json:"parent_id"`
	Sort        int    `json:"sort"`
}

// CategoryUpdateRequest 更新分类请求
type CategoryUpdateRequest struct {
	Name        string         `json:"name" validate:"max=100"`
	Code        string         `json:"code" validate:"max=50"`
	Description string         `json:"description" validate:"max=255"`
	ParentID    *uint          `json:"parent_id"`
	Sort        int            `json:"sort"`
	Status      CategoryStatus `json:"status" validate:"oneof=0 1"`
}

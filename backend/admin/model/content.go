package model

import (
	"bico-admin/core/model"
)

// Content 内容模型
type Content struct {
	model.BaseModel
	Title       string        `json:"title" gorm:"size:255;not null" validate:"required,max=255"`
	Content     string        `json:"content" gorm:"type:text"`
	Summary     string        `json:"summary" gorm:"size:500"`
	CategoryID  uint          `json:"category_id" gorm:"index"`
	AuthorID    uint          `json:"author_id" gorm:"index"`
	Status      ContentStatus `json:"status" gorm:"default:1"`
	ViewCount   int           `json:"view_count" gorm:"default:0"`
	LikeCount   int           `json:"like_count" gorm:"default:0"`
	CommentCount int          `json:"comment_count" gorm:"default:0"`
	
	// 关联关系
	Category *Category    `json:"category" gorm:"foreignKey:CategoryID"`
	Author   *model.User  `json:"author" gorm:"foreignKey:AuthorID"`
}

// ContentStatus 内容状态
type ContentStatus int

const (
	ContentStatusDraft     ContentStatus = 1 // 草稿
	ContentStatusPublished ContentStatus = 2 // 已发布
	ContentStatusArchived  ContentStatus = 3 // 已归档
)

// ContentCreateRequest 创建内容请求
type ContentCreateRequest struct {
	Title      string `json:"title" validate:"required,max=255"`
	Content    string `json:"content" validate:"required"`
	Summary    string `json:"summary" validate:"max=500"`
	CategoryID uint   `json:"category_id" validate:"required"`
	Status     ContentStatus `json:"status" validate:"oneof=1 2 3"`
}

// ContentUpdateRequest 更新内容请求
type ContentUpdateRequest struct {
	Title      string        `json:"title" validate:"max=255"`
	Content    string        `json:"content"`
	Summary    string        `json:"summary" validate:"max=500"`
	CategoryID uint          `json:"category_id"`
	Status     ContentStatus `json:"status" validate:"oneof=1 2 3"`
}

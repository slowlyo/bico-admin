package model

// BaseModel 公共字段
type BaseModel struct {
	ID        uint     `gorm:"primarykey" json:"id"`
	CreatedAt JSONTime `json:"created_at"`
	UpdatedAt JSONTime `json:"updated_at"`
}

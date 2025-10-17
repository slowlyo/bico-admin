package model

// Role 角色模型
type Role struct {
	BaseModel
	Name   string `gorm:"size:64;uniqueIndex;not null" json:"name"`
	Code   string `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Remark string `gorm:"size:255" json:"remark"`
	Status int8   `gorm:"default:1;comment:状态 1正常 0禁用" json:"status"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

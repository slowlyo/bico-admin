package model

// User 用户模型
type User struct {
	BaseModel
	Username string `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`
	Nickname string `gorm:"size:64" json:"nickname"`
	Email    string `gorm:"size:128" json:"email"`
	Phone    string `gorm:"size:32" json:"phone"`
	Avatar   string `gorm:"size:255" json:"avatar"`
	Status   int8   `gorm:"default:1;comment:状态 1正常 0禁用" json:"status"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

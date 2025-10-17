package model

// Log 日志模型
type Log struct {
	BaseModel
	UserID   uint   `gorm:"index" json:"user_id"`
	Module   string `gorm:"size:64" json:"module"`
	Action   string `gorm:"size:64" json:"action"`
	Content  string `gorm:"type:text" json:"content"`
	IP       string `gorm:"size:64" json:"ip"`
	UserAgent string `gorm:"size:255" json:"user_agent"`
}

// TableName 指定表名
func (Log) TableName() string {
	return "logs"
}

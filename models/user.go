package models

// User 对应数据库中的 user 表
type User struct {
	UserID   int64  `gorm:"column:user_id;uniqueIndex:idx_user_id;not null" json:"user_id,string"`              // 用户逻辑ID
	Username string `gorm:"column:username;type:varchar(64);uniqueIndex:idx_username;not null" json:"username"` // 用户名
	Password string `gorm:"column:password;type:varchar(64);not null" json:"-"`                                 // 密码（json化时隐藏）
	Token    string `json:"token"`
}

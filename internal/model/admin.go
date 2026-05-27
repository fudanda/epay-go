// internal/model/admin.go
package model

import "time"

// Admin 管理员
type Admin struct {
	BaseModel
	Username    string     `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password    string     `gorm:"size:128;not null" json:"-"`
	Role        string     `gorm:"size:20;default:admin" json:"role"` // super, admin
	LastLoginAt *time.Time `json:"last_login_at"`
}

func (Admin) TableName() string {
	return prefixedTableName("admins")
}

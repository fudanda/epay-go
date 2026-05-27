// internal/model/config.go
package model

import "time"

// Config 系统配置
type Config struct {
	Key       string    `gorm:"size:64;primarykey" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Config) TableName() string {
	return prefixedTableName("configs")
}

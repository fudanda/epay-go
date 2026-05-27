// internal/model/merchant.go
package model

import (
	"github.com/shopspring/decimal"
)

// Merchant 商户
type Merchant struct {
	BaseModel
	Username      string          `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password      string          `gorm:"size:128;not null" json:"-"`
	Email         string          `gorm:"size:128" json:"email"`
	Phone         string          `gorm:"size:20" json:"phone"`
	ApiKey        string          `gorm:"size:64;uniqueIndex;not null" json:"api_key"`
	PublicKey     string          `gorm:"type:text" json:"-"`
	PrivateKey    string          `gorm:"type:text" json:"-"`
	Balance       decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"balance"`
	FrozenBalance decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"frozen_balance"`
	Status        int8            `gorm:"default:1" json:"status"` // 0禁用 1正常
}

func (Merchant) TableName() string {
	return prefixedTableName("merchants")
}

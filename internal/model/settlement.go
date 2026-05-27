// internal/model/settlement.go
package model

import (
	"github.com/shopspring/decimal"
)

// Settlement 结算记录
type Settlement struct {
	BaseModel
	SettleNo     string          `gorm:"size:32;uniqueIndex;not null" json:"settle_no"`
	MerchantID   int64           `gorm:"index;not null" json:"merchant_id"`
	Amount       decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"amount"`
	Fee          decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"fee"`
	ActualAmount decimal.Decimal `gorm:"type:decimal(12,2)" json:"actual_amount"`
	AccountType  string          `gorm:"size:20" json:"account_type"` // alipay, bank
	AccountNo    string          `gorm:"size:64" json:"account_no"`
	AccountName  string          `gorm:"size:64" json:"account_name"`
	Status       int8            `gorm:"default:0" json:"status"` // 0待审核 1处理中 2已完成 3已驳回
	Remark       string          `gorm:"size:255" json:"remark"`

	// 关联
	Merchant *Merchant `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
}

func (Settlement) TableName() string {
	return prefixedTableName("settlements")
}

// 结算状态常量
const (
	SettleStatusPending    = 0
	SettleStatusProcessing = 1
	SettleStatusCompleted  = 2
	SettleStatusRejected   = 3
)

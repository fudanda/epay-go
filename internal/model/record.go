// internal/model/record.go
package model

import (
	"github.com/shopspring/decimal"
)

// BalanceRecord 资金变动记录
type BalanceRecord struct {
	BaseModel
	MerchantID    int64           `gorm:"index;not null" json:"merchant_id"`
	Action        int8            `gorm:"not null" json:"action"` // 1收入 2支出
	Amount        decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"amount"`
	BeforeBalance decimal.Decimal `gorm:"type:decimal(12,2)" json:"before_balance"`
	AfterBalance  decimal.Decimal `gorm:"type:decimal(12,2)" json:"after_balance"`
	Type          string          `gorm:"size:32" json:"type"` // order_income, fee, settle, refund
	TradeNo       string          `gorm:"size:64" json:"trade_no"`

	// 关联
	Merchant *Merchant `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
}

func (BalanceRecord) TableName() string {
	return prefixedTableName("balance_records")
}

// 资金变动类型常量
const (
	RecordActionIncome  = 1
	RecordActionExpense = 2
)

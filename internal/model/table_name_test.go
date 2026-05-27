package model

import "testing"

func TestPrefixedTableNames(t *testing.T) {
	testCases := map[string]string{
		Admin{}.TableName():         "epay_admins",
		Merchant{}.TableName():      "epay_merchants",
		Channel{}.TableName():       "epay_channels",
		Order{}.TableName():         "epay_orders",
		Settlement{}.TableName():    "epay_settlements",
		BalanceRecord{}.TableName(): "epay_balance_records",
		Config{}.TableName():        "epay_configs",
		Refund{}.TableName():        "epay_refunds",
	}

	for got, want := range testCases {
		if got != want {
			t.Fatalf("table name mismatch: got %s want %s", got, want)
		}
	}
}

package model

const tablePrefix = "epay_"

func prefixedTableName(base string) string {
	return tablePrefix + base
}

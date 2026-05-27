// internal/database/migrate.go
package database

import (
	"fmt"
	"log"
	"strings"

	"github.com/example/epay-go/internal/model"
)

// Migrate 自动迁移数据库表
func Migrate() error {
	log.Println("Running database migrations...")

	if err := renameLegacyTables(); err != nil {
		return err
	}

	err := DB.AutoMigrate(
		&model.Admin{},
		&model.Merchant{},
		&model.Channel{},
		&model.Order{},
		&model.Settlement{},
		&model.BalanceRecord{},
		&model.Config{},
		&model.Refund{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed")
	return nil
}

type tableRename struct {
	legacy  string
	current string
}

func legacyTableRenames() []tableRename {
	return []tableRename{
		{legacy: "admins", current: model.Admin{}.TableName()},
		{legacy: "merchants", current: model.Merchant{}.TableName()},
		{legacy: "channels", current: model.Channel{}.TableName()},
		{legacy: "orders", current: model.Order{}.TableName()},
		{legacy: "settlements", current: model.Settlement{}.TableName()},
		{legacy: "balance_records", current: model.BalanceRecord{}.TableName()},
		{legacy: "configs", current: model.Config{}.TableName()},
		{legacy: "refunds", current: model.Refund{}.TableName()},
	}
}

func renameLegacyTables() error {
	migrator := DB.Migrator()

	for _, rename := range legacyTableRenames() {
		hasLegacy := migrator.HasTable(rename.legacy)
		hasCurrent := migrator.HasTable(rename.current)

		if hasLegacy && hasCurrent {
			return fmt.Errorf("both legacy and prefixed tables exist: %s and %s", rename.legacy, rename.current)
		}

		if !hasLegacy || hasCurrent {
			continue
		}

		if err := migrator.RenameTable(rename.legacy, rename.current); err != nil {
			return fmt.Errorf("rename table %s -> %s: %w", rename.legacy, rename.current, err)
		}

		if err := renameTableIndexes(rename.legacy, rename.current); err != nil {
			return err
		}

		if err := renameTableConstraints(rename.legacy, rename.current); err != nil {
			return err
		}

		log.Printf("Renamed legacy table %s to %s", rename.legacy, rename.current)
	}

	return nil
}

func renameTableIndexes(legacyTableName, currentTableName string) error {
	oldPrefix := "idx_" + legacyTableName + "_"
	newPrefix := "idx_" + currentTableName + "_"

	var indexNames []string
	if err := DB.Raw(`
		SELECT indexname
		FROM pg_indexes
		WHERE schemaname = CURRENT_SCHEMA()
		  AND tablename = ?
		  AND indexname LIKE ?
	`, currentTableName, oldPrefix+"%").Scan(&indexNames).Error; err != nil {
		return fmt.Errorf("load indexes for %s: %w", currentTableName, err)
	}

	for _, indexName := range indexNames {
		newName := strings.Replace(indexName, oldPrefix, newPrefix, 1)
		if err := DB.Exec(
			"ALTER INDEX " + quoteIdentifier(indexName) + " RENAME TO " + quoteIdentifier(newName),
		).Error; err != nil {
			return fmt.Errorf("rename index %s -> %s: %w", indexName, newName, err)
		}
	}

	return nil
}

func renameTableConstraints(legacyTableName, currentTableName string) error {
	oldFKPrefix := "fk_" + legacyTableName + "_"
	newFKPrefix := "fk_" + currentTableName + "_"
	oldPrimaryKey := legacyTableName + "_pkey"
	newPrimaryKey := currentTableName + "_pkey"

	var constraintNames []string
	if err := DB.Raw(`
		SELECT c.conname
		FROM pg_constraint c
		JOIN pg_class t ON t.oid = c.conrelid
		JOIN pg_namespace n ON n.oid = t.relnamespace
		WHERE n.nspname = CURRENT_SCHEMA()
		  AND t.relname = ?
		  AND (c.conname = ? OR c.conname LIKE ?)
	`, currentTableName, oldPrimaryKey, oldFKPrefix+"%").Scan(&constraintNames).Error; err != nil {
		return fmt.Errorf("load constraints for %s: %w", currentTableName, err)
	}

	for _, constraintName := range constraintNames {
		newName := constraintName
		switch {
		case constraintName == oldPrimaryKey:
			newName = newPrimaryKey
		case strings.HasPrefix(constraintName, oldFKPrefix):
			newName = strings.Replace(constraintName, oldFKPrefix, newFKPrefix, 1)
		default:
			continue
		}

		if err := DB.Exec(
			"ALTER TABLE " + quoteIdentifier(currentTableName) +
				" RENAME CONSTRAINT " + quoteIdentifier(constraintName) +
				" TO " + quoteIdentifier(newName),
		).Error; err != nil {
			return fmt.Errorf("rename constraint %s -> %s: %w", constraintName, newName, err)
		}
	}

	return nil
}

func quoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

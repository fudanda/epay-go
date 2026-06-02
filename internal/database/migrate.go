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
			if err := reconcileLegacyAndPrefixedTable(rename.legacy, rename.current); err != nil {
				return err
			}
			continue
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

func reconcileLegacyAndPrefixedTable(legacyTable, currentTable string) error {
	legacyRows, err := tableRowCount(legacyTable)
	if err != nil {
		log.Printf("Warning: failed to count legacy table %s: %v. Keeping prefixed table %s.", legacyTable, err, currentTable)
		return nil
	}
	currentRows, err := tableRowCount(currentTable)
	if err != nil {
		log.Printf("Warning: failed to count prefixed table %s: %v. Keeping prefixed table.", currentTable, err)
		return nil
	}

	if legacyRows == 0 {
		log.Printf("Legacy and prefixed tables both exist: %s and %s (legacy rows=0). Keeping prefixed table.", legacyTable, currentTable)
		return nil
	}

	// If both tables exist and only legacy table has rows, backfill to prefixed table
	// so runtime always reads the canonical prefixed table.
	if currentRows == 0 {
		if err := copyLegacyRowsToPrefixed(legacyTable, currentTable); err != nil {
			log.Printf(
				"Warning: failed to backfill prefixed table %s from legacy table %s: %v. "+
					"Service will continue with prefixed table; please reconcile data manually if needed.",
				currentTable,
				legacyTable,
				err,
			)
			return nil
		}
		log.Printf("Backfilled %s from %s (%d rows).", currentTable, legacyTable, legacyRows)
		return nil
	}

	log.Printf(
		"Legacy and prefixed tables both exist and contain data: %s(%d) and %s(%d). Keeping prefixed table for runtime.",
		legacyTable,
		legacyRows,
		currentTable,
		currentRows,
	)
	return nil
}

func tableRowCount(tableName string) (int64, error) {
	var count int64
	if err := DB.Raw("SELECT COUNT(*) FROM " + quoteIdentifier(tableName)).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func copyLegacyRowsToPrefixed(legacyTable, currentTable string) error {
	// Tables are structurally equivalent (legacy name migration). ON CONFLICT DO NOTHING
	// ensures this remains idempotent if rerun.
	return DB.Exec(
		"INSERT INTO " + quoteIdentifier(currentTable) +
			" SELECT * FROM " + quoteIdentifier(legacyTable) +
			" ON CONFLICT DO NOTHING",
	).Error
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

package r2dbc

import (
	"context"
	"fmt"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken/persistence"
	"gorm.io/gorm"
)

type R2dbcEventStreamPersistence struct {
	db        *gorm.DB
	tableName string
}

func NewR2dbcEventStreamPersistence(db *gorm.DB, tableName string) *R2dbcEventStreamPersistence {
	return &R2dbcEventStreamPersistence{
		db:        db,
		tableName: tableName,
	}
}

func (r *R2dbcEventStreamPersistence) Append(ctx context.Context, name string, data []byte, expectedVersion int64, newVersion int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var actualVersion int64

		err := tx.Raw(
			fmt.Sprintf("SELECT COALESCE(MAX(version), %d) FROM %s WHERE name = ?", 0, r.tableName),
			name,
		).Scan(&actualVersion).Error
		if err != nil {
			return fmt.Errorf("failed to determine actual version: %w", err)
		}

		if actualVersion != expectedVersion {
			return &persistence.PersistenceConcurrencyError{
				Name:            name,
				ActualVersion:   actualVersion,
				ExpectedVersion: expectedVersion,
			}
		}

		err = tx.Exec(
			fmt.Sprintf("INSERT INTO %s (name, version, data) VALUES (?, ?, ?)", r.tableName),
			name, newVersion, data,
		).Error
		if err != nil {
			return fmt.Errorf("failed to append data: %w", err)
		}

		return nil
	})
}

func (r *R2dbcEventStreamPersistence) ReadStream(ctx context.Context, name string, afterVersion *int64) ([]persistence.BinaryData, error) {
	query := fmt.Sprintf("SELECT name, version, data FROM %s WHERE name = ? ORDER BY version", r.tableName)
	var params []interface{}

	params = append(params, name)

	if afterVersion != nil {
		query += " AND version > ?"
		params = append(params, *afterVersion)
	}

	var results []persistence.BinaryData
	err := r.db.WithContext(ctx).Raw(query, params...).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to read stream: %w", err)
	}

	return results, nil
}

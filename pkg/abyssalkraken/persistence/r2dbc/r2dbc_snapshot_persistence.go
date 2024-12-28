package r2dbc

import (
	"context"
	"errors"
	"fmt"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken/persistence"
	"gorm.io/gorm"
)

type R2dbcSnapshotPersistence struct {
	db        *gorm.DB
	tableName string
}

func NewR2dbcSnapshotPersistence(db *gorm.DB, tableName string) *R2dbcSnapshotPersistence {
	return &R2dbcSnapshotPersistence{
		db:        db,
		tableName: tableName,
	}
}

func (r *R2dbcSnapshotPersistence) Upsert(ctx context.Context, name string, data []byte, expectedVersion int64, newVersion int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var actualVersion *int64

		err := tx.Raw(
			fmt.Sprintf("SELECT MAX(version) FROM %s WHERE name = ?", r.tableName),
			name,
		).Scan(&actualVersion).Error
		if err != nil {
			return fmt.Errorf("failed to determine actual version: %w", err)
		}

		if actualVersion == nil {
			err = tx.Exec(
				fmt.Sprintf("INSERT INTO %s (name, version, data) VALUES (?, ?, ?)", r.tableName),
				name, newVersion, data,
			).Error
			if err != nil {
				return fmt.Errorf("failed to insert snapshot: %w", err)
			}
		} else {
			if *actualVersion != expectedVersion {
				return &persistence.PersistenceConcurrencyError{
					Name:            name,
					ActualVersion:   *actualVersion,
					ExpectedVersion: expectedVersion,
				}
			}

			err = tx.Exec(
				fmt.Sprintf("UPDATE %s SET data = ?, version = ? WHERE name = ?", r.tableName),
				data, newVersion, name,
			).Error
			if err != nil {
				return fmt.Errorf("failed to update snapshot: %w", err)
			}
		}

		return nil
	})
}

func (r *R2dbcSnapshotPersistence) ReadRecord(ctx context.Context, name string) (*persistence.BinaryData, error) {
	var result persistence.BinaryData

	err := r.db.WithContext(ctx).Raw(
		fmt.Sprintf("SELECT name, version, data FROM %s WHERE name = ?", r.tableName),
		name,
	).Scan(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot: %w", err)
	}

	return &result, nil
}

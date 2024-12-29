package r2dbc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken/config"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken/persistence"
)

type R2dbcSnapshotPersistence struct {
	db        *sql.DB
	tableName string
}

func NewR2dbcSnapshotPersistence(db *sql.DB, config config.R2dbcPersistenceConfig) *R2dbcSnapshotPersistence {
	return &R2dbcSnapshotPersistence{
		db:        db,
		tableName: config.TableName,
	}
}

func (r *R2dbcSnapshotPersistence) Upsert(ctx context.Context, name string, data []byte, expectedVersion int64, newVersion int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var actualVersion *int64
	versionQuery := fmt.Sprintf("SELECT MAX(version) FROM %s WHERE name = ?", r.tableName)
	err = tx.QueryRowContext(ctx, versionQuery, name).Scan(&actualVersion)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to determine actual version: %w", err)
	}

	if actualVersion == nil {
		insertQuery := fmt.Sprintf("INSERT INTO %s (name, version, data) VALUES (?, ?, ?)", r.tableName)
		_, err = tx.ExecContext(ctx, insertQuery, name, newVersion, data)
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

		updateQuery := fmt.Sprintf("UPDATE %s SET data = ?, version = ? WHERE name = ?", r.tableName)
		_, err = tx.ExecContext(ctx, updateQuery, data, newVersion, name)
		if err != nil {
			return fmt.Errorf("failed to update snapshot: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *R2dbcSnapshotPersistence) ReadRecord(ctx context.Context, name string) (*persistence.BinaryData, error) {
	query := fmt.Sprintf("SELECT version, data FROM %s WHERE name = ?", r.tableName)
	var result persistence.BinaryData

	err := r.db.QueryRowContext(ctx, query, name).Scan(&result.Version, &result.Data)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot: %w", err)
	}

	return &result, nil
}

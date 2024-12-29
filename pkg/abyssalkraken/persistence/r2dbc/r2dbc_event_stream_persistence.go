package r2dbc

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken/config"

	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken/persistence"
)

type R2dbcEventStreamPersistence struct {
	db        *sql.DB
	tableName string
}

func NewR2dbcEventStreamPersistence(db *sql.DB, config config.R2dbcPersistenceConfig) *R2dbcEventStreamPersistence {
	return &R2dbcEventStreamPersistence{
		db:        db,
		tableName: config.TableName,
	}
}

func (r *R2dbcEventStreamPersistence) Append(ctx context.Context, name string, data []byte, expectedVersion int64, newVersion int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var actualVersion int64
	query := fmt.Sprintf("SELECT COALESCE(MAX(version), %d) FROM %s WHERE name = ?", 0, r.tableName)
	err = tx.QueryRowContext(ctx, query, name).Scan(&actualVersion)
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

	insertQuery := fmt.Sprintf("INSERT INTO %s (name, version, data) VALUES (?, ?, ?)", r.tableName)
	_, err = tx.ExecContext(ctx, insertQuery, name, newVersion, data)
	if err != nil {
		return fmt.Errorf("failed to append data: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *R2dbcEventStreamPersistence) ReadStream(ctx context.Context, name string, afterVersion *int64) ([]persistence.BinaryData, error) {
	query := fmt.Sprintf("SELECT version, data FROM %s WHERE name = ? ORDER BY version", r.tableName)
	params := []interface{}{name}

	if afterVersion != nil {
		query += " AND version > ?"
		params = append(params, *afterVersion)
	}

	rows, err := r.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to read stream: %w", err)
	}
	defer rows.Close()

	var results []persistence.BinaryData
	for rows.Next() {
		var result persistence.BinaryData
		if err := rows.Scan(&result.Version, &result.Data); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while reading rows: %w", err)
	}

	return results, nil
}

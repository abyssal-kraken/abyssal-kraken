package r2dbc

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken/config"
)

var (
	r2dbcSnapshotRegistries     sync.Map
	r2dbcSnapshotRegistryErrors sync.Map
)

func R2dbcSnapshotRegistryConfig(db *sql.DB, config config.R2dbcPersistenceConfig) (*R2dbcSnapshotPersistence, error) {
	key := config.TableName

	if registry, exists := r2dbcSnapshotRegistries.Load(key); exists {
		if err, ok := r2dbcSnapshotRegistryErrors.Load(key); ok && err != nil {
			return nil, err.(error)
		}
		return registry.(*R2dbcSnapshotPersistence), nil
	}

	var once sync.Once
	var persistence *R2dbcSnapshotPersistence
	var initErr error

	once.Do(func() {
		if db == nil {
			initErr = fmt.Errorf("database connection is not configured")
			r2dbcSnapshotRegistryErrors.Store(key, initErr)
			return
		}
		if config.TableName == "" {
			initErr = fmt.Errorf("tableName is not configured")
			r2dbcSnapshotRegistryErrors.Store(key, initErr)
			return
		}

		persistence = NewR2dbcSnapshotPersistence(db, config)

		if config.ValidateSchemaOnInit {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := persistence.ValidateSchema(ctx); err != nil {
				initErr = fmt.Errorf("failed to validate schema for %s: %w", config.TableName, err)
				r2dbcSnapshotRegistryErrors.Store(key, initErr)
				return
			}
		}

		r2dbcSnapshotRegistries.Store(key, persistence)
		r2dbcSnapshotRegistryErrors.Store(key, nil)
	})

	return persistence, initErr
}

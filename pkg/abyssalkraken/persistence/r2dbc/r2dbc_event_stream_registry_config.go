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
	r2dbcEventStreamRegistries     sync.Map
	r2dbcEventStreamRegistryErrors sync.Map
)

func R2dbcEventStreamRegistryConfig(db *sql.DB, config config.R2dbcPersistenceConfig) (*R2dbcEventStreamPersistence, error) {
	key := config.TableName

	if registry, exists := r2dbcEventStreamRegistries.Load(key); exists {
		if err, ok := r2dbcEventStreamRegistryErrors.Load(key); ok && err != nil {
			return nil, err.(error)
		}
		return registry.(*R2dbcEventStreamPersistence), nil
	}

	var once sync.Once
	var persistence *R2dbcEventStreamPersistence
	var initErr error

	once.Do(func() {
		if db == nil {
			initErr = fmt.Errorf("database connection is not configured")
			r2dbcEventStreamRegistryErrors.Store(key, initErr)
			return
		}
		if config.TableName == "" {
			initErr = fmt.Errorf("tableName is not configured")
			r2dbcEventStreamRegistryErrors.Store(key, initErr)
			return
		}

		persistence = NewR2dbcEventStreamPersistence(db, config)

		if config.ValidateSchemaOnInit {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := persistence.ValidateSchema(ctx); err != nil {
				initErr = fmt.Errorf("failed to validate schema for %s: %w", config.TableName, err)
				r2dbcEventStreamRegistryErrors.Store(key, initErr)
				return
			}
		}

		r2dbcEventStreamRegistries.Store(key, persistence)
		r2dbcEventStreamRegistryErrors.Store(key, nil)
	})

	return persistence, initErr
}

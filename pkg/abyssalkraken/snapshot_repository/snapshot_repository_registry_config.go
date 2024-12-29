package snapshot_repository

import (
	"errors"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"reflect"
	"sync"

	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken/persistence"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken/serialization"
)

var (
	snapshotRepositoryRegistries     sync.Map
	snapshotRepositoryRegistryErrors sync.Map
)

func SnapshotRepositoryRegistryConfig[
	ID abyssalkraken.AggregateID,
	E abyssalkraken.DomainEvent[ID],
	A abyssalkraken.AggregateRoot[ID, E],
](
	persistenceLayer persistence.SnapshotPersistence,
	serializationLayer serialization.SnapshotSerialization[ID, E, A],
) (*SnapshotRepository[ID, E, A], error) {
	typeKey := [3]reflect.Type{
		reflect.TypeOf((*ID)(nil)).Elem(),
		reflect.TypeOf((*E)(nil)).Elem(),
		reflect.TypeOf((*A)(nil)).Elem(),
	}

	if registry, exists := snapshotRepositoryRegistries.Load(typeKey); exists {
		if err, ok := snapshotRepositoryRegistryErrors.Load(typeKey); ok && err != nil {
			return nil, err.(error)
		}
		return registry.(*SnapshotRepository[ID, E, A]), nil
	}

	var once sync.Once
	var singleton *SnapshotRepository[ID, E, A]
	var initErr error

	once.Do(func() {
		if persistenceLayer == nil {
			initErr = errors.New("snapshotPersistence is not configured")
			snapshotRepositoryRegistryErrors.Store(typeKey, initErr)
			return
		}
		if serializationLayer == nil {
			initErr = errors.New("snapshotPersistence is not configured")
			snapshotRepositoryRegistryErrors.Store(typeKey, initErr)
			return
		}

		singleton = NewSnapshotRepository(persistenceLayer, serializationLayer)
		snapshotRepositoryRegistries.Store(typeKey, singleton)
		snapshotRepositoryRegistryErrors.Store(typeKey, nil)
	})

	return singleton, initErr
}

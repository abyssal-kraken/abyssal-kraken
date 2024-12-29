package avro

import (
	"errors"
	"reflect"
	"sync"

	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
)

var (
	snapshotRegistries     sync.Map
	snapshotRegistryErrors sync.Map
)

func AvroSnapshotConverterRegistryConfig[
	ID abyssalkraken.AggregateID,
	E abyssalkraken.DomainEvent[ID],
	A abyssalkraken.AggregateRoot[ID, E],
	GC any,
](
	converters []AvroSnapshotConverter[ID, E, A, GC],
) (*AvroSnapshotConverterRegistry[ID, E, A, GC], error) {
	typeKey := [4]reflect.Type{
		reflect.TypeOf((*ID)(nil)).Elem(),
		reflect.TypeOf((*E)(nil)).Elem(),
		reflect.TypeOf((*A)(nil)).Elem(),
		reflect.TypeOf((*GC)(nil)).Elem(),
	}

	if registry, exists := snapshotRegistries.Load(typeKey); exists {
		if err, ok := snapshotRegistryErrors.Load(typeKey); ok && err != nil {
			return nil, err.(error)
		}
		return registry.(*AvroSnapshotConverterRegistry[ID, E, A, GC]), nil
	}

	var once sync.Once
	var registryInstance *AvroSnapshotConverterRegistry[ID, E, A, GC]
	var initError error

	once.Do(func() {
		if len(converters) == 0 {
			initError = errors.New("no AvroSnapshotConverters provided")
			snapshotRegistryErrors.Store(typeKey, initError)
			return
		}

		registryInstance = NewAvroSnapshotConverterRegistry(converters)
		snapshotRegistries.Store(typeKey, registryInstance)
		snapshotRegistryErrors.Store(typeKey, nil)
	})

	return registryInstance, initError
}

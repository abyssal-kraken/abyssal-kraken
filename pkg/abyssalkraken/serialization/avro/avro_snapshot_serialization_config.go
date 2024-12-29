package avro

import (
	"errors"
	"reflect"
	"sync"

	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
)

var (
	snapshotSerializationInstances sync.Map
	snapshotSerializationErrors    sync.Map
)

func AvroSnapshotSerializationConfig[
	ID abyssalkraken.AggregateID,
	E abyssalkraken.DomainEvent[ID],
	A abyssalkraken.AggregateRoot[ID, E],
](
	converterRegistry *AvroSnapshotConverterRegistry[ID, E, A, map[string]interface{}],
) (*AvroSnapshotSerialization[ID, E, A], error) {
	if converterRegistry == nil {
		return nil, errors.New("AvroSnapshotConverterRegistry is not configured")
	}

	typeKey := [3]reflect.Type{
		reflect.TypeOf((*ID)(nil)).Elem(),
		reflect.TypeOf((*E)(nil)).Elem(),
		reflect.TypeOf((*A)(nil)).Elem(),
	}

	if instance, exists := snapshotSerializationInstances.Load(typeKey); exists {
		if err, ok := snapshotSerializationErrors.Load(typeKey); ok && err != nil {
			return nil, err.(error)
		}
		return instance.(*AvroSnapshotSerialization[ID, E, A]), nil
	}

	var once sync.Once
	var instance *AvroSnapshotSerialization[ID, E, A]
	var initError error

	once.Do(func() {
		instance = NewAvroSnapshotSerialization(converterRegistry)
		snapshotSerializationInstances.Store(typeKey, instance)
		snapshotSerializationErrors.Store(typeKey, nil)
	})

	return instance, initError
}

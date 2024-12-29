package avro

import (
	"errors"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"reflect"
	"sync"
)

var (
	serializationInstances sync.Map
	serializationErrors    sync.Map
)

func AvroEventSerializationConfig[
	ID abyssalkraken.AggregateID,
	E abyssalkraken.DomainEvent[ID],
](
	registry *AvroEventConverterRegistry[ID, E, map[string]interface{}],
) (*AvroEventSerialization[ID, E], error) {
	if registry == nil {
		return nil, errors.New("AvroEventConverterRegistry is not configured")
	}

	typeKey := [2]reflect.Type{
		reflect.TypeOf((*ID)(nil)).Elem(),
		reflect.TypeOf((*E)(nil)).Elem(),
	}

	if instance, exists := serializationInstances.Load(typeKey); exists {
		if err, ok := serializationErrors.Load(typeKey); ok && err != nil {
			return nil, err.(error)
		}
		return instance.(*AvroEventSerialization[ID, E]), nil
	}

	var once sync.Once
	var instance *AvroEventSerialization[ID, E]
	var initErr error

	once.Do(func() {
		if registry == nil {
			initErr = errors.New("no AvroEventConverterRegistry provided for serialization")
			serializationErrors.Store(typeKey, initErr)
			return
		}

		instance = NewAvroEventSerialization(registry)
		serializationInstances.Store(typeKey, instance)
		serializationErrors.Store(typeKey, nil)
	})

	return instance, initErr
}

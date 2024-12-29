package avro

import (
	"errors"
	"reflect"
	"sync"

	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
)

var (
	eventStreamSerializationInstances sync.Map
	eventStreamSerializationErrors    sync.Map
)

func AvroEventStreamSerializationConfig[
	ID abyssalkraken.AggregateID,
	E abyssalkraken.DomainEvent[ID],
](
	converterRegistry *AvroEventStreamConverterRegistry[ID, E, []map[string]interface{}],
) (*AvroEventStreamSerialization[ID, E], error) {
	if converterRegistry == nil {
		return nil, errors.New("AvroEventStreamConverterRegistry is not configured")
	}

	typeKey := [2]reflect.Type{
		reflect.TypeOf((*ID)(nil)).Elem(),
		reflect.TypeOf((*E)(nil)).Elem(),
	}

	if instance, exists := eventStreamSerializationInstances.Load(typeKey); exists {
		if err, ok := eventStreamSerializationErrors.Load(typeKey); ok && err != nil {
			return nil, err.(error)
		}
		return instance.(*AvroEventStreamSerialization[ID, E]), nil
	}

	var once sync.Once
	var instance *AvroEventStreamSerialization[ID, E]
	var initError error

	once.Do(func() {
		instance = NewAvroEventStreamSerialization(converterRegistry)
		eventStreamSerializationInstances.Store(typeKey, instance)
		eventStreamSerializationErrors.Store(typeKey, nil)
	})

	return instance, initError
}

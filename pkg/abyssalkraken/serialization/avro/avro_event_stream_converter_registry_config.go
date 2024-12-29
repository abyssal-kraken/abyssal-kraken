package avro

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
)

var (
	eventStreamRegistries     sync.Map
	eventStreamRegistryErrors sync.Map
)

func AvroEventStreamConverterRegistryConfig[
	ID abyssalkraken.AggregateID,
	E abyssalkraken.DomainEvent[ID],
	GC any,
](
	converters []AvroEventStreamConverter[ID, E, GC],
) (*AvroEventStreamConverterRegistry[ID, E, GC], error) {
	typeKey := [3]reflect.Type{
		reflect.TypeOf((*ID)(nil)).Elem(),
		reflect.TypeOf((*E)(nil)).Elem(),
		reflect.TypeOf((*GC)(nil)).Elem(),
	}

	if registry, exists := eventStreamRegistries.Load(typeKey); exists {
		if err, ok := eventStreamRegistryErrors.Load(typeKey); ok && err != nil {
			return nil, err.(error)
		}
		return registry.(*AvroEventStreamConverterRegistry[ID, E, GC]), nil
	}

	var once sync.Once
	var registryInstance *AvroEventStreamConverterRegistry[ID, E, GC]
	var initError error

	once.Do(func() {
		if len(converters) == 0 {
			initError = errors.New("no AvroEventStreamConverters provided")
			eventStreamRegistryErrors.Store(typeKey, initError)
			return
		}

		registryInstance = NewAvroEventStreamConverterRegistry[ID, E, GC]()

		for _, converter := range converters {
			if err := registryInstance.Register(typeKey[1], converter); err != nil {
				initError = fmt.Errorf("failed to register event stream converter: %w", err)
				registryErrors.Store(typeKey, initError)
				return
			}
		}

		eventStreamRegistries.Store(typeKey, registryInstance)
		eventStreamRegistryErrors.Store(typeKey, nil)
	})

	return registryInstance, initError
}

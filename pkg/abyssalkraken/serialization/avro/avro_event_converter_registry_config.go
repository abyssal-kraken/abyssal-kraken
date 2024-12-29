package avro

import (
	"errors"
	"fmt"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"reflect"
	"sync"
)

var (
	registries     sync.Map
	registryErrors sync.Map
)

func AvroEventConverterRegistryConfig[
	ID abyssalkraken.AggregateID,
	E abyssalkraken.DomainEvent[ID],
	GC any,
](
	converters []AvroEventConverter[ID, E, GC],
) (*AvroEventConverterRegistry[ID, E, GC], error) {
	typeKey := [3]reflect.Type{
		reflect.TypeOf((*ID)(nil)).Elem(),
		reflect.TypeOf((*E)(nil)).Elem(),
		reflect.TypeOf((*GC)(nil)).Elem(),
	}

	if registry, exists := registries.Load(typeKey); exists {
		if err, ok := registryErrors.Load(typeKey); ok && err != nil {
			return nil, err.(error)
		}
		return registry.(*AvroEventConverterRegistry[ID, E, GC]), nil
	}

	var once sync.Once
	var singleton *AvroEventConverterRegistry[ID, E, GC]
	var initErr error

	once.Do(func() {
		if len(converters) == 0 {
			initErr = errors.New("no AvroEventConverters provided")
			registryErrors.Store(typeKey, initErr)
			return
		}

		singleton = NewAvroEventConverterRegistry[ID, E, GC]()

		for _, converter := range converters {
			if err := singleton.Register(typeKey[1], converter); err != nil {
				initErr = fmt.Errorf("failed to register event converter: %w", err)
				registryErrors.Store(typeKey, initErr)
				return
			}
		}

		registries.Store(typeKey, singleton)
		registryErrors.Store(typeKey, nil)
	})

	return singleton, initErr
}

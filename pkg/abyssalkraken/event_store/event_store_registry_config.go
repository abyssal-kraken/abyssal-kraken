package event_store

import (
	"errors"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken/persistence"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken/serialization"
	"reflect"
	"sync"
)

var (
	eventStoreRegistries     sync.Map
	eventStoreRegistryErrors sync.Map
)

func EventStoreRegistryConfig[
	ID abyssalkraken.AggregateID,
	E abyssalkraken.DomainEvent[ID],
](
	persistenceLayer persistence.EventStreamPersistence,
	serializationLayer serialization.EventStreamSerialization[ID, E],
) (*EventStore[ID, E], error) {
	typeKey := [2]reflect.Type{
		reflect.TypeOf((*ID)(nil)).Elem(),
		reflect.TypeOf((*E)(nil)).Elem(),
	}

	if registry, exists := eventStoreRegistries.Load(typeKey); exists {
		if err, ok := eventStoreRegistryErrors.Load(typeKey); ok && err != nil {
			return nil, err.(error)
		}
		return registry.(*EventStore[ID, E]), nil
	}

	var once sync.Once
	var singleton *EventStore[ID, E]
	var initErr error

	once.Do(func() {
		if persistenceLayer == nil {
			initErr = errors.New("eventStreamPersistence is not configured")
			eventStoreRegistryErrors.Store(typeKey, initErr)
			return
		}
		if serializationLayer == nil {
			initErr = errors.New("eventStreamSerialization is not configured")
			eventStoreRegistryErrors.Store(typeKey, initErr)
			return
		}

		singleton = NewEventStore[ID, E](persistenceLayer, serializationLayer)
		eventStoreRegistries.Store(typeKey, singleton)
		eventStoreRegistryErrors.Store(typeKey, nil)
	})

	return singleton, initErr
}

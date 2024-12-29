package eventbus

import (
	"errors"
	"reflect"
	"sync"
)

var (
	eventHandlerRegistries     sync.Map
	eventHandlerRegistryErrors sync.Map
)

func EventHandlerRegistryConfig(eventType reflect.Type) (*EventHandlerRegistry, error) {
	if eventType == nil {
		return nil, errors.New("eventType não pode ser nil")
	}

	if registry, exists := eventHandlerRegistries.Load(eventType); exists {
		if err, ok := eventHandlerRegistryErrors.Load(eventType); ok && err != nil {
			return nil, err.(error)
		}
		return registry.(*EventHandlerRegistry), nil
	}

	var once sync.Once
	var singleton *EventHandlerRegistry
	var initErr error

	once.Do(func() {
		singleton = &EventHandlerRegistry{}
		eventHandlerRegistries.Store(eventType, singleton)
		eventHandlerRegistryErrors.Store(eventType, nil)
	})

	return singleton, initErr
}

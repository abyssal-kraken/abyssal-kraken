package commandbus

import (
	"errors"
	"reflect"
	"sync"
)

var (
	handlerRegistries     sync.Map
	handlerRegistryErrors sync.Map
)

func CommandHandlerRegistryConfig(commandType reflect.Type) (*CommandHandlerRegistry, error) {
	if commandType == nil {
		return nil, errors.New("commandType cannot be nil")
	}

	if registry, exists := handlerRegistries.Load(commandType); exists {
		if err, ok := handlerRegistryErrors.Load(commandType); ok && err != nil {
			return nil, err.(error)
		}
		return registry.(*CommandHandlerRegistry), nil
	}

	var once sync.Once
	var singleton *CommandHandlerRegistry
	var initErr error

	once.Do(func() {
		singleton = &CommandHandlerRegistry{}
		handlerRegistries.Store(commandType, singleton)
		handlerRegistryErrors.Store(commandType, nil)
	})

	return singleton, initErr
}

package decorators

import (
	"fmt"
	"sync"
)

var (
	loggingEventDecoratorSingleton    *LoggingEventBusDecorator
	loggingEventDecoratorRegistryOnce sync.Once
	loggingEventDecoratorRegistryErr  error
)

func LoggingEventBusDecoratorRegistry() (*LoggingEventBusDecorator, error) {
	loggingEventDecoratorRegistryOnce.Do(func() {
		loggingEventDecoratorSingleton = &LoggingEventBusDecorator{}
	})

	if loggingEventDecoratorRegistryErr != nil {
		return nil, fmt.Errorf("erro ao registrar LoggingEventBusDecorator: %w", loggingEventDecoratorRegistryErr)
	}

	return loggingEventDecoratorSingleton, nil
}

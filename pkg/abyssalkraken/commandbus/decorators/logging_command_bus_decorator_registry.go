package decorators

import (
	"fmt"
	"sync"
)

var (
	loggingDecoratorSingleton     *LoggingCommandBusDecorator
	loggingDecoratorRegistryOnce  sync.Once
	loggingDecoratorRegistryError error
)

func LoggingCommandBusDecoratorRegistry() (*LoggingCommandBusDecorator, error) {
	loggingDecoratorRegistryOnce.Do(func() {
		loggingDecoratorSingleton = &LoggingCommandBusDecorator{}
	})

	if loggingDecoratorRegistryError != nil {
		return nil, fmt.Errorf("error to register LoggingCommandBusDecorator: %w", loggingDecoratorRegistryError)
	}

	return loggingDecoratorSingleton, nil
}

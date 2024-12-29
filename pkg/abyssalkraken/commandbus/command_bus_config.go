package commandbus

import (
	"errors"
	"sync"
)

var (
	commandBusOnce  sync.Once
	commandBus      *CommandBus
	commandBusError error
)

func CommandBusConfig(decorators ...CommandBusDecorator) (*CommandBus, error) {
	commandBusOnce.Do(func() {
		registry := &CommandHandlerRegistry{}

		commandBus = NewCommandBus(registry, decorators)
		if commandBus == nil {
			commandBusError = errors.New("falha ao inicializar CommandBus")
		}
	})

	return commandBus, commandBusError
}

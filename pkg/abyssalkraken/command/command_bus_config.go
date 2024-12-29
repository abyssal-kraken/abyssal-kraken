package command

import (
	"errors"
	"sync"
)

var (
	commandBusOnce  sync.Once
	commandBus      *CommandBus
	commandBusError error
)

func CommandBusConfig() (*CommandBus, error) {
	commandBusOnce.Do(func() {
		registry := &CommandHandlerRegistry{}

		commandBus = NewCommandBus(registry)
		if commandBus == nil {
			commandBusError = errors.New("failed to initialize CommandBus")
		}
	})

	return commandBus, commandBusError
}

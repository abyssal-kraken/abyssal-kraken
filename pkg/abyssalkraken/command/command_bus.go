package command

import (
	"context"
	"reflect"
)

type CommandBus struct {
	registry *CommandHandlerRegistry
}

func NewCommandBus(registry *CommandHandlerRegistry) *CommandBus {
	return &CommandBus{registry: registry}
}

func (cb *CommandBus) Execute(ctx context.Context, command interface{}) (interface{}, error) {
	commandType := reflect.TypeOf(command)

	handler, err := cb.registry.FindCommandHandler(commandType)
	if err != nil {
		return nil, err
	}

	commandHandler, ok := handler.(CommandHandler[interface{}, interface{}])
	if !ok {
		return nil, NewCommandHandlerNotFoundException(commandType)
	}

	return commandHandler.Handle(ctx, command)
}

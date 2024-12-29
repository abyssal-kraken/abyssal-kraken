package commandbus

import (
	"context"
	"reflect"
)

type CommandBus struct {
	registry   *CommandHandlerRegistry
	decorators []CommandBusDecorator
}

func NewCommandBus(registry *CommandHandlerRegistry, decorators []CommandBusDecorator) *CommandBus {
	return &CommandBus{registry: registry, decorators: decorators}
}

func (cb *CommandBus) Execute(ctx context.Context, command Command[any]) (any, error) {
	commandType := reflect.TypeOf(command)

	handler, err := cb.registry.FindCommandHandler(commandType)
	if err != nil {
		return nil, err
	}

	commandHandler, ok := handler.(CommandHandler[interface{}, interface{}])
	if !ok {
		return nil, NewCommandHandlerNotFoundException(commandType)
	}

	execution := func() (interface{}, error) {
		return commandHandler.Handle(ctx, command)
	}

	for i := len(cb.decorators) - 1; i >= 0; i-- {
		decorator := cb.decorators[i]
		execution = decorator.Decorate(command, execution)
	}

	return execution()
}

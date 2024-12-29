package commandbus

import (
	"fmt"
	"reflect"
	"sync"
)

type CommandHandlerRegistry struct {
	commandHandlers sync.Map
}

func (r *CommandHandlerRegistry) Register(handler interface{}) error {
	handlerType := reflect.TypeOf(handler)

	if handlerType.Kind() != reflect.Ptr {
		return fmt.Errorf("handler must be a pointer to struct implementing CommandHandler, got: %s", handlerType.String())
	}

	if handlerType.Implements(reflect.TypeOf((*CommandHandler[any, any])(nil)).Elem()) {
		commandType := handlerType.In(0).Elem()
		r.commandHandlers.Store(commandType, handler)
		return nil
	}

	return fmt.Errorf("handler does not implement CommandHandler interface: %s", handlerType.String())
}

func (r *CommandHandlerRegistry) FindCommandHandler(commandType reflect.Type) (interface{}, error) {
	handler, exists := r.commandHandlers.Load(commandType)
	if !exists {
		return nil, NewCommandHandlerNotFoundException(commandType)
	}
	return handler, nil
}

type CommandHandlerNotFoundException struct {
	commandType reflect.Type
}

func NewCommandHandlerNotFoundException(commandType reflect.Type) error {
	return &CommandHandlerNotFoundException{commandType: commandType}
}

func (e *CommandHandlerNotFoundException) Error() string {
	return fmt.Sprintf("no commandbus handler found for commandbus %s", e.commandType.String())
}

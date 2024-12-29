package eventbus

import (
	"fmt"
	"reflect"
	"sync"
)

type EventHandlerRegistry struct {
	eventHandlers sync.Map
}

func (r *EventHandlerRegistry) Register(
	eventType reflect.Type,
	eventHandlers []EventHandler,
) error {
	if eventType == nil {
		return fmt.Errorf("eventType não pode ser nil")
	}

	r.eventHandlers.Store(eventType, eventHandlers)
	return nil
}

func (r *EventHandlerRegistry) FindEventHandlers(
	eventType reflect.Type,
) ([]EventHandler, error) {
	handlers, exists := r.eventHandlers.Load(eventType)
	if !exists {
		return nil, fmt.Errorf("nenhum handler registrado para o tipo de evento: %s", eventType.String())
	}

	eventHandlers, ok := handlers.([]EventHandler)
	if !ok {
		return nil, fmt.Errorf("handlers registrados não são compatíveis com o tipo: %s", eventType.String())
	}

	return eventHandlers, nil
}

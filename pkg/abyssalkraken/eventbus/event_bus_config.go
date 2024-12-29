package eventbus

import (
	"errors"
	"sync"
)

var (
	eventBusOnce  sync.Once
	eventBus      *EventBus
	eventBusError error
)

func EventBusConfig() (*EventBus, error) {
	eventBusOnce.Do(func() {
		registry := &EventHandlerRegistry{}

		eventBus = &EventBus{registry: registry}
		if eventBus == nil {
			eventBusError = errors.New("falha ao inicializar EventBus")
		}
	})

	return eventBus, eventBusError
}

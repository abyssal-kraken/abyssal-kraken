package eventbus

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
)

type EventBus struct {
	registry   *EventHandlerRegistry
	decorators []EventBusDecorator
}

func NewEventBus(registry *EventHandlerRegistry, decorators []EventBusDecorator) *EventBus {
	return &EventBus{registry: registry, decorators: decorators}
}

func (bus *EventBus) PublishEvent(
	ctx context.Context,
	event abyssalkraken.DomainEvent[abyssalkraken.AggregateID],
	parallel bool,
) error {
	eventType := reflect.TypeOf(event)
	handlers, err := bus.registry.FindEventHandlers(eventType)
	if err != nil {
		return fmt.Errorf("falha ao localizar handlers para o evento: %w", err)
	}

	var execute func() error
	if parallel {
		execute = func() error {
			var wg sync.WaitGroup
			errChan := make(chan error, len(handlers))

			for _, handler := range handlers {
				wg.Add(1)

				handler := handler
				go func() {
					defer wg.Done()
					if execErr := handler.Handle(ctx, event); execErr != nil {
						errChan <- execErr
					}
				}()
			}

			wg.Wait()
			close(errChan)

			if len(errChan) > 0 {
				return <-errChan
			}
			return nil
		}
	} else {
		execute = func() error {
			for _, handler := range handlers {
				if err := handler.Handle(ctx, event); err != nil {
					return err
				}
			}
			return nil
		}
	}

	for i := len(bus.decorators) - 1; i >= 0; i-- {
		decorator := bus.decorators[i]
		execute = decorator.DecorateMultiple(ctx, event, handlers, execute)
	}

	return execute()
}

func (bus *EventBus) PublishEvents(ctx context.Context, events []abyssalkraken.DomainEvent[abyssalkraken.AggregateID], parallel bool) error {
	for _, event := range events {
		if err := bus.PublishEvent(ctx, event, parallel); err != nil {
			return err
		}
	}
	return nil
}

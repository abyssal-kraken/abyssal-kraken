package decorators

import (
	"context"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"log"
	"reflect"
	"time"

	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken/eventbus"
)

type LoggingEventBusDecorator struct{}

func (d LoggingEventBusDecorator) DecorateMultiple(
	ctx context.Context,
	event abyssalkraken.DomainEvent[abyssalkraken.AggregateID],
	eventHandlers []eventbus.EventHandler,
	execution func() error,
) func() error {
	return func() error {
		if len(eventHandlers) == 0 {
			log.Printf("No event handlers found for event: %T", event)
		} else {
			handlerNames := ""
			for _, handler := range eventHandlers {
				handlerNames += reflect.TypeOf(handler).String() + ", "
			}
			log.Printf("Publishing event %T to %d handlers: [%s]", event, len(eventHandlers), handlerNames)
		}

		return execution()
	}
}

func (d LoggingEventBusDecorator) DecorateSingle(
	ctx context.Context,
	event abyssalkraken.DomainEvent[abyssalkraken.AggregateID],
	eventHandler eventbus.EventHandler,
	execution func() error,
) func() error {
	return func() error {
		handlerName := reflect.TypeOf(eventHandler).String()
		log.Printf("Starting event handler: %s for event: %T", handlerName, event)

		startTime := time.Now()
		err := execution()
		duration := time.Since(startTime)

		if err != nil {
			log.Printf("Error in event handler: %s for event: %T, duration: %v, error: %v", handlerName, event, duration, err)
		} else {
			log.Printf("Finished event handler: %s for event: %T in %v", handlerName, event, duration)
		}

		return err
	}
}

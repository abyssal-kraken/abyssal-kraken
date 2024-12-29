package abyssalkraken

import (
	"errors"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"reflect"
	"sync"
)

var (
	eventMetricsHandlers      sync.Map
	eventMetricsHandlerErrors sync.Map
)

func EventMetricsHandlerRegistryConfig[
	ID abyssalkraken.AggregateID,
	E abyssalkraken.DomainEvent[ID],
](
	metricsPublisher EventMetricsPublisher[ID, E],
) (*EventMetricsHandler[ID, E], error) {
	typeKey := reflect.TypeOf((*EventMetricsPublisher[ID, E])(nil)).Elem()

	if handler, exists := eventMetricsHandlers.Load(typeKey); exists {
		if err, ok := eventMetricsHandlerErrors.Load(typeKey); ok && err != nil {
			return nil, err.(error)
		}
		return handler.(*EventMetricsHandler[ID, E]), nil
	}

	var once sync.Once
	var singleton *EventMetricsHandler[ID, E]
	var initErr error

	once.Do(func() {
		if metricsPublisher == nil {
			initErr = errors.New("eventMetricsPublisher is not configured")
			eventMetricsHandlerErrors.Store(typeKey, initErr)
			return
		}

		singleton = &EventMetricsHandler[ID, E]{
			MetricsPublisher: metricsPublisher,
		}

		eventMetricsHandlers.Store(typeKey, singleton)
		eventMetricsHandlerErrors.Store(typeKey, nil)
	})

	return singleton, initErr
}

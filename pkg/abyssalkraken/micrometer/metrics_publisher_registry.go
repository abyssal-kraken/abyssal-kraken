package micrometer

import (
	"errors"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"github.com/prometheus/client_golang/prometheus"
	"reflect"
	"sync"
)

var (
	micrometerMetricsPublisherRegistries     sync.Map
	micrometerMetricsPublisherRegistryErrors sync.Map
)

func MicrometerEventMetricsPublisherRegistry[
	ID abyssalkraken.AggregateID,
	E abyssalkraken.DomainEvent[ID],
](
	registry *prometheus.Registry,
) (*PrometheusEventMetricsPublisher[ID, E], error) {
	typeKey := reflect.TypeOf((*PrometheusEventMetricsPublisher[ID, E])(nil)).Elem()

	if publisher, exists := micrometerMetricsPublisherRegistries.Load(typeKey); exists {
		if err, ok := micrometerMetricsPublisherRegistryErrors.Load(typeKey); ok && err != nil {
			return nil, err.(error)
		}
		return publisher.(*PrometheusEventMetricsPublisher[ID, E]), nil
	}

	var once sync.Once
	var singleton *PrometheusEventMetricsPublisher[ID, E]
	var initErr error

	once.Do(func() {
		if registry == nil {
			initErr = errors.New("meterRegistry is not configured")
			micrometerMetricsPublisherRegistryErrors.Store(typeKey, initErr)
			return
		}

		singleton = NewPrometheusEventMetricsPublisher[ID, E](registry)

		micrometerMetricsPublisherRegistries.Store(typeKey, singleton)
		micrometerMetricsPublisherRegistryErrors.Store(typeKey, nil)
	})

	return singleton, initErr
}

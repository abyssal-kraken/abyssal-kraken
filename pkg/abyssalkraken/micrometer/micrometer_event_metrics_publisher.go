package micrometer

import (
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusEventMetricsPublisher[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID]] struct {
	totalEvents prometheus.Counter
}

func NewPrometheusEventMetricsPublisher[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID]](registry *prometheus.Registry) *PrometheusEventMetricsPublisher[ID, E] {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "abyssal_kraken_total_events",
		Help: "Total domain events emitted by this bounded context",
	})

	registry.MustRegister(counter)

	return &PrometheusEventMetricsPublisher[ID, E]{totalEvents: counter}
}

func (p *PrometheusEventMetricsPublisher[ID, E]) PublishMetrics(_ E) error {
	p.totalEvents.Inc()
	return nil
}

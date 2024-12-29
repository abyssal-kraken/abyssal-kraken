package abyssalkraken

import "github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"

type EventMetricsHandler[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID]] struct {
	MetricsPublisher EventMetricsPublisher[ID, E]
}

func (h EventMetricsHandler[ID, E]) HandleEvent(event E) error {
	return h.MetricsPublisher.PublishMetrics(event)
}

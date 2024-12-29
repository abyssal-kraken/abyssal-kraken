package abyssalkraken

import "github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"

type EventMetricsPublisher[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID]] interface {
	PublishMetrics(event E) error
}

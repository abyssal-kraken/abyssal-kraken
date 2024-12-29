package eventbus

import (
	"context"

	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
)

type EventBusDecorator interface {
	DecorateMultiple(
		ctx context.Context,
		event abyssalkraken.DomainEvent[abyssalkraken.AggregateID],
		eventHandlers []EventHandler,
		execution func() error,
	) func() error

	DecorateSingle(
		ctx context.Context,
		event abyssalkraken.DomainEvent[abyssalkraken.AggregateID],
		eventHandler EventHandler,
		execution func() error,
	) func() error
}

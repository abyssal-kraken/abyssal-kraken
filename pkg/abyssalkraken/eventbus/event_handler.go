package eventbus

import (
	"context"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
)

type EventHandler interface {
	Handle(ctx context.Context, event abyssalkraken.DomainEvent[abyssalkraken.AggregateID]) error
}

package serialization

import (
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"reflect"
)

type SnapshotSerialization[
	ID abyssalkraken.AggregateID,
	E abyssalkraken.DomainEvent[ID],
	A abyssalkraken.AggregateRoot[ID, E],
] interface {
	Serialize(aggregateRoot A, aggregateRootClass reflect.Type) ([]byte, error)

	Deserialize(data []byte, aggregateRootClass reflect.Type) (A, error)
}

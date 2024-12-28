package serialization

import "github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"

type SnapshotSerialization[
	ID abyssalkraken.AggregateID,
	E abyssalkraken.DomainEvent[ID],
	A abyssalkraken.AggregateRoot[ID, E],
] interface {
	Serialize(aggregateRoot A) ([]byte, error)

	Deserialize(data []byte) (A, error)
}

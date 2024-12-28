package serialization

import "github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"

type SnapshotSerialization[
	ID abyssalkraken.AggregateID,
	E abyssalkraken.DomainEvent[ID],
	A abyssalkraken.AggregateRoot[ID, E],
] interface {
	Serialize(snapshot abyssalkraken.Snapshot[ID, E, A]) ([]byte, error)
	Deserialize(data []byte) (abyssalkraken.Snapshot[ID, E, A], error)
}

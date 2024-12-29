package serialization

import "github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"

type EventStreamSerialization[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID]] interface {
	Serialize(events []E) ([]byte, error)
	Deserialize(data []byte) ([]E, error)
}

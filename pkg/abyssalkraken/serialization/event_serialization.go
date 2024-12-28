package serialization

import "github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"

type EventSerialization[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID]] interface {
	Serialize(events []E) ([]byte, error)
	Deserialize(data []byte) ([]E, error)
}

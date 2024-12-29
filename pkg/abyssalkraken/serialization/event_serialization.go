package serialization

import (
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"reflect"
)

type EventSerialization[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID]] interface {
	Serialize(events E, eventClass reflect.Type) ([]byte, error)
	Deserialize(data []byte, eventClass reflect.Type) (E, error)
}

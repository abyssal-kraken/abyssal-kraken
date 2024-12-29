package abyssalkraken

import "time"

type DomainEvent[ID AggregateID] interface {
	GetAggregateID() ID
	GetEventID() EventID
	GetEventType() EventType
	GetOccurredOn() time.Time
	GetMetadata() map[string]string
}

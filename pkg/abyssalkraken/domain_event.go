package abyssalkraken

import "time"

type DomainEvent[ID AggregateID] interface {
	AggregateID() ID
	EventID() EventID
	EventType() EventType
	OccurredOn() time.Time
	Metadata() map[string]string
}

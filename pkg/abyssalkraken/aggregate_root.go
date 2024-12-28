package abyssalkraken

type AggregateID interface {
	String() string
}

type AggregateRoot[ID AggregateID, E DomainEvent[ID]] interface {
	ID() ID
	AddEvent(event E)
	HasPendingEvents() bool
	CollectPendingEvents() []E
}

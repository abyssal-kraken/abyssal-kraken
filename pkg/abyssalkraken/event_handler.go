package abyssalkraken

type EventHandler[ID AggregateID, E DomainEvent[ID]] interface {
	HandleEvent(event E) error
}

package abyssalkraken

type EventPublisher[ID AggregateID, E DomainEvent[ID]] interface {
	PublishEvent(event E) error
}

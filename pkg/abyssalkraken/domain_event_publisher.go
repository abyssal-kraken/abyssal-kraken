package abyssalkraken

type DomainEventPublisher[ID AggregateID, E DomainEvent[ID]] interface {
	Publish(event E) error
}

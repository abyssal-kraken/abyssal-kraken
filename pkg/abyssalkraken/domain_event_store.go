package abyssalkraken

type DomainEventStore[ID AggregateID, E DomainEvent[ID]] interface {
	FindByAggregateID(aggregateID ID) ([]E, error)
	FindByEventID(eventID string) (E, error)
	Save(event E) error
}

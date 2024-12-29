package abyssalkraken

type AggregateID interface {
	String() string
}

type AggregateRoot[ID AggregateID, E DomainEvent[ID]] interface {
	ID() ID
	Type() AggregateType
	HasChanges() bool
	CollectChanges() []E
	Mutate(event E)
	Apply(event E)
	ReplayEvents(events []E)
}

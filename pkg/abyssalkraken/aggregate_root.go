package abyssalkraken

type AggregateID interface {
	String() string
}

type AggregateRoot[ID AggregateID, E DomainEvent[ID]] interface {
	ID() ID
	Type() AggregateType
	HasChanges() bool
	CollectChanges() []E
	Mutate(event E) *AggregateRoot[ID, E[ID]]
	Apply(event E) *AggregateRoot[ID, E[ID]]
	ReplayEvents(events []E) *AggregateRoot[ID, E[ID]]
}

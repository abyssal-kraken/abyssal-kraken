package abyssalkraken

type Snapshot[
	ID AggregateID,
	E DomainEvent[ID],
	A AggregateRoot[ID, E],
] struct {
	AggregateRoot A
	Version       Version
}

func TakeSnapshot[
	ID AggregateID,
	E DomainEvent[ID],
	A AggregateRoot[ID, E],
](aggregateRoot A, version Version) Snapshot[ID, E, A] {
	return Snapshot[ID, E, A]{
		AggregateRoot: aggregateRoot,
		Version:       version,
	}
}

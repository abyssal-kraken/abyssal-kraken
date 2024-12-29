package abyssalkraken

type AggregateVersion[ID AggregateID, E DomainEvent[ID], A AggregateRoot[ID, E]] struct {
	AggregateRoot A
	Version       Version
}

func WithVersion[ID AggregateID, E DomainEvent[ID], A AggregateRoot[ID, E]](aggregateRoot A, version Version) AggregateVersion[ID, E, A] {
	return AggregateVersion[ID, E, A]{
		AggregateRoot: aggregateRoot,
		Version:       version,
	}
}

func (av AggregateVersion[ID, E, A]) Map(block func(A) A) AggregateVersion[ID, E, A] {
	return WithVersion(block(av.AggregateRoot), av.Version)
}

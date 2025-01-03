package abyssalkraken

import (
	"sort"
)

type SimpleAggregateRoot[ID AggregateID, E DomainEvent[ID]] struct {
	id            ID
	aggregateType AggregateType
	changes       []E
}

func NewSimpleAggregateRoot[ID AggregateID, E DomainEvent[ID]](id ID, aggregateType AggregateType) *SimpleAggregateRoot[ID, E] {
	return &SimpleAggregateRoot[ID, E]{
		id:            id,
		aggregateType: aggregateType,
		changes:       make([]E, 0),
	}
}

func (a *SimpleAggregateRoot[ID, E]) ID() ID {
	return a.id
}

func (a *SimpleAggregateRoot[ID, E]) Type() AggregateType {
	return a.aggregateType
}

func (a *SimpleAggregateRoot[ID, E]) HasChanges() bool {
	return len(a.changes) > 0
}

func (a *SimpleAggregateRoot[ID, E]) CollectChanges() []E {
	sort.Slice(a.changes, func(i, j int) bool {
		return a.changes[i].GetOccurredOn().Before(a.changes[j].GetOccurredOn())
	})

	collected := a.changes
	a.changes = nil
	return collected
}

func (a *SimpleAggregateRoot[ID, E]) Mutate(event E) {
	panic("Mutate must be implemented by the aggregate root")
}

func (a *SimpleAggregateRoot[ID, E]) ApplyChange(event E) {
	a.changes = append(a.changes, event)

	a.Mutate(event)
}

func (a *SimpleAggregateRoot[ID, E]) ReplayEvents(events []E) {
	for _, event := range events {
		a.Mutate(event)
	}
}

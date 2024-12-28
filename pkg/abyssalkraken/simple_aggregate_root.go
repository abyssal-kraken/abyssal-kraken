package abyssalkraken

import (
	"sort"
	"sync"
)

type SimpleAggregateRoot[ID AggregateID, E DomainEvent[ID]] struct {
	mu            sync.Mutex
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
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.changes) > 0
}

func (a *SimpleAggregateRoot[ID, E]) CollectChanges() []E {
	a.mu.Lock()
	defer a.mu.Unlock()

	sort.Slice(a.changes, func(i, j int) bool {
		return a.changes[i].OccurredOn().Before(a.changes[j].OccurredOn())
	})

	collected := a.changes
	a.changes = nil
	return collected
}

func (a *SimpleAggregateRoot[ID, E]) Mutate(event E) {
	panic("Mutate must be implemented by the aggregate root")
}

func (a *SimpleAggregateRoot[ID, E]) Apply(event E) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.changes = append(a.changes, event)

	a.Mutate(event)
}

func (a *SimpleAggregateRoot[ID, E]) ReplayEvents(events []E) {
	for _, event := range events {
		a.Mutate(event)
	}
}

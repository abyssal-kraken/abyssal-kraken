package abyssalkraken

import (
	"sync"
)

type SimpleAggregateRoot[ID AggregateID, E DomainEvent[ID]] struct {
	mu            sync.Mutex
	aggregateID   ID
	pendingEvents []E
}

func NewSimpleAggregateRoot[ID AggregateID, E DomainEvent[ID]](id ID) *SimpleAggregateRoot[ID, E] {
	return &SimpleAggregateRoot[ID, E]{
		aggregateID:   id,
		pendingEvents: []E{},
	}
}

func (a *SimpleAggregateRoot[ID, E]) ID() ID {
	return a.aggregateID
}

func (a *SimpleAggregateRoot[ID, E]) AddEvent(event E) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.pendingEvents = append(a.pendingEvents, event)
}

func (a *SimpleAggregateRoot[ID, E]) HasPendingEvents() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.pendingEvents) > 0
}

func (a *SimpleAggregateRoot[ID, E]) CollectPendingEvents() []E {
	a.mu.Lock()
	defer a.mu.Unlock()

	events := a.pendingEvents
	a.pendingEvents = nil
	return events
}

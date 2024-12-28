package in_memory

import (
	"account/account/core/domain/account"
	"errors"
)

type InMemoryDomainEventStore[ID account.AccountID, E account.AccountEvent] struct {
	eventsByAggregateID map[string][]E
	eventsByEventID     map[string]E
}

func NewInMemoryDomainEventStore[ID account.AccountID, E account.AccountEvent]() *InMemoryDomainEventStore[ID, E] {
	return &InMemoryDomainEventStore[ID, E]{
		eventsByAggregateID: make(map[string][]E),
		eventsByEventID:     make(map[string]E),
	}
}

func (s *InMemoryDomainEventStore[ID, E]) FindByAggregateID(aggregateID ID) ([]E, error) {
	events, exists := s.eventsByAggregateID[string(aggregateID)]
	if !exists {
		return nil, errors.New("no events found for the given aggregate ID")
	}
	return events, nil
}

func (s *InMemoryDomainEventStore[ID, E]) FindByEventID(eventID string) (E, error) {
	event, exists := s.eventsByEventID[eventID]
	if !exists {
		var zeroValue E
		return zeroValue, errors.New("event not found for the given event ID")
	}
	return event, nil
}

func (s *InMemoryDomainEventStore[ID, E]) Save(event E) error {
	aggregateID := event.AggregateID().String()
	s.eventsByAggregateID[aggregateID] = append(s.eventsByAggregateID[aggregateID], event)

	eventID := event.EventID()
	if _, exists := s.eventsByEventID[eventID.String()]; exists {
		return errors.New("event with the same ID already exists")
	}
	s.eventsByEventID[eventID.String()] = event

	return nil
}

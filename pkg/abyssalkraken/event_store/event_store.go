package event_store

import (
	"context"
	"errors"
	"fmt"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken/persistence"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken/serialization"
	"sync"
)

type EventStore[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID]] struct {
	eventPersistence   persistence.EventPersistence
	eventSerialization serialization.EventSerialization[ID, E]
	locks              sync.Map
}

func NewEventStore[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID]](
	eventPersistence persistence.EventPersistence,
	eventSerialization serialization.EventSerialization[ID, E],
) *EventStore[ID, E] {
	return &EventStore[ID, E]{
		eventPersistence:   eventPersistence,
		eventSerialization: eventSerialization,
	}
}

func (es *EventStore[ID, E]) AppendToStream(ctx context.Context, aggregateID ID, expectedVersion abyssalkraken.Version, events []E) error {
	if len(events) == 0 {
		return nil
	}

	data, err := es.eventSerialization.Serialize(events)
	if err != nil {
		return fmt.Errorf("failed to serialize events: %w", err)
	}

	lock := es.lockAggregate(aggregateID)
	lock.Lock()
	defer lock.Unlock()

	expectedVersionInt := expectedVersion.ToInt()
	newVersion := expectedVersionInt + 1

	err = es.eventPersistence.Append(ctx, aggregateID.String(), data, expectedVersionInt, newVersion)
	if err != nil {
		var concurrencyErr *persistence.PersistenceConcurrencyError
		if errors.As(err, &concurrencyErr) {
			return &EventStoreConcurrencyException[ID]{
				AggregateID:     aggregateID,
				ExpectedVersion: int(concurrencyErr.ExpectedVersion),
				ActualVersion:   int(concurrencyErr.ActualVersion),
				InnerError:      err,
			}
		}
		return fmt.Errorf("failed to append events: %w", err)
	}

	return nil
}

func (es *EventStore[ID, E]) LoadEventStream(ctx context.Context, aggregateID ID) (abyssalkraken.EventStream[ID, E], error) {
	return es.loadEventStream(ctx, aggregateID, nil)
}

func (es *EventStore[ID, E]) LoadEventStreamAfterVersion(ctx context.Context, aggregateID ID, version abyssalkraken.Version) (abyssalkraken.EventStream[ID, E], error) {
	return es.loadEventStream(ctx, aggregateID, &version)
}

func (es *EventStore[ID, E]) loadEventStream(ctx context.Context, aggregateID ID, version *abyssalkraken.Version) (abyssalkraken.EventStream[ID, E], error) {
	var afterVersion *int64
	if version != nil {
		v := version.ToInt()
		afterVersion = &v
	}

	records, err := es.eventPersistence.ReadRecords(ctx, aggregateID.String(), afterVersion)
	if err != nil {
		return abyssalkraken.EventStream[ID, E]{}, fmt.Errorf("failed to read event records: %w", err)
	}

	stream := abyssalkraken.EmptyStream[ID, E]()
	results := make(chan abyssalkraken.EventStream[ID, E], len(records))
	wg := sync.WaitGroup{}
	wg.Add(len(records))

	for _, record := range records {
		go func(record persistence.BinaryData) {
			defer wg.Done()
			events, err := es.eventSerialization.Deserialize(record.Data)
			if err == nil {
				v, err := abyssalkraken.ToVersion(record.Version)
				if err == nil {
					results <- abyssalkraken.StreamOf(v, events)
				}
			}
		}(record)
	}

	wg.Wait()
	close(results)

	for partial := range results {
		stream = stream.Plus(partial)
	}

	return stream, nil
}

func (es *EventStore[ID, E]) lockAggregate(aggregateID ID) *sync.Mutex {
	lock, _ := es.locks.LoadOrStore(aggregateID.String(), &sync.Mutex{})
	return lock.(*sync.Mutex)
}

type EventStoreConcurrencyException[ID abyssalkraken.AggregateID] struct {
	AggregateID     ID
	ExpectedVersion int
	ActualVersion   int
	InnerError      error
}

func (e *EventStoreConcurrencyException[ID]) Error() string {
	return fmt.Sprintf(
		"Concurrency failure when appending events to the stream for aggregate ID %s. Expected version: %d, Actual version: %d",
		e.AggregateID.String(),
		e.ExpectedVersion,
		e.ActualVersion,
	)
}

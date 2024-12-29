package repository

import (
	"context"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken/event_store"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken/snapshot_repository"
	"reflect"
)

type AggregateRepository[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID], A abyssalkraken.AggregateRoot[ID, E]] struct {
	eventStore         event_store.EventStore[ID, E]
	snapshotRepository snapshot_repository.SnapshotRepository[ID, E, A]
}

func (r *AggregateRepository[ID, E, A]) CreateAggregateFromFirstEvent(event E) A {
	panic("createAggregateFromFirstEvent must be implemented by the concrete repository")
}

func (r *AggregateRepository[ID, E, A]) LoadAggregateById(
	ctx context.Context,
	aggregateId ID,
	aggregateRootClass reflect.Type,
	eventClass reflect.Type,
) (*abyssalkraken.AggregateVersion[ID, E, A], error) {
	snapshot, err := r.snapshotRepository.FindSnapshot(ctx, aggregateId, aggregateRootClass)
	if err != nil {
		return nil, err
	}
	if snapshot != nil {
		eventStream, err := r.eventStore.LoadEventStreamAfterVersion(ctx, aggregateId, eventClass, snapshot.Version)
		if err != nil {
			return nil, err
		}
		updatedAggregate := snapshot.AggregateRoot
		updatedAggregate.ReplayEvents(eventStream.Events)

		return &abyssalkraken.AggregateVersion[ID, E, A]{
			AggregateRoot: updatedAggregate,
			Version:       eventStream.Version,
		}, nil
	}

	eventStream, err := r.eventStore.LoadEventStream(ctx, aggregateId, eventClass)
	if err != nil {
		return nil, err
	}
	if len(eventStream.Events) == 0 {
		return nil, nil
	}

	aggregateRoot := r.CreateAggregateFromFirstEvent(eventStream.Events[0])
	aggregateRoot.ReplayEvents(eventStream.Events)

	return &abyssalkraken.AggregateVersion[ID, E, A]{
		AggregateRoot: aggregateRoot,
		Version:       eventStream.Version,
	}, nil
}

func (r *AggregateRepository[ID, E, A]) LoadAggregateByIdShortcut(
	ctx context.Context,
	aggregateId ID,
) (*abyssalkraken.AggregateVersion[ID, E, A], error) {
	eventType := reflect.TypeOf((*E)(nil)).Elem()
	aggregateType := reflect.TypeOf((*A)(nil)).Elem()
	return r.LoadAggregateById(ctx, aggregateId, aggregateType, eventType)
}

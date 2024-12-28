package repository

import (
	"context"
	"fmt"

	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken/persistence"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken/serialization"
)

type SnapshotRepository[
	ID abyssalkraken.AggregateID,
	E abyssalkraken.DomainEvent[ID],
	A abyssalkraken.AggregateRoot[ID, E],
] struct {
	snapshotPersistence   persistence.SnapshotPersistence
	snapshotSerialization serialization.SnapshotSerialization[ID, E, A]
}

func NewSnapshotRepository[
	ID abyssalkraken.AggregateID,
	E abyssalkraken.DomainEvent[ID],
	A abyssalkraken.AggregateRoot[ID, E],
](
	snapshotPersistence persistence.SnapshotPersistence,
	snapshotSerialization serialization.SnapshotSerialization[ID, E, A],
) *SnapshotRepository[ID, E, A] {
	return &SnapshotRepository[ID, E, A]{
		snapshotPersistence:   snapshotPersistence,
		snapshotSerialization: snapshotSerialization,
	}
}

func (r *SnapshotRepository[ID, E, A]) SaveSnapshot(
	ctx context.Context,
	snapshot abyssalkraken.Snapshot[ID, E, A],
	expectedVersion abyssalkraken.Version,
) error {
	data, err := r.snapshotSerialization.Serialize(snapshot.AggregateRoot)
	if err != nil {
		return fmt.Errorf("failed to serialize snapshot: %w", err)
	}

	err = r.snapshotPersistence.Upsert(
		ctx,
		snapshot.AggregateRoot.ID().String(),
		data,
		expectedVersion.ToInt(),
		snapshot.Version.ToInt(),
	)
	if err != nil {
		var concurrencyErr *persistence.PersistenceConcurrencyError
		if ok := persistence.AsConcurrencyError(err, &concurrencyErr); ok {
			return &SnapshotConcurrencyException[ID]{
				AggregateID:     snapshot.AggregateRoot.ID(),
				ExpectedVersion: abyssalkraken.Version{Value: concurrencyErr.ExpectedVersion},
				ActualVersion:   abyssalkraken.Version{Value: concurrencyErr.ActualVersion},
				InnerError:      concurrencyErr,
			}
		}
		return fmt.Errorf("failed to upsert snapshot: %w", err)
	}

	return nil
}

func (r *SnapshotRepository[ID, E, A]) FindSnapshot(
	ctx context.Context,
	aggregateID ID,
) (*abyssalkraken.Snapshot[ID, E, A], error) {
	record, err := r.snapshotPersistence.ReadRecord(ctx, aggregateID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot record: %w", err)
	}
	if record == nil {
		return nil, nil
	}

	aggregateRoot, err := r.snapshotSerialization.Deserialize(record.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize snapshot: %w", err)
	}

	return &abyssalkraken.Snapshot[ID, E, A]{
		AggregateRoot: aggregateRoot,
		Version:       abyssalkraken.Version{Value: record.Version},
	}, nil
}

type SnapshotConcurrencyException[ID abyssalkraken.AggregateID] struct {
	AggregateID     ID
	ExpectedVersion abyssalkraken.Version
	ActualVersion   abyssalkraken.Version
	InnerError      error
}

func (e *SnapshotConcurrencyException[ID]) Error() string {
	return fmt.Sprintf(
		"Concurrency failure when saving snapshot for aggregate id %s. Expected version: %d, Actual version: %d",
		e.AggregateID.String(),
		e.ExpectedVersion.ToInt(),
		e.ActualVersion.ToInt(),
	)
}

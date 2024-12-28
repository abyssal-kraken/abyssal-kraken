package persistence

import "context"

type SnapshotPersistence interface {
	Upsert(ctx context.Context, name string, data []byte, expectedVersion int64, newVersion int64) error
	ReadRecord(ctx context.Context, name string) (*BinaryData, error)
}

package persistence

import "context"

type EventPersistence interface {
	Append(ctx context.Context, name string, data []byte, expectedVersion int64, newVersion int64) error
	ReadRecords(ctx context.Context, name string, afterVersion *int64) ([]BinaryData, error)
}

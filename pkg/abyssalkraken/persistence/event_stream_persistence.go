package persistence

import "context"

type EventStreamPersistence interface {
	Append(ctx context.Context, name string, data []byte, expectedVersion int64, newVersion int64) error
	ReadStream(ctx context.Context, name string, afterVersion *int64) ([]BinaryData, error)
}

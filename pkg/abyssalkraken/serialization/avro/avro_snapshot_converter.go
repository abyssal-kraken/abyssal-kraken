package avro

import (
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"github.com/hamba/avro/v2"
	"reflect"
)

type AvroSnapshotConverter[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID], A abyssalkraken.AggregateRoot[ID, E], GC any] interface {
	AggregateRootType() reflect.Type

	AvroSchema() avro.Schema

	ToAvroSchema(aggregateRoot A) (GC, error)

	FromAvroSchema(avroContainer GC) (A, error)
}

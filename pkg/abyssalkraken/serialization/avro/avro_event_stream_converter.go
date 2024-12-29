package avro

import (
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"github.com/hamba/avro/v2"
)

type AvroEventStreamConverter[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID], GC any] interface {
	AvroSchema() avro.Schema

	ToAvroSchema(events []E) (GC, error)

	FromAvroSchema(avroContainer GC) ([]E, error)
}

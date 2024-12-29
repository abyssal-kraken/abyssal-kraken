package avro

import (
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"github.com/hamba/avro/v2"
)

type AvroEventConverter[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID], GC any] interface {
	AvroSchema() avro.Schema

	ToAvroSchema(event E) (GC, error)

	FromAvroSchema(avroContainer GC) (E, error)
}

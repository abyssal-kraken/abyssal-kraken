package avro

import (
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"reflect"
)

type AvroEventConverterRegistry[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID], GC any] struct {
	convertersMap map[reflect.Type]AvroEventConverter[ID, E, GC]
}

func NewAvroEventConverterRegistry[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID], GC any](
	converters []AvroEventConverter[ID, E, GC],
) *AvroEventConverterRegistry[ID, E, GC] {
	convertersMap := make(map[reflect.Type]AvroEventConverter[ID, E, GC])
	for _, converter := range converters {
		convertersMap[converter.EventType()] = converter
	}
	return &AvroEventConverterRegistry[ID, E, GC]{convertersMap: convertersMap}
}

func (r *AvroEventConverterRegistry[ID, E, GC]) FindConverter(
	eventClass reflect.Type,
) (AvroEventConverter[ID, E, GC], error) {
	converter, exists := r.convertersMap[eventClass]
	if !exists {
		return nil, &AvroEventConverterNotFoundException{
			EventClass: eventClass,
		}
	}
	return converter, nil
}

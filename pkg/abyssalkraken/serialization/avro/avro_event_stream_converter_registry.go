package avro

import (
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"reflect"
)

type AvroEventStreamConverterRegistry[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID], GC any] struct {
	convertersMap map[reflect.Type]AvroEventStreamConverter[ID, E, GC]
}

func NewAvroEventStreamConverterRegistry[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID], GC any](
	converters []AvroEventStreamConverter[ID, E, GC],
) *AvroEventStreamConverterRegistry[ID, E, GC] {
	convertersMap := make(map[reflect.Type]AvroEventStreamConverter[ID, E, GC])
	for _, converter := range converters {
		convertersMap[converter.EventType()] = converter
	}
	return &AvroEventStreamConverterRegistry[ID, E, GC]{convertersMap: convertersMap}
}

func (r *AvroEventStreamConverterRegistry[ID, E, GC]) FindConverter(
	eventClass reflect.Type,
) (AvroEventStreamConverter[ID, E, GC], error) {
	converter, exists := r.convertersMap[eventClass]
	if !exists {
		return nil, &AvroEventStreamConverterNotFoundException{
			EventClass: eventClass,
		}
	}
	return converter, nil
}

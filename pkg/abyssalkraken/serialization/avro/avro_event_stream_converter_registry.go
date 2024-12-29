package avro

import (
	"fmt"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"reflect"
)

type AvroEventStreamConverterRegistry[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID], GC any] struct {
	converters map[reflect.Type]AvroEventStreamConverter[ID, E, GC]
}

func NewAvroEventStreamConverterRegistry[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID], GC any]() *AvroEventStreamConverterRegistry[ID, E, GC] {
	return &AvroEventStreamConverterRegistry[ID, E, GC]{
		converters: make(map[reflect.Type]AvroEventStreamConverter[ID, E, GC]),
	}
}

func (r *AvroEventStreamConverterRegistry[ID, E, GC]) Register(eventType reflect.Type, converter AvroEventStreamConverter[ID, E, GC]) error {
	if _, exists := r.converters[eventType]; exists {
		return fmt.Errorf("event stream converter already registered for event type %s", eventType.String())
	}

	r.converters[eventType] = converter
	return nil
}

func (r *AvroEventStreamConverterRegistry[ID, E, GC]) FindConverter(
	eventType reflect.Type,
) (AvroEventStreamConverter[ID, E, GC], error) {
	converter, exists := r.converters[eventType]
	if !exists {
		return nil, &AvroEventStreamConverterNotFoundException{
			EventType: eventType,
		}
	}
	return converter, nil
}

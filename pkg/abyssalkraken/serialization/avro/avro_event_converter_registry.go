package avro

import (
	"fmt"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"reflect"
)

type AvroEventConverterRegistry[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID], GC any] struct {
	converters map[reflect.Type]AvroEventConverter[ID, E, GC]
}

func NewAvroEventConverterRegistry[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID], GC any]() *AvroEventConverterRegistry[ID, E, GC] {
	return &AvroEventConverterRegistry[ID, E, GC]{
		converters: make(map[reflect.Type]AvroEventConverter[ID, E, GC]),
	}
}

func (r *AvroEventConverterRegistry[ID, E, GC]) Register(eventType reflect.Type, converter AvroEventConverter[ID, E, GC]) error {
	if _, exists := r.converters[eventType]; exists {
		return fmt.Errorf("event converter already registered for event type %s", eventType.String())
	}

	r.converters[eventType] = converter
	return nil
}

func (r *AvroEventConverterRegistry[ID, E, GC]) FindConverter(eventType reflect.Type) (AvroEventConverter[ID, E, GC], error) {
	converter, exists := r.converters[eventType]
	if !exists {
		return nil, &AvroEventConverterNotFoundException{
			EventType: eventType,
		}
	}
	return converter, nil
}

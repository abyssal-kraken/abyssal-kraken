package avro

import (
	"fmt"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"reflect"
)

type AvroSnapshotConverterRegistry[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID], A abyssalkraken.AggregateRoot[ID, E], GC any] struct {
	converters map[reflect.Type]AvroSnapshotConverter[ID, E, A, GC]
}

func NewAvroSnapshotConverterRegistry[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID], A abyssalkraken.AggregateRoot[ID, E], GC any]() *AvroSnapshotConverterRegistry[ID, E, A, GC] {
	return &AvroSnapshotConverterRegistry[ID, E, A, GC]{
		converters: make(map[reflect.Type]AvroSnapshotConverter[ID, E, A, GC]),
	}
}

func (r *AvroSnapshotConverterRegistry[ID, E, A, GC]) Register(aggregateRootType reflect.Type, converter AvroSnapshotConverter[ID, E, A, GC]) error {
	if _, exists := r.converters[aggregateRootType]; exists {
		return fmt.Errorf("snapshot converter already registered for aggregate root type %s", aggregateRootType.String())
	}

	r.converters[aggregateRootType] = converter
	return nil
}

func (r *AvroSnapshotConverterRegistry[ID, E, A, GC]) FindConverter(
	aggregateRootType reflect.Type,
) (AvroSnapshotConverter[ID, E, A, GC], error) {
	converter, exists := r.converters[aggregateRootType]
	if !exists {
		return nil, &AvroSnapshotConverterNotFoundException{
			AggregateRootType: aggregateRootType,
		}
	}
	return converter, nil
}

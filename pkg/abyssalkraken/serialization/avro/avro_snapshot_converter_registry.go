package avro

import (
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"reflect"
)

type AvroSnapshotConverterRegistry[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID], A abyssalkraken.AggregateRoot[ID, E], GC any] struct {
	convertersMap map[reflect.Type]AvroSnapshotConverter[ID, E, A, GC]
}

func NewAvroSnapshotConverterRegistry[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID], A abyssalkraken.AggregateRoot[ID, E], GC any](
	converters []AvroSnapshotConverter[ID, E, A, GC],
) *AvroSnapshotConverterRegistry[ID, E, A, GC] {
	convertersMap := make(map[reflect.Type]AvroSnapshotConverter[ID, E, A, GC])
	for _, converter := range converters {
		convertersMap[converter.AggregateRootType()] = converter
	}
	return &AvroSnapshotConverterRegistry[ID, E, A, GC]{convertersMap: convertersMap}
}

func (r *AvroSnapshotConverterRegistry[ID, E, A, GC]) FindConverter(
	aggregateRootClass reflect.Type,
) (AvroSnapshotConverter[ID, E, A, GC], error) {
	converter, exists := r.convertersMap[aggregateRootClass]
	if !exists {
		return nil, &AvroSnapshotConverterNotFoundException{
			AggregateRootClass: aggregateRootClass,
		}
	}
	return converter, nil
}

package avro

import (
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"github.com/linkedin/goavro/v2"
	"reflect"
)

type AvroSnapshotSerialization[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID], A abyssalkraken.AggregateRoot[ID, E]] struct {
	converterRegistry *AvroSnapshotConverterRegistry[ID, E, A, map[string]interface{}]
}

func NewAvroSnapshotSerialization[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID], A abyssalkraken.AggregateRoot[ID, E]](
	converterRegistry *AvroSnapshotConverterRegistry[ID, E, A, map[string]interface{}],
) *AvroSnapshotSerialization[ID, E, A] {
	return &AvroSnapshotSerialization[ID, E, A]{converterRegistry: converterRegistry}
}

func (s *AvroSnapshotSerialization[ID, E, A]) Serialize(aggregateRoot A, aggregateRootType reflect.Type) ([]byte, error) {
	converter, err := s.converterRegistry.FindConverter(aggregateRootType)
	if err != nil {
		return nil, err
	}

	avroData, err := converter.ToAvroSchema(aggregateRoot)
	if err != nil {
		return nil, err
	}

	codec, err := goavro.NewCodec(converter.AvroSchema().String())
	if err != nil {
		return nil, err
	}

	binaryData, err := codec.BinaryFromNative(nil, avroData)
	if err != nil {
		return nil, err
	}

	return binaryData, nil
}

func (s *AvroSnapshotSerialization[ID, E, A]) Deserialize(data []byte, aggregateRootType reflect.Type) (A, error) {
	converter, err := s.converterRegistry.FindConverter(aggregateRootType)
	if err != nil {
		var zeroAggregateRoot A
		return zeroAggregateRoot, err
	}

	codec, err := goavro.NewCodec(converter.AvroSchema().String())
	if err != nil {
		var zeroAggregateRoot A
		return zeroAggregateRoot, err
	}

	nativeData, _, err := codec.NativeFromBinary(data)
	if err != nil {
		var zeroAggregateRoot A
		return zeroAggregateRoot, err
	}

	aggregateRoot, err := converter.FromAvroSchema(nativeData.(map[string]interface{}))
	if err != nil {
		var zeroAggregateRoot A
		return zeroAggregateRoot, err
	}

	return aggregateRoot, nil
}

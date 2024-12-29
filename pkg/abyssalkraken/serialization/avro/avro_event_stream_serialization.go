package avro

import (
	"errors"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"github.com/linkedin/goavro/v2"
	"reflect"
)

type AvroEventStreamSerialization[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID]] struct {
	converterRegistry *AvroEventStreamConverterRegistry[ID, E, []map[string]interface{}]
}

func NewAvroEventStreamSerialization[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID]](
	converterRegistry *AvroEventStreamConverterRegistry[ID, E, []map[string]interface{}],
) *AvroEventStreamSerialization[ID, E] {
	return &AvroEventStreamSerialization[ID, E]{converterRegistry: converterRegistry}
}

func (s *AvroEventStreamSerialization[ID, E]) Serialize(events []E, eventType reflect.Type) ([]byte, error) {
	if len(events) == 0 {
		return nil, errors.New("no events provided for serialization")
	}

	converter, err := s.converterRegistry.FindConverter(eventType)
	if err != nil {
		return nil, err
	}

	avroData, err := converter.ToAvroSchema(events)
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

func (s *AvroEventStreamSerialization[ID, E]) Deserialize(data []byte, eventType reflect.Type) ([]E, error) {
	converter, err := s.converterRegistry.FindConverter(eventType)
	if err != nil {
		return nil, err
	}

	codec, err := goavro.NewCodec(converter.AvroSchema().String())
	if err != nil {
		return nil, err
	}

	nativeData, _, err := codec.NativeFromBinary(data)
	if err != nil {
		return nil, err
	}

	events, err := converter.FromAvroSchema(nativeData.([]map[string]interface{}))
	if err != nil {
		return nil, err
	}

	return events, nil
}

package avro

import (
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"github.com/linkedin/goavro/v2"
	"reflect"
)

type AvroEventSerialization[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID]] struct {
	converterRegistry *AvroEventConverterRegistry[ID, E, map[string]interface{}]
}

func NewAvroEventSerialization[ID abyssalkraken.AggregateID, E abyssalkraken.DomainEvent[ID]](
	converterRegistry *AvroEventConverterRegistry[ID, E, map[string]interface{}],
) *AvroEventSerialization[ID, E] {
	return &AvroEventSerialization[ID, E]{converterRegistry: converterRegistry}
}

func (s *AvroEventSerialization[ID, E]) Serialize(event E, eventType reflect.Type) ([]byte, error) {
	converter, err := s.converterRegistry.FindConverter(eventType)
	if err != nil {
		return nil, err
	}

	avroData, err := converter.ToAvroSchema(event)
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

func (s *AvroEventSerialization[ID, E]) Deserialize(data []byte, eventType reflect.Type) (E, error) {
	converter, err := s.converterRegistry.FindConverter(eventType)
	if err != nil {
		var zeroEvent E
		return zeroEvent, err
	}

	codec, err := goavro.NewCodec(converter.AvroSchema().String())
	if err != nil {
		var zeroEvent E
		return zeroEvent, err
	}

	avroData, _, err := codec.NativeFromBinary(data)
	if err != nil {
		var zeroEvent E
		return zeroEvent, err
	}

	event, err := converter.FromAvroSchema(avroData.(map[string]interface{}))
	if err != nil {
		var zeroEvent E
		return zeroEvent, err
	}

	return event, nil
}

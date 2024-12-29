package avro

import (
	"fmt"
	"reflect"
)

type AvroEventConverterNotFoundException struct {
	EventType reflect.Type
}

func (e *AvroEventConverterNotFoundException) Error() string {
	return fmt.Sprintf("No Avro event converter found for event type %s", e.EventType.String())
}

type AvroEventStreamConverterNotFoundException struct {
	EventType reflect.Type
}

func (e *AvroEventStreamConverterNotFoundException) Error() string {
	return fmt.Sprintf("No Avro event stream converter found for event type %s", e.EventType.String())
}

type AvroSnapshotConverterNotFoundException struct {
	AggregateRootType reflect.Type
}

func (e *AvroSnapshotConverterNotFoundException) Error() string {
	return fmt.Sprintf("No Avro snapshot converter found for aggregate root type %s", e.AggregateRootType.String())
}

package avro

import (
	"fmt"
	"reflect"
)

type AvroEventConverterNotFoundException struct {
	EventClass reflect.Type
}

func (e *AvroEventConverterNotFoundException) Error() string {
	return fmt.Sprintf("No Avro event converter found for event class %s", e.EventClass.String())
}

type AvroEventStreamConverterNotFoundException struct {
	EventClass reflect.Type
}

func (e *AvroEventStreamConverterNotFoundException) Error() string {
	return fmt.Sprintf("No Avro event stream converter found for event class %s", e.EventClass.String())
}

type AvroSnapshotConverterNotFoundException struct {
	AggregateRootClass reflect.Type
}

func (e *AvroSnapshotConverterNotFoundException) Error() string {
	return fmt.Sprintf("No Avro snapshot converter found for aggregate root class %s", e.AggregateRootClass.String())
}

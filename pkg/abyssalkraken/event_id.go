package abyssalkraken

import (
	"github.com/google/uuid"
)

type EventID struct {
	id string
}

func NewEventID(id string) (*EventID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, err
	}
	return &EventID{id: id}, nil
}

func RandomEventID() EventID {
	return EventID{id: uuid.NewString()}
}

func (e *EventID) ToUUID() uuid.UUID {
	return uuid.MustParse(e.id)
}

func (e *EventID) String() string {
	return e.id
}

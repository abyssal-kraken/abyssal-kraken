package abyssalkraken

type EventStream[ID AggregateID, E DomainEvent[ID]] struct {
	Version Version
	Events  []E
}

func (es EventStream[ID, E]) Plus(another EventStream[ID, E]) EventStream[ID, E] {
	combinedEvents := make([]E, len(es.Events)+len(another.Events))
	copy(combinedEvents, es.Events)
	copy(combinedEvents[len(es.Events):], another.Events)

	return EventStream[ID, E]{
		Version: another.Version,
		Events:  combinedEvents,
	}
}

func (es EventStream[ID, E]) IsEmpty() bool {
	return len(es.Events) == 0
}

func EmptyStream[ID AggregateID, E DomainEvent[ID]](version *Version) EventStream[ID, E] {
	if version != nil {
		return EventStream[ID, E]{
			Version: *version,
			Events:  make([]E, 0),
		}
	}

	return EventStream[ID, E]{
		Version: MinVersion,
		Events:  make([]E, 0),
	}
}

func StreamOf[ID AggregateID, E DomainEvent[ID]](version Version, events []E) EventStream[ID, E] {
	return EventStream[ID, E]{
		Version: version,
		Events:  events,
	}
}

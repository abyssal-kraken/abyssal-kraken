package abyssalkraken

type EventStream[ID AggregateID, E DomainEvent[ID]] struct {
	Version Version
	Events  map[string]E
}

func (es EventStream[ID, E]) Plus(another EventStream[ID, E]) EventStream[ID, E] {
	newEvents := make(map[string]E, len(es.Events))
	for k, v := range es.Events {
		newEvents[k] = v
	}

	for k, v := range another.Events {
		newEvents[k] = v
	}

	return EventStream[ID, E]{
		Version: another.Version,
		Events:  newEvents,
	}
}

func EmptyStream[ID AggregateID, E DomainEvent[ID]]() EventStream[ID, E] {
	return EventStream[ID, E]{
		Version: MinVersion,
		Events:  make(map[string]E),
	}
}

func StreamOf[ID AggregateID, E DomainEvent[ID]](version Version, events []E) EventStream[ID, E] {
	eventSet := make(map[string]E, len(events))
	for _, event := range events {
		id := event.GetEventID()
		eventSet[id.String()] = event // Usamos o GetEventID como chave única.
	}
	return EventStream[ID, E]{
		Version: version,
		Events:  eventSet,
	}
}

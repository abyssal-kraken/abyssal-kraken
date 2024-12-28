package abyssalkraken

import "sync"

type DomainEventHandler[
	A AggregateRoot[ID, E],
	ID AggregateID,
	E DomainEvent[ID],
] interface {
	Transactional() bool
	Handle(aggregateRoot *A, domainEvent E) error
}

func HandlePendingEvent[
	A AggregateRoot[ID, E],
	ID AggregateID,
	E DomainEvent[ID],
](
	domainEvent E,
	domainEventHandlers []DomainEventHandler[A, ID, E],
	aggregateRoot *A,
) error {
	var transactionalHandlers, nonTransactionalHandlers []DomainEventHandler[A, ID, E]

	for _, handler := range domainEventHandlers {
		if handler.Transactional() {
			transactionalHandlers = append(transactionalHandlers, handler)
		} else {
			nonTransactionalHandlers = append(nonTransactionalHandlers, handler)
		}
	}

	for _, handler := range transactionalHandlers {
		if err := handler.Handle(aggregateRoot, domainEvent); err != nil {
			return err
		}
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(nonTransactionalHandlers))

	for _, handler := range nonTransactionalHandlers {
		wg.Add(1)
		go func(h DomainEventHandler[A, ID, E]) {
			defer wg.Done()
			if err := h.Handle(aggregateRoot, domainEvent); err != nil {
				errors <- err
			}
		}(handler)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		return err
	}

	return nil
}

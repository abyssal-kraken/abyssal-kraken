package persistence

import "fmt"

type PersistenceConcurrencyError struct {
	Name            string
	ExpectedVersion int64
	ActualVersion   int64
}

func (e *PersistenceConcurrencyError) Error() string {
	return fmt.Sprintf(
		"concurrency error on %s: expected version %d but got %d",
		e.Name, e.ExpectedVersion, e.ActualVersion,
	)
}

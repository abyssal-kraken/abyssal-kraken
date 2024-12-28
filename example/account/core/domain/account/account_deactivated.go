package account

import (
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"time"
)

type AccountDeactivatedEvent struct {
	accountID    AccountID
	eventID      abyssalkraken.EventID
	eventType    abyssalkraken.EventType
	eventVersion int
	occurredOn   time.Time
	metadata     map[string]string
	Active       bool
}

func NewAccountDeactivatedEvent(accountID AccountID) *AccountDeactivatedEvent {
	return &AccountDeactivatedEvent{
		accountID:    accountID,
		eventID:      abyssalkraken.RandomEventID(),
		eventType:    AccountDeactivated,
		eventVersion: 1,
		occurredOn:   time.Now(),
		metadata:     map[string]string{},
		Active:       false,
	}
}

func (e AccountDeactivatedEvent) AggregateID() AccountID             { return e.accountID }
func (e AccountDeactivatedEvent) EventID() abyssalkraken.EventID     { return e.eventID }
func (e AccountDeactivatedEvent) EventType() abyssalkraken.EventType { return e.eventType }
func (e AccountDeactivatedEvent) EventVersion() int                  { return e.eventVersion }
func (e AccountDeactivatedEvent) OccurredOn() time.Time              { return e.occurredOn }
func (e AccountDeactivatedEvent) Metadata() map[string]string        { return e.metadata }

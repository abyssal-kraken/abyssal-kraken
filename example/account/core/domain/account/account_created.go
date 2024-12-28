package account

import (
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"time"
)

type AccountCreatedEvent struct {
	accountID    AccountID
	eventID      abyssalkraken.EventID
	eventType    abyssalkraken.EventType
	eventVersion int
	occurredOn   time.Time
	metadata     map[string]string
	Name         string
	Email        string
	Active       bool
}

func NewAccountCreatedEvent(accountID AccountID, name, email string) *AccountCreatedEvent {
	return &AccountCreatedEvent{
		accountID:    accountID,
		eventID:      abyssalkraken.RandomEventID(),
		eventType:    AccountCreated,
		eventVersion: 1,
		occurredOn:   time.Now(),
		metadata:     map[string]string{},
		Name:         name,
		Email:        email,
		Active:       true,
	}
}

func (e AccountCreatedEvent) AggregateID() AccountID             { return e.accountID }
func (e AccountCreatedEvent) EventID() abyssalkraken.EventID     { return e.eventID }
func (e AccountCreatedEvent) EventType() abyssalkraken.EventType { return e.eventType }
func (e AccountCreatedEvent) EventVersion() int                  { return e.eventVersion }
func (e AccountCreatedEvent) OccurredOn() time.Time              { return e.occurredOn }
func (e AccountCreatedEvent) Metadata() map[string]string        { return e.metadata }

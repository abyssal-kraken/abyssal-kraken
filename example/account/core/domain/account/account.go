package account

import (
	"errors"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"time"
)

type AccountEvent interface {
	abyssalkraken.DomainEvent[AccountID]
}

type Account struct {
	root      *abyssalkraken.SimpleAggregateRoot[AccountID, AccountEvent]
	Name      string
	Email     string
	Active    bool
	CreatedAt time.Time
}

func (a *Account) ID() AccountID {
	return a.root.ID()
}

func (a *Account) AddEvent(event AccountEvent) {
	a.root.AddEvent(event)
}

func (a *Account) HasPendingEvents() bool {
	return a.root.HasPendingEvents()
}

func (a *Account) CollectPendingEvents() []AccountEvent {
	return a.root.CollectPendingEvents()
}

func CreateAccount(id AccountID, name, email string) (*Account, error) {
	if name == "" || email == "" {
		return nil, errors.New("nome e email são obrigatórios")
	}

	raw := abyssalkraken.NewSimpleAggregateRoot[AccountID, AccountEvent](id)

	account := &Account{
		root:      raw,
		Name:      name,
		Email:     email,
		Active:    true,
		CreatedAt: time.Now(),
	}

	account.AddEvent(NewAccountCreatedEvent(id, name, email))

	return account, nil
}

func (a *Account) Deactivate() error {
	if !a.Active {
		return errors.New("a conta já está inativa")
	}
	a.Active = false

	a.AddEvent(NewAccountDeactivatedEvent(a.ID()))

	return nil
}

package event

import (
	"account/account/core/domain/account"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
)

type AccountEventStoreHandler struct {
	EventStore abyssalkraken.DomainEventStore[account.AccountID, account.AccountEvent]
}

func NewAccountEventStoreHandler(eventStore abyssalkraken.DomainEventStore[account.AccountID, account.AccountEvent]) *AccountEventStoreHandler {
	return &AccountEventStoreHandler{EventStore: eventStore}
}

func (h *AccountEventStoreHandler) Transactional() bool {
	return true
}

func (h *AccountEventStoreHandler) Handle(_ **account.Account, event account.AccountEvent) error {
	return h.EventStore.Save(event)
}

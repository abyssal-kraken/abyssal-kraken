package event

import (
	domain "account/account/core/domain/account"
)

type AccountEventRepositoryHandler struct {
	AccountRepository domain.AccountRepository
}

func NewAccountEventRepositoryHandler(accountRepository domain.AccountRepository) *AccountEventRepositoryHandler {
	return &AccountEventRepositoryHandler{AccountRepository: accountRepository}
}

func (h *AccountEventRepositoryHandler) Transactional() bool {
	return true
}

func (h *AccountEventRepositoryHandler) Handle(account **domain.Account, event domain.AccountEvent) error {
	if event.EventType().Name() == domain.AccountCreated.Name() {
		return h.AccountRepository.Save(*account)
	}

	return h.AccountRepository.Update(*account)
}

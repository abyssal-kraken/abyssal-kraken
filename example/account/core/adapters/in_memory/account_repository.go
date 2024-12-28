package in_memory

import (
	"account/account/core/domain/account"
	"errors"
	"sync"
)

type InMemoryAccountRepository struct {
	data map[string]*account.Account
	mu   sync.RWMutex
}

func NewInMemoryAccountRepository() *InMemoryAccountRepository {
	return &InMemoryAccountRepository{
		data: make(map[string]*account.Account),
	}
}

func (r *InMemoryAccountRepository) Save(account *account.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[account.ID().String()] = account
	return nil
}

func (r *InMemoryAccountRepository) Update(account *account.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	account, exists := r.data[account.ID().String()]
	if !exists {
		return errors.New("conta não encontrada")
	}

	r.data[account.ID().String()] = account

	return nil
}

func (r *InMemoryAccountRepository) FindByID(accountID account.AccountID) (*account.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	account, exists := r.data[accountID.String()]
	if !exists {
		return nil, errors.New("conta não encontrada")
	}
	return account, nil
}

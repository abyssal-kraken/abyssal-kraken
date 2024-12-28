package query

import (
	"account/account/core/domain/account"
)

type FindAccountByIDQuery struct {
	AccountID account.AccountID
}

type FindAccountByIDQueryProjection struct {
	ID     account.AccountID
	Name   string
	Email  string
	Active bool
}

type FindAccountByIDQueryHandler struct {
	Repo account.AccountRepository
}

func (h *FindAccountByIDQueryHandler) Handle(query FindAccountByIDQuery) (*FindAccountByIDQueryProjection, error) {
	acc, err := h.Repo.FindByID(query.AccountID)
	if err != nil {
		return nil, err
	}

	return &FindAccountByIDQueryProjection{
		ID:     acc.ID(),
		Name:   acc.Name,
		Email:  acc.Email,
		Active: acc.Active,
	}, nil
}

package account

import (
	"fmt"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"github.com/google/uuid"
)

type AccountID string

func (a AccountID) String() string {
	return string(a)
}

func NewAccountID() AccountID {
	return AccountID(uuid.New().String())
}

func FromString(id string) (AccountID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return "", fmt.Errorf("invalid AccountID, must be a valid UUID: %w", err)
	}
	return AccountID(id), nil
}

func ToAccountId(ID abyssalkraken.AggregateID) AccountID {
	return AccountID(ID.String())
}

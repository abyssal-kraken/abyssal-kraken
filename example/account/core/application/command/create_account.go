package command

import (
	"account/account/core/domain/account"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
)

type CreateAccountCommand struct {
	ID    account.AccountID
	Name  string
	Email string
}

type CreateAccountCommandHandler struct {
	EventHandlers []abyssalkraken.DomainEventHandler[*account.Account, account.AccountID, account.AccountEvent]
}

func (h *CreateAccountCommandHandler) Handle(command CreateAccountCommand) error {
	account, err := account.CreateAccount(command.ID, command.Name, command.Email)

	if err != nil {
		return err
	}

	for _, event := range account.CollectPendingEvents() {
		err := abyssalkraken.HandlePendingEvent(event, h.EventHandlers, &account)

		if err != nil {
			return err
		}
	}

	return nil
}

package command

import (
	"account/account/core/domain/account"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
)

type DeactivateAccountCommand struct {
	ID account.AccountID
}

type DeactivateAccountCommandHandler struct {
	Repo          account.AccountRepository
	EventHandlers []abyssalkraken.DomainEventHandler[*account.Account, account.AccountID, account.AccountEvent]
}

func (h *DeactivateAccountCommandHandler) Handle(command DeactivateAccountCommand) error {
	account, err := h.Repo.FindByID(command.ID)

	if err != nil {
		return err
	}

	err = account.Deactivate()

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

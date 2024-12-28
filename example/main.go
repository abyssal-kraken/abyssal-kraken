package main

import (
	http "account/account/core/adapters/http"
	repository "account/account/core/adapters/in_memory"
	"account/account/core/application/command"
	"account/account/core/application/event"
	queries "account/account/core/application/query"
	"account/account/core/domain/account"
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	repo := repository.NewInMemoryAccountRepository()
	eventStore := repository.NewInMemoryDomainEventStore[account.AccountID, account.AccountEvent]()

	eventStoreHandler := event.NewAccountEventStoreHandler(eventStore)
	eventRepositoryHandler := event.NewAccountEventRepositoryHandler(repo)

	createUserHandler := &command.CreateAccountCommandHandler{
		EventHandlers: []abyssalkraken.DomainEventHandler[*account.Account, account.AccountID, account.AccountEvent]{
			eventRepositoryHandler,
			eventStoreHandler,
		},
	}
	deactivateUserHandler := &command.DeactivateAccountCommandHandler{Repo: repo}
	getUserByIDHandler := &queries.FindAccountByIDQueryHandler{Repo: repo}

	accountCommandHandler := http.NewAccountHttpCommandHandler(createUserHandler, deactivateUserHandler)
	accountQueryHandler := http.NewAccountHttpQueryHandler(getUserByIDHandler)

	api := e.Group("/api/v1")

	http.AccountCommandRoutes(api, accountCommandHandler)
	http.AccountQueryRoutes(api, accountQueryHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

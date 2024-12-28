package http

import (
	"account/account/core/application/command"
	"account/account/core/domain/account"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AccountHttpCommandHandler struct {
	CreateAccountCommandHandler     *command.CreateAccountCommandHandler
	DeactivateAccountCommandHandler *command.DeactivateAccountCommandHandler
}

func NewAccountHttpCommandHandler(
	createUserHandler *command.CreateAccountCommandHandler,
	deactivateUserHandler *command.DeactivateAccountCommandHandler,
) *AccountHttpCommandHandler {
	return &AccountHttpCommandHandler{
		CreateAccountCommandHandler:     createUserHandler,
		DeactivateAccountCommandHandler: deactivateUserHandler,
	}
}

func (h *AccountHttpCommandHandler) CreateAccount(c echo.Context) error {
	var req CreateAccountRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "JSON inválido"})
	}

	cmd := command.CreateAccountCommand{
		ID:    account.NewAccountID(),
		Name:  req.Name,
		Email: req.Email,
	}

	if err := h.CreateAccountCommandHandler.Handle(cmd); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"id": "123"})
}

func (h *AccountHttpCommandHandler) Deactivate(c echo.Context) error {
	id := c.Param("id")

	accountId, err := account.FromString(id)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "id inválido"})
	}

	cmd := command.DeactivateAccountCommand{ID: accountId}

	if err := h.DeactivateAccountCommandHandler.Handle(cmd); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

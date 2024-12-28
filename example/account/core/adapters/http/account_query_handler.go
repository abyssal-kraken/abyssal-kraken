package http

import (
	"account/account/core/application/query"
	"account/account/core/domain/account"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AccountHttpQueryHandler struct {
	GetUserByIDHandler *query.FindAccountByIDQueryHandler
}

func NewAccountHttpQueryHandler(getUserByIDHandler *query.FindAccountByIDQueryHandler) *AccountHttpQueryHandler {
	return &AccountHttpQueryHandler{
		GetUserByIDHandler: getUserByIDHandler,
	}
}

func (h *AccountHttpQueryHandler) GetAccountByID(c echo.Context) error {
	id := c.Param("id")

	accountId, err := account.FromString(id)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "id inválido"})
	}

	query := query.FindAccountByIDQuery{AccountID: accountId}

	result, err := h.GetUserByIDHandler.Handle(query)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Conta não encontrada"})
	}

	accountResponse := AccountResponse{
		ID:     result.ID.String(),
		Name:   result.Name,
		Email:  result.Email,
		Active: result.Active,
	}

	return c.JSON(http.StatusOK, accountResponse)
}

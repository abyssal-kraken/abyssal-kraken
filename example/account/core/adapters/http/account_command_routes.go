package http

import (
	"github.com/labstack/echo/v4"
)

func AccountCommandRoutes(e *echo.Group, handler *AccountHttpCommandHandler) {
	e.POST("/accounts", handler.CreateAccount)
	e.PUT("/accounts/:id/deactivate", handler.Deactivate)
}

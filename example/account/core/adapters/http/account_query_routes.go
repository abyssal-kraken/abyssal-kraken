package http

import "github.com/labstack/echo/v4"

func AccountQueryRoutes(e *echo.Group, handler *AccountHttpQueryHandler) {
	e.GET("/accounts/:id", handler.GetAccountByID)
}

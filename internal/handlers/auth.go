package handlers

import (
	"net/http"
	"pastebin/internal/models"
	"pastebin/internal/services"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authSvc *services.AuthService
}

func (h *AuthHandler) Register(c echo.Context) error {
	var RegisterInput *models.RegisterInput
	if err := c.Bind(RegisterInput); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	ctx := c.Request().Context()
	return h.authSvc.Register(ctx, RegisterInput)
}

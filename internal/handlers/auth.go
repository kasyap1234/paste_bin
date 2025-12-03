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

func NewAuthHandler(authSvc *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var RegisterInput models.RegisterInput
	if err := c.Bind(RegisterInput); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	ctx := c.Request().Context()
	if err := h.authSvc.Register(ctx, &RegisterInput); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to register"})

	}
	return c.JSON(http.StatusCreated, echo.Map{"error": "user registered"})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var loginInput models.LoginInput
	if err := c.Bind(loginInput); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	ctx := c.Request().Context()
	resp, err := h.authSvc.Login(ctx, &loginInput)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized "})

	}
	return c.JSON(http.StatusOK, resp)
}

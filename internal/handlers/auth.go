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

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Register a new user with email and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.RegisterInput	true	"User registration data"
//	@Success		201		{object}	map[string]string	"User registered successfully"
//	@Failure		400		{object}	map[string]string	"Invalid request"
//	@Failure		500		{object}	map[string]string	"Failed to register"
//	@Router			/register [post]
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

// Login godoc
//
//	@Summary		Login user
//	@Description	Authenticate user and return JWT token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.LoginInput		true	"User login credentials"
//	@Success		200		{object}	models.LoginResponse	"Login successful"
//	@Failure		400		{object}	map[string]string		"Invalid request"
//	@Failure		401		{object}	map[string]string		"Unauthorized"
//	@Router			/login [post]
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

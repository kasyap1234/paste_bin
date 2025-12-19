package handlers

import (
	"fmt"
	"net/http"
	"pastebin/internal/models"
	"pastebin/internal/services"
	"strings"

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
	if err := c.Bind(&RegisterInput); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	// Validate input
	if RegisterInput.Name == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "name is required"})
	}
	if len(RegisterInput.Name) < 2 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "name must be at least 2 characters"})
	}
	if RegisterInput.Email == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "email is required"})
	}
	if !strings.Contains(RegisterInput.Email, "@") {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid email format"})
	}
	if len(RegisterInput.Password) < 6 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "password must be at least 6 characters"})
	}

	ctx := c.Request().Context()
	if err := h.authSvc.Register(ctx, &RegisterInput); err != nil {
		fmt.Printf("Register error: %v\n", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to register"})
	}
	return c.JSON(http.StatusCreated, echo.Map{"message": "user registered"})
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
	if err := c.Bind(&loginInput); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	// Validate input
	if loginInput.Email == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "email is required"})
	}
	if !strings.Contains(loginInput.Email, "@") {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid email format"})
	}
	if len(loginInput.Password) < 6 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "password must be at least 6 characters"})
	}

	ctx := c.Request().Context()
	resp, err := h.authSvc.Login(ctx, &loginInput)
	if err != nil {
		fmt.Printf("Login error: %v\n", err)
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	return c.JSON(http.StatusOK, resp)
}

package handlers

import (
	"net/http"
	"pastebin/internal/models"
	"pastebin/internal/services"
	"pastebin/pkg/utils"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	authSvc *services.AuthService
	userSvc *services.UserService
	logger  zerolog.Logger
}

func NewAuthHandler(authSvc *services.AuthService, userSvc *services.UserService, logger zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
		userSvc: userSvc,
		logger:  logger,
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
		return utils.SendError(c, http.StatusBadRequest, "invalid request")
	}

	// Validate input
	if RegisterInput.Name == "" {
		return utils.SendError(c, http.StatusBadRequest, "name is required")
	}
	if len(RegisterInput.Name) < 2 {
		return utils.SendError(c, http.StatusBadRequest, "name must be at least 2 characters")
	}
	if RegisterInput.Email == "" {
		return utils.SendError(c, http.StatusBadRequest, "email is required")
	}
	if !strings.Contains(RegisterInput.Email, "@") {
		return utils.SendError(c, http.StatusBadRequest, "invalid email format")
	}
	if len(RegisterInput.Password) < 6 {
		return utils.SendError(c, http.StatusBadRequest, "password must be at least 6 characters")
	}

	ctx := c.Request().Context()
	if err := h.authSvc.Register(ctx, &RegisterInput); err != nil {
		h.logger.Error().Err(err).Msg("failed to register user")
		return utils.SendError(c, http.StatusInternalServerError, "failed to register user")
	}
	return utils.SendSuccess(c, http.StatusCreated, nil, "user registered successfully")
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
		return utils.SendError(c, http.StatusBadRequest, "invalid request")
	}

	// Validate input
	if loginInput.Email == "" {
		return utils.SendError(c, http.StatusBadRequest, "email is required")
	}
	if !strings.Contains(loginInput.Email, "@") {
		return utils.SendError(c, http.StatusBadRequest, "invalid email format")
	}
	if len(loginInput.Password) < 6 {
		return utils.SendError(c, http.StatusBadRequest, "password must be at least 6 characters")
	}

	ctx := c.Request().Context()
	resp, err := h.authSvc.Login(ctx, &loginInput)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to login")
		return utils.SendError(c, http.StatusUnauthorized, "invalid email or password")
	}
	return utils.SendSuccess(c, http.StatusOK, resp, "login successful")
}

// ForgotPassword godoc
//
//	@Summary		Request password reset
//	@Description	Send a password reset link to the user's email
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		map[string]string	true	"Email payload"
//	@Success		200		{object}	map[string]string	"Reset link sent if email exists"
//	@Failure		400		{object}	map[string]string	"Invalid request"
//	@Router			/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	var payload struct {
		Email string `json:"email"`
	}

	if err := c.Bind(&payload); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "invalid request")
	}

	if payload.Email == "" {
		return utils.SendError(c, http.StatusBadRequest, "email is required")
	}
	if !strings.Contains(payload.Email, "@") {
		return utils.SendError(c, http.StatusBadRequest, "invalid email format")
	}

	ctx := c.Request().Context()
	if err := h.userSvc.ForgotPassword(ctx, payload.Email); err != nil {
		h.logger.Error().Err(err).Msg("failed to process forgot password")
		// Return generic success to avoid email enumeration.
	}

	return utils.SendSuccess(c, http.StatusOK, nil, "if the email exists, a reset link has been sent")
}

// ResetPassword godoc
//
//	@Summary		Reset password using token
//	@Description	Reset user's password with a valid reset token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		map[string]string	true	"Reset payload"
//	@Success		200		{object}	map[string]string	"Password reset successfully"
//	@Failure		400		{object}	map[string]string	"Invalid request"
//	@Failure		400		{object}	map[string]string	"Invalid or expired token"
//	@Router			/reset-password [post]
func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var payload struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.Bind(&payload); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "invalid request")
	}

	if payload.Token == "" {
		return utils.SendError(c, http.StatusBadRequest, "token is required")
	}
	if len(payload.NewPassword) < 6 {
		return utils.SendError(c, http.StatusBadRequest, "password must be at least 6 characters")
	}

	ctx := c.Request().Context()
	if err := h.userSvc.ResetPasswordWithToken(ctx, payload.Token, payload.NewPassword); err != nil {
		h.logger.Error().Err(err).Msg("failed to reset password with token")
		return utils.SendError(c, http.StatusBadRequest, "invalid or expired token")
	}

	return utils.SendSuccess(c, http.StatusOK, nil, "password reset successfully")
}

// VerifyEmail godoc
//
//	@Summary		Verify user email
//	@Description	Verify user email using verification token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string				true	"Verification token"
//	@Success		200		{object}	map[string]string	"Email verified successfully"
//	@Failure		400		{object}	map[string]string	"Invalid or expired token"
//	@Failure		500		{object}	map[string]string	"Failed to verify email"
//	@Router			/verify-email/{token} [get]
func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	token := c.Param("token")
	if token == "" {
		return utils.SendError(c, http.StatusBadRequest, "verification token is required")
	}

	ctx := c.Request().Context()
	if err := h.authSvc.VerifyEmail(ctx, token); err != nil {
		h.logger.Error().Err(err).Msg("failed to verify email")
		return utils.SendError(c, http.StatusBadRequest, "invalid or expired verification token")
	}

	return utils.SendSuccess(c, http.StatusOK, nil, "email verified successfully")
}

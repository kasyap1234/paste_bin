package services

import (
	"context"
	"errors"
	"fmt"
	"pastebin/internal/auth"
	"pastebin/internal/models"
	"pastebin/internal/repositories"
	"pastebin/pkg/utils"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type AuthService struct {
	authRepo   *repositories.AuthRepository
	jwtManager *auth.JWTManager
	userRepo   *repositories.UserRepository
	logger     zerolog.Logger
}

func NewAuthService(authRepo *repositories.AuthRepository, userRepo *repositories.UserRepository, jwtMgr *auth.JWTManager, logger zerolog.Logger) *AuthService {
	return &AuthService{
		authRepo:   authRepo,
		jwtManager: jwtMgr,
		userRepo:   userRepo,
		logger:     logger,
	}
}

func (a *AuthService) Register(ctx context.Context, registerInput *models.RegisterInput) error {
	// check if user already exists
	// only proceed if user does not exists
	ok, err := a.userRepo.ExistsUser(ctx, registerInput.Email)
	if err != nil {
		a.logger.Error().Err(err).Msg("error checking existing user")
		return fmt.Errorf("error checking existing user: %w", err)
	}

	if ok {
		return fmt.Errorf("user already exists with email: %s", registerInput.Email)
	}
	regErr := a.authRepo.Register(ctx, registerInput)
	if regErr != nil {
		a.logger.Error().Err(regErr).Msg("error registering user")
		return regErr
	}
	return nil
}

func (a *AuthService) Login(ctx context.Context, loginInput *models.LoginInput) (*models.LoginResponse, error) {
	user, err := a.userRepo.GetUserByEmail(ctx, loginInput.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		a.logger.Error().Msg("user not found")
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if err != nil {
		a.logger.Error().Err(err).Msg("failed to get user by email")
		return nil, fmt.Errorf("invalid email or password: %w", err)
	}

	if !utils.VerifyPassword(user.PasswordHash, loginInput.Password) {
		a.logger.Error().Msg("invalid email or password")
		return nil, fmt.Errorf("invalid email or password: %w", err)
	}
	token, err := a.jwtManager.GenerateToken(user.ID, user.Email, 24*time.Hour)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to generate token")
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

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
	"github.com/jackc/pgx/v5/pgconn"
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
	user, regErr := a.authRepo.Register(ctx, registerInput)
	if regErr != nil {
		var pgErr *pgconn.PgError
		if errors.As(regErr, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("user already exists with email: %s", registerInput.Email)
		}
		a.logger.Error().Err(regErr).Msg("error registering user")
		return regErr
	}
	verificationToken := utils.GenerateResetToken()
	if err := a.userRepo.SaveVerifyToken(ctx, user.ID, verificationToken, time.Now().Add(24*time.Hour)); err != nil {
		a.logger.Error().Err(err).Msg("error updating verification token")
		return err
	}
	// Don't fail registration if email fails - just log the error
	if err := utils.SendVerifyEmail(registerInput.Email, verificationToken); err != nil {
		a.logger.Error().Err(err).Msg("failed to send verification email - registration still successful")
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
		return nil, fmt.Errorf("invalid email or password")
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

func (a *AuthService) VerifyEmail(ctx context.Context, token string) error {
	user, err := a.userRepo.GetUserByVerifyToken(ctx, token)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to get user by token")
		return fmt.Errorf("failed to get user by token: %w", err)
	}
	if user.VerifyToken == nil || user.VerifyTokenExpiresAt == nil || time.Now().After(*user.VerifyTokenExpiresAt) {
		return fmt.Errorf("invalid or expired verify token")
	}

	user.IsVerified = true
	err = a.userRepo.UpdateUserIsVerified(ctx, user.ID, user.IsVerified)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to update user")
		return fmt.Errorf("failed to update user: %w", err)
	}
	err = a.userRepo.ClearVerifyToken(ctx, user.ID)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to clear verify token")
		return fmt.Errorf("failed to clear verify token: %w", err)
	}
	a.logger.Info().Msg("email verified successfully")
	return nil
}

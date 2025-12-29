package services

import (
	"context"
	"errors"
	"fmt"
	"pastebin/internal/models"
	"pastebin/internal/repositories"
	"pastebin/pkg/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type UserService struct {
	userRepo *repositories.UserRepository
	logger   zerolog.Logger
}

func NewUserService(userRepo *repositories.UserRepository, logger zerolog.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// CheckUserExists returns true if a user with the given userID exists.
// If the user is not found it returns (false, nil).
// If an unexpected error occurs while checking, it returns (false, error).

func (u *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("unable to get user by email id : %w", err)
	}
	return user, nil
}
func (u *UserService) CheckUserExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	if u == nil || u.userRepo == nil {
		u.logger.Error().Msg("user service or repository is not initialized")
		return false, fmt.Errorf("user service or repository is not initialized")
	}

	user, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		// If repository wrapped a pgx.ErrNoRows, treat that as "not exists".
		if errors.Is(err, pgx.ErrNoRows) {
			u.logger.Error().Msg("user not found")
			return false, nil
		}
		// Propagate unexpected errors.
		u.logger.Error().Err(err).Msg("failed to check user existence")
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	// Defensive: if repository returned a nil pointer but no error, treat as not exists.
	if user == nil {
		u.logger.Error().Msg("user is nil")
		return false, nil
	}
	return true, nil
}

func (u *UserService) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	_, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("unable to check email existence :%w", err)
	}
	return true, nil
}

func (u *UserService) UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	if len(newPassword) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	err := u.userRepo.UpdatePassword(ctx, userID, newPassword)
	if err != nil {
		u.logger.Error().Err(err).Msg("failed to update password")
		return fmt.Errorf("failed to update password :%w", err)
	}
	return nil
}

func (u *UserService) ForgotPassword(ctx context.Context, email string) error {
	// check if email exists
	user, err := u.GetUserByEmail(ctx, email)
	if err != nil {
		// Do not reveal whether the email exists to the caller; just log.
		u.logger.Debug().Err(err).Str("email", email).Msg("forgot password requested for non-existent email")
		return nil
	}

	// Generate a reset token and store it with expiry
	resetToken := utils.GenerateResetToken()
	expiresAt := time.Now().Add(1 * time.Hour)

	if err := u.userRepo.SavePasswordResetToken(ctx, user.ID, resetToken, expiresAt); err != nil {
		u.logger.Error().Err(err).Msg("failed to save password reset token")
		return fmt.Errorf("failed to save password reset token: %w", err)
	}

	if err = utils.SendPasswordResetEmail(email, resetToken); err != nil {
		u.logger.Error().Err(err).Msg("failed to send password reset link")
		return fmt.Errorf("failed to send password from email: %w", err)
	}

	return nil
}

// ResetPasswordWithToken verifies the provided reset token and, if valid,
// updates the user's password and clears the token.
func (u *UserService) ResetPasswordWithToken(ctx context.Context, token, newPassword string) error {
	if len(newPassword) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}

	user, err := u.userRepo.GetUserByResetToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("invalid or expired reset token")
		}
		return fmt.Errorf("failed to verify reset token: %w", err)
	}

	if err := u.userRepo.UpdatePassword(ctx, user.ID, newPassword); err != nil {
		u.logger.Error().Err(err).Msg("failed to update password via reset token")
		return fmt.Errorf("failed to update password: %w", err)
	}

	if err := u.userRepo.ClearPasswordResetToken(ctx, user.ID); err != nil {
		u.logger.Error().Err(err).Msg("failed to clear password reset token")
		return fmt.Errorf("failed to clear password reset token: %w", err)
	}

	return nil
}

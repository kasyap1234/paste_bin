package services

import (
	"context"
	"errors"
	"fmt"
	"pastebin/internal/repositories"

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

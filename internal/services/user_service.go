package services

import (
	"context"
	"errors"
	"fmt"
	"pastebin/internal/repositories"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CheckUserExists returns true if a user with the given userID exists.
// If the user is not found it returns (false, nil).
// If an unexpected error occurs while checking, it returns (false, error).
func (u *UserService) CheckUserExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	if u == nil || u.userRepo == nil {
		return false, fmt.Errorf("user service or repository is not initialized")
	}

	user, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		// If repository wrapped a pgx.ErrNoRows, treat that as "not exists".
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		// Propagate unexpected errors.
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	// Defensive: if repository returned a nil pointer but no error, treat as not exists.
	if user == nil {
		return false, nil
	}
	return true, nil
}

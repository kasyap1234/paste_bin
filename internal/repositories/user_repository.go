package repositories

import (
	"context"
	"fmt"
	"pastebin/internal/models"
	"pastebin/pkg/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
func (u *UserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	query := `SELECT id, name, email, password_hash FROM users WHERE id=$1`
	rows, err := u.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	defer rows.Close()
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, fmt.Errorf("unable to collect row: %w", err)
	}
	return &user, nil
}
func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, name, email, password_hash FROM users WHERE email=$1`
	rows, err := u.db.Query(ctx, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to query user by email: %w", err)
	}
	defer rows.Close()
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}
	return &user, nil
}

func (u *UserRepository) ExistsUser(ctx context.Context, email string) (bool, error) {
	_, err := u.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil // User doesn't exist, which is not an error
		}
		return false, fmt.Errorf("failed to get user by email %s: %w", email, err)
	}
	return true, nil // User exists
}

func (u *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, name, email, password_hash) VALUES ($1, $2, $3, $4)`
	_, err := u.db.Exec(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET name=$2,email=$3,password_hash=$4 WHERE id=$1`
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("unable to begin transaction :%w", err)
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction %w", err)
	}

	return nil
}

func (u *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	query := `UPDATE users SET password_hash=$1 WHERE id=$2`
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction : %w", err)
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, query, hashedPassword, userID)
	if err != nil {
		return fmt.Errorf("unable to update password :%w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction:%w", err)
	}
	return nil

}

// SavePasswordResetToken stores a reset token and its expiry for the given user.
func (u *UserRepository) SavePasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	query := `UPDATE users SET reset_token=$1, reset_token_expires_at=$2 WHERE id=$3`
	_, err := u.db.Exec(ctx, query, token, expiresAt, userID)
	if err != nil {
		return fmt.Errorf("failed to save password reset token: %w", err)
	}
	return nil
}

// GetUserByResetToken returns the user associated with a valid, non-expired reset token.
func (u *UserRepository) GetUserByResetToken(ctx context.Context, token string) (*models.User, error) {
	query := `SELECT id, name, email, password_hash FROM users WHERE reset_token=$1 AND reset_token_expires_at > NOW()`
	rows, err := u.db.Query(ctx, query, token)
	if err != nil {
		return nil, fmt.Errorf("failed to query user by reset token: %w", err)
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("failed to scan user by reset token: %w", err)
	}
	return &user, nil
}

// ClearPasswordResetToken removes the reset token and its expiry for a user after successful reset.
func (u *UserRepository) ClearPasswordResetToken(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET reset_token=NULL, reset_token_expires_at=NULL WHERE id=$1`
	_, err := u.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to clear password reset token: %w", err)
	}
	return nil
}

func (u *UserRepository) GetVerificationToken(ctx context.Context, userID uuid.UUID) (string, error) {
	query := `SELECT verification_token FROM users WHERE id = $1`
	var token *string
	err := u.db.QueryRow(ctx, query, userID).Scan(&token)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil // Return empty string if user not found, or maybe error? existing logic suggests returning empty string or specific error. Usually if user not found it's an error, but if token is null it's empty string.
		}
		return "", fmt.Errorf("failed to scan verification token: %w", err)
	}
	if token == nil {
		return "", nil
	}
	return *token, nil
}

func (u *UserRepository) ClearVerificationToken(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET verification_token = NULL, verification_token_expires_at = NULL WHERE id = $1`
	_, err := u.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to clear verification token: %w", err)
	}
	return nil
}

func (u *UserRepository) SaveVerificationToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	query := `UPDATE users SET verification_token = $1, verification_token_expires_at = $2 WHERE id = $3`
	_, err := u.db.Exec(ctx, query, token, expiresAt, userID)
	if err != nil {
		return fmt.Errorf("failed to save verification token: %w", err)
	}
	return nil
}

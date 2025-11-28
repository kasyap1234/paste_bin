package repositories

import (
	"context"
	"fmt"
	"pastebin/internal/models"

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

func (a *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT * FROM users WHERE email=$1`
	rows, err := a.db.Query(ctx, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to query user by email: %w", err)
	}
	defer rows.Close()
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}
	return &user, nil
}

func (u *UserRepository) ExistsUser(ctx context.Context, email string) (bool, error) {
	_, err := u.GetUserByEmail(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to get user by email %s", email)
	}
	return true, nil
}

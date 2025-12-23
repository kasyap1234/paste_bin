package repositories

import (
	"context"
	"fmt"
	"pastebin/internal/models"

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


func(u*UserRepository)UpdateUser(ctx context.Context,user *models.User)error{
	query:=`UPDATE users SET name=$2,email=$3,password_hash=$4 WHERE id=$1`
	_,err:=u.db.Exec(ctx,query,user.ID,user.Name,user.Email,user.PasswordHash)
	if err!=nil{
		return fmt.Errorf("failed to update user: %w",err)
	}
	return nil
}



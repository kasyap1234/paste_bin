package repositories

import (
	"context"
	"fmt"
	"pastebin/internal/models"
	"pastebin/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (a *AuthRepository) Register(ctx context.Context, registerInput *models.RegisterInput) error {
	hashedPassword, err := utils.HashPassword(registerInput.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user := &models.User{
		ID:       uuid.New(),
		Name:     registerInput.Name,
		Email:    registerInput.Email,
		Password: hashedPassword,
	}
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	query := `INSERT INTO users(id,name,email,password_hash) VALUES ($1,$2,$3,$4)`

	_, err = tx.Exec(ctx, query, user.ID, user.Name, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

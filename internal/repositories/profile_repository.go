package repositories

import (
	"context"
	"errors"
	"fmt"
	"pastebin/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{
		db: db,
	}
}

func (p *ProfileRepository) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	row, err := p.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	defer row.Close()
	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to collect user: %w", err)
	}
	return &user, nil
}

func (p *ProfileRepository) UpdateProfile(ctx context.Context, userID uuid.UUID, patch *models.PatchProfile) (*models.User, error) {
	// Build dynamic update query based on provided fields
	updateBuilder := sq.Update("users").Where(sq.Eq{"id": userID}).PlaceholderFormat(sq.Dollar)

	hasUpdates := false

	if patch.Name != nil {
		updateBuilder = updateBuilder.Set("name", *patch.Name)
		hasUpdates = true
	}

	if patch.Avatar != nil {
		updateBuilder = updateBuilder.Set("avatar", *patch.Avatar)
		hasUpdates = true
	}

	// If no fields to update, just return the current profile
	if !hasUpdates {
		return p.GetProfile(ctx, userID)
	}

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build update query: %w", err)
	}

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return p.GetProfile(ctx, userID)
}

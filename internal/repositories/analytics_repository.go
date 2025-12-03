package repositories

import (
	"context"
	"fmt"
	"pastebin/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalyticsRepository struct {
	db *pgxpool.Pool
}

func NewAnalyticsRepository(db *pgxpool.Pool) *AnalyticsRepository {
	return &AnalyticsRepository{
		db: db,
	}
}

func (a *AnalyticsRepository) CreateAnalytics(ctx context.Context, pasteID uuid.UUID, url string) error {
	query := `INSERT INTO pastes_analytics(paste_id,url) VALUES ($1,$2)`
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, query, pasteID, url)
	if err != nil {
		return fmt.Errorf("failed to create analytics: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (a *AnalyticsRepository) GetAnalyticsByPasteID(ctx context.Context, pasteID uuid.UUID) (*models.Analytics, error) {
	query := `SELECT * FROM pastes_analytics WHERE pasteID=$1`
	row, err := a.db.Query(ctx, query, pasteID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no analytics found for pasteID: %s", pasteID)
		}
		return nil, fmt.Errorf("failed to get analytics by pasteID: %w", err)
	}
	defer row.Close()
	analytics, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Analytics])
	if err != nil {
		return nil, fmt.Errorf("failed to collect analytics: %w", err)
	}
	return &analytics, nil
}

func (a *AnalyticsRepository) IncrementViews(ctx context.Context, pasteID uuid.UUID) error {
	query := `UPDATE pastes_analytics SET views=views+1 WHERE pasteID=$1`
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, query, pasteID)
	if err != nil {
		return fmt.Errorf("failed to increment views: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (a *AnalyticsRepository) GetAnalyticsByURL(ctx context.Context, url string) (*models.Analytics, error) {
	query := `SELECT * FROM pastes_analytics WHERE url=$1`
	row, err := a.db.Query(ctx, query, url)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no analytics found for url: %s", url)
		}
		return nil, fmt.Errorf("failed to get analytics by url: %w", err)
	}
	defer row.Close()
	analytics, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Analytics])
	if err != nil {
		return nil, fmt.Errorf("failed to collect analytics: %w", err)
	}
	return &analytics, nil
}

func (a *AnalyticsRepository) GetAllAnalytics(ctx context.Context, order string, limit int, offset int) ([]models.Analytics, error) {
	query := `SELECT * FROM pastes_analytics  ORDER BY $1 LIMIT $2 OFFSET $3`
	rows, err := a.db.Query(ctx, query, order, limit, offset)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no analytics found")
		}
		return nil, fmt.Errorf("failed to get all analytics: %w", err)
	}
	defer rows.Close()
	analytics, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Analytics])
	if err != nil {
		return nil, fmt.Errorf("failed to collect analytics: %w", err)
	}
	return analytics, nil
}

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

type AnalyticsRepository struct {
	db *pgxpool.Pool
}

func NewAnalyticsRepository(db *pgxpool.Pool) *AnalyticsRepository {
	return &AnalyticsRepository{
		db: db,
	}
}

func (a *AnalyticsRepository) CreateAnalytics(ctx context.Context, pasteID uuid.UUID, url string) error {
	query := `INSERT INTO pastes_analytics (paste_id, url) VALUES ($1, $2)`
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
	query := `SELECT * FROM pastes_analytics WHERE paste_id = $1`
	row, err := a.db.Query(ctx, query, pasteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics by paste_id: %w", err)
	}
	defer row.Close()
	analytics, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Analytics])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to collect analytics: %w", err)
	}
	return &analytics, nil
}

func (a *AnalyticsRepository) IncrementViews(ctx context.Context, pasteID uuid.UUID) error {
	// Build the update query using squirrel
	query := sq.Update("pastes_analytics").
		Set("views", sq.Expr("views+1")).
		Where(sq.Eq{"paste_id": pasteID}).
		PlaceholderFormat(sq.Dollar)

	queryStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build increment views query: %w", err)
	}

	tx, err := a.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, queryStr, args...)
	if err != nil {
		return fmt.Errorf("failed to increment views: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (a *AnalyticsRepository) GetAnalyticsByURL(ctx context.Context, url string) (*models.Analytics, error) {
	query := `SELECT * FROM pastes_analytics WHERE url = $1`
	row, err := a.db.Query(ctx, query, url)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics by url: %w", err)
	}
	defer row.Close()
	analytics, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Analytics])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to collect analytics: %w", err)
	}
	return &analytics, nil
}

func (a *AnalyticsRepository) GetAllAnalytics(ctx context.Context, order string, limit int, offset int) ([]models.Analytics, error) {
	// ORDER BY cannot use parameters, so we need to validate and use string interpolation carefully
	// Only allow safe column names
	allowedOrders := map[string]string{
		"created_at": "created_at",
		"updated_at": "updated_at",
		"views":      "views",
	}
	orderBy := "created_at DESC" // default
	if orderCol, ok := allowedOrders[order]; ok {
		orderBy = orderCol + " DESC"
	}
	query := fmt.Sprintf(`SELECT * FROM pastes_analytics ORDER BY %s LIMIT $1 OFFSET $2`, orderBy)
	rows, err := a.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all analytics: %w", err)
	}
	defer rows.Close()
	analytics, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Analytics])
	if err != nil {
		return nil, fmt.Errorf("failed to collect analytics: %w", err)
	}
	return analytics, nil
}

func (a *AnalyticsRepository) GetAllAnalyticsByUser(ctx context.Context, userID uuid.UUID, order string, limit, offset int) ([]models.Analytics, error) {
	// ORDER BY cannot use parameters, so we need to validate and use string interpolation carefully
	// Only allow safe column names
	allowedOrders := map[string]string{
		"created_at": "a.created_at",
		"updated_at": "a.updated_at",
		"views":      "a.views",
	}
	orderBy := "a.created_at DESC" // default
	if orderCol, ok := allowedOrders[order]; ok {
		orderBy = orderCol + " DESC"
	}
	query := fmt.Sprintf(`SELECT a.* FROM pastes_analytics a JOIN pastes p ON a.paste_id = p.id WHERE p.user_id = $1 ORDER BY %s LIMIT $2 OFFSET $3`, orderBy)
	rows, err := a.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all analytics by user: %w", err)
	}
	defer rows.Close()
	analytics, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Analytics])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}
	return analytics, nil
}

func (a *AnalyticsRepository) GetAnalyticsByID(ctx context.Context, id uuid.UUID) (*models.Analytics, error) {
	query := `SELECT * FROM pastes_analytics WHERE id = $1`
	rows, err := a.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics for id: %w", err)
	}
	defer rows.Close()
	analytic, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Analytics])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row: %w", err)
	}
	return &analytic, nil
}

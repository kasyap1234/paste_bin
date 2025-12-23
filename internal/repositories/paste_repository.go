package repositories

import (
	"context"
	"fmt"
	"os"
	"pastebin/internal/models"
	"pastebin/pkg/utils"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PasteRepository struct {
	db *pgxpool.Pool
}

func NewPasteRepository(db *pgxpool.Pool) *PasteRepository {
	return &PasteRepository{
		db: db,
	}
}

func (p *PasteRepository) CreatePaste(ctx context.Context, userID uuid.UUID, pasteInput *models.PasteInput) (*models.PasteOutput, error) {
	query := `INSERT INTO pastes (user_id, title, is_private, content, language, url, password, expires_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	title := pasteInput.Title
	if title == "" {
		title = "Untitled"
	}
	urlSlug := uuid.New().String()[:8]
	var isPrivate bool
	var passwordHash string
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	url := baseURL + "/p/" + urlSlug
	if pasteInput.Password == "" {
		isPrivate = false
		passwordHash = ""
	} else {
		isPrivate = true
		hashedPassword, err := utils.HashPassword(pasteInput.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		passwordHash = hashedPassword
	}
	language := pasteInput.Language
	content := pasteInput.Content
	expiresAt := pasteInput.ExpiresAt
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, query, userID, title, isPrivate, content, language, url, passwordHash, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert paste: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Retrieve the created paste to return it
	getQuery := `SELECT id, user_id, title, is_private, content, language, url, expires_at, created_at, updated_at FROM pastes WHERE url = $1`
	row, err := p.db.Query(ctx, getQuery, url)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created paste: %w", err)
	}
	defer row.Close()
	paste, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[models.PasteOutput])
	if err != nil {
		return nil, fmt.Errorf("failed to collect created paste: %w", err)
	}

	return &paste, nil
}

func (p *PasteRepository) UpdatePaste(ctx context.Context, pasteID uuid.UUID, patchInput *models.PatchPaste) error {
	// Convert patch input to a map of updates, skipping nil fields
	updates := utils.StructToMap(patchInput, "db")

	// Handle special password logic
	if patchInput.Password != nil {
		password := *patchInput.Password
		if password == "" {
			// Empty password means removing password protection
			updates["password"] = ""
			// If is_private is not explicitly set, set it to false when password is removed
			if patchInput.IsPrivate == nil {
				updates["is_private"] = false
			}
		} else {
			// Hash the password before storing
			hashedPassword, err := utils.HashPassword(password)
			if err != nil {
				return fmt.Errorf("failed to hash password: %w", err)
			}
			updates["password"] = hashedPassword
			// If is_private is not explicitly set, set it to true when password is provided
			if patchInput.IsPrivate == nil {
				updates["is_private"] = true
			}
		}
	}

	// Always update the updated_at timestamp
	updates["updated_at"] = time.Now()

	// Build the update query using squirrel
	updateBuilder := sq.Update("pastes").
		SetMap(updates).
		Where(sq.Eq{"id": pasteID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	// Begin transaction
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Execute the update query
	cmdTag, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update paste: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("paste not found with id: %s", pasteID.String())
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (p *PasteRepository) GetPasteByID(ctx context.Context, pasteID uuid.UUID, isAuthenticated bool, userID uuid.UUID) (*models.PasteOutput, error) {
	query := `SELECT p.id, p.user_id, p.title, p.is_private, p.content, p.language, p.url, p.expires_at, p.created_at, p.updated_at, COALESCE(a.views, 0) as views FROM pastes p LEFT JOIN pastes_analytics a ON p.id = a."pasteID" WHERE p.id = $1`
	row, err := p.db.Query(ctx, query, pasteID)
	if err != nil {
		return nil, fmt.Errorf("failed to query paste: %w", err)
	}
	defer row.Close()
	paste, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[models.PasteOutput])
	if err != nil {
		return nil, fmt.Errorf("failed to collect paste: %w", err)
	}

	// Check if paste has expired
	if paste.ExpiresAt != nil && paste.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("paste has expired")
	}
	// Check if user is the owner
	isOwner := isAuthenticated && paste.UserID == userID
	if !isOwner {
		// Increment view count for non-owner views
		if err := p.incrementViewCount(ctx, pasteID); err != nil {
			// Log the error but don't fail the paste retrieval
			// since view counting is not critical
			fmt.Printf("Failed to increment view count for paste %s: %v\n", pasteID, err)
		}
	}
	return &paste, nil
}

func (p *PasteRepository) incrementViewCount(ctx context.Context, pasteID uuid.UUID) error {
	query := `INSERT INTO pastes_analytics ("pasteID", views, updated_at) VALUES ($1, 1, NOW())
ON CONFLICT ("pasteID")
DO UPDATE SET views = pastes_analytics.views + 1, updated_at = NOW()`
	_, err := p.db.Exec(ctx, query, pasteID)
	return err
}

func (p *PasteRepository) GetAllPastes(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.PasteOutput, int, error) {
	// First, get the total count of non-expired pastes for the user
	countQuery := `SELECT COUNT(*) FROM pastes WHERE user_id = $1 AND (expires_at IS NULL OR expires_at > NOW())`
	var total int
	err := p.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Then get the paginated results
	query := `SELECT p.id, p.user_id, p.title, p.is_private, p.language, p.url, p.expires_at, p.created_at, COALESCE(a.views, 0) as views 
		FROM pastes p 
		LEFT JOIN pastes_analytics a ON p.id = a."pasteID" 
		WHERE p.user_id = $1 AND (p.expires_at IS NULL OR p.expires_at > NOW())
		ORDER BY p.created_at DESC 
		LIMIT $2 OFFSET $3`
	row, err := p.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get pastes for user ID: %w", err)
	}
	defer row.Close()
	pastes, err := pgx.CollectRows(row, pgx.RowToStructByName[models.PasteOutput])
	if err != nil {
		return nil, 0, fmt.Errorf("failed to collect pastes: %w", err)
	}

	return pastes, total, nil
}

func (p *PasteRepository) DeletePasteByID(ctx context.Context, pasteID uuid.UUID) error {
	// Build the delete query using squirrel
	query := sq.Delete("pastes").
		Where(sq.Eq{"id": pasteID}).
		PlaceholderFormat(sq.Dollar)

	queryStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	// Begin transaction
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Execute the delete query
	cmdTag, err := tx.Exec(ctx, queryStr, args...)
	if err != nil {
		return fmt.Errorf("failed to delete paste by id: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("paste not found with id: %s", pasteID)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (p *PasteRepository) GetPasteBySlug(ctx context.Context, slug string) (*models.PasteOutput, error) {
	// Query for paste where URL ends with /p/slug
	query := `SELECT p.id, p.user_id, p.title, p.is_private, p.content, p.language, p.url, p.expires_at, p.created_at, p.updated_at, COALESCE(a.views, 0) as views FROM pastes p LEFT JOIN pastes_analytics a ON p.id = a."pasteID" WHERE p.url LIKE $1`
	row, err := p.db.Query(ctx, query, "%/p/"+slug)
	if err != nil {
		return nil, fmt.Errorf("failed to query paste by slug: %w", err)
	}
	defer row.Close()
	paste, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[models.PasteOutput])
	if err != nil {
		return nil, fmt.Errorf("failed to collect paste: %w", err)
	}

	// Check if paste has expired
	if paste.ExpiresAt != nil && paste.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("paste has expired")
	}

	// Check if paste is private (should not be accessible publicly)
	if paste.IsPrivate {
		return nil, fmt.Errorf("paste is private")
	}

	// Increment view count for public access
	err = p.incrementViewCount(ctx, paste.ID)
	if err != nil {
		// Log the error but don't fail the paste retrieval
		// since view counting is not critical
		fmt.Printf("Failed to increment view count for paste %s: %v\n", paste.ID, err)
	}

	return &paste, nil
}

func (p *PasteRepository) FilterPastes(ctx context.Context, userID uuid.UUID, pasteFilter *models.PasteFilters) (*[]models.PasteOutput, error) {
	// Build base select query - only select columns that exist in PasteOutput
	builder := sq.Select(
		"p.id",
		"p.user_id",
		"p.title",
		"p.is_private",
		"p.content",
		"p.language",
		"p.url",
		"p.expires_at",
		"COALESCE(a.views, 0) as views",
	).From("pastes p").
		LeftJoin("pastes_analytics a ON p.id = a.\"pasteID\"").
		PlaceholderFormat(sq.Dollar)

	// Exclude expired pastes
	builder = builder.Where(sq.Or{
		sq.Eq{"p.expires_at": nil},
		sq.Gt{"p.expires_at": time.Now()},
	})
	// Only selected users
	builder = builder.Where(sq.Eq{"p.user_id": userID})
	// Apply filters if provided
	if pasteFilter != nil {
		if len(pasteFilter.Languages) > 0 {
			builder = builder.Where(sq.Eq{"p.language": pasteFilter.Languages})
		}
		if pasteFilter.DateFrom != nil {
			builder = builder.Where(sq.GtOrEq{"p.created_at": *pasteFilter.DateFrom})
		}
		if pasteFilter.DateTo != nil {
			builder = builder.Where(sq.LtOrEq{"p.created_at": *pasteFilter.DateTo})
		}

		// Handle sorting with allow-list for security
		sortBy := "p.created_at"
		switch pasteFilter.SortBy {
		case "created_at":
			sortBy = "p.created_at"
		case "updated_at":
			sortBy = "p.updated_at"
		case "views":
			sortBy = "views"
		case "title":
			sortBy = "p.title"
		}

		// Handle sort order
		sortOrder := "DESC"
		if pasteFilter.SortOrder == "asc" || pasteFilter.SortOrder == "ASC" {
			sortOrder = "ASC"
		}
		builder = builder.OrderBy(fmt.Sprintf("%s %s", sortBy, sortOrder))
	} else {
		// Default sorting when no filter is provided
		builder = builder.OrderBy("p.created_at DESC")
	}

	// Build the SQL query
	queryStr, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build filter query: %w", err)
	}

	// Execute the query
	rows, err := p.db.Query(ctx, queryStr, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to filter pastes: %w", err)
	}
	defer rows.Close()

	// Collect rows into structs using pgx
	pastes, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.PasteOutput])
	if err != nil {
		return nil, fmt.Errorf("failed to collect filtered pastes: %w", err)
	}

	return &pastes, nil
}

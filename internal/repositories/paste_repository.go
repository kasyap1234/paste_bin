package repositories

import (
	"context"
	"fmt"
	"os"
	"pastebin/internal/models"
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
	query := `INSERT INTO pastes (user_id,title,is_private,content,language,url,expires_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	title := pasteInput.Title
	if title == "" {
		title = "Untitled"
	}
	urlSlug := uuid.New().String()[:8]
	var isPrivate bool
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	url := baseURL + "/p/" + urlSlug
	if pasteInput.Password == "" {
		isPrivate = false
	} else {
		isPrivate = true
	}
	language := pasteInput.Language
	content := pasteInput.Content
	expiresAt := pasteInput.ExpiresAt
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, query, userID, title, isPrivate, content, language, url, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert paste: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Retrieve the created paste to return it
	getQuery := `SELECT id,user_id,title,is_private,content,language,url,expires_at,created_at,updated_at FROM pastes WHERE url = $1`
	row, err := p.db.Query(ctx, getQuery, url)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created paste: %w", err)
	}
	paste, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[models.PasteOutput])
	if err != nil {
		return nil, fmt.Errorf("failed to collect created paste: %w", err)
	}

	return &paste, nil
}

func (p *PasteRepository) UpdatePaste(ctx context.Context, pasteID uuid.UUID, patchInput *models.PatchPaste) error {
	// Build the update query using squirrel
	updateBuilder := sq.Update("pastes").PlaceholderFormat(sq.Dollar)

	// Add fields to update if they are provided (not nil)
	if patchInput.Title != nil {
		updateBuilder = updateBuilder.Set("title", *patchInput.Title)
	}
	if patchInput.Content != nil {
		updateBuilder = updateBuilder.Set("content", *patchInput.Content)
	}
	if patchInput.Language != nil {
		updateBuilder = updateBuilder.Set("language", *patchInput.Language)
	}
	if patchInput.IsPrivate != nil {
		updateBuilder = updateBuilder.Set("is_private", *patchInput.IsPrivate)
	}
	if patchInput.Password != nil {
		updateBuilder = updateBuilder.Set("password", *patchInput.Password)
	}
	if patchInput.ExpiresAt != nil {
		updateBuilder = updateBuilder.Set("expires_at", *patchInput.ExpiresAt)
	}

	// Always update the updated_at timestamp
	updateBuilder = updateBuilder.Set("updated_at", time.Now()).Where(sq.Eq{"id": pasteID})

	// Check if any fields were provided to update
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
	query := `SELECT p.id,p.user_id,p.title,p.is_private,p.content,p.language,p.url,p.expires_at ,p.created_at,p.updated_at ,COALESCE(a.views,0) as views FROM pastes p LEFT JOIN pastes_analytics a ON p.id=a.paste_id WHERE p.id=$1`
	row, err := p.db.Query(ctx, query, pasteID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no pastes with pasteID=%s", pasteID)
		}
		return nil, fmt.Errorf("failed to query paste: %w", err)
	}
	paste, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[models.PasteOutput])
	if err != nil {
		return nil, fmt.Errorf("failed to collect paste: %w", err)
	}

	// Check if paste has expired
	if paste.ExpiresAt != nil && paste.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("paste has expired")
	}
	// check if it is public
	// if yes
	isOwner := isAuthenticated && paste.UserID == userID
	if !isOwner {
		// it is public view
		err := p.incrementViewCount(ctx, pasteID)
		if err != nil {
			// Log the error but don't fail the paste retrieval
			// since view counting is not critical
			fmt.Printf("Failed to increment view count for paste %s: %v\n", pasteID, err)
		}
	}
	return &paste, nil
}

func (p *PasteRepository) incrementViewCount(ctx context.Context, pasteID uuid.UUID) error {
	query := `INSERT INTO pastes_analytics (pasteid,views,updated_at) VALUES($1,1,NOW())
ON CONFLICT (pasteid)
DO UPDATE SET
views=paste_analytics.views +1 , updated_at=NOW()
`
	_, err := p.db.Exec(ctx, query, pasteID)
	return err

}

func (p *PasteRepository) GetAllPastes(ctx context.Context, userID uuid.UUID) (*[]models.PasteOutput, error) {
	sql := `SELECT p.id,p.user_id,p.title,p.is_private,p.language,p.url,p.expires_at,p.crated_at ,COALESCE(a.views,0) as views  FROM pastes p LEFT JOIN pastes_analytics a ON p.id=a.paste_id WHERE user_id=$1 ORDER BY p.created_at DESC`
	row, err := p.db.Query(ctx, sql, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no pastes found for userID : %s", userID)
		}
		return nil, fmt.Errorf("failed to get pastes for user ID : %s", userID)
	}
	pastes, err := pgx.CollectRows(row, pgx.RowToStructByName[models.PasteOutput])
	if err != nil {
		return nil, fmt.Errorf("failed to collect pastes :  %w", err)
	}

	// Filter out expired pastes
	var validPastes []models.PasteOutput
	for _, paste := range pastes {
		if paste.ExpiresAt == nil || !paste.ExpiresAt.Before(time.Now()) {
			validPastes = append(validPastes, paste)
		}
	}

	return &validPastes, nil
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
	query := `SELECT p.id,p.user_id,p.title,p.is_private,p.content,p.language,p.url,p.expires_at,p.created_at,p.updated_at,COALESCE(a.views,0) as views FROM pastes p LEFT JOIN pastes_analytics a ON p.id=a.paste_id WHERE p.url LIKE $1`
	row, err := p.db.Query(ctx, query, "%/p/"+slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no paste found with slug: %s", slug)
		}
		return nil, fmt.Errorf("failed to query paste by slug: %w", err)
	}
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
		LeftJoin("pastes_analytics a ON p.id = a.paste_id").
		PlaceholderFormat(sq.Dollar)

	// Exclude expired pastes
	builder = builder.Where(sq.Or{
		sq.Eq{"p.expires_at": nil},
		sq.Gt{"p.expires_at": time.Now()},
	})
	// only selected users
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
		if err == pgx.ErrNoRows {
			empty := []models.PasteOutput{}
			return &empty, nil
		}
		return nil, fmt.Errorf("failed to filter pastes: %w", err)
	}

	// Collect rows into structs using pgx
	pastes, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.PasteOutput])
	if err != nil {
		return nil, fmt.Errorf("failed to collect filtered pastes: %w", err)
	}

	return &pastes, nil
}

package repositories

import (
	"context"
	"fmt"
	"pastebin/internal/models"
	"pastebin/pkg/utils"
	"strings"
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

func (p *PasteRepository) CreatePaste(ctx context.Context, userID uuid.UUID, pasteInput *models.PasteInput) error {
	query := `INSERT INTO pastes (user_id,title,is_private,content,language,url,expires_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	title := pasteInput.Title
	if title == "" {
		title = "Untitled"
	}
	urlSlug := uuid.New().String()[:8]
	var isPrivate bool
	url := "https://pastebin.com/" + urlSlug
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
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, query, userID, title, isPrivate, content, language, url, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to insert paste: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (p *PasteRepository) UpdatePaste(ctx context.Context, pasteID uuid.UUID, patchInput *models.PatchPaste) error {
	sets, values, nextIndex := utils.BuildSets(patchInput)
	if len(sets) == 0 {
		return fmt.Errorf("no fields to update")
	}
	values = append(values, pasteID)
	query := fmt.Sprintf(`UPDATE pastes SET %s ,updated_at=NOW() WHERE id=$%d`, strings.Join(sets, ","), nextIndex)
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction : %w", err)
	}
	defer tx.Rollback(ctx)
	cmdTag, err := tx.Exec(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("failed to udpate paste: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("paste not found with id : %s", pasteID.String())

	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction : %w", err)

	}
	return nil
}

func (p *PasteRepository) GetPasteByID(ctx context.Context, pasteID uuid.UUID) (*models.PasteOutput, error) {
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

	return &paste, nil
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
	sql := `DELETE * FROM pastes WHERE paste_id=$1`
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction %w", err)
	}
	defer tx.Rollback(ctx)
	cmdTag, err := tx.Exec(ctx, sql, pasteID)
	if err != nil {
		return fmt.Errorf("unable to delete paste by id ")
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("paste not found with id %s", pasteID)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction %w", err)
	}
	return nil
}

// func (p *PasteRepository) FilterPastes(ctx context.Context, userID uuid.UUID, pasteFilter *models.PasteFilters) (*[]models.PasteOutput, error) {
// 	qb := sq.Select("*").From("pastes").Where(sq)
// }

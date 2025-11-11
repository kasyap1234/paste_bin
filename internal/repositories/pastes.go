package repositories

import (
	"context"
	"pastebin/internal/models"

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

func (p *PasteRepository) CreatePaste(ctx context.Context, pasteInput *models.PasteInput) error {
sql :=`INSERT `
}

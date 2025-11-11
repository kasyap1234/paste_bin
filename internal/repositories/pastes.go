package repositories

import "github.com/jackc/pgx/v5/pgxpool"

type PasteRepository struct {
	db *pgxpool.Pool
}

func NewPasteRepository(db *pgxpool.Pool) *PasteRepository {
	return &PasteRepository{
		db: db,
	}
}

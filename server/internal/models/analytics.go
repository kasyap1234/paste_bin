package models

import (
	"time"

	"github.com/google/uuid"
)

type Analytics struct {
	ID        uuid.UUID `json:"id" db:"id"`
	PasteID   uuid.UUID `json:"paste_id" db:"paste_id"`
	URL       string    `json:"url" db:"url"`
	Views     int       `json:"views" db:"views"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type AnalyticsInput struct {
	PasteID uuid.UUID `json:"paste_id"`
	URL     string    `json:"url"`
}

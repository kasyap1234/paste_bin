package models

import (
	"time"

	"github.com/google/uuid"
)

type Analytics struct {
	ID        uuid.UUID `json:"id"`
	PasteID   uuid.UUID `json:"paste_id"`
	URL       string    `json:"url"`
	Views     int       `json:"views"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}



type AnalyticsInput struct {
	PasteID uuid.UUID `json:"paste_id"`
	URL     string    `json:"url"`
}

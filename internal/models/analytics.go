package models

import "github.com/google/uuid"

type Analytics struct {
	ID      uuid.UUID `json:"id"`
	PasteID uuid.UUID `json:"paste_id"`
	URL     string    `json:"url"`
	Views   int       `json:"views"`
}

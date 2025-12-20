package models

import (
	"time"

	"github.com/google/uuid"
)

type PasteInput struct {
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Language  string     `json:"language"`
	Password  string     `json:"password"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type PasteOutput struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	Title     string     `json:"title"`
	IsPrivate bool       `json:"is_private"`
	Content   string     `json:"content"`
	Language  string     `json:"language"`
	URL       string     `json:"url"`
	Views     int        `json:"views"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type PatchPaste struct {
	ID        *uuid.UUID `json:"id" db:"id"`
	UserID    *uuid.UUID `json:"user_id" db:"user_id"`
	Title     *string    `json:"title" db:"title"`
	Content   *string    `json:"content" db:"content"`
	Language  *string    `json:"language" db:"language"`
	IsPrivate *bool      `json:"is_private" db:"is_private"`
	Password  *string    `json:"password" db:"password"`
	ExpiresAt *time.Time `json:"expires_at" db:"expires_at"`
}

type PaginatedPastesResponse struct {
	Pastes    []PasteOutput `json:"pastes"`
	Total     int           `json:"total"`
	Limit     int           `json:"limit"`
	Offset    int           `json:"offset"`
	HasMore   bool          `json:"has_more"`
}

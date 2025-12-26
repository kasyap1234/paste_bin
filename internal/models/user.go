package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	Avatar       string    `json:"avatar" db:"avatar"`
	PasswordHash string    `json:"-" db:"password_hash"`
}

// PatchProfile represents optional fields for partial profile updates
type PatchProfile struct {
	Name   *string `json:"name,omitempty"`
	Avatar *string `json:"avatar,omitempty"`
}

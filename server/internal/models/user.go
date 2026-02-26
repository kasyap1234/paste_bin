package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                   uuid.UUID  `json:"id" db:"id"`
	Name                 string     `json:"name" db:"name"`
	Email                string     `json:"email" db:"email"`
	Avatar               *string    `json:"avatar" db:"avatar"`
	PasswordHash         string     `json:"-" db:"password_hash"`
	VerifyToken          *string    `json:"-" db:"verify_token"`
	VerifyTokenExpiresAt *time.Time `json:"-" db:"verify_token_expires_at"`
	IsVerified           bool       `json:"is_verified" db:"is_verified"`
	ResetToken           *string    `json:"-" db:"reset_token"`
	ResetTokenExpiresAt  *time.Time `json:"-" db:"reset_token_expires_at"`
}

// PatchProfile represents optional fields for partial profile updates
type PatchProfile struct {
	Name   *string `json:"name,omitempty"`
	Avatar *string `json:"avatar,omitempty"`
}

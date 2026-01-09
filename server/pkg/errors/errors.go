package errors

import "errors"

// Sentinel errors for paste operations
var (
	ErrPasswordRequired = errors.New("password required")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrPasteNotFound    = errors.New("paste not found")
	ErrPasteExpired     = errors.New("paste has expired")
	ErrPermissionDenied = errors.New("permission denied")
)

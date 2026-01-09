package utils

import "github.com/google/uuid"

// GenerateResetToken creates a random token string that can be used
// for password reset or email verification flows.
func GenerateResetToken() string {
	return uuid.New().String()
}

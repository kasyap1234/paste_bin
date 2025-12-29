package utils

import (
	"fmt"
	"os"
)

// SendPasswordResetEmail builds a reset URL with the given token and sends it
// to the provided email address. For now this just logs/prints the URL; in a
// real system you would integrate with an email provider here.
func SendPasswordResetEmail(email, token string) error {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", baseURL, token)

	// TODO: integrate with real email provider.
	// For now we just print the URL so it is visible in logs.
	fmt.Printf("Password reset link for %s: %s\n", email, resetURL)

	return nil
}

func SendVerifyEmail(email string, token string) error {
	// send email for verification
	baseURL := os.Getenv("EMAIL_VERIFICATION_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8001"
	}
	verifyURL := fmt.Sprintf("%s/verify-email/%s", baseURL, token)

	// send the verify url via email to the email id
	fmt.Print(verifyURL)
	return nil
}

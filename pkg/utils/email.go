package utils

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
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
	if err := SendEmail(email, "Password Reset", resetURL); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

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
	if err := SendEmail(email, "Verify Email ", verifyURL); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func SendEmail(to string, subject, body string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY not set")
	}

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "pastebin@gmail.com", // Default sender for testing
		To:      []string{to},
		Subject: subject,
		Html:    body,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	fmt.Printf("Email sent successfully! ID: %s\n", sent.Id)
	return nil
}


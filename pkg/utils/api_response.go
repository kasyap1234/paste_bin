package utils

import (
	"github.com/labstack/echo/v4"
)

// APIResponse represents a standardized API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SendSuccess sends a success response with optional data and message
func SendSuccess(c echo.Context, code int, data interface{}, message string) error {
	return c.JSON(code, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendError sends an error response with a specific error message
func SendError(c echo.Context, code int, message string) error {
	return c.JSON(code, APIResponse{
		Success: false,
		Error:   message,
	})
}

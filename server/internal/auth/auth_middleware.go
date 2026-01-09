package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ContextKey is a typed key for storing values in context.Context to avoid collisions.
type ContextKey string

const (
	userIDCtxKey    ContextKey = "userID"
	userEmailCtxKey ContextKey = "userEmail"
)

// AuthMiddleware validates the Authorization header using the provided JWTManager.
// On success it injects the user's ID and email into the request's context.Context
// using typed context keys. It does NOT use echo.Context's Set/Get map.
func AuthMiddleware(jwtManager *JWTManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := extractToken(c.Request().Header.Get("Authorization"))
			if err != nil {
				return echo.NewHTTPError(401, "missing or invalid authorization")
			}

			claims, err := jwtManager.VerifyToken(token)
			if err != nil {
				return echo.NewHTTPError(401, "invalid token")
			}

			// Put values into the request's context.Context using typed keys.
			req := c.Request()
			ctx := context.WithValue(req.Context(), userIDCtxKey, claims.UserID)
			ctx = context.WithValue(ctx, userEmailCtxKey, claims.Email)
			c.SetRequest(req.WithContext(ctx))

			return next(c)
		}
	}
}

// extractToken extracts a bearer token from an Authorization header value.
func extractToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("auth header missing")
	}
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", errors.New("authorization scheme must be Bearer")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[1] == "" {
		return "", errors.New("authorization header format wrong")
	}
	return parts[1], nil
}

// GetUserIDFromContext reads the user ID from a standard context.Context.
// Returns uuid.Nil with an error if the value is missing or has the wrong type.
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	if ctx == nil {
		return uuid.Nil, errors.New("context is nil")
	}
	val := ctx.Value(userIDCtxKey)
	if val == nil {
		return uuid.Nil, errors.New("user id not found in context")
	}
	id, ok := val.(uuid.UUID)
	if !ok {
		// Sometimes IDs might be stored as strings by other code; try parsing.
		if s, sok := val.(string); sok {
			parsed, err := uuid.Parse(s)
			if err != nil {
				return uuid.Nil, errors.New("user id in context has invalid string format")
			}
			return parsed, nil
		}
		return uuid.Nil, errors.New("user id in context has unexpected type")
	}
	return id, nil
}

// GetUserEmailFromContext reads the user email from a standard context.Context.
func GetUserEmailFromContext(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", errors.New("context is nil")
	}
	val := ctx.Value(userEmailCtxKey)
	if val == nil {
		return "", errors.New("user email not found in context")
	}
	email, ok := val.(string)
	if !ok {
		return "", errors.New("user email in context has unexpected type")
	}
	return email, nil
}

// GetUserIDFromEchoContext reads the user ID from an echo.Context by delegating to the
// request's context.Context. It is a convenience wrapper for handler code.
func GetUserIDFromEchoContext(c echo.Context) (uuid.UUID, error) {
	if c == nil || c.Request() == nil {
		return uuid.Nil, errors.New("echo context or request is nil")
	}
	return GetUserIDFromContext(c.Request().Context())
}

// GetUserEmailFromEchoContext reads the user email from an echo.Context by delegating to the
// request's context.Context. It is a convenience wrapper for handler code.
func GetUserEmailFromEchoContext(c echo.Context) (string, error) {
	if c == nil || c.Request() == nil {
		return "", errors.New("echo context or request is nil")
	}
	return GetUserEmailFromContext(c.Request().Context())
}

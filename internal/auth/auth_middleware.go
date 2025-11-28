package auth

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var userIDKey string
var userEmailKey string

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
			c.Set(string(userIDKey), claims.UserID)
			c.Set(string(userEmailKey), claims.Email)
			return next(c)
		}
	}
}

func extractToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("auth Header missing")
	}
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", errors.New("authorization scheme must be bearer")
	}
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 {
		return "", errors.New("authorization header format wrong")
	}
	token := tokenParts[1]
	return token, nil
}

func GetUserIDFromContext(c echo.Context) (uuid.UUID, error) {
	userID := c.Get(string(userIDKey))
	if userID == nil {
		return uuid.Nil, errors.New("user id not found in context")
	}
	id, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("invalid user id type")
	}
	return id, nil
}

func GetUserEmailFromContext(c echo.Context) (string, error) {
	userEmail := c.Get(string(userEmailKey))
	if userEmail == nil {
		return "", errors.New("email nto found in the context")
	}
	email, ok := userEmail.(string)
	if !ok {
		return "", errors.New("invalid user email type ")
	}
	return email, nil
}

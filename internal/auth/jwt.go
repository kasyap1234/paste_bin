package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the JWT payload used by the application.
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// JWTManager handles creation and verification of JWT tokens.
type JWTManager struct {
	secretKey []byte
}

// NewJWTManager constructs a JWTManager. secretKey should be a sufficiently long random string.
func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{
		secretKey: []byte(secretKey),
	}
}

// GenerateToken creates a signed JWT containing the user's ID and email.
// expirationTime is a duration from now after which the token is invalid.
func (j *JWTManager) GenerateToken(userID uuid.UUID, email string, expirationTime time.Duration) (string, error) {
	if len(j.secretKey) == 0 {
		return "", errors.New("jwt secret key is empty")
	}
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expirationTime)),
			ID:        uuid.NewString(), // jti
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signed, nil
}

// VerifyToken parses and validates a token string and returns the Claims if valid.
func (j *JWTManager) VerifyToken(tokenStr string) (*Claims, error) {
	if tokenStr == "" {
		return nil, errors.New("token is empty")
	}
	if len(j.secretKey) == 0 {
		return nil, errors.New("jwt secret key is empty")
	}

	claims := &Claims{}
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	_, err := parser.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		// ensure signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Explicitly validate expiration time. The jwt library's RegisteredClaims.Valid
	// method signature/behavior can vary between versions, so perform the expiration
	// check directly to avoid depending on that helper.
	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}


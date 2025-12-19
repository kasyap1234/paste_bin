package services

import (
	"context"
	"fmt"
	"log"
	"pastebin/internal/auth"
	"pastebin/internal/models"
	"pastebin/internal/repositories"
	"pastebin/pkg/utils"
	"time"
)

type AuthService struct {
	authRepo   *repositories.AuthRepository
	jwtManager *auth.JWTManager
	userRepo   *repositories.UserRepository
}

func NewAuthService(authRepo *repositories.AuthRepository, userRepo *repositories.UserRepository, jwtMgr *auth.JWTManager) *AuthService {
	return &AuthService{
		authRepo:   authRepo,
		jwtManager: jwtMgr,
		userRepo:   userRepo,
	}
}

func (a *AuthService) Register(ctx context.Context, registerInput *models.RegisterInput) error {
	// check if user already exists
	// only proceed if user does not exists
	ok, err := a.userRepo.ExistsUser(ctx, registerInput.Email)
	if err != nil {
		log.Printf("ExistsUser error: %v", err)
		return fmt.Errorf("error checking existing user: %w", err)
	}
	if ok {
		return fmt.Errorf("user already exists with email: %s", registerInput.Email)
	}
	regErr := a.authRepo.Register(ctx, registerInput)
	if regErr != nil {
		log.Printf("Register error: %v", regErr)
		return regErr
	}
	return nil
}

func (a *AuthService) Login(ctx context.Context, loginInput *models.LoginInput) (*models.LoginResponse, error) {
	user, err := a.userRepo.GetUserByEmail(ctx, loginInput.Email)

	if err != nil {
		log.Printf("GetUserByEmail error: %v", err)
		return nil, fmt.Errorf("invalid email or password")
	}

	if !utils.VerifyPassword(user.Password, loginInput.Password) {
		return nil, fmt.Errorf("invalid email or password")
	}
	token, err := a.jwtManager.GenerateToken(user.ID, user.Email, 24*time.Hour)
	if err != nil {
		log.Printf("GenerateToken error: %v", err)
		return nil, fmt.Errorf("failed to generate token")
	}
	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

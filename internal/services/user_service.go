package services

import (
	"context"
	"pastebin/internal/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}


func(u*UserService)CheckUserExists(ctx context.Context,userID)
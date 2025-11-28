package services

import (
	"context"
	"fmt"
	"pastebin/internal/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}



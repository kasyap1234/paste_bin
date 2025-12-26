package services

import (
	"context"
	"pastebin/internal/models"
	"pastebin/internal/repositories"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type ProfileService struct {
	profileRepo *repositories.ProfileRepository
	logger      zerolog.Logger
}

func NewProfileService(profileRepo *repositories.ProfileRepository, logger zerolog.Logger) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
		logger:      logger,
	}
}

func (p *ProfileService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := p.profileRepo.GetProfile(ctx, userID)
	if err != nil {
		p.logger.Err(err).Msg("failed to get profile")
		return nil, err
	}
	return user, nil
}

func (p *ProfileService) UpdateProfile(ctx context.Context, userID uuid.UUID, patch *models.PatchProfile) (*models.User, error) {
	user, err := p.profileRepo.UpdateProfile(ctx, userID, patch)
	if err != nil {
		p.logger.Err(err).Msg("failed to update profile")
		return nil, err
	}
	return user, nil
}

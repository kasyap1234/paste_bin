package services

import (
	"context"
	"fmt"
	"pastebin/internal/models"
	"pastebin/internal/repositories"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type AnalyticsService struct {
	analyticsRepo *repositories.AnalyticsRepository
	logger        zerolog.Logger
}

func NewAnalyticsService(analyticsRepo *repositories.AnalyticsRepository, logger zerolog.Logger) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo: analyticsRepo,
		logger:        logger,
	}
}

func (s *AnalyticsService) CreateAnalytics(ctx context.Context, pasteID uuid.UUID, url string) error {
	if pasteID == uuid.Nil {
		return fmt.Errorf("unable to create analytics for nil pasteID")
	}
	if url == "" {
		return fmt.Errorf("unable to create analytics for empty url")
	}
	analytics, err := s.analyticsRepo.GetAnalyticsByPasteID(ctx, pasteID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get analytics by pasteID")
		return fmt.Errorf("failed to get analytics by pasteID: %w", err)
	}
	if analytics == nil {
		err = s.analyticsRepo.CreateAnalytics(ctx, pasteID, url)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to create analytics")
			return fmt.Errorf("failed to create analytics: %w", err)
		}
		return nil
	}
	return nil
}

func (s *AnalyticsService) GetAnalyticsByPasteID(ctx context.Context, pasteID uuid.UUID) (*models.Analytics, error) {
	if pasteID == uuid.Nil {
		return nil, fmt.Errorf("unable to get analytics for nil pasteID %s", pasteID)
	}
	return s.analyticsRepo.GetAnalyticsByPasteID(ctx, pasteID)
}

func (s *AnalyticsService) GetAnalyticsByID(ctx context.Context, ID uuid.UUID) (*models.Analytics, error) {
	if ID == uuid.Nil {
		return nil, fmt.Errorf("unable to get analytics for nil analytics id ")
	}
	return s.analyticsRepo.GetAnalyticsByID(ctx, ID)

}

func (s *AnalyticsService) IncrementViews(ctx context.Context, pasteID uuid.UUID) error {
	if pasteID == uuid.Nil {
		return fmt.Errorf("unable to increment views for nil pasteID")
	}

	return s.analyticsRepo.IncrementViews(ctx, pasteID)
}

func (s *AnalyticsService) GetAnalyticsByURL(ctx context.Context, url string) (*models.Analytics, error) {
	if url == "" {
		return nil, fmt.Errorf("unable to get analytics for empty url %s ", url)
	}
	analytics, err := s.analyticsRepo.GetAnalyticsByURL(ctx, url)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get analytics by url")
		return nil, fmt.Errorf("failed to get analytics by url: %w", err)
	}
	if analytics == nil {
		return nil, fmt.Errorf("analytics not found for url %s", url)
	}
	return analytics, nil
}

func (s *AnalyticsService) GetAllAnalytics(ctx context.Context, order string, limit int, offset int) ([]models.Analytics, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	if order == "" {
		order = "created_at DESC"
	}
	return s.analyticsRepo.GetAllAnalytics(ctx, order, limit, offset)
}

func (s *AnalyticsService) GetAllAnalyticsByUser(ctx context.Context, userID uuid.UUID, order string, limit, offset int) ([]models.Analytics, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("unable to get analytics for nil userID")
	}
	if order == "" {
		// by default use order created_at desc
		order = "created_at DESC"
	}

	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}
	analytics, err := s.analyticsRepo.GetAllAnalyticsByUser(ctx, userID, order, limit, offset)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get all analytics by user")
		return nil, fmt.Errorf("unable to get all analytics by user: %w", err)
	}
	return analytics, nil
}

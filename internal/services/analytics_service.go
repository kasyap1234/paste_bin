package services

import (
	"context"
	"fmt"
	"pastebin/internal/models"
	"pastebin/internal/repositories"

	"github.com/google/uuid"
)

type AnalyticsService struct {
	analyticsRepo *repositories.AnalyticsRepository
}

func NewAnalyticsService(analyticsRepo *repositories.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo: analyticsRepo,
	}
}

func (s *AnalyticsService) CreateAnalytics(ctx context.Context, pasteID uuid.UUID, url string) error {
	if url == "" {
		return fmt.Errorf("unable to create analytics for empty url %s ", url)
	}
	_, err := s.analyticsRepo.GetAnalyticsByPasteID(ctx, pasteID)
	if err != nil {
		return fmt.Errorf("unable to get analytics by pasteID %s : %w", pasteID, err)
	}

	return s.analyticsRepo.CreateAnalytics(ctx, pasteID, url)
}

func (s *AnalyticsService) GetAnalyticsByPasteID(ctx context.Context, pasteID uuid.UUID) (*models.Analytics, error) {
	if pasteID == uuid.Nil {
		return nil, fmt.Errorf("unable to get analytics for nil pasteID")
	}
	return s.analyticsRepo.GetAnalyticsByPasteID(ctx, pasteID)
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

	return s.analyticsRepo.GetAnalyticsByURL(ctx, url)
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

func (s *AnalyticsService) GetAllAnalyticsByUser(ctx context.Context, userID uuid.UUID, order string, limit, offset int) (*[]models.Analytics, error) {
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
	return s.analyticsRepo.GetAllAnalyticsByUser(ctx, userID, order, limit, offset)

}

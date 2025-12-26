package services

import (
	"context"
	"fmt"
	"pastebin/internal/auth"
	"pastebin/internal/models"
	"pastebin/internal/repositories"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type PasteService struct {
	pasteRepo *repositories.PasteRepository
	logger    zerolog.Logger
}

func NewPasteService(pasteRepo *repositories.PasteRepository, logger zerolog.Logger) *PasteService {
	return &PasteService{
		pasteRepo: pasteRepo,
		logger: logger,
	}
}
func (p *PasteService) CreatePaste(ctx context.Context, createPaste *models.PasteInput) (*models.PasteOutput, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to get userID from context")
		return nil, fmt.Errorf("unable to get userID from context: %w", err)
	}
	paste, err := p.pasteRepo.CreatePaste(ctx, userID, createPaste)
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to create paste")
		return nil, fmt.Errorf("unable to create paste: %w", err)
	}
	return paste, nil
}

func (p *PasteService) UpdatePaste(ctx context.Context, pasteID uuid.UUID, patchPaste *models.PatchPaste) error {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to get userID from context")
		return fmt.Errorf("unable to get userID from context: %w", err)
	}

	_, err = p.GetPasteByID(ctx, pasteID, true, userID, "")
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to get paste by ID")
		return fmt.Errorf("unable to find paste with ID: %s ", pasteID)
	}
	err = p.pasteRepo.UpdatePaste(ctx, pasteID, patchPaste)
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to update paste")
		return fmt.Errorf("unable to update paste: %w", err)
	}
	return nil
}

func (p *PasteService) GetPasteByID(ctx context.Context, pasteID uuid.UUID, isAuthenticated bool, userID uuid.UUID, password string) (*models.PasteOutput, error) {

	paste, err := p.pasteRepo.GetPasteByID(ctx, pasteID, isAuthenticated, userID, password)
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to get paste by ID")
		return nil, fmt.Errorf("unable to get paste by ID: %w", err)
	}
	return paste, nil
}

func (p *PasteService) GetAllPastes(ctx context.Context, userID uuid.UUID, limit, offset int) (*models.PaginatedPastesResponse, error) {
	// Validate and set defaults
	if limit <= 0 {
		limit = 10 // default limit
	}
	if limit > 100 {
		limit = 100 // max limit to prevent abuse
	}
	if offset < 0 {
		offset = 0
	}

	pastes, total, err := p.pasteRepo.GetAllPastes(ctx, userID, limit, offset)
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to get pastes")
		return nil, fmt.Errorf("unable to get pastes: %w", err)
	}

	hasMore := offset+limit < total

	return &models.PaginatedPastesResponse{
		Pastes:  pastes,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: hasMore,
	}, nil
}

func (p *PasteService) DeletePasteByID(ctx context.Context, pasteID uuid.UUID) error {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to get userID from context")
		return fmt.Errorf("unable to get userID from context : %w", err)
	}
	paste, err := p.pasteRepo.GetPasteByID(ctx, pasteID, true, userID, "")
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to get paste by ID")
		return fmt.Errorf("unable to get paste by ID: %w", err)
	}
	if paste.UserID != userID {
		p.logger.Error().Msg("user does not have permission to delete this paste")
		return fmt.Errorf("user does not have permission to delete this paste")
	}

	err = p.pasteRepo.DeletePasteByID(ctx, pasteID)
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to delete paste by ID")
		return fmt.Errorf("unable to delete paste by ID: %w", err)
	}
	return nil					
}

func (p *PasteService) FilterPastes(ctx context.Context, filter *models.PasteFilters) (*[]models.PasteOutput, error) {
	if filter == nil {
		p.logger.Error().Msg("filter is nil")
		return nil, fmt.Errorf("filter is nil")
	}

	
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to get userID from context")
		return nil, fmt.Errorf("unable to get userID from context: %w", err)
	}
	if userID == uuid.Nil {
		p.logger.Error().Msg("userID is nil")
		return nil, fmt.Errorf("userID is nil")
	}
	
	pastes, err := p.pasteRepo.FilterPastes(ctx, userID, filter)
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to filter pastes")
		return nil, fmt.Errorf("unable to filter pastes: %w", err)
	}
	return pastes, nil
}

func (p *PasteService) GetPasteBySlug(ctx context.Context, slug, password string) (*models.PasteOutput, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to get userID from context")
		return nil, fmt.Errorf("unable to get userID from context: %w", err)
	}
	if userID == uuid.Nil {
		p.logger.Error().Msg("userID is nil")
		return nil, fmt.Errorf("userID is nil")
	}
	paste, err := p.pasteRepo.GetPasteBySlug(ctx, slug, password)
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to get paste by slug")
		return nil, fmt.Errorf("unable to get paste by slug: %w", err)
	}
	if paste == nil {
		p.logger.Error().Msg("paste is nil")
		return nil, fmt.Errorf("unable to get paste by slug: %w", err)
	}
	if paste.IsPrivate && paste.UserID != userID {
		p.logger.Error().Msg("user does not have permission to view this paste")
		return nil, fmt.Errorf("user does not have permission to view this paste")
	}
	return paste, nil
}

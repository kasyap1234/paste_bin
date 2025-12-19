package services

import (
	"context"
	"fmt"
	"pastebin/internal/auth"
	"pastebin/internal/models"
	"pastebin/internal/repositories"

	"github.com/google/uuid"
)

type PasteService struct {
	pasteRepo *repositories.PasteRepository
}

func NewPasteService(pasteRepo *repositories.PasteRepository) *PasteService {
	return &PasteService{
		pasteRepo: pasteRepo,
	}
}
func (s *PasteService) CreatePaste(ctx context.Context, createPaste *models.PasteInput) (*models.PasteOutput, error) {
	val := ctx.Value("userIDKey")
	if val == nil {
		return nil, fmt.Errorf("missing userID in the context")
	}

	userID, ok := val.(uuid.UUID)
	if !ok {
		str, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("user ID is not a valid string")
		}
		var err error
		userID, err = uuid.Parse(str)
		if err != nil {
			return nil, fmt.Errorf("unable to parse userID : %w", err)
		}
	}
	return s.pasteRepo.CreatePaste(ctx, userID, createPaste)
}

func (p *PasteService) UpdatePaste(ctx context.Context, pasteID uuid.UUID, patchPaste *models.PatchPaste) error {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("unable to get userID from context: %w", err)
	}

	_, err = p.GetPasteByID(ctx, pasteID, true, userID)
	if err != nil {
		return fmt.Errorf("unable to find paste with ID: %s ", pasteID)
	}
	err = p.pasteRepo.UpdatePaste(ctx, pasteID, patchPaste)
	return err
}

func (p *PasteService) GetPasteByID(ctx context.Context, pasteID uuid.UUID, isAuthenticated bool, userID uuid.UUID) (*models.PasteOutput, error) {

	return p.pasteRepo.GetPasteByID(ctx, pasteID, isAuthenticated, userID)
}

func (p *PasteService) GetAllPastes(ctx context.Context, userID uuid.UUID) (*[]models.PasteOutput, error) {

	pastes, err := p.pasteRepo.GetAllPastes(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("unable to get pastes: %w", err)
	}
	return pastes, nil
}

func (p *PasteService) DeletePasteByID(ctx context.Context, pasteID uuid.UUID) error {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("unable to get userID from context : %w", err)
	}
	paste, err := p.pasteRepo.GetPasteByID(ctx, pasteID, true, userID)
	if err != nil {
		return fmt.Errorf("unable to get paste by ID: %w", err)
	}
	if paste.UserID != userID {
		return fmt.Errorf("user does not have permission to delete this paste")
	}

	err = p.pasteRepo.DeletePasteByID(ctx, pasteID)
	return err
}

func (p *PasteService) FilterPastes(ctx context.Context, filter *models.PasteFilters) (*[]models.PasteOutput, error) {

	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get userID from context: %w", err)
	}
	pastes, err := p.pasteRepo.FilterPastes(ctx, userID, filter)
	if err != nil {
		return nil, fmt.Errorf("unable to filter pastes: %w", err)
	}
	return pastes, nil
}

func (p *PasteService) GetPasteBySlug(ctx context.Context, slug string) (*models.PasteOutput, error) {
	return p.pasteRepo.GetPasteBySlug(ctx, slug)
}

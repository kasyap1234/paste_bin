package services

import (
	"context"
	"fmt"
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
func (s *PasteService) CreatePaste(ctx context.Context, createPaste *models.PasteInput) error {
	val := ctx.Value("userIDKey")
	if val == nil {
		return fmt.Errorf("missing userID in the context")
	}

	userID, ok := val.(uuid.UUID)
	if !ok {
		str, ok := val.(string)
		if !ok {
			return fmt.Errorf("user ID is not a valid string")
		}
		var err error
		userID, err = uuid.Parse(str)
		if err != nil {
			return fmt.Errorf("unable to parse userID : %w", err)
		}
	}
	return s.pasteRepo.CreatePaste(ctx, userID, createPaste)
}

func (p *PasteService) UpdatePaste(ctx context.Context, pasteID uuid.UUID, patchPaste *models.PatchPaste) error {
	_, err := p.GetPasteByID(ctx, pasteID)
	if err != nil {
		return fmt.Errorf("unable to find paste with ID: %s ", pasteID)
	}
	err = p.pasteRepo.UpdatePaste(ctx, pasteID, patchPaste)
	return err
}

func (p *PasteService) GetPasteByID(ctx context.Context, pasteID uuid.UUID) (*models.PasteOutput, error) {

	return p.pasteRepo.GetPasteByID(ctx, pasteID)
}


func(p*PasteService)GetAllPastes(ctx context.Context,userID uuid.UUID)(*[]models.PasteOutput,error){
	// check if user exists or not 
	
	pastes,err :=p.pasteRepo.GetAllPastes(ctx,userID)
	return pastes,nil 
}
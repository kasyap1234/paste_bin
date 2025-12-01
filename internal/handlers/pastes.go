package handlers

import (
	"net/http"
	"pastebin/internal/auth"
	"pastebin/internal/models"
	"pastebin/internal/services"

	"github.com/labstack/echo/v4"
)

type PasteHandler struct {
	pasteSvc *services.PasteService
}

func NewPasteHandler(pasteSvc *services.PasteService) *PasteHandler {
	return &PasteHandler{
		pasteSvc: pasteSvc,
	}
}

func (p *PasteHandler) CreatePaste(c echo.Context) error {
	var createPaste *models.PasteInput
	if err := c.Bind(createPaste); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "please use valid paste"})
	}

	ctx := c.Request().Context()
	if err := p.pasteSvc.CreatePaste(ctx, createPaste); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "unable to create a paste"})

	}
	return c.JSON(201, echo.Map{"error": "paste created"})

}

func (p *PasteHandler) UpdatePaste(c echo.Context) error {
	var patchPaste *models.PatchPaste
	if err := c.Bind(patchPaste); err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}
	ctx := c.Request().Context()
	if err := p.pasteSvc.UpdatePaste(ctx, *patchPaste.ID, patchPaste); err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{"error": "unable to update paste"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "paste updated "})

}

func (p *PasteHandler) ListAllPastes(c echo.Context) error {
	userID, err := auth.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unable to get userID "})
	}
	
}

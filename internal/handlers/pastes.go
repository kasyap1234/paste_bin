package handlers

import (
	"net/http"
	"pastebin/internal/auth"
	"pastebin/internal/models"
	"pastebin/internal/services"

	"github.com/google/uuid"
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
	var createPaste models.PasteInput
	if err := c.Bind(&createPaste); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "please use valid paste"})
	}

	ctx := c.Request().Context()
	if err := p.pasteSvc.CreatePaste(ctx, &createPaste); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "unable to create a paste"})

	}
	return c.JSON(201, echo.Map{"error": "paste created"})

}

func (p *PasteHandler) UpdatePaste(c echo.Context) error {
	var patchPaste models.PatchPaste
	if err := c.Bind(&patchPaste); err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}
	ctx := c.Request().Context()
	if err := p.pasteSvc.UpdatePaste(ctx, *patchPaste.ID, &patchPaste); err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{"error": "unable to update paste"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "paste updated "})

}

func (p *PasteHandler) GetAllPastes(c echo.Context) error {
	userID, err := auth.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unable to get userID "})
	}

	pastes, err := p.pasteSvc.GetAllPastes(c.Request().Context(), userID)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unable to get userID "})
	}

	return c.JSON(http.StatusOK, pastes)
}

func (p *PasteHandler) GetPasteByID(c echo.Context) error {
	pasteIDParam := c.Param("id")
	if pasteIDParam == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "paste ID is required"})
	}
	pasteID, err := uuid.Parse(pasteIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "paste ID is invalid"})
	}
	ctx := c.Request().Context()
	paste, err := p.pasteSvc.GetPasteByID(ctx, pasteID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "unable to get paste"})
	}
	return c.JSON(http.StatusOK, paste)
}

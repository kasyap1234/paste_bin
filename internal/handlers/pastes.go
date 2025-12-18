package handlers

import (
	"net/http"
	"pastebin/internal/auth"
	"pastebin/internal/models"
	"pastebin/internal/services"
	"time"

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

// CreatePaste godoc
//
//	@Summary		Create a new paste
//	@Description	Create a new paste with optional expiration
//	@Tags			pastes
//	@Accept			json
//	@Produce		json
//	@Param			request		body		models.PasteInput		true	"Paste data"
//	@Param			expires_in	query		string					false	"Expiration duration (e.g., '24h', '7d')"
//	@Success		201			{object}	map[string]string		"Paste created successfully"
//	@Failure		400			{object}	map[string]string		"Invalid request"
//	@Failure		500			{object}	map[string]string		"Unable to create paste"
//	@Security		BearerAuth
//	@Router			/paste [post]
func (p *PasteHandler) CreatePaste(c echo.Context) error {
	var createPaste models.PasteInput
	if err := c.Bind(&createPaste); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "please use valid paste"})
	}

	// Handle expiry parameter from query string
	expiresIn := c.QueryParam("expires_in")
	if expiresIn != "" {
		duration, err := time.ParseDuration(expiresIn)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid expires_in format, use duration format like '24h', '7d', etc."})
		}
		expiresAt := time.Now().Add(duration)
		createPaste.ExpiresAt = &expiresAt
	}

	ctx := c.Request().Context()
	if err := p.pasteSvc.CreatePaste(ctx, &createPaste); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "unable to create a paste"})

	}
	return c.JSON(201, echo.Map{"error": "paste created"})

}

// UpdatePaste godoc
//
//	@Summary		Update a paste
//	@Description	Update an existing paste by ID
//	@Tags			pastes
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Paste ID"
//	@Param			request	body		models.PatchPaste	true	"Paste update data"
//	@Success		200		{object}	map[string]string	"Paste updated successfully"
//	@Failure		400		{object}	map[string]string	"Invalid request"
//	@Failure		500		{object}	map[string]string	"Unable to update paste"
//	@Security		BearerAuth
//	@Router			/paste/{id} [put]
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

// GetAllPastes godoc
//
//	@Summary		Get all pastes for user
//	@Description	Retrieve all pastes for the authenticated user
//	@Tags			pastes
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.PasteOutput	"List of pastes"
//	@Failure		401	{object}	map[string]string	"Unauthorized"
//	@Security		BearerAuth
//	@Router			/pastes [get]
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

// GetPasteByID godoc
//
//	@Summary		Get paste by ID
//	@Description	Retrieve a specific paste by its ID
//	@Tags			pastes
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string				true	"Paste ID"
//	@Success		200	{object}	models.PasteOutput	"Paste data"
//	@Failure		400	{object}	map[string]string	"Invalid paste ID"
//	@Failure		500	{object}	map[string]string	"Unable to get paste"
//	@Router			/paste/{id} [get]
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
	userID, err := auth.GetUserIDFromContext(ctx)
	isAuthenticated := err == nil
	var requestUserID uuid.UUID
	if isAuthenticated {
		requestUserID = userID
	}
	paste, err := p.pasteSvc.GetPasteByID(ctx, pasteID, isAuthenticated, requestUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "unable to get paste"})
	}
	return c.JSON(http.StatusOK, paste)
}

// DeletePasteByID godoc
//
//	@Summary		Delete paste by ID
//	@Description	Delete a specific paste by its ID
//	@Tags			pastes
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string				true	"Paste ID"
//	@Success		200	{object}	map[string]string	"Paste deleted successfully"
//	@Failure		400	{object}	map[string]string	"Invalid paste ID"
//	@Failure		500	{object}	map[string]string	"Cannot delete paste"
//	@Security		BearerAuth
//	@Router			/paste/{id} [delete]
func (p *PasteHandler) DeletePasteByID(c echo.Context) error {
	pasteIDParam := c.Param("id")
	if pasteIDParam == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "paste id is required"})

	}
	pasteID, err := uuid.Parse(pasteIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "paste id is invalid"})
	}
	ctx := c.Request().Context()
	err = p.pasteSvc.DeletePasteByID(ctx, pasteID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "cannot delete by paste id "})

	}
	return c.JSON(http.StatusOK, echo.Map{"error": "deleted"})
}

package handlers

import (
	"net/http"
	"pastebin/internal/auth"
	"pastebin/internal/models"
	"pastebin/internal/services"
	"pastebin/pkg/utils"
	"strconv"
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
//	@Success		201			{object}	models.PasteOutput		"Created paste with shareable URL"
//	@Failure		400			{object}	map[string]string		"Invalid request"
//	@Failure		500			{object}	map[string]string		"Unable to create paste"
//	@Security		BearerAuth
//	@Router			/paste [post]
func (p *PasteHandler) CreatePaste(c echo.Context) error {
	var createPaste models.PasteInput
	if err := c.Bind(&createPaste); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "please use valid paste")
	}

	// Handle expiry parameter from query string
	expiresIn := c.QueryParam("expires_in")
	if expiresIn != "" {
		duration, err := time.ParseDuration(expiresIn)
		if err != nil {
			return utils.SendError(c, http.StatusBadRequest, "invalid expires_in format, use duration format like '24h', '7d', etc.")
		}
		expiresAt := time.Now().Add(duration)
		createPaste.ExpiresAt = &expiresAt
	}

	ctx := c.Request().Context()
	paste, err := p.pasteSvc.CreatePaste(ctx, &createPaste)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "unable to create a paste")
	}

	return utils.SendSuccess(c, http.StatusCreated, paste, "paste created successfully")
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
	// Get paste ID from URL path parameter
	pasteIDParam := c.Param("id")
	if pasteIDParam == "" {
		return utils.SendError(c, http.StatusBadRequest, "paste ID is required")
	}
	pasteID, err := uuid.Parse(pasteIDParam)
	if err != nil {
		return utils.SendError(c, http.StatusBadRequest, "invalid paste ID")
	}

	var patchPaste models.PatchPaste
	if err := c.Bind(&patchPaste); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "invalid request body")
	}

	ctx := c.Request().Context()
	if err := p.pasteSvc.UpdatePaste(ctx, pasteID, &patchPaste); err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "unable to update paste")
	}
	return utils.SendSuccess(c, http.StatusOK, nil, "paste updated successfully")
}

// GetAllPastes godoc
//
//	@Summary		Get all pastes for user
//	@Description	Retrieve all pastes for the authenticated user with pagination
//	@Tags			pastes
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Number of pastes to return (default: 10, max: 100)"
//	@Param			offset	query		int		false	"Number of pastes to skip (default: 0)"
//	@Success		200		{object}	models.PaginatedPastesResponse	"Paginated list of pastes"
//	@Failure		400		{object}	map[string]string	"Invalid pagination parameters"
//	@Failure		500		{object}	map[string]string	"Unable to get pastes"
//	@Security		BearerAuth
//	@Router			/pastes [get]
func (p *PasteHandler) GetAllPastes(c echo.Context) error {
	userID, _ := auth.GetUserIDFromEchoContext(c) // Middleware ensures this succeeds

	// Parse pagination parameters
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit := 10 // default limit
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			return utils.SendError(c, http.StatusBadRequest, "invalid limit parameter")
		}
	}

	offset := 0 // default offset
	if offsetStr != "" {
		var err error
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			return utils.SendError(c, http.StatusBadRequest, "invalid offset parameter")
		}
	}

	result, err := p.pasteSvc.GetAllPastes(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "unable to get pastes")
	}

	return utils.SendSuccess(c, http.StatusOK, result, "pastes retrieved")
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
	password := c.QueryParam("password")
	paste, err := p.pasteSvc.GetPasteByID(ctx, pasteID, isAuthenticated, requestUserID, password)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "unable to get paste")
	}
	return utils.SendSuccess(c, http.StatusOK, paste, "paste details")
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
//	@Failure		500	{object}	map[string]string	"Unable to delete paste"
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
		return utils.SendError(c, http.StatusInternalServerError, "cannot delete by paste id")
	}
	return utils.SendSuccess(c, http.StatusOK, nil, "paste deleted successfully")
}

// GetPublicPaste godoc
//
//	@Summary		Get public paste by slug
//	@Description	Retrieve a public paste by its URL slug
//	@Tags			pastes
//	@Accept			json
//	@Produce		json
//	@Param			slug	path		string				true	"Paste slug"
//	@Success		200		{object}	models.PasteOutput	"Paste data"
//	@Failure		400		{object}	map[string]string	"Invalid slug"
//	@Failure		404		{object}	map[string]string	"Paste not found"
//	@Failure		500		{object}	map[string]string	"Unable to get paste"
//	@Router			/p/{slug} [get]
func (p *PasteHandler) GetPublicPaste(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "paste slug is required"})
	}

	ctx := c.Request().Context()
	paste, err := p.pasteSvc.GetPasteBySlug(ctx, slug)
	if err != nil {
		return utils.SendError(c, http.StatusNotFound, "paste not found")
	}

	return utils.SendSuccess(c, http.StatusOK, paste, "paste details")
}

// GetRawPaste godoc
//
//	@Summary		Get raw paste content by slug
//	@Description	Retrieve raw text content of a public paste by its URL slug
//	@Tags			pastes
//	@Accept			json
//	@Produce		text/plain
//	@Param			slug	path		string				true	"Paste slug"
//	@Success		200		{string}	string				"Raw paste content"
//	@Failure		400		{object}	map[string]string	"Invalid slug"
//	@Failure		404		{object}	map[string]string	"Paste not found"
//	@Failure		500		{object}	map[string]string	"Unable to get paste"
//	@Router			/raw/{slug} [get]
func (p *PasteHandler) GetRawPaste(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "paste slug is required"})
	}

	ctx := c.Request().Context()
	paste, err := p.pasteSvc.GetPasteBySlug(ctx, slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "paste not found"})
	}

	c.Response().Header().Set("Content-Type", "text/plain")
	return c.String(http.StatusOK, paste.Content)
}

func (p *PasteHandler) FilterPastes(c echo.Context) error {
	var filter models.PasteFilters
	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}
	ctx := c.Request().Context()
	pastes, err := p.pasteSvc.FilterPastes(ctx, &filter)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "unable to filter pastes")
	}
	return utils.SendSuccess(c, http.StatusOK, pastes, "filtered pastes")
}

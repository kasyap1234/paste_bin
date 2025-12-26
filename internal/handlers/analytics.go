package handlers

import (
	"fmt"
	"net/http"
	"pastebin/internal/models"
	"pastebin/internal/services"
	"pastebin/pkg/utils"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type AnalyticsHandler struct {
	analyticsSvc *services.AnalyticsService
	logger       zerolog.Logger
}

func NewAnalyticsHandler(analyticsSvc *services.AnalyticsService, logger zerolog.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsSvc: analyticsSvc,
		logger:       logger,
	}
}

// CreateAnalytics godoc
//
//	@Summary		Create analytics entry
//	@Description	Create a new analytics entry for a paste
//	@Tags			analytics
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.AnalyticsInput	true	"Analytics data"
//	@Success		201		{object}	map[string]string		"Analytics created successfully"
//	@Failure		400		{object}	map[string]string		"Invalid request"
//	@Failure		500		{object}	map[string]string		"Unable to create analytics"
//	@Security		BearerAuth
//	@Router			/create-analytics [post]
func (h *AnalyticsHandler) CreateAnalytics(c echo.Context) error {
	var createAnalytics models.AnalyticsInput
	if err := c.Bind(&createAnalytics); err != nil {
		h.logger.Error().Err(err).Msg("failed to bind create analytics")
		return utils.SendError(c, http.StatusBadRequest, "invalid request body")
	}
	ctx := c.Request().Context()
	if err := h.analyticsSvc.CreateAnalytics(ctx, createAnalytics.PasteID, createAnalytics.URL); err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "failed to create analytics")
	}
	return utils.SendSuccess(c, http.StatusCreated, nil, "analytics created successfully")
}

// GetAllAnalytics godoc
//
//	@Summary		Get all analytics
//	@Description	Retrieve all analytics with pagination
//	@Tags			analytics
//	@Accept			json
//	@Produce		json
//	@Param			order	query		string	false	"Sort order"
//	@Param			limit	query		int		false	"Limit number of results"
//	@Param			offset	query		int		false	"Offset for pagination"
//	@Success		200		{array}		models.Analytics	"List of analytics"
//	@Failure		400		{object}	map[string]string	"Invalid parameters"
//	@Failure		500		{object}	map[string]string	"Unable to get analytics"
//	@Security		BearerAuth
//	@Router			/analytics [get]
func (h *AnalyticsHandler) GetAllAnalytics(c echo.Context) error {
	order := c.QueryParam("order")
	limitStr := c.QueryParam("limit")
	limit := 10 // default limit
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return utils.SendError(c, http.StatusBadRequest, "invalid limit")
		}
	}
	offsetStr := c.QueryParam("offset")
	offset := 0 // default offset
	if offsetStr != "" {
		var err error
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return utils.SendError(c, http.StatusBadRequest, "invalid offset")
		}
	}
	ctx := c.Request().Context()
	analytics, err := h.analyticsSvc.GetAllAnalytics(ctx, order, limit, offset)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "failed to retrieve analytics")
	}
	return utils.SendSuccess(c, http.StatusOK, analytics, "analytics retrieved successfully")
}

// GetAllAnalyticsByUser godoc
//
//	@Summary		Get analytics by user
//	@Description	Retrieve analytics for a specific user with pagination
//	@Tags			analytics
//	@Accept			json
//	@Produce		json
//	@Param			userID	query		string			true	"User ID"
//	@Param			order	query		string			false	"Sort order"
//	@Param			limit	query		int				false	"Limit number of results"
//	@Param			offset	query		int				false	"Offset for pagination"
//	@Success		200		{array}		models.Analytics	"List of analytics"
//	@Failure		400		{object}	map[string]string	"Invalid parameters"
//	@Failure		500		{object}	map[string]string	"Unable to get analytics"
//	@Security		BearerAuth
//	@Router			/analytics/user [get]
func (h *AnalyticsHandler) GetAllAnalyticsByUser(c echo.Context) error {
	userIDStr := c.QueryParam("userID")
	if userIDStr == "" {
		return utils.SendError(c, http.StatusBadRequest, "user id is required")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.SendError(c, http.StatusBadRequest, fmt.Sprintf("invalid user id: %s", userIDStr))
	}
	order := c.QueryParam("order")
	limitStr := c.QueryParam("limit")
	limit := 10 // default limit
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return utils.SendError(c, http.StatusBadRequest, "invalid limit")
		}
	}
	offsetStr := c.QueryParam("offset")
	offset := 0 // default offset
	if offsetStr != "" {
		var err error
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return utils.SendError(c, http.StatusBadRequest, "invalid offset")
		}
	}
	ctx := c.Request().Context()
	analytics, err := h.analyticsSvc.GetAllAnalyticsByUser(ctx, userID, order, limit, offset)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "failed to retrieve analytics by user")
	}
	return utils.SendSuccess(c, http.StatusOK, analytics, "user analytics retrieved successfully")
}

// GetAnalyticsByID godoc
//
//	@Summary		Get analytics by ID
//	@Description	Retrieve specific analytics entry by its ID
//	@Tags			analytics
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"Analytics ID"
//	@Success		200	{object}	models.Analytics	"Analytics data"
//	@Failure		400	{object}	map[string]string	"Invalid ID"
//	@Failure		500	{object}	map[string]string	"Unable to get analytics"
//	@Security		BearerAuth
//	@Router			/analytics/{id} [get]
func (h *AnalyticsHandler) GetAnalyticsByID(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return utils.SendError(c, http.StatusBadRequest, "id is required")
	}
	ID, err := uuid.Parse(idStr)
	if err != nil {
		return utils.SendError(c, http.StatusBadRequest, fmt.Sprintf("invalid analytics id: %s", idStr))
	}
	ctx := c.Request().Context()
	analytic, err := h.analyticsSvc.GetAnalyticsByID(ctx, ID)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "failed to retrieve analytics")
	}
	return utils.SendSuccess(c, http.StatusOK, analytic, "analytics retrieved successfully")
}

// GetAnalyticsByPasteID godoc
//
//	@Summary		Get analytics by paste ID
//	@Description	Retrieve analytics for a specific paste
//	@Tags			analytics
//	@Accept			json
//	@Produce		json
//	@Param			pasteID	query		string			true	"Paste ID"
//	@Success		200		{object}	models.Analytics	"Analytics data"
//	@Failure		400		{object}	map[string]string	"Invalid paste ID"
//	@Failure		500		{object}	map[string]string	"Unable to get analytics"
//	@Security		BearerAuth
//	@Router			/analytics/paste [get]
func (h *AnalyticsHandler) GetAnalyticsByPasteID(c echo.Context) error {
	pasteIDStr := c.QueryParam("pasteID")
	if pasteIDStr == "" {
		return utils.SendError(c, http.StatusBadRequest, "paste id is required")
	}
	pasteID, err := uuid.Parse(pasteIDStr)
	if err != nil {
		return utils.SendError(c, http.StatusBadRequest, fmt.Sprintf("invalid paste id: %s", pasteIDStr))
	}
	ctx := c.Request().Context()
	analytic, err := h.analyticsSvc.GetAnalyticsByPasteID(ctx, pasteID)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "failed to retrieve analytics by paste id")
	}
	return utils.SendSuccess(c, http.StatusOK, analytic, "paste analytics retrieved successfully")
}

package handlers

import (
	"fmt"
	"net/http"
	"pastebin/internal/models"
	"pastebin/internal/services"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AnalyticsHandler struct {
	analyticsSvc *services.AnalyticsService
}

func NewAnalyticsHandler(analyticsSvc *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsSvc: analyticsSvc,
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
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}
	ctx := c.Request().Context()
	if err := h.analyticsSvc.CreateAnalytics(ctx, createAnalytics.PasteID, createAnalytics.URL); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "unable to create analytics"})
	}
	return c.JSON(http.StatusCreated, echo.Map{"message": "analytics created"})
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
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid limit"})
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid offset"})
	}
	ctx := c.Request().Context()
	analytics, err := h.analyticsSvc.GetAllAnalytics(ctx, order, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "unable to get all analytics"})
	}
	return c.JSON(http.StatusOK, analytics)
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
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "user id is required "})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf("invalid user id :%s", userIDStr)
	}
	order := c.QueryParam("order")
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid limit"})

	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid offset"})
	}
	ctx := c.Request().Context()
	analytics, err := h.analyticsSvc.GetAllAnalyticsByUser(ctx, userID, order, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "unable to get all analytics by user "})
	}
	return c.JSON(http.StatusOK, analytics)
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
	idStr := c.Param("ID")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "id is required"})
	}
	ID, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("invalid analytics ID %s", idStr)
	}
	ctx := c.Request().Context()
	analytic, err := h.analyticsSvc.GetAnalyticsByID(ctx, ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "unable to get analytics by ID"})
	}
	return c.JSON(http.StatusOK, analytic)
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
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "pasteID is required "})
	}
	pasteID, err := uuid.Parse(pasteIDStr)
	if err != nil {
		return fmt.Errorf("invalid pasteID %s", pasteIDStr)
	}
	ctx := c.Request().Context()
	analytic, err := h.analyticsSvc.GetAnalyticsByPasteID(ctx, pasteID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "unable to get analytics by pasteID "})
	}
	return c.JSON(http.StatusOK, analytic)
}

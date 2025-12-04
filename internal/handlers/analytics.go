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
	return c.JSON(http.StatusOK,analytics)
}

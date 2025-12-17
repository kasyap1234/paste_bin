package handlers

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Handlers struct {
	authHandler      *AuthHandler
	pasteHandler     *PasteHandler
	analyticsHandler *AnalyticsHandler
}

func NewHandlers(authHandler *AuthHandler, pasteHandler *PasteHandler, analyticsHandler *AnalyticsHandler) *Handlers {
	return &Handlers{
		authHandler:      authHandler,
		pasteHandler:     pasteHandler,
		analyticsHandler: analyticsHandler,
	}
}

func (h *Handlers) RegisterRoutes(e *echo.Echo) {
	e.POST("/register", h.authHandler.Register)
	e.POST("/login", h.authHandler.Login)
	e.POST("/paste", h.pasteHandler.CreatePaste)
	e.PUT("/paste/:id", h.pasteHandler.UpdatePaste)
	e.GET("/paste/:id", h.pasteHandler.GetPasteByID)
	e.DELETE("/paste/:id", h.pasteHandler.DeletePasteByID)
	e.GET("/pastes", h.pasteHandler.GetAllPastes)
	e.GET("/analytics", h.analyticsHandler.GetAllAnalytics)
	e.GET("/analytics/user", h.analyticsHandler.GetAllAnalyticsByUser)
	e.GET("/analytics/paste", h.analyticsHandler.GetAnalyticsByPasteID)
	e.POST("/create-analytics", h.analyticsHandler.CreateAnalytics)
	e.GET("/analytics/:id", h.analyticsHandler.GetAnalyticsByID)

	// Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}

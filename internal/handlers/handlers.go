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

func (h *Handlers) RegisterRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	// Public routes (no authentication required)
	e.POST("/register", h.authHandler.Register)
	e.POST("/login", h.authHandler.Login)
	e.GET("/paste/:id", h.pasteHandler.GetPasteByID) // Allow public viewing by UUID
	e.GET("/p/:slug", h.pasteHandler.GetPublicPaste) // Public sharing by slug
	e.GET("/raw/:slug", h.pasteHandler.GetRawPaste)  // Raw content by slug

	// Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Handle favicon.ico
	e.GET("/favicon.ico", func(c echo.Context) error {
		return c.NoContent(204)
	})

	// Protected routes (require authentication)
	protected := e.Group("", authMiddleware)
	protected.POST("/paste", h.pasteHandler.CreatePaste)
	protected.PUT("/paste/:id", h.pasteHandler.UpdatePaste)
	protected.DELETE("/paste/:id", h.pasteHandler.DeletePasteByID)
	protected.GET("/pastes", h.pasteHandler.GetAllPastes)
	protected.GET("/paste/filter", h.pasteHandler.FilterPastes)
	protected.GET("/analytics", h.analyticsHandler.GetAllAnalytics)
	protected.GET("/analytics/user", h.analyticsHandler.GetAllAnalyticsByUser)
	protected.GET("/analytics/paste", h.analyticsHandler.GetAnalyticsByPasteID)
	protected.POST("/create-analytics", h.analyticsHandler.CreateAnalytics)
	protected.GET("/analytics/:id", h.analyticsHandler.GetAnalyticsByID)
}

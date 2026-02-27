package handlers

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type HealthHandler struct {
	logger zerolog.Logger
	db     *pgxpool.Pool
}

func NewHealthHandler(logger zerolog.Logger, db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{
		logger: logger,
		db:     db,
	}
}

func (h *HealthHandler) Health(c echo.Context) error {
	// check if db is connected and backend is running
	if err := h.db.Ping(context.Background()); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database connection failed"})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})

}

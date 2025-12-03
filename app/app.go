package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"

	"pastebin/internal/auth"
	"pastebin/internal/config"
	"pastebin/internal/database"
	"pastebin/internal/handlers"
	"pastebin/internal/repositories"
	"pastebin/internal/services"
)

type App struct {
	server   *echo.Echo
	logger   zerolog.Logger
	addr     string
	db       *pgxpool.Pool
	handlers *handlers.Handlers
}

// New initializes the entire application graph: logger, database connections,
// repositories, services, handlers, and the Echo server. Everything outside of
// wiring lives in their respective packages.
func New() (*App, error) {
	logger := initLogger()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET env var is required")
	}
	jwtMgr := auth.NewJWTManager(jwtSecret)

	db, err := database.InitDB(&config.DBConfig{})
	if err != nil {
		logger.Error().Err(err).Msg("failed to initialize database")
		return nil, fmt.Errorf("init db: %w", err)
	}

	authRepo := repositories.NewAuthRepository(db)
	userRepo := repositories.NewUserRepository(db)
	pasteRepo := repositories.NewPasteRepository(db)

	authSvc := services.NewAuthService(authRepo, userRepo, jwtMgr)
	pasteSvc := services.NewPasteService(pasteRepo)

	authHandler := handlers.NewAuthHandler(authSvc)
	pasteHandler := handlers.NewPasteHandler(pasteSvc)
	handlerSet := handlers.NewHandlers(authHandler, pasteHandler)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	handlerSet.RegisterRoutes(e)

	addr := resolveAddr()

	return &App{
		server:   e,
		logger:   logger,
		addr:     addr,
		db:       db,
		handlers: handlerSet,
	}, nil
}

// Run starts the HTTP server and blocks until it exits.
func (a *App) Run() error {
	defer a.db.Close()
	a.logger.Info().Str("addr", a.addr).Msg("starting pastebin api")
	return a.server.Start(a.addr)
}

func resolveAddr() string {
	addr := os.Getenv("PORT")
	if addr == "" {
		return ":8080"
	}
	if !strings.HasPrefix(addr, ":") {
		return ":" + addr
	}
	return addr
}

func initLogger() zerolog.Logger {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "production"
	}

	if env == "development" {
		return zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
			Level(zerolog.DebugLevel).
			With().
			Timestamp().
			Str("env", env).
			Logger()
	}

	return zerolog.New(os.Stdout).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Str("env", env).
		Logger()
}

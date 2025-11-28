package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"pastebin/internal/auth"
	"pastebin/internal/config"
	"pastebin/internal/database"
)

var Log zerolog.Logger // Value, not pointer

func initLogger() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "production"
	}

	if env == "development" {
		Log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
			Level(zerolog.DebugLevel).
			With().
			Timestamp().
			Str("env", env).
			Logger()
	} else {
		Log = zerolog.New(os.Stdout).
			Level(zerolog.NoLevel).
			With().
			Timestamp().
			Str("env", env).
			Logger()
	}
}

func main() {
	initLogger()

	Log.Info().Msg("starting pastebin api")
	jwtMgr :=auth.NewJWTManager(os.Getenv("JWT_SECRET"))
	e :=echo.New()
	
	api :=e.Group("/api")
	api.Use(auth.AuthMiddleware(jwtMgr))
	

	dbConfig := &config.DBConfig{}
	pool, err := database.InitDB(dbConfig)
	if err != nil {
		Log.Fatal().Err(err).Msg("failed to initialize database")
	}
	defer pool.Close() // No ctx arg

	Log.Info().Msg("database initialized successfully")

	// Keep running
	select {}
}

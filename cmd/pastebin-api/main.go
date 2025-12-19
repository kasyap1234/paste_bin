// Package main Pastebin API
//
//	@title			Pastebin API
//	@version		1.0
//	@description	A simple pastebin service with analytics
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
//	@host		localhost:8080
//	@BasePath	/
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
package main

import (
	"log"

	"pastebin/app"
	_ "pastebin/docs" // This is required for swagger

	"github.com/joho/godotenv"
)

func main() {
	// Try to load .env file from multiple possible locations
	envPaths := []string{
		".env",       // Current directory
		"../../.env", // From cmd/pastebin-api/ relative to project root
		"../.env",    // From bin/ directory
		"./.env",     // Explicit current directory
	}

	envLoaded := false
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			log.Printf("Loaded .env file from: %s", path)
			envLoaded = true
			break
		}
	}

	if !envLoaded {
		log.Printf("Warning: .env file not found in any of the expected locations, using system env vars")
	}

	application, appErr := app.New()
	if appErr != nil {
		log.Fatalf("failed to initialize app: %v", appErr)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}

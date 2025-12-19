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
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Printf("Warning: .env file not found, using system env vars: %v", err)
	}

	application, appErr := app.New()
	if appErr != nil {
		log.Fatalf("failed to initialize app: %v", appErr)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}

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
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}

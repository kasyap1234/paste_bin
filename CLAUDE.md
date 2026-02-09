# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Backend (Go)
- **Run Application**: `go run cmd/pastebin-api/main.go` (from `server/` directory)
- **Build**: `go build -o pastebin-api cmd/pastebin-api/main.go`
- **Run Tests**: `go test ./...`
- **Run Single Test**: `go test -v -run TestName ./path/to/package`
- **Lint**: `golangci-lint run` (if installed)
- **Generate Swagger Docs**: `swag init -g cmd/pastebin-api/main.go -o docs`

### Database Migrations (using Task & Goose)
- **Status**: `task migrate:status`
- **Up**: `task migrate:up`
- **Down**: `task migrate:down`
- **Create Migration**: `task migrate:create name=<migration_name>`

### Docker
- **Start Services**: `docker-compose up -d` (from `server/` directory)

## Architecture and Structure

The project is organized into a `server` (Go backend) and `client` (frontend, currently empty). The backend follows a layered architecture:

- `server/cmd/pastebin-api`: Entry point of the application.
- `server/app`: Application wiring and dependency injection.
- `server/internal/handlers`: HTTP request handlers (Echo framework).
- `server/internal/services`: Business logic layer.
- `server/internal/repositories`: Data access layer (PostgreSQL).
- `server/internal/models`: Data structures and database entities.
- `server/internal/auth`: JWT management and authentication middleware.
- `server/pkg`: Shared utilities (errors, password hashing, email, etc.).
- `server/db/migrations`: SQL migration files managed by Goose.

### Design Patterns
- **Dependency Injection**: Dependencies are manually wired in `server/app/app.go`.
- **Repository Pattern**: Data access is abstracted behind repository interfaces.
- **Service Layer**: Business logic is decoupled from transport and storage.
- **Fluent SQL**: `squirrel` is used for building SQL queries.
- **Structured Logging**: `zerolog` is used throughout the application.
- **Error Handling**: Custom errors are defined in `server/pkg/errors` using sentinel patterns.

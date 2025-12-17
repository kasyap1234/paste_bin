# Swagger Documentation

This project now includes Swagger documentation for the Pastebin API.

## Accessing Swagger UI

Once the application is running, you can access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

## API Documentation

The Swagger documentation includes all endpoints with:

- Request/response schemas
- Authentication requirements (Bearer token)
- Parameter descriptions
- Example requests

## Endpoints Covered

### Authentication
- `POST /register` - Register a new user
- `POST /login` - Login user

### Pastes
- `POST /paste` - Create a new paste
- `PUT /paste/{id}` - Update a paste
- `GET /paste/{id}` - Get paste by ID
- `DELETE /paste/{id}` - Delete paste by ID
- `GET /pastes` - Get all pastes for authenticated user

### Analytics
- `POST /create-analytics` - Create analytics entry
- `GET /analytics` - Get all analytics
- `GET /analytics/user` - Get analytics by user
- `GET /analytics/paste` - Get analytics by paste ID
- `GET /analytics/{id}` - Get analytics by ID

## Authentication

Most endpoints require authentication using Bearer tokens. After logging in via `POST /login`, include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Regenerating Documentation

If you modify the API endpoints or models, regenerate the Swagger documentation:

```bash
go run github.com/swaggo/swag/cmd/swag init -g cmd/pastebin-api/main.go --output docs
```
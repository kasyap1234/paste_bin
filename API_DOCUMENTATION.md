# Pastebin API Documentation

## Overview

The Pastebin API is a RESTful service that allows users to create, manage, and share text pastes with analytics tracking. The API supports user authentication via JWT tokens and provides endpoints for paste management, user authentication, and analytics.

## Base URL

- **Development**: `http://localhost:8080`
- **Production**: `https://api.example.com`

## Authentication

The API uses Bearer token authentication with JWT (JSON Web Tokens) for protected endpoints.

### How to Authenticate

1. Register a new user using the `/register` endpoint
2. Login using the `/login` endpoint to receive a JWT token
3. Include the token in the `Authorization` header for protected endpoints:

```
Authorization: Bearer <your_jwt_token>
```

## API Endpoints

### Authentication Endpoints

#### Register User
- **POST** `/register`
- **Description**: Register a new user account
- **Request Body**:
  ```json
  {
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }
  ```
- **Response** (201 Created):
  ```json
  {
    "message": "user registered"
  }
  ```
- **Validations**:
  - Name: minimum 2 characters
  - Email: valid email format
  - Password: minimum 6 characters

#### Login User
- **POST** `/login`
- **Description**: Authenticate user and receive JWT token
- **Request Body**:
  ```json
  {
    "email": "john@example.com",
    "password": "password123"
  }
  ```
- **Response** (200 OK):
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
  ```

### Paste Endpoints

#### Create Paste
- **POST** `/paste` (Requires Authentication)
- **Description**: Create a new paste with optional expiration
- **Query Parameters**:
  - `expires_in` (optional): Expiration duration (e.g., '24h', '7d', '1m')
- **Request Body**:
  ```json
  {
    "title": "My Code Snippet",
    "content": "function hello() { console.log(\"world\"); }",
    "language": "javascript",
    "password": "optional_password"
  }
  ```
- **Response** (201 Created):
  ```json
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "title": "My Code Snippet",
    "is_private": false,
    "content": "function hello() { console.log(\"world\"); }",
    "language": "javascript",
    "url": "http://localhost:8080/p/abc123def",
    "views": 0,
    "expires_at": null
  }
  ```

#### Get All Pastes
- **GET** `/pastes` (Requires Authentication)
- **Description**: Retrieve all pastes for the authenticated user
- **Response** (200 OK):
  ```json
  [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "user_id": "550e8400-e29b-41d4-a716-446655440001",
      "title": "My Code Snippet",
      "is_private": false,
      "content": "function hello() { console.log(\"world\"); }",
      "language": "javascript",
      "url": "http://localhost:8080/p/abc123def",
      "views": 5,
      "expires_at": null
    }
  ]
  ```

#### Get Paste by ID
- **GET** `/paste/{id}`
- **Description**: Retrieve a specific paste by its ID
- **Path Parameters**:
  - `id` (required): Paste ID (UUID)
- **Response** (200 OK): Returns paste object
- **Error Responses**:
  - 400: Invalid paste ID
  - 500: Unable to get paste

#### Update Paste
- **PUT** `/paste/{id}` (Requires Authentication)
- **Description**: Update an existing paste
- **Path Parameters**:
  - `id` (required): Paste ID (UUID)
- **Request Body** (all fields optional):
  ```json
  {
    "title": "Updated Title",
    "content": "Updated content",
    "language": "python",
    "is_private": true,
    "password": "new_password",
    "expires_at": "2024-12-31T23:59:59Z"
  }
  ```
- **Response** (200 OK):
  ```json
  {
    "message": "paste updated"
  }
  ```

#### Delete Paste
- **DELETE** `/paste/{id}` (Requires Authentication)
- **Description**: Delete a specific paste by its ID
- **Path Parameters**:
  - `id` (required): Paste ID (UUID)
- **Response** (200 OK):
  ```json
  {
    "error": "deleted"
  }
  ```

#### Get Public Paste by Slug
- **GET** `/p/{slug}`
- **Description**: Retrieve a public paste by its URL slug
- **Path Parameters**:
  - `slug` (required): Paste slug
- **Response** (200 OK): Returns paste object
- **Error Responses**:
  - 400: Invalid slug
  - 404: Paste not found
  - 500: Unable to get paste

#### Get Raw Paste Content
- **GET** `/raw/{slug}`
- **Description**: Retrieve raw text content of a public paste
- **Path Parameters**:
  - `slug` (required): Paste slug
- **Response** (200 OK): Raw text content with Content-Type: text/plain
- **Error Responses**:
  - 400: Invalid slug
  - 404: Paste not found
  - 500: Unable to get paste

### Analytics Endpoints

#### Create Analytics Entry
- **POST** `/create-analytics` (Requires Authentication)
- **Description**: Create a new analytics entry for a paste
- **Request Body**:
  ```json
  {
    "paste_id": "550e8400-e29b-41d4-a716-446655440000",
    "url": "http://localhost:8080/p/abc123def"
  }
  ```
- **Response** (201 Created):
  ```json
  {
    "message": "analytics created"
  }
  ```

#### Get All Analytics
- **GET** `/analytics` (Requires Authentication)
- **Description**: Retrieve all analytics with pagination
- **Query Parameters**:
  - `order` (optional): Sort order
  - `limit` (optional): Limit number of results
  - `offset` (optional): Offset for pagination
- **Response** (200 OK):
  ```json
  [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "paste_id": "550e8400-e29b-41d4-a716-446655440001",
      "url": "http://localhost:8080/p/abc123def",
      "views": 10
    }
  ]
  ```

#### Get Analytics by User
- **GET** `/analytics/user` (Requires Authentication)
- **Description**: Retrieve analytics for a specific user with pagination
- **Query Parameters**:
  - `userID` (required): User ID (UUID)
  - `order` (optional): Sort order
  - `limit` (optional): Limit number of results
  - `offset` (optional): Offset for pagination
- **Response** (200 OK): Array of analytics objects

#### Get Analytics by Paste ID
- **GET** `/analytics/paste` (Requires Authentication)
- **Description**: Retrieve analytics for a specific paste
- **Query Parameters**:
  - `pasteID` (required): Paste ID (UUID)
- **Response** (200 OK):
  ```json
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "paste_id": "550e8400-e29b-41d4-a716-446655440001",
    "url": "http://localhost:8080/p/abc123def",
    "views": 10
  }
  ```

#### Get Analytics by ID
- **GET** `/analytics/{id}` (Requires Authentication)
- **Description**: Retrieve specific analytics entry by its ID
- **Path Parameters**:
  - `id` (required): Analytics ID (UUID)
- **Response** (200 OK): Analytics object
- **Error Responses**:
  - 400: Invalid ID
  - 500: Unable to get analytics

## Error Handling

The API returns appropriate HTTP status codes and error messages:

- **400 Bad Request**: Invalid input or request parameters
- **401 Unauthorized**: Missing or invalid authentication token
- **404 Not Found**: Resource not found
- **500 Internal Server Error**: Server error

Error Response Format:
```json
{
  "error": "Description of the error"
}
```

## Using the OpenAPI Specification

The API includes a complete OpenAPI 3.0 specification in `openapi.yaml`. You can:

1. **Import into Postman**: Import the OpenAPI file to test endpoints
2. **Generate Client Libraries**: Use OpenAPI generators to create client code
3. **View Interactive Documentation**: Use tools like Swagger UI or ReDoc
4. **API Validation**: Validate your API implementation against the spec

## Common Use Cases

### 1. Register and Login
```bash
# Register
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 2. Create and Share a Paste
```bash
# Create paste (replace TOKEN with your JWT)
curl -X POST http://localhost:8080/paste \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Script",
    "content": "echo Hello World",
    "language": "bash"
  }'

# Get public paste by slug
curl http://localhost:8080/p/abc123def
```

### 3. Create Analytics
```bash
curl -X POST http://localhost:8080/create-analytics \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "paste_id": "550e8400-e29b-41d4-a716-446655440000",
    "url": "http://localhost:8080/p/abc123def"
  }'
```

## Rate Limiting

Currently, the API does not implement rate limiting. Future versions may include this feature.

## Versioning

Current API Version: 1.0.0

The API may introduce breaking changes in future major versions. It is recommended to implement version-aware client code.

## Support

For issues, bugs, or feature requests, please contact the API support team:
- Email: support@swagger.io
- URL: http://www.swagger.io/support
# Chirpy - Go HTTP Server

A RESTful HTTP server for a Twitter-like social media platform built with Go, PostgreSQL, and JWT authentication.

## Overview

Chirpy is a backend API server that allows users to create accounts, post short messages (chirps), and manage their content. The application features user authentication with JWT tokens, refresh token management, and premium user upgrades via webhook integration.

## Features

- **User Management**: Create accounts, login, and update credentials
- **Authentication**: JWT-based authentication with refresh tokens
- **Chirps**: Create, read, and delete short messages (max 140 characters)
- **Profanity Filter**: Automatic filtering of inappropriate content
- **Premium Upgrades**: Webhook integration for Chirpy Red premium accounts
- **Metrics**: Request tracking and admin dashboard

## Technologies

- **Go 1.25.1**: Core programming language
- **PostgreSQL**: Database for persistent storage
- **sqlc**: Type-safe SQL code generation
- **JWT (golang-jwt/jwt)**: Token-based authentication
- **Argon2id**: Secure password hashing
- **godotenv**: Environment variable management

## Installation

### Prerequisites

- Go 1.25.1 or higher
- PostgreSQL database
- Git

### Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/ManoloEsS/go_http_server.git
   cd go_http_server
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Set up PostgreSQL database**:
   ```bash
   createdb chirpy
   ```

4. **Run database migrations**:
   Execute the SQL files in `sql/schema/` in order:
   ```bash
   psql -d chirpy -f sql/schema/001_users.sql
   psql -d chirpy -f sql/schema/002_chirps.sql
   psql -d chirpy -f sql/schema/003_passwords.sql
   psql -d chirpy -f sql/schema/004_refresh_tokens.sql
   psql -d chirpy -f sql/schema/005_is_chirpy_red.sql
   ```

5. **Create `.env` file**:
   ```env
   DB_URL="postgres://username:password@localhost:5432/chirpy?sslmode=disable"
   PLATFORM="dev"
   SECRET="your-jwt-secret-key-here"
   POLKA_KEY="your-webhook-api-key-here"
   ```

6. **Build and run**:
   ```bash
   go build -o out ./cmd/go_http_server
   ./out
   ```

The server will start on `http://localhost:8080`.

## API Endpoints

### Health Check

#### `GET /api/healthz`
Check server health status.

**Response**: `200 OK`

---

### User Management

#### `POST /api/users`
Create a new user account.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response**: `201 Created`
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

#### `PUT /api/users`
Update user email and/or password (requires authentication).

**Headers**: `Authorization: Bearer <jwt_token>`

**Request Body**:
```json
{
  "email": "newemail@example.com",
  "password": "newpassword"
}
```

**Response**: `200 OK`
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "newemail@example.com",
  "is_chirpy_red": false
}
```

---

### Authentication

#### `POST /api/login`
Authenticate user and receive JWT and refresh tokens.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securepassword",
  "expires_in_seconds": 3600
}
```

**Response**: `200 OK`
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "token": "jwt_access_token",
  "refresh_token": "refresh_token_string"
}
```

#### `POST /api/refresh`
Exchange refresh token for a new JWT access token.

**Headers**: `Authorization: Bearer <refresh_token>`

**Response**: `200 OK`
```json
{
  "token": "new_jwt_access_token"
}
```

#### `POST /api/revoke`
Revoke a refresh token (logout).

**Headers**: `Authorization: Bearer <refresh_token>`

**Response**: `204 No Content`

---

### Chirps

#### `POST /api/chirps`
Create a new chirp (requires authentication).

**Headers**: `Authorization: Bearer <jwt_token>`

**Request Body**:
```json
{
  "body": "This is my first chirp!"
}
```

**Response**: `201 Created`
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "body": "This is my first chirp!",
  "user_id": "uuid"
}
```

#### `GET /api/chirps`
Retrieve all chirps with optional filtering and sorting.

**Query Parameters**:
- `author_id` (optional): Filter chirps by user UUID
- `sort` (optional): Sort order - `asc` (default) or `desc`

**Response**: `200 OK`
```json
[
  {
    "id": "uuid",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "body": "This is a chirp!",
    "user_id": "uuid"
  }
]
```

#### `GET /api/chirps/{chirpID}`
Retrieve a specific chirp by ID.

**Response**: `200 OK`
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "body": "This is a chirp!",
  "user_id": "uuid"
}
```

#### `DELETE /api/chirps/{chirpID}`
Delete a chirp (requires authentication and ownership).

**Headers**: `Authorization: Bearer <jwt_token>`

**Response**: `204 No Content`

---

### Premium Features

#### `POST /api/polka/webhooks`
Webhook endpoint for upgrading users to Chirpy Red (requires API key).

**Headers**: `Authorization: ApiKey <polka_api_key>`

**Request Body**:
```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "uuid"
  }
}
```

**Response**: `204 No Content`

---

### Admin Endpoints

#### `GET /admin/metrics`
View server metrics (file server hit count).

**Response**: `200 OK` (HTML page)

#### `POST /admin/reset`
Reset all users and metrics (development environment only).

**Response**: `200 OK`

---

### Static Files

#### `GET /app/*`
Serves static files from the root directory.

## Project Structure

```
go_http_server/
├── cmd/
│   └── go_http_server/
│       └── main.go              # Application entry point
├── internal/
│   ├── auth/                    # Authentication utilities
│   │   ├── get_bearer_token.go  # JWT token extraction
│   │   ├── get_polka_api_key.go # API key extraction
│   │   ├── jwt.go               # JWT creation and validation
│   │   ├── make_refresh_token.go # Refresh token generation
│   │   └── password_hashing.go  # Argon2id password hashing
│   ├── config/
│   │   └── constants.go         # Application constants
│   └── database/                # Database layer (generated by sqlc)
│       ├── chirps.sql.go        # Chirp queries
│       ├── db.go                # Database interface
│       ├── models.go            # Data models
│       ├── refresh_tokens.sql.go # Token queries
│       └── users.sql.go         # User queries
├── server/
│   ├── handlers/                # HTTP request handlers
│   │   ├── api_config.go        # Shared configuration
│   │   ├── handler_*.go         # Individual endpoint handlers
│   │   ├── middleware.go        # Metrics middleware
│   │   └── structs.go           # Response structures
│   └── json_responses.go        # JSON response utilities
├── sql/
│   ├── queries/                 # SQL queries for sqlc
│   └── schema/                  # Database migrations
├── .env                         # Environment variables
├── go.mod                       # Go module dependencies
├── sqlc.yaml                    # sqlc configuration
└── README.md                    # This file
```

## Authentication Flow

1. **User Registration**: User creates account via `POST /api/users`
2. **Login**: User authenticates via `POST /api/login`, receives JWT (1 hour) and refresh token (60 days)
3. **API Requests**: User includes JWT in `Authorization: Bearer <token>` header
4. **Token Refresh**: When JWT expires, use refresh token via `POST /api/refresh` to get new JWT
5. **Logout**: Revoke refresh token via `POST /api/revoke`

## Security Features

- **Argon2id Password Hashing**: Industry-standard secure password storage
- **JWT Authentication**: Stateless token-based authentication
- **Refresh Tokens**: Long-lived tokens stored in database with revocation support
- **Profanity Filtering**: Automatic content moderation
- **Authorization Checks**: Ownership validation for resource deletion

## Development

### Running Tests
```bash
go test ./...
```

### Generate Database Code (after SQL changes)
```bash
sqlc generate
```

### Environment Modes
- `PLATFORM=dev`: Enables admin reset endpoint
- `PLATFORM=prod`: Disables destructive admin operations

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is part of a learning exercise and is provided as-is for educational purposes.

## Authors

- ManoloEsS

## Acknowledgments

- Built as part of the Boot.dev backend development course
- Uses sqlc for type-safe SQL queries
- JWT authentication with golang-jwt
- Password hashing with Argon2id

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

### Health & Metrics

#### Check Server Health
```
GET /api/healthz
```
Returns server health status.

**Response**: `200 OK`

#### Get Metrics
```
GET /admin/metrics
```
Returns server hit metrics.

**Response**: Metrics page with request count

#### Reset Database
```
POST /admin/reset
```
Resets all users in the database (admin only).

**Response**: `200 OK`

---

### User Management

#### Create User
```
POST /api/users
```

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
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

#### User Login
```
POST /api/login
```

Authenticate a user and receive access and refresh tokens.

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
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "token": "jwt-access-token",
  "refresh_token": "refresh-token"
}
```

**Notes**:
- `expires_in_seconds` is optional (max 1 hour)
- Default expiration is 1 hour if not specified

#### Update User
```
PUT /api/users
```

Update user email and/or password.

**Headers**:
```
Authorization: Bearer <jwt-token>
```

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
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "email": "newemail@example.com",
  "is_chirpy_red": false
}
```

### Token Management

#### Refresh Access Token
```
POST /api/refresh
```

Generate a new access token using a refresh token.

**Headers**:
```
Authorization: Bearer <refresh-token>
```

**Response**: `200 OK`
```json
{
  "token": "new-jwt-access-token"
}
```

#### Revoke Refresh Token
```
POST /api/revoke
```

Revoke a refresh token (logout).

**Headers**:
```
Authorization: Bearer <refresh-token>
```

**Response**: `200 OK`

---

### Chirps (Messages)

#### Validate Chirp
```
POST /api/validate_chirp
```

Validate chirp content without saving it.

**Request Body**:
```json
{
  "body": "This is a test chirp message"
}
```

**Response**: `200 OK`

#### Create Chirp
```
POST /api/chirps
```

Create a new chirp (authenticated users only).

**Headers**:
```
Authorization: Bearer <jwt-token>
```

**Request Body**:
```json
{
  "body": "This is my chirp message",
  "user_id": "user-uuid"
}
```

**Response**: `201 Created`
```json
{
  "id": "chirp-uuid",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "body": "This is my chirp message",
  "user_id": "user-uuid"
}
```

**Notes**:
- Maximum chirp length: 140 characters
- Profane words are automatically replaced with `****`
- Filtered words: kerfuffle, sharbert, fornax

#### Get All Chirps
```
GET /api/chirps
```

Retrieve all chirps from the database.

**Response**: `200 OK`
```json
[
  {
    "id": "chirp-uuid",
    "created_at": "timestamp",
    "updated_at": "timestamp",
    "body": "Chirp message",
    "user_id": "user-uuid"
  }
]
```

#### Get Chirp by ID
```
GET /api/chirps/{chirpID}
```

Retrieve a specific chirp by its ID.

**Parameters**:
- `chirpID`: UUID of the chirp

**Response**: `200 OK`
```json
{
  "id": "chirp-uuid",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "body": "Chirp message",
  "user_id": "user-uuid"
}
```

#### Delete Chirp
```
DELETE /api/chirps/{chirpID}
```

Delete a chirp (only the chirp's author can delete it).

**Headers**:
```
Authorization: Bearer <jwt-token>
```

**Parameters**:
- `chirpID`: UUID of the chirp

**Response**: `204 No Content`

**Error**: `403 Forbidden` if user is not the chirp's author

---

### Webhooks

#### Polka Upgrade Webhook
```
POST /api/polka/webhooks
```

Webhook for upgrading users to Chirpy Red premium membership.

**Headers**:
```
Authorization: ApiKey <polka-api-key>
```

**Request Body**:
```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "user-uuid"
  }
}
```

**Response**: `204 No Content`

**Notes**:
- Requires valid Polka API key
- Only processes `user.upgraded` events
- Returns `204` for other event types

---

### Static Files

#### Application Files
```
GET /app/*
```

Serves static files from the root directory.

**Example**: 
- `/app/` → serves `index.html`
- `/app/assets/style.css` → serves static assets

---

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

## License

This project is part of a learning exercise and is provided as-is for educational purposes.

## Authors

- ManoloEsS

## Acknowledgments

- Built as part of the Boot.dev backend development course
- Uses sqlc for type-safe SQL queries
- JWT authentication with golang-jwt
- Password hashing with Argon2id

# Chirpy - RESTful API Server

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A complete RESTful API implementing user authentication, content management, and third-party integrations. Features JWT-based auth with refresh token rotation, Argon2id password hashing, and webhook integration patterns—demonstrating production-grade backend architecture.

## Features

| Feature | Details |
|---------|---------|
| **JWT Authentication** | JWT with HS256 signing, 1-hour expiration, and refresh token rotation for secure stateless auth |
| **Password Security** | Argon2id password hashing—OWASP-recommended, memory-hard, resistant to GPU attacks |
| **Type-Safe Database** | sqlc generates type-safe database queries at compile time, eliminating runtime SQL errors |
| **Refresh Tokens** | 64-char cryptographic refresh tokens stored in PostgreSQL with revocation support |
| **Webhook Integration** | External webhook integration with API key verification for third-party services |
| **Content Moderation** | Configurable profanity filter with word replacement |
| **Request Metrics** | Atomic request counter with admin metrics endpoint |

## Prerequisites

- [Go](https://go.dev/) 1.21 or higher
- [PostgreSQL](https://www.postgresql.org/) 12 or higher
- [goose](https://github.com/pressly/goose) for database migrations

## Installation

### Clone and Build

```bash
git clone https://github.com/ManoloEsS/go_http_server.git
cd go_http_server
go build -o chirpy ./cmd/go_http_server
```

### Database Setup

```bash
# Create the database
createdb chirpy

# Run migrations
goose -dir sql/schema postgres "postgres://user:pass@localhost/chirpy?sslmode=disable" up
```

### Environment Configuration

Create a `.env` file in the project root:

```bash
cat > .env << 'EOF'
DB_URL="postgres://user:pass@localhost/chirpy?sslmode=disable"
PLATFORM="dev"
SECRET="your-32-char-secret-key-here"
POLKA_KEY="your-polka-api-key"
EOF
```

### Run the Server

```bash
./chirpy
```

The server starts at `http://localhost:8080`

## Usage

### Register a New User

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"securepass123"}'
```

**Response:**
```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "created_at": "2024-01-15T10:30:00Z",
  "email": "alice@example.com",
  "is_chirpy_red": false
}
```

### Login

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"securepass123"}'
```

**Response:**
```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "email": "alice@example.com",
  "is_chirpy_red": false,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "abc123def456..."
}
```

### Create a Chirp (Authenticated)

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{"body":"Hello, world!","user_id":"a1b2c3d4-e5f6-7890-abcd-ef1234567890"}'
```

**Response:**
```json
{
  "id": "chirp-uuid",
  "created_at": "2024-01-15T10:31:00Z",
  "body": "Hello, world!",
  "user_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

### Get All Chirps

```bash
curl http://localhost:8080/api/chirps
```

### Delete Own Chirp

```bash
curl -X DELETE http://localhost:8080/api/chirps/{chirp_id} \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Refresh Token

```bash
curl -X POST http://localhost:8080/api/refresh \
  -H "Authorization: Bearer <refresh_token>"
```

### Revoke Token (Logout)

```bash
curl -X POST http://localhost:8080/api/revoke \
  -H "Authorization: Bearer <refresh_token>"
```

## Configuration

| Variable | Description |
|----------|-------------|
| `DB_URL` | PostgreSQL connection string |
| `PLATFORM` | `dev` or `prod` (enables/disables admin endpoints) |
| `SECRET` | JWT signing key (minimum 32 characters) |
| `POLKA_KEY` | API key for webhook verification |

## API Reference

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/healthz` | No | Health check |
| `POST` | `/api/users` | No | User registration |
| `POST` | `/api/login` | No | Login (returns JWT + refresh token) |
| `PUT` | `/api/users` | JWT | Update email/password |
| `POST` | `/api/chirps` | JWT | Create chirp (140 char max) |
| `GET` | `/api/chirps` | No | List all chirps |
| `GET` | `/api/chirps/{id}` | No | Get single chirp |
| `DELETE` | `/api/chirps/{id}` | JWT | Delete own chirp |
| `POST` | `/api/refresh` | Refresh | Rotate JWT |
| `POST` | `/api/revoke` | Refresh | Logout |
| `POST` | `/api/polka/webhooks` | API Key | Premium upgrade webhook |


## Architecture

```
Request → JWT Middleware → Handler → sqlc Queries → PostgreSQL
                ↓
         Password validation
         Token rotation
         Ownership checks
```

## Running Tests

```bash
go test ./...
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you'd like to change.

## License

Distributed under the MIT License. See `LICENSE` for more information.

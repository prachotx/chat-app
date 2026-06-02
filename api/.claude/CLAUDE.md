# API — Real-Time Chat

Go REST API built with Fiber v3, GORM, and PostgreSQL.

## Tech Stack

- **Framework:** [Fiber v3](https://github.com/gofiber/fiber)
- **ORM:** GORM with PostgreSQL driver
- **Auth:** JWT (HS256) stored in HttpOnly cookie (`access_token`, 24h TTL)
- **Validation:** go-playground/validator v10 via Fiber's `StructValidator`
- **Hot reload:** Air (`.air.toml`)

## Project Layout

```
api/
├── cmd/server/main.go        # Entrypoint — wires dependencies, registers routes
├── config/config.go          # Loads env vars into Config struct (singleton)
├── internal/
│   ├── dto/                  # Request/response shapes (no sensitive fields in *Response types)
│   ├── handler/              # HTTP layer — binds input, calls service, returns response
│   ├── middleware/           # Fiber middleware (auth JWT validation)
│   ├── model/                # GORM models (GormModel base: ID, CreatedAt, UpdatedAt)
│   ├── repository/           # DB queries — interface + implementation
│   └── service/              # Business logic — interface + implementation
└── pkg/
    ├── crypto/               # bcrypt helpers (HashPassword, CheckPasswordHash)
    ├── jwt/                  # GenerateToken / ValidateToken
    └── response/             # response.Send(c, status, message, data) → BaseResponse JSON
```

## Architecture Pattern

All features follow the same layered pattern. Each layer depends only on the layer below via an interface:

```
Handler → Service (interface) → Repository (interface) → GORM
```

When adding a new feature:
1. Add model in `internal/model/`
2. Add repository interface + implementation in `internal/repository/`
3. Add service interface + implementation in `internal/service/`
4. Add DTOs in `internal/dto/` — **never expose sensitive fields** in `*Response` types
5. Add handler in `internal/handler/`
6. Register routes in `cmd/server/main.go`

## Response Format

All responses use `response.Send()` which returns:

```json
{ "message": "...", "data": ... }
```

`data` is omitted when `nil`.

## Auth Flow

- `POST /api/auth/register` — hash password with bcrypt, store user
- `POST /api/auth/login` — verify password, set `access_token` cookie (HttpOnly, Secure, SameSite=Lax)
- `GET /api/auth/me` — protected by `middleware.AuthMiddleware`; reads `user_id` from JWT claims and stores it in `c.Locals("user_id")`

## Environment Variables

| Variable      | Description               |
|---------------|---------------------------|
| `PORT`        | Server port               |
| `DB_HOST`     | Postgres host             |
| `DB_PORT`     | Postgres port             |
| `DB_USER`     | Postgres user             |
| `DB_PASSWORD` | Postgres password         |
| `DB_NAME`     | Postgres database name    |
| `SECRET_KEY`  | JWT signing secret (HMAC) |

## Running Locally

```bash
# from /api
air          # hot reload via .air.toml
# or
go run ./cmd/server/main.go
```

Database and pgAdmin are managed via Docker Compose at the repo root (`docker-compose.dev.yml`).

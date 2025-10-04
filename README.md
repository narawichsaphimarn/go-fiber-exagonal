# Go Fiber + Hexagonal Architecture (Latest)

This project follows Hexagonal Architecture with Go (Fiber). It separates core business logic (domain + usecase + ports) from technologies (HTTP, DB, auth), uses `pgx v5` for PostgreSQL, JWT for auth, and validator for request validation.

## Overview
- Core (domain/usecase/ports): business logic and contracts (interfaces).
- Adapters (http/repo/auth): Fiber handlers/middleware, PostgreSQL repo via `pgx`, JWT provider.
- Composition root: `cmd/api/main.go` wires everything, starts the server.

## Folder Structure (Current)
```plaintext
├── cmd/
│   └── api/
│       └── main.go                 # Fiber app, pgx pool, route wiring
├── configs/
│   └── app.yaml                    # App/DB/Auth configuration
├── internal/
│   ├── adapters/
│   │   ├── auth/                   # JWT provider
│   │   ├── http/
│   │   │   ├── handlers/           # Book/User handlers
│   │   │   └── middleware/         # Auth middleware (Bearer token)
│   │   └── repo/                   # pgx v5 repositories (book, user)
│   ├── domain/
│   │   ├── auth.go                 # Auth DTOs (e.g., login request)
│   │   ├── book.go                 # Book entity
│   │   └── user.go                 # User entity + ComparePassword
│   ├── ports/
│   │   ├── book_repository.go      # Book repo interface
│   │   ├── token_provider.go       # Token provider interface
│   │   └── user_repository.go      # User repo interface
│   └── usecase/
│       ├── book/                   # BookService
│       └── user/                   # UserService (register/login/etc.)
├── migrations/
│   ├── 001_init_books.sql
│   └── 002_inti_user.sql           # Users table, constraints, seeds
├── pkg/
│   ├── config.go                   # Load config from app.yaml/env
│   ├── hash.go                     # bcrypt hash helpers
│   └── validate_struct.go          # go-playground/validator v10 helper
├── client.http                     # API examples for testing
├── docker-compose.yml              # (Optional) Postgres service
├── go.mod
└── go.sum
```

## Configuration
- `configs/app.yaml` expected keys:
  - `app.port`: server port (e.g., 3000)
  - `db`: `host`, `port`, `user`, `password`, `dbName`, `options`
  - `auth.jwt.secret`: JWT signing secret
- Loaded by `pkg/config.go`. Server listens on `":" + cfg.App.Port`.

## Database (pgx v5)
- Connection via `pgxpool` in `cmd/api/main.go`.
- Repositories implement `ports.*Repository` and use context-aware queries.
- Migrations:
  - `001_init_books.sql`: books schema.
  - `002_inti_user.sql`: users schema, constraints, and potential seed.
    - Note: `ON CONFLICT (username) DO NOTHING;` will ignore duplicate inserts.
      - For strict uniqueness, handle duplicates in repo/usecase with proper errors.
- Recommended environment:
  ```bash
  export DATABASE_URL="postgres://user:pass@localhost:5432/dbname?sslmode=disable"
  ```

## Auth
- JWT provider: `internal/adapters/auth/jwt/provider.go`
  - Methods: `GenerateToken(userId)`, `ValidateToken(tokenString)`
- Middleware: `internal/adapters/http/middleware/auth.go`
  - Checks `Authorization: Bearer <token>`; rejects with `401` if missing/invalid.
- Route groups:
  - Public: `/v1` (no middleware) → register/login
  - Protected: `/v1` with middleware → users/books CRUD
  - Tip: Prefer nested groups to avoid prefix conflicts:
    - `v1 := app.Group("/v1")`
    - `v1Protected := v1.Group("", ProtectMiddleware)`

## Validation & Passwords
- Validation: `github.com/go-playground/validator/v10`
  - Tags on `domain.User` (e.g., `validate:"required,email,min=8"`).
  - Used via `pkg.ValidateStruct(ctx, obj)`.
- Password hashing: `golang.org/x/crypto/bcrypt`
  - `pkg.HashPassword` hashes on register/update.
  - `domain.User.ComparePassword` verifies on login.

## Context
- Handlers create per-request `context.WithTimeout` (e.g., 5s).
- Usecase and repo receive `ctx` for cancellation and deadlines.
- Avoid `context.Background()` in repo/usecase to respect request lifecycle.

## API Endpoints
- Public (no token):
  - `POST /v1/register` → create user
  - `POST /v1/login` → returns `{ "token": "<jwt>" }`
- Protected (Bearer token required):
  - Users:
    - `GET /v1/user/:id`
    - `GET /v1/user/email/:email`
    - `GET /v1/users`
    - `PUT /v1/user/:id`
    - `PUT /v1/user/:id/password`
    - `DELETE /v1/user/:id`
  - Books:
    - `GET /v1/books`
    - `GET /v1/books/:id`
    - `POST /v1/books`
    - `PUT /v1/books/:id`
    - `DELETE /v1/books/:id`

## Usage
- Run server:
  ```bash
  go run ./cmd/api
  ```
- Test with `client.http` (VS Code/JetBrains HTTP client) or curl.
- Example login then call protected:
  1. `POST /v1/login` → get token
  2. Use header `Authorization: Bearer <token>` on protected routes

## Troubleshooting
- 401 on `POST /v1/register`:
  - Ensure it’s registered under the public group `/v1` (no middleware).
  - If both public and protected groups use the same prefix (`/v1`) from `app`, nesting is safer:
    - Register protected routes under `v1Protected := v1.Group("", Protect)`.
  - Confirm client doesn’t send a stale `Authorization` header inadvertently.
- CORS:
  - Current config allows `http://localhost:3000` with standard methods/headers. Adjust if needed.

## Transactions (Quick Guide)
- `database/sql`: use `BeginTx(ctx, opts)`; defer `Rollback()`; `Commit()` on success.
- `pgx v5`: use `pool.Begin(ctx)`; `tx.Exec/Query`; `tx.Commit(ctx)` or `tx.Rollback(ctx)`.
- GORM: use `db.Transaction(func(tx *gorm.DB) error { ... })` or manual `Begin/Commit/Rollback`.
- Prefer wrapping transactional operations with helpers and propagate `ctx`.

## Development Notes
- Keep domain/usecase pure (no Fiber, pgx, JWT).
- Ports define contracts; adapters implement them.
- Lifecycle (DB pool, providers) managed in `cmd/api/main.go`.

---
This README reflects the latest structure, routing model (public vs protected), configuration, validation, auth, DB integration, and common pitfalls such as 401 on public endpoints with middleware.
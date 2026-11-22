# Go Clean Architecture Boilerplate

A robust, production-ready REST API boilerplate with Go, Echo, and GORM following Clean Architecture principles.

---

## Architecture Overview

This project strictly follows **Clean Architecture**, dividing responsibilities into distinct layers:

```
┌────────────────────────────────────────────────────────┐
│                        app/                            │
│         (Entry Point, DI & Graceful Shutdown)          │
└───────────────────────────┬────────────────────────────┘
                            │
┌───────────────────────────▼────────────────────────────┐
│                    domain/user/                        │
│  ┌──────────────────────────────────────────────────┐  │
│  │  delivery/http/  (Echo Handlers & DTO Binding)   │  │
│  ├──────────────────────────────────────────────────┤  │
│  │  usecase/        (Business Logic & Validation)   │  │
│  ├──────────────────────────────────────────────────┤  │
│  │  repository/     (GORM Database Operations)      │  │
│  ├──────────────────────────────────────────────────┤  │
│  │  interfaces/     (Contract Port Definitions)     │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────┘
```

- **`app/`**: Application bootstrap, dependency injection, middleware setup, and graceful shutdown handling.
- **`domain/`**: Core business domains containing interfaces, use cases, data repositories, and HTTP delivery handlers.
- **`model/`**: Shared entities, request DTOs (with struct validation), and response DTOs.
- **`shared/`**: Common utilities including configuration (`config`), database connections (`database`), migrations (`migration`), structured logging (`logger`), custom error definitions (`errors`), JSON response formatter (`json`), and JWT security middleware (`middleware`).
- **`db/migrations/`**: Versioned SQL migration files (`up` and `down`).

---

## Tech Stack

- **Go 1.25+**
- **Echo v4** — High-performance web framework
- **GORM v2** — Advanced ORM (`gorm.io/gorm` with PostgreSQL and SQLite support)
- **golang-migrate** — Robust database schema migrations
- **Uber Zap** — Blazing-fast structured JSON/development logging
- **go-playground/validator** — Request payload validation
- **golang-jwt/jwt/v5** — JWT authentication middleware
- **Testify** — Unit testing assertion and mocking framework

---

## Configuration

Create or update your environment configuration in `conf/config.env`:

```env
APP_NAME=boilerplate
APP_ENV=development
APP_PORT=8080

# Database Configuration (PostgreSQL / SQLite)
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=user_db
DB_SSLMODE=disable

# JWT Secret Key
JWT_SECRET=super-secret-key-change-me
```

*(Note: Set `DB_DRIVER=sqlite3` and `DB_NAME=app.db` to run without PostgreSQL).*

---

## Getting Started

### 1. Install Dependencies
```bash
go mod download
```

### 2. Run Database Migrations (Optional / Automatic)
Migrations will run automatically upon startup, or you can manage them manually using `golang-migrate`.

### 3. Run the Server
```bash
go run app/main.go
```
The server will start at `http://localhost:8080`.

---

## Testing

Run all unit tests with verbose output:
```bash
go test -v ./...
```

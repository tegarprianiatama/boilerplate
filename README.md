# Go Clean Architecture Boilerplate

A robust, production-ready REST API boilerplate with Go, Echo, and GORM following Clean Architecture principles.

---

## Architecture Overview

This project strictly follows **Clean Architecture**, dividing responsibilities into distinct layers:

- **`app/`**: Application bootstrap, dependency injection, middleware setup, and graceful shutdown handling.
- **`domain/`**: Core business domains containing interfaces, use cases, data repositories, and HTTP delivery handlers.
- **`model/`**: Shared entities, request DTOs (with struct validation), and response DTOs.
- **`shared/`**: Common utilities including configuration (`config`), database connections (`database`), migrations (`migration`), Redis client (`cache`), RabbitMQ client (`queue`), structured logging (`logger`), custom error definitions (`errors`), JSON response formatter (`json`), and JWT security middleware (`middleware`).
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
- **go-redis/v9** — Redis client and connection pooling
- **rabbitmq/amqp091-go** — RabbitMQ (AMQP 0-9-1) client
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
JWT_EXPIRE_HOURS=24

# Redis Configuration
REDIS_ENABLED=true
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# RabbitMQ Configuration
RABBITMQ_ENABLED=true
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_VHOST=/
```

*(Note: Set `DB_DRIVER=sqlite3` and `DB_NAME=app.db` to run without PostgreSQL. Set `REDIS_ENABLED=false` or `RABBITMQ_ENABLED=false` to skip those connections.)*

---

## Getting Started

### 1. Install Dependencies
```bash
go mod download
```

### 2. Start Infrastructure (Optional)
Spin up PostgreSQL, Redis, and RabbitMQ locally:
```bash
docker compose up -d
```
RabbitMQ's management UI is available at `http://localhost:15672` (default `guest` / `guest`).

### 3. Run Database Migrations (Optional / Automatic)
Migrations will run automatically upon startup, or you can manage them manually using `golang-migrate`.

### 4. Run the Server
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

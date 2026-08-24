# Go Clean Architecture Boilerplate

REST API boilerplate untuk User CRUD menggunakan Go, Echo, dan GORM dengan pola Clean Architecture.

## Tech Stack

- **Go 1.24**
- **Echo v4** — HTTP framework
- **GORM** — ORM (PostgreSQL / SQLite)
- **validator** — validasi request
- **godotenv** — load config env

## Struktur Proyek

```
app/                  → entry point & dependency injection
conf/                 → file konfigurasi (config.env)
domain/user/          → domain user (interfaces, usecase, repository, handler HTTP)
model/                → entity, request DTO, response DTO
shared/               → config, database, helper response JSON
```

## Cara Menjalankan

```bash
# 1. Install dependency
go mod download

# 2. Sesuaikan konfigurasi database di conf/config.env
#    (set DB_DRIVER=sqlite3 untuk jalan tanpa PostgreSQL)

# 3. Jalankan server
make run
```

Server berjalan di `http://localhost:8080`.

Menjalankan test:

```bash
make test
```

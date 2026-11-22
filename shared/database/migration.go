package database

import (
	"go-boilerplate/shared/config"
	"go-boilerplate/shared/logger"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

// RunMigrations runs pending SQL migrations located in `migrationsPath`.
func RunMigrations(migrationsPath string) error {
	cfg := config.GetConfig()
	log := logger.Log()

	var dbURL string
	if cfg.DBDriver == "postgres" {
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)
	} else {
		dbURL = fmt.Sprintf("sqlite3://%s", cfg.DBName)
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		dbURL,
	)
	if err != nil {
		log.Error("Failed to create migration instance", zap.Error(err))
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error("Migration execution failed", zap.Error(err))
		return err
	}

	log.Info("Database migrations applied successfully")
	return nil
}

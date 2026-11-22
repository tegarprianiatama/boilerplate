package database

import (
	"boilerplate/model/entity"
	"boilerplate/shared/config"
	"boilerplate/shared/logger"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func OpenDB() (*gorm.DB, error) {
	cfg := config.GetConfig()
	log := logger.Log()

	var dialector gorm.Dialector

	if cfg.DBDriver == "postgres" {
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
		dialector = postgres.Open(dsn)
	} else {
		dsn := cfg.DBName
		if dsn == "" {
			dsn = "app.db"
		}
		dialector = sqlite.Open(dsn)
	}

	gormLogLevel := gormLogger.Info
	if cfg.AppEnv == "production" {
		gormLogLevel = gormLogger.Error
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database (%s): %w", cfg.DBDriver, err)
	}

	// Auto Migration for User entity
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		log.Error("Failed auto-migrating tables", zap.Error(err))
	} else {
		log.Info("Database migration completed successfully")
	}

	return db, nil
}

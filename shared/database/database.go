package database

import (
	"boilerplate/model/entity"
	"boilerplate/shared/config"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/mattn/go-sqlite3"
)

func OpenDB() (*gorm.DB, error) {
	driver := config.Env("DB_DRIVER", "postgres")

	var dsn string
	if driver == "postgres" {
		host := config.Env("DB_HOST", "localhost")
		port := config.Env("DB_PORT", "5432")
		user := config.Env("DB_USER", "postgres")
		pass := config.Env("DB_PASSWORD", "postgres")
		dbname := config.Env("DB_NAME", "user_db")
		sslmode := config.Env("DB_SSLMODE", "disable")

		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, pass, dbname, sslmode)
	} else {
		dsn = config.Env("DB_NAME", "app.db")
	}

	db, err := gorm.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database (%s): %w", driver, err)
	}

	db.LogMode(true)

	// Auto Migration for User entity
	if err := db.AutoMigrate(&entity.User{}).Error; err != nil {
		log.Printf("Failed auto-migrating tables: %v", err)
	}

	return db, nil
}

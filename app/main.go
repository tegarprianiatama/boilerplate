package main

import (
	"log"
	"net/http"

	userHttpDelivery "boilerplate/domain/user/delivery/http"
	userRepository "boilerplate/domain/user/repository"
	userUsecase "boilerplate/domain/user/usecase"

	jsonResponse "boilerplate/shared/common/json"
	"boilerplate/shared/config"
	"boilerplate/shared/database"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func main() {
	config.LoadEnv()

	db, err := database.OpenDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to Go Clean Architecture Boilerplate")
	})

	jsonHelper := jsonResponse.NewResponse()

	apiV1 := e.Group("/api/v1")

	// Dependency Injection
	userRepo := userRepository.NewRepository(db)
	userUse := userUsecase.NewUsecase(userRepo)
	userHttpDelivery.UserHandler(apiV1, userUse, jsonHelper)

	port := config.Env("APP_PORT", "8080")
	log.Printf("Server starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
}

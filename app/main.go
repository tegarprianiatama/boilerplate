package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	userHttpDelivery "boilerplate/domain/user/delivery/http"
	userRepository "boilerplate/domain/user/repository"
	userUsecase "boilerplate/domain/user/usecase"

	"boilerplate/shared/cache"
	jsonResponse "boilerplate/shared/common/json"
	"boilerplate/shared/config"
	"boilerplate/shared/database"
	"boilerplate/shared/logger"
	"boilerplate/shared/queue"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
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

// @title Go Clean Architecture Boilerplate API
// @version 1.0
// @description REST API boilerplate for User CRUD using Go, Echo, and GORM with Clean Architecture.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// 7. Type-safe configuration
	cfg := config.LoadEnv()

	// 4. Structured Logging (Zap)
	log := logger.InitLogger(cfg.AppEnv)
	defer logger.Sync()

	// 2. Database Migrations & 1. GORM v2 connection
	if err := database.RunMigrations("db/migrations"); err != nil {
		log.Warn("Golang-migrate encountered an issue or skipped, falling back to GORM AutoMigrate", zap.Error(err))
	}

	db, err := database.OpenDB()
	if err != nil {
		log.Fatal("Failed to initialize database", zap.Error(err))
	}
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	// Redis connection
	if cfg.RedisEnabled {
		redisClient, err := cache.OpenRedis()
		if err != nil {
			log.Fatal("Failed to initialize Redis", zap.Error(err))
		}
		defer func() {
			if err := redisClient.Close(); err != nil {
				log.Warn("Failed to close Redis connection", zap.Error(err))
			}
		}()
	}

	// RabbitMQ connection
	if cfg.RabbitMQEnabled {
		rabbitConn, err := queue.OpenRabbitMQ()
		if err != nil {
			log.Fatal("Failed to initialize RabbitMQ", zap.Error(err))
		}
		defer func() {
			if err := rabbitConn.Close(); err != nil {
				log.Warn("Failed to close RabbitMQ connection", zap.Error(err))
			}
		}()
	}

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Zap structured request logging middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			stop := time.Since(start)

			req := c.Request()
			res := c.Response()

			log.Info("HTTP Request",
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.Int("status", res.Status),
				zap.Duration("latency", stop),
				zap.String("remote_ip", c.RealIP()),
			)
			return err
		}
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to Go Clean Architecture Boilerplate (Production-Ready)")
	})

	jsonHelper := jsonResponse.NewResponse()
	apiV1 := e.Group("/api/v1")

	// Dependency Injection
	userRepo := userRepository.NewRepository(db)
	userUse := userUsecase.NewUsecase(userRepo)
	userHttpDelivery.UserHandler(apiV1, userUse, jsonHelper)

	// 9. Graceful Shutdown
	port := cfg.AppPort
	go func() {
		log.Info("Server starting", zap.String("port", port))
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server shutdown error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Info("Shutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exiting")
}

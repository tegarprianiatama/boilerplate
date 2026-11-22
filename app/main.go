package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	userHttpDelivery "go-boilerplate/domain/user/delivery/http"
	userRepository "go-boilerplate/domain/user/repository"
	userUsecase "go-boilerplate/domain/user/usecase"

	authHttpDelivery "go-boilerplate/domain/auth/delivery/http"
	authRepository "go-boilerplate/domain/auth/repository"
	authUsecase "go-boilerplate/domain/auth/usecase"

	"go-boilerplate/shared/cache"
	jsonResponse "go-boilerplate/shared/common/json"
	"go-boilerplate/shared/config"
	"go-boilerplate/shared/database"
	"go-boilerplate/shared/logger"
	"go-boilerplate/shared/queue"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
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
// @description REST API go-boilerplate for User CRUD using Go, Echo, GORM, Redis, and RabbitMQ with Clean Architecture.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Type-safe configuration
	cfg := config.LoadEnv()

	// Structured Logging (Zap)
	log := logger.InitLogger(cfg.AppEnv)
	defer logger.Sync()

	// Database Migrations & GORM v2 connection
	if err := database.RunMigrations("db/migrations"); err != nil {
		log.Warn("Golang-migrate encountered an issue or skipped, falling back to GORM AutoMigrate", zap.Error(err))
	}

	db, err := database.OpenDB()
	if err != nil {
		log.Fatal("Failed to initialize database", zap.Error(err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB", zap.Error(err))
	}
	defer sqlDB.Close()

	// Redis connection
	var redisClient *redis.Client
	if cfg.RedisEnabled {
		redisClient, err = cache.OpenRedis()
		if err != nil {
			log.Fatal("Failed to initialize Redis", zap.Error(err))
		}
		defer redisClient.Close()
	}

	// RabbitMQ connection
	var rabbitConn *amqp.Connection
	if cfg.RabbitMQEnabled {
		rabbitConn, err = queue.OpenRabbitMQ()
		if err != nil {
			log.Fatal("Failed to initialize RabbitMQ", zap.Error(err))
		}
		defer rabbitConn.Close()
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

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
		defer cancel()

		healthStatus := map[string]string{
			"status":   "ok",
			"database": "connected",
		}
		statusCode := http.StatusOK

		// Check database
		if err := sqlDB.PingContext(ctx); err != nil {
			healthStatus["database"] = "disconnected"
			healthStatus["status"] = "error"
			statusCode = http.StatusServiceUnavailable
		}

		// Check Redis if enabled
		if cfg.RedisEnabled {
			if redisClient != nil {
				if err := redisClient.Ping(ctx).Err(); err != nil {
					healthStatus["redis"] = "disconnected"
					healthStatus["status"] = "error"
					statusCode = http.StatusServiceUnavailable
				} else {
					healthStatus["redis"] = "connected"
				}
			} else {
				healthStatus["redis"] = "disabled"
			}
		}

		// Check RabbitMQ if enabled
		if cfg.RabbitMQEnabled {
			if rabbitConn != nil {
				if rabbitConn.IsClosed() {
					healthStatus["rabbitmq"] = "disconnected"
					healthStatus["status"] = "error"
					statusCode = http.StatusServiceUnavailable
				} else {
					healthStatus["rabbitmq"] = "connected"
				}
			} else {
				healthStatus["rabbitmq"] = "disabled"
			}
		}

		return c.JSON(statusCode, healthStatus)
	})

	jsonHelper := jsonResponse.NewResponse()
	apiV1 := e.Group("/api/v1")

	// Dependency Injection
	authRepo := authRepository.NewRepository(db)
	authUse := authUsecase.NewUsecase(authRepo)
	authHttpDelivery.AuthHandler(apiV1, authUse, jsonHelper)

	userRepo := userRepository.NewRepository(db)
	userUse := userUsecase.NewUsecase(userRepo)
	userHttpDelivery.UserHandler(apiV1, userUse, jsonHelper)

	// Graceful Shutdown
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

	// 1. Shutdown HTTP server first
	if err := e.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
	}

	// 2. Close RabbitMQ connection
	if rabbitConn != nil && !rabbitConn.IsClosed() {
		if err := rabbitConn.Close(); err != nil {
			log.Warn("Failed to close RabbitMQ connection", zap.Error(err))
		} else {
			log.Info("RabbitMQ connection closed successfully")
		}
	}

	// 3. Close Redis connection
	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			log.Warn("Failed to close Redis connection", zap.Error(err))
		} else {
			log.Info("Redis connection closed successfully")
		}
	}

	// 4. Database connection closed via defer sqlDB.Close()

	log.Info("Server exiting")
}

package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv         string
	AppPort        string
	DBDriver       string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	JWTSecret      string
	JWTExpireHours int

	RedisEnabled  bool
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	RabbitMQEnabled  bool
	RabbitMQHost     string
	RabbitMQPort     string
	RabbitMQUser     string
	RabbitMQPassword string
	RabbitMQVHost    string
}

var globalConfig *Config

func LoadEnv() *Config {
	err := godotenv.Load("conf/config.env")
	if err != nil {
		log.Println("Notice: conf/config.env file not found, fallback to system env")
	}

	globalConfig = &Config{
		AppEnv:         getEnv("APP_ENV", "development"),
		AppPort:        getEnv("APP_PORT", "8080"),
		DBDriver:       getEnv("DB_DRIVER", "postgres"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "user_db"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		JWTSecret:      getEnv("JWT_SECRET", "super-secret-key-change-me"),
		JWTExpireHours: getEnvAsInt("JWT_EXPIRE_HOURS", 24),

		RedisEnabled:  getEnvAsBool("REDIS_ENABLED", true),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		RabbitMQEnabled:  getEnvAsBool("RABBITMQ_ENABLED", true),
		RabbitMQHost:     getEnv("RABBITMQ_HOST", "localhost"),
		RabbitMQPort:     getEnv("RABBITMQ_PORT", "5672"),
		RabbitMQUser:     getEnv("RABBITMQ_USER", "guest"),
		RabbitMQPassword: getEnv("RABBITMQ_PASSWORD", "guest"),
		RabbitMQVHost:    getEnv("RABBITMQ_VHOST", "/"),
	}

	if globalConfig.JWTExpireHours <= 0 {
		globalConfig.JWTExpireHours = 24
	}

	return globalConfig
}

func GetConfig() *Config {
	if globalConfig == nil {
		return LoadEnv()
	}
	return globalConfig
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return fallback
}

// Env helper maintained for backward compatibility if needed
func Env(key, fallback string) string {
	return getEnv(key, fallback)
}

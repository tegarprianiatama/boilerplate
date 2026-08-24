package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load("conf/config.env")
	if err != nil {
		log.Println("Notice: conf/config.env file not found, fallback to system env")
	}
}

func Env(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

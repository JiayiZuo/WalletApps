package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_DSN     string
	PORT       string
	JWT_SECRET string
	REDIS_ADDR string
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	return &Config{
		DB_DSN:     getEnv("DB_DSN", "postgres://postgres:postgres@localhost:5432/wallet?sslmode=disable"),
		PORT:       getEnv("PORT", "8080"),
		JWT_SECRET: getEnv("JWT_SECRET", "mysecret"),
		REDIS_ADDR: getEnv("REDIS_ADDR", "localhost:6379"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	log.Printf("WARN: env %s not set, using default: %s", key, fallback)
	return fallback
}

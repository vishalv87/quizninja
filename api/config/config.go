package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	Port           string
	GinMode        string
	AllowedOrigins string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "quizninja"),
		JWTSecret:      getEnv("JWT_SECRET", "your_jwt_secret"),
		Port:           getEnv("PORT", "8080"),
		GinMode:        getEnv("GIN_MODE", "debug"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
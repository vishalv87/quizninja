package config

import (
	"log"
	"os"
	"strconv"
	"strings"

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

	// Supabase configuration
	UseSupabase        bool
	SupabaseURL        string
	SupabaseAnonKey    string
	SupabaseDBHost     string
	SupabaseDBPort     string
	SupabaseDBUser     string
	SupabaseDBPassword string
	SupabaseDBName     string
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

		// Supabase configuration
		UseSupabase:        getBoolEnv("USE_SUPABASE", false),
		SupabaseURL:        getEnv("SUPABASE_URL", ""),
		SupabaseAnonKey:    getEnv("SUPABASE_ANON_KEY", ""),
		SupabaseDBHost:     getEnv("SUPABASE_DB_HOST", ""),
		SupabaseDBPort:     getEnv("SUPABASE_DB_PORT", "5432"),
		SupabaseDBUser:     getEnv("SUPABASE_DB_USER", ""),
		SupabaseDBPassword: getEnv("SUPABASE_DB_PASSWORD", ""),
		SupabaseDBName:     getEnv("SUPABASE_DB_NAME", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		// Parse common truthy/falsy values
		switch strings.ToLower(value) {
		case "true", "1", "yes", "on":
			return true
		case "false", "0", "no", "off":
			return false
		default:
			// Try parsing as boolean
			if parsed, err := strconv.ParseBool(value); err == nil {
				return parsed
			}
		}
	}
	return defaultValue
}
package config

import (
	"fmt"
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
	Port           string
	GinMode        string
	AllowedOrigins string

	// Supabase configuration
	UseSupabase        bool
	UseSupabaseAuth    bool
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
		Port:           getEnv("PORT", "8080"),
		GinMode:        getEnv("GIN_MODE", "debug"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),

		// Supabase configuration
		UseSupabase:        getBoolEnv("USE_SUPABASE", false),
		UseSupabaseAuth:    getBoolEnv("USE_SUPABASE_AUTH", false),
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

// ValidateSupabaseConfig validates Supabase configuration when enabled
func (c *Config) ValidateSupabaseConfig() error {
	if !c.UseSupabaseAuth {
		return nil // No validation needed when disabled
	}

	var errors []string

	if c.SupabaseURL == "" {
		errors = append(errors, "SUPABASE_URL is required when USE_SUPABASE_AUTH=true")
	}

	if c.SupabaseAnonKey == "" {
		errors = append(errors, "SUPABASE_ANON_KEY is required when USE_SUPABASE_AUTH=true")
	}

	if len(errors) > 0 {
		return fmt.Errorf("supabase configuration errors: %v", errors)
	}

	return nil
}

// GetAuthStrategy returns the current authentication strategy
func (c *Config) GetAuthStrategy() string {
	if c.UseSupabaseAuth {
		return "supabase-only"
	}
	return "supabase-only"
}

// IsSupabaseEnabled returns true if Supabase features are enabled
func (c *Config) IsSupabaseEnabled() bool {
	return c.UseSupabase
}

// IsSupabaseAuthEnabled returns true if Supabase authentication is enabled
func (c *Config) IsSupabaseAuthEnabled() bool {
	return c.UseSupabaseAuth
}

// ValidateConfig validates the entire configuration
func (c *Config) ValidateConfig() error {
	var errors []string

	// Basic validation

	if c.DBHost == "" {
		errors = append(errors, "DB_HOST is required")
	}

	if c.DBName == "" {
		errors = append(errors, "DB_NAME is required")
	}

	// Supabase validation
	if err := c.ValidateSupabaseConfig(); err != nil {
		errors = append(errors, err.Error())
	}

	// Database configuration validation
	if c.UseSupabase {
		if c.SupabaseDBHost == "" {
			errors = append(errors, "SUPABASE_DB_HOST is required when USE_SUPABASE=true")
		}
		if c.SupabaseDBUser == "" {
			errors = append(errors, "SUPABASE_DB_USER is required when USE_SUPABASE=true")
		}
		if c.SupabaseDBPassword == "" {
			errors = append(errors, "SUPABASE_DB_PASSWORD is required when USE_SUPABASE=true")
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration validation failed: %v", errors)
	}

	return nil
}
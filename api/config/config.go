package config

import (
	"flag"
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
	SupabaseURL        string
	SupabaseAnonKey    string
	SupabaseServiceKey string // For admin operations and testing
	SupabaseDBHost     string
	SupabaseDBPort     string
	SupabaseDBUser     string
	SupabaseDBPassword string
	SupabaseDBName     string

	// Test configuration
	UseMockAuth bool // Use mock auth manager for tests instead of real Supabase

	// Rate Limiting Configuration
	RateLimitEnabled bool
	RateLimitGlobal  int64 // requests per minute per IP
	RateLimitAuth    int64 // requests per minute per IP for auth endpoints
	RateLimitWrite   int64 // requests per minute per IP for write operations
	RateLimitPerUser int64 // requests per minute per authenticated user
}

func Load() *Config {
	// Detect if we're running tests
	isTestEnv := isTestEnvironment()

	if isTestEnv {
		// Load test configuration - no fallback
		err := godotenv.Load(".env.test")
		if err != nil {
			// Try from parent directory (for tests running from tests/ dir)
			err = godotenv.Load("../.env.test")
			if err != nil {
				log.Println("Warning: No .env.test file found, using environment variables for tests")
			}
		}
		log.Println("Loaded test environment configuration (.env.test)")
	} else {
		// Load production/development configuration
		err := godotenv.Load(".env")
		if err != nil {
			// If not found, try parent directory
			err = godotenv.Load("../.env")
			if err != nil {
				log.Println("No .env file found, using environment variables")
			}
		}
		log.Println("Loaded development/production environment configuration (.env)")
	}

	cfg := &Config{
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
		SupabaseURL:        getEnv("SUPABASE_URL", ""),
		SupabaseAnonKey:    getEnv("SUPABASE_ANON_KEY", ""),
		SupabaseServiceKey: getEnv("SUPABASE_SERVICE_KEY", ""),
		SupabaseDBHost:     getEnv("SUPABASE_DB_HOST", ""),
		SupabaseDBPort:     getEnv("SUPABASE_DB_PORT", "5432"),
		SupabaseDBUser:     getEnv("SUPABASE_DB_USER", ""),
		SupabaseDBPassword: getEnv("SUPABASE_DB_PASSWORD", ""),
		SupabaseDBName:     getEnv("SUPABASE_DB_NAME", ""),

		// Test configuration
		UseMockAuth: getBoolEnv("USE_MOCK_AUTH", false),

		// Rate Limiting Configuration
		RateLimitEnabled: getBoolEnv("RATE_LIMIT_ENABLED", true),
		RateLimitGlobal:  getInt64Env("RATE_LIMIT_GLOBAL", 100),
		RateLimitAuth:    getInt64Env("RATE_LIMIT_AUTH", 5),
		RateLimitWrite:   getInt64Env("RATE_LIMIT_WRITE", 20),
		RateLimitPerUser: getInt64Env("RATE_LIMIT_PER_USER", 60),
	}

	// ✅ SECURITY: Prevent mock auth in production
	if cfg.GinMode == "release" && cfg.UseMockAuth {
		log.Fatal("SECURITY ERROR: Mock authentication cannot be enabled in release mode")
	}

	return cfg
}

// isTestEnvironment detects if we're running in test mode
func isTestEnvironment() bool {
	// Check if we're running with go test
	if flag.Lookup("test.v") != nil {
		return true
	}

	// Check if the executable name contains ".test"
	if strings.Contains(os.Args[0], ".test") {
		return true
	}

	// Check GO_ENV environment variable
	if os.Getenv("GO_ENV") == "test" {
		return true
	}

	// Check if any test-related flags are set
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}

	return false
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

func getInt64Env(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// ValidateSupabaseConfig validates Supabase configuration for authentication
func (c *Config) ValidateSupabaseConfig() error {
	// Skip validation if using mock auth for testing
	if c.UseMockAuth {
		return nil
	}

	var errors []string

	if c.SupabaseURL == "" {
		errors = append(errors, "SUPABASE_URL is required for authentication (unless USE_MOCK_AUTH=true)")
	}

	if c.SupabaseAnonKey == "" {
		errors = append(errors, "SUPABASE_ANON_KEY is required for authentication (unless USE_MOCK_AUTH=true)")
	}

	if len(errors) > 0 {
		return fmt.Errorf("supabase configuration errors: %v", errors)
	}

	return nil
}

// GetAuthStrategy returns the current authentication strategy
func (c *Config) GetAuthStrategy() string {
	if c.UseMockAuth {
		return "mock-auth"
	}
	return "supabase-only"
}

// IsSupabaseEnabled returns true if Supabase features are enabled
func (c *Config) IsSupabaseEnabled() bool {
	return c.UseSupabase
}

// IsMockAuthEnabled returns true if mock authentication is enabled for testing
func (c *Config) IsMockAuthEnabled() bool {
	return c.UseMockAuth
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

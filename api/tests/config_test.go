package tests

import (
	"os"
	"testing"

	"quizninja-api/config"

	"github.com/stretchr/testify/assert"
)

// TestConfigEnvironmentDetection verifies that the config loader correctly
// detects test environment and loads appropriate configuration files
func TestConfigEnvironmentDetection(t *testing.T) {
	t.Run("TestEnvironmentConfiguration", func(t *testing.T) {
		// This test should automatically load .env.test because we're running with go test
		cfg := config.Load()

		// Verify test environment configuration
		assert.True(t, cfg.IsMockAuthEnabled(), "Mock auth should be enabled in test environment")
		assert.False(t, cfg.IsSupabaseEnabled(), "Supabase database should be disabled in test environment")

		// Verify test database configuration
		assert.Equal(t, "localhost", cfg.DBHost, "Test should use local database host")
		assert.Equal(t, "quizninja", cfg.DBName, "Test should use test database name")
		assert.Equal(t, "test", cfg.GinMode, "Test should use test gin mode")

		// Verify auth strategy
		assert.Equal(t, "mock-auth", cfg.GetAuthStrategy(), "Test should use mock auth strategy")

		// Verify empty Supabase credentials (not needed for mock auth)
		assert.Empty(t, cfg.SupabaseURL, "Supabase URL should be empty in test environment")
		assert.Empty(t, cfg.SupabaseAnonKey, "Supabase anon key should be empty in test environment")
		assert.Empty(t, cfg.SupabaseServiceKey, "Supabase service key should be empty in test environment")
	})
}

// TestConfigProductionEnvironment verifies production configuration
// This test temporarily simulates non-test environment
func TestConfigProductionEnvironment(t *testing.T) {
	t.Run("ProductionEnvironmentConfiguration", func(t *testing.T) {
		// Save original args
		originalArgs := os.Args

		// Temporarily modify os.Args to simulate non-test environment
		os.Args = []string{"quizninja-api"}

		// Temporarily set production environment variables
		os.Setenv("USE_SUPABASE", "true")
		os.Setenv("USE_MOCK_AUTH", "false")
		os.Setenv("DB_NAME", "quizninja")
		os.Setenv("GIN_MODE", "debug")

		// Load config (should load .env, not .env.test)
		cfg := config.Load()

		// Verify production environment configuration
		assert.False(t, cfg.IsMockAuthEnabled(), "Mock auth should be disabled in production environment")
		assert.True(t, cfg.IsSupabaseEnabled(), "Supabase database should be enabled in production")

		// Verify production database configuration
		assert.Equal(t, "quizninja", cfg.DBName, "Production should use production database name")
		assert.Equal(t, "debug", cfg.GinMode, "Production should use debug gin mode")

		// Verify auth strategy
		assert.Equal(t, "supabase-only", cfg.GetAuthStrategy(), "Production should use real Supabase auth strategy")

		// Restore original args and clean up environment
		os.Args = originalArgs
		os.Unsetenv("USE_SUPABASE")
		os.Unsetenv("USE_MOCK_AUTH")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("GIN_MODE")
	})
}

// TestConfigValidation verifies that configuration validation works correctly
func TestConfigValidation(t *testing.T) {
	t.Run("MockAuthValidation", func(t *testing.T) {
		cfg := config.Load()

		// Mock auth configuration should pass validation
		err := cfg.ValidateConfig()
		assert.NoError(t, err, "Mock auth configuration should be valid")

		// Supabase validation should be skipped when mock auth is enabled
		err = cfg.ValidateSupabaseConfig()
		assert.NoError(t, err, "Supabase validation should be skipped for mock auth")
	})
}

// TestEnvironmentVariableDifferences documents the expected differences between environments
func TestEnvironmentVariableDifferences(t *testing.T) {
	cfg := config.Load()

	// Document expected test environment values
	expectedTestConfig := map[string]interface{}{
		"USE_MOCK_AUTH":     true,
		"USE_SUPABASE":      false,
		"DB_NAME":           "quizninja",
		"GIN_MODE":          "test",
		"SUPABASE_URL":      "",
		"SUPABASE_ANON_KEY": "",
	}

	// Verify each expected configuration
	assert.Equal(t, expectedTestConfig["USE_MOCK_AUTH"], cfg.UseMockAuth)
	assert.Equal(t, expectedTestConfig["USE_SUPABASE"], cfg.UseSupabase)
	assert.Equal(t, expectedTestConfig["DB_NAME"], cfg.DBName)
	assert.Equal(t, expectedTestConfig["GIN_MODE"], cfg.GinMode)
	assert.Equal(t, expectedTestConfig["SUPABASE_URL"], cfg.SupabaseURL)
	assert.Equal(t, expectedTestConfig["SUPABASE_ANON_KEY"], cfg.SupabaseAnonKey)

	// Log configuration for debugging
	t.Logf("Test Configuration Summary:")
	t.Logf("- Auth Strategy: %s", cfg.GetAuthStrategy())
	t.Logf("- Mock Auth: %v", cfg.IsMockAuthEnabled())
	t.Logf("- Supabase DB: %v", cfg.IsSupabaseEnabled())
	t.Logf("- Database Name: %s", cfg.DBName)
	t.Logf("- Gin Mode: %s", cfg.GinMode)
}

// TestConfigFileLoading verifies that the correct config files are loaded
func TestConfigFileLoading(t *testing.T) {
	// This test verifies that .env.test is loaded instead of .env during tests
	cfg := config.Load()

	// These values should come from .env.test, not .env
	assert.Equal(t, "quizninja", cfg.DBName, "Should load test database name from .env.test")
	assert.Equal(t, "test", cfg.GinMode, "Should load test gin mode from .env.test")
	assert.True(t, cfg.UseMockAuth, "Should load mock auth setting from .env.test")
	assert.False(t, cfg.UseSupabase, "Should load local DB setting from .env.test")

	// Verify that test gin mode is loaded (not production debug mode)
	assert.NotEqual(t, "debug", cfg.GinMode, "Should not load production gin mode")
	assert.Equal(t, "test", cfg.GinMode, "Should specifically load test gin mode")

	// Verify DB_NAME is the expected value (same for both test and production environments)
	assert.Equal(t, "quizninja", cfg.DBName, "DB name should be 'quizninja' for both test and production")
}

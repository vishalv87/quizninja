package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMockJWTGeneration verifies that mock JWTs can be generated without JWT_SECRET env var
func TestMockJWTGeneration(t *testing.T) {
	userID := "test-user-123"
	email := "test@example.com"
	name := "Test User"

	// Generate mock JWT using hardcoded config (no env var needed)
	token, err := GenerateMockJWT(userID, email, name, DefaultMockJWTConfig)
	require.NoError(t, err, "Should generate mock JWT successfully")
	assert.NotEmpty(t, token, "Generated token should not be empty")

	// Validate the generated token
	user, err := ValidateMockJWT(token, DefaultMockJWTConfig)
	require.NoError(t, err, "Should validate mock JWT successfully")

	assert.Equal(t, userID, user.ID, "User ID should match")
	assert.Equal(t, email, user.Email, "Email should match")
	assert.Equal(t, name, user.UserMetadata["name"], "Name should match")
	assert.True(t, user.UserMetadata["is_test"].(bool), "Should be marked as test")
}

// TestMockAuthManager verifies the mock auth manager works without JWT_SECRET
func TestMockAuthManager(t *testing.T) {
	manager := NewMockSupabaseTestAuthManager()

	// Create a test user
	testUser, err := manager.CreateUniqueTestUser("Mock Test")
	require.NoError(t, err, "Should create test user successfully")
	assert.NotEmpty(t, testUser.AccessToken, "Access token should be generated")

	// Validate the token
	user, err := manager.ValidateTestToken(testUser.AccessToken)
	require.NoError(t, err, "Should validate test token successfully")
	assert.Equal(t, testUser.ID, user.ID, "User ID should match")

	// Cleanup
	err = manager.CleanupTestUser(testUser.ID)
	assert.NoError(t, err, "Should cleanup test user successfully")
}

// TestMockJWTWithoutEnvVar verifies mock auth doesn't depend on JWT_SECRET env var
func TestMockJWTWithoutEnvVar(t *testing.T) {
	// This test runs without any JWT_SECRET environment variable
	// and should still work because mock auth uses hardcoded secret

	config := DefaultMockJWTConfig
	assert.Equal(t, "mock_test_secret_key_for_quizninja_testing_only", config.SecretKey)
	assert.Equal(t, "mock-supabase-auth", config.Issuer)
	assert.Equal(t, "authenticated", config.Audience)

	// Generate and validate token
	token, err := GenerateMockJWT("user-123", "test@example.com", "Test User", config)
	require.NoError(t, err)

	user, err := ValidateMockJWT(token, config)
	require.NoError(t, err)
	assert.Equal(t, "user-123", user.ID)
}
package utils

// TestAuthManager defines the interface for authentication managers used in testing
type TestAuthManager interface {
	// CreateUniqueTestUser creates a test user with a unique email
	CreateUniqueTestUser(namePrefix string) (*SupabaseTestUser, error)

	// CreateTestUser creates a test user with specific email and name
	CreateTestUser(email, name string) (*SupabaseTestUser, error)

	// CleanupTestUser removes a test user
	CleanupTestUser(userID string) error

	// CleanupAllTestUsers removes all tracked test users
	CleanupAllTestUsers()

	// ValidateTestToken validates that a test token is still valid
	ValidateTestToken(token string) (*SupabaseUser, error)

	// RefreshTestToken refreshes a test user's token
	RefreshTestToken(refreshToken string) (*SupabaseAuthResponse, error)
}

// TokenValidator defines the interface for token validation
type TokenValidator interface {
	// ValidateToken validates a token and returns user information
	ValidateToken(token string) (*SupabaseUser, error)
}

// Ensure existing SupabaseTestAuthManager implements the interface
var _ TestAuthManager = (*SupabaseTestAuthManager)(nil)
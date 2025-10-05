package utils

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MockSupabaseTestAuthManager provides a mock implementation of TestAuthManager
// that doesn't require real Supabase connections for testing
type MockSupabaseTestAuthManager struct {
	jwtConfig MockJWTConfig
	users     map[string]*MockUser // userID -> user
	usersByEmail map[string]*MockUser // email -> user
	mutex     sync.RWMutex
}

// MockUser represents a mock user for testing
type MockUser struct {
	ID           string
	Email        string
	Name         string
	AccessToken  string
	RefreshToken string
	CreatedAt    time.Time
}

// NewMockSupabaseTestAuthManager creates a new mock auth manager
func NewMockSupabaseTestAuthManager() *MockSupabaseTestAuthManager {
	return &MockSupabaseTestAuthManager{
		jwtConfig:    DefaultMockJWTConfig,
		users:        make(map[string]*MockUser),
		usersByEmail: make(map[string]*MockUser),
	}
}

// CreateTestUser creates a mock test user with specific email and name
func (m *MockSupabaseTestAuthManager) CreateTestUser(email, name string) (*SupabaseTestUser, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if user already exists
	if _, exists := m.usersByEmail[email]; exists {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Generate new user ID
	userID := uuid.New().String()

	// Generate tokens
	accessToken, err := GenerateMockJWT(userID, email, name, m.jwtConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := GenerateMockRefreshToken(userID, m.jwtConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create mock user
	mockUser := &MockUser{
		ID:           userID,
		Email:        email,
		Name:         name,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now(),
	}

	// Store user
	m.users[userID] = mockUser
	m.usersByEmail[email] = mockUser

	// Return SupabaseTestUser
	return &SupabaseTestUser{
		ID:           userID,
		Email:        email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour).Unix(),
	}, nil
}

// CreateUniqueTestUser creates a test user with a unique email
func (m *MockSupabaseTestAuthManager) CreateUniqueTestUser(namePrefix string) (*SupabaseTestUser, error) {
	uniqueID := uuid.New().String()[:8]
	cleanPrefix := strings.ReplaceAll(strings.ToLower(namePrefix), " ", "")
	email := fmt.Sprintf("test%s%s@example.com", cleanPrefix, uniqueID)
	name := fmt.Sprintf("Test %s User", namePrefix)

	return m.CreateTestUser(email, name)
}

// CleanupTestUser removes a test user
func (m *MockSupabaseTestAuthManager) CleanupTestUser(userID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	user, exists := m.users[userID]
	if !exists {
		// Not an error - user might already be cleaned up
		return nil
	}

	// Remove from both maps
	delete(m.users, userID)
	delete(m.usersByEmail, user.Email)

	return nil
}

// CleanupAllTestUsers removes all tracked test users
func (m *MockSupabaseTestAuthManager) CleanupAllTestUsers() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Clear all users
	m.users = make(map[string]*MockUser)
	m.usersByEmail = make(map[string]*MockUser)
}

// ValidateTestToken validates that a test token is still valid
func (m *MockSupabaseTestAuthManager) ValidateTestToken(token string) (*SupabaseUser, error) {
	return ValidateMockJWT(token, m.jwtConfig)
}

// RefreshTestToken refreshes a test user's token
func (m *MockSupabaseTestAuthManager) RefreshTestToken(refreshToken string) (*SupabaseAuthResponse, error) {
	// Validate the refresh token
	userID, err := ValidateMockRefreshToken(refreshToken, m.jwtConfig)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	m.mutex.RLock()
	user, exists := m.users[userID]
	m.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("user not found for refresh token")
	}

	// Generate new tokens
	newAccessToken, err := GenerateMockJWT(user.ID, user.Email, user.Name, m.jwtConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	newRefreshToken, err := GenerateMockRefreshToken(user.ID, m.jwtConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	// Update stored tokens
	m.mutex.Lock()
	user.AccessToken = newAccessToken
	user.RefreshToken = newRefreshToken
	m.mutex.Unlock()

	// Create user object
	now := time.Now()
	supabaseUser := SupabaseUser{
		ID:               user.ID,
		Email:            user.Email,
		EmailConfirmedAt: &now,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        now,
		LastSignInAt:     &now,
		Role:             "authenticated",
		UserMetadata: map[string]interface{}{
			"name":    user.Name,
			"is_test": true,
		},
		AppMetadata:  map[string]interface{}{},
		IdentityData: []interface{}{},
		Aud:          m.jwtConfig.Audience,
	}

	return &SupabaseAuthResponse{
		User:         supabaseUser,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    86400, // 24 hours
		ExpiresAt:    time.Now().Add(24 * time.Hour).Unix(),
		TokenType:    "Bearer",
	}, nil
}

// GetUserByID retrieves a mock user by ID (utility method for testing)
func (m *MockSupabaseTestAuthManager) GetUserByID(userID string) (*MockUser, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	user, exists := m.users[userID]
	return user, exists
}

// GetUserByEmail retrieves a mock user by email (utility method for testing)
func (m *MockSupabaseTestAuthManager) GetUserByEmail(email string) (*MockUser, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	user, exists := m.usersByEmail[email]
	return user, exists
}

// Ensure MockSupabaseTestAuthManager implements TestAuthManager interface
var _ TestAuthManager = (*MockSupabaseTestAuthManager)(nil)
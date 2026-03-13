package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SupabaseTestUser represents a test user created in Supabase for testing
type SupabaseTestUser struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

// CreateSupabaseTestUserRequest represents the request to create a test user via admin API
type CreateSupabaseTestUserRequest struct {
	Email        string                 `json:"email"`
	Password     string                 `json:"password"`
	EmailConfirm bool                   `json:"email_confirm"`
	UserMetadata map[string]interface{} `json:"user_metadata,omitempty"`
	AppMetadata  map[string]interface{} `json:"app_metadata,omitempty"`
}

// SupabaseTestAuthManager handles test user creation and cleanup
type SupabaseTestAuthManager struct {
	supabaseURL string
	serviceKey  string
	anonKey     string
	testUsers   []string // Track created users for cleanup
}

// NewSupabaseTestAuthManager creates a new test auth manager
func NewSupabaseTestAuthManager(supabaseURL, serviceKey, anonKey string) *SupabaseTestAuthManager {
	return &SupabaseTestAuthManager{
		supabaseURL: supabaseURL,
		serviceKey:  serviceKey,
		anonKey:     anonKey,
		testUsers:   make([]string, 0),
	}
}

// CreateTestUser creates a test user in Supabase and returns access token
func (m *SupabaseTestAuthManager) CreateTestUser(email, name string) (*SupabaseTestUser, error) {
	if m.serviceKey == "" || m.supabaseURL == "" {
		return nil, fmt.Errorf("Supabase service key and URL required for test user creation")
	}

	// Create user via Supabase Admin API
	createReq := CreateSupabaseTestUserRequest{
		Email:        email,
		Password:     "test_password_123", // Fixed password for all test users
		EmailConfirm: true,                // Auto-confirm email for testing
		UserMetadata: map[string]interface{}{
			"name":        name,
			"is_test":     true,
			"created_for": "integration_test",
		},
	}

	reqBody, err := json.Marshal(createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create user request: %w", err)
	}

	// Create user via admin API
	url := fmt.Sprintf("%s/auth/v1/admin/users", m.supabaseURL)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", m.serviceKey),
		"apikey":        m.serviceKey,
	}

	resp, err := makeSupabaseRequest("POST", url, headers, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create test user: %w", err)
	}
	defer resp.Body.Close()

	var responseBody bytes.Buffer
	responseBody.ReadFrom(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create test user: HTTP %d: %s", resp.StatusCode, responseBody.String())
	}

	var supabaseUser SupabaseUser
	if err := json.Unmarshal(responseBody.Bytes(), &supabaseUser); err != nil {
		return nil, fmt.Errorf("failed to parse create user response: %w", err)
	}

	// Track the user for cleanup
	m.testUsers = append(m.testUsers, supabaseUser.ID)

	// Now generate an access token for this user by signing them in
	accessToken, refreshToken, expiresAt, err := m.generateTokenForUser(email, "test_password_123")
	if err != nil {
		// If token generation fails, still return the user info but log the error
		return &SupabaseTestUser{
			ID:           supabaseUser.ID,
			Email:        supabaseUser.Email,
			AccessToken:  fmt.Sprintf("mock_token_for_%s", supabaseUser.ID),
			RefreshToken: "",
			ExpiresAt:    time.Now().Add(time.Hour).Unix(),
		}, fmt.Errorf("user created but token generation failed: %w", err)
	}

	return &SupabaseTestUser{
		ID:           supabaseUser.ID,
		Email:        supabaseUser.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// generateTokenForUser signs in the user to get a valid access token
func (m *SupabaseTestAuthManager) generateTokenForUser(email, password string) (string, string, int64, error) {
	url := fmt.Sprintf("%s/auth/v1/token?grant_type=password", m.supabaseURL)

	signInReq := map[string]string{
		"email":    email,
		"password": password,
	}

	reqBody, err := json.Marshal(signInReq)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to marshal sign-in request: %w", err)
	}

	headers := map[string]string{
		"apikey": m.anonKey,
	}

	resp, err := makeSupabaseRequest("POST", url, headers, reqBody)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to sign in test user: %w", err)
	}
	defer resp.Body.Close()

	var responseBody bytes.Buffer
	responseBody.ReadFrom(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", "", 0, fmt.Errorf("failed to sign in test user: HTTP %d: %s", resp.StatusCode, responseBody.String())
	}

	var authResp SupabaseAuthResponse
	if err := json.Unmarshal(responseBody.Bytes(), &authResp); err != nil {
		return "", "", 0, fmt.Errorf("failed to parse sign-in response: %w", err)
	}

	return authResp.AccessToken, authResp.RefreshToken, authResp.ExpiresAt, nil
}

// CleanupTestUser removes a test user from Supabase
func (m *SupabaseTestAuthManager) CleanupTestUser(userID string) error {
	if m.serviceKey == "" || m.supabaseURL == "" {
		return nil // Skip cleanup if not configured
	}

	url := fmt.Sprintf("%s/auth/v1/admin/users/%s", m.supabaseURL, userID)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", m.serviceKey),
		"apikey":        m.serviceKey,
	}

	resp, err := makeSupabaseRequest("DELETE", url, headers, nil)
	if err != nil {
		return fmt.Errorf("failed to delete test user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		var responseBody bytes.Buffer
		responseBody.ReadFrom(resp.Body)
		return fmt.Errorf("failed to delete test user: HTTP %d: %s", resp.StatusCode, responseBody.String())
	}

	// HTTP 404 means user already deleted - this is fine for cleanup purposes

	return nil
}

// CleanupAllTestUsers removes all tracked test users
func (m *SupabaseTestAuthManager) CleanupAllTestUsers() {
	for _, userID := range m.testUsers {
		if err := m.CleanupTestUser(userID); err != nil {
			// Log error but continue cleanup
			fmt.Printf("Warning: Failed to cleanup test user %s: %v\n", userID, err)
		}
	}
	m.testUsers = make([]string, 0)
}

// CreateUniqueTestUser creates a test user with a unique email
func (m *SupabaseTestAuthManager) CreateUniqueTestUser(namePrefix string) (*SupabaseTestUser, error) {
	uniqueID := uuid.New().String()[:8]
	// Use a simpler email format that Supabase will accept
	cleanPrefix := strings.ReplaceAll(strings.ToLower(namePrefix), " ", "")
	email := fmt.Sprintf("test%s%s@example.com", cleanPrefix, uniqueID)
	name := fmt.Sprintf("Test %s User", namePrefix)

	return m.CreateTestUser(email, name)
}

// ValidateTestToken validates that a test token is still valid
func (m *SupabaseTestAuthManager) ValidateTestToken(token string) (*SupabaseUser, error) {
	user, supabaseErr := ValidateSupabaseTokenHTTP(token, m.supabaseURL, m.anonKey)
	if supabaseErr != nil {
		return nil, supabaseErr // SupabaseAuthError implements error interface
	}
	return user, nil
}

// RefreshTestToken refreshes a test user's token
func (m *SupabaseTestAuthManager) RefreshTestToken(refreshToken string) (*SupabaseAuthResponse, error) {
	authResp, supabaseErr := RefreshSupabaseTokenHTTP(refreshToken, m.supabaseURL, m.anonKey)
	if supabaseErr != nil {
		return nil, supabaseErr // SupabaseAuthError implements error interface
	}
	return authResp, nil
}

package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// MockJWTConfig holds configuration for mock JWT generation
type MockJWTConfig struct {
	SecretKey string
	Issuer    string
	Audience  string
}

// Default mock JWT configuration for testing
var DefaultMockJWTConfig = MockJWTConfig{
	SecretKey: "mock_test_secret_key_for_quizninja_testing_only",
	Issuer:    "mock-supabase-auth",
	Audience:  "authenticated",
}

// MockJWTClaims represents the claims in a mock JWT token
type MockJWTClaims struct {
	UserID   string `json:"sub"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Audience string `json:"aud"`
	IsTest   bool   `json:"is_test"`
	jwt.RegisteredClaims
}

// GenerateMockJWT creates a mock JWT token for testing
func GenerateMockJWT(userID, email, name string, config MockJWTConfig) (string, error) {
	now := time.Now()
	expirationTime := now.Add(24 * time.Hour) // 24 hour expiry

	claims := &MockJWTClaims{
		UserID:   userID,
		Email:    email,
		Name:     name,
		Audience: config.Audience,
		IsTest:   true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    config.Issuer,
			Subject:   userID,
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey))
}

// ValidateMockJWT validates a mock JWT token and extracts user information
func ValidateMockJWT(tokenString string, config MockJWTConfig) (*SupabaseUser, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MockJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.SecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(*MockJWTClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse claims")
	}

	// Check if token is expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("token is expired")
	}

	// Create SupabaseUser from claims
	now := time.Now()
	user := &SupabaseUser{
		ID:               claims.UserID,
		Email:            claims.Email,
		EmailConfirmedAt: &now, // Mock as confirmed
		CreatedAt:        now,
		UpdatedAt:        now,
		LastSignInAt:     &now,
		Role:             "authenticated",
		UserMetadata: map[string]interface{}{
			"name":    claims.Name,
			"is_test": true,
		},
		AppMetadata:  map[string]interface{}{},
		IdentityData: []interface{}{},
		Aud:          claims.Audience,
	}

	return user, nil
}

// GenerateMockRefreshToken creates a mock refresh token
func GenerateMockRefreshToken(userID string, config MockJWTConfig) (string, error) {
	now := time.Now()
	expirationTime := now.Add(30 * 24 * time.Hour) // 30 day expiry

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    config.Issuer,
		Subject:   userID,
		ID:        uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey + "_refresh"))
}

// ValidateMockRefreshToken validates a mock refresh token
func ValidateMockRefreshToken(tokenString string, config MockJWTConfig) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.SecretKey + "_refresh"), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse refresh token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("refresh token is not valid")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", fmt.Errorf("failed to parse refresh token claims")
	}

	// Check if token is expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return "", fmt.Errorf("refresh token is expired")
	}

	return claims.Subject, nil
}
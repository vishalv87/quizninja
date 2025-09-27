package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"quizninja-api/config"
	"quizninja-api/database"
	"quizninja-api/models"
	"quizninja-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TestConfig holds test configuration
type TestConfig struct {
	Server *gin.Engine
	Config *config.Config
}

// SetupTestServer initializes the test server with real database connection
func SetupTestServer(t *testing.T) *TestConfig {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Load configuration
	cfg := config.Load()

	// Initialize database connection (includes migrations)
	database.Connect(cfg)

	// Create gin engine and setup routes
	server := gin.New()
	routes.SetupRoutes(server, cfg)

	return &TestConfig{
		Server: server,
		Config: cfg,
	}
}

// CreateTestUser creates a test user and returns user ID and auth token
func CreateTestUser(t *testing.T, tc *TestConfig) (uuid.UUID, string) {
	// Create a test user
	registerReq := models.RegisterRequest{
		Email:    fmt.Sprintf("test_%s@example.com", uuid.New().String()[:8]),
		Password: "testpassword123",
		Name:     "Test User",
		Age:      intPtr(25),
	}

	reqBody, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	tc.Server.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create test user: %d %s", w.Code, w.Body.String())
	}

	var response models.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse register response: %v", err)
	}

	// Mark the user as test data
	_, err = database.DB.Exec("UPDATE users SET is_test_data = true WHERE id = $1", response.User.ID)
	if err != nil {
		t.Fatalf("Failed to mark user as test data: %v", err)
	}

	return response.User.ID, response.AccessToken
}

// MakeAuthenticatedRequest makes an HTTP request with authentication
func MakeAuthenticatedRequest(t *testing.T, tc *TestConfig, method, path, token string, body []byte) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	tc.Server.ServeHTTP(w, req)
	return w
}

// MakeRequest makes a simple HTTP request without authentication
func MakeRequest(t *testing.T, tc *TestConfig, method, path string) *httptest.ResponseRecorder {
	return MakeAuthenticatedRequest(t, tc, method, path, "", nil)
}

// VerifyIsTestDataField verifies that the is_test_data field exists and has the expected value
func VerifyIsTestDataField(t *testing.T, data map[string]interface{}, expectedValue bool, fieldPath string) {
	value, exists := data["is_test_data"]
	if !exists {
		t.Errorf("Field 'is_test_data' not found in %s", fieldPath)
		return
	}

	boolValue, ok := value.(bool)
	if !ok {
		t.Errorf("Field 'is_test_data' in %s is not a boolean, got %T", fieldPath, value)
		return
	}

	if boolValue != expectedValue {
		t.Errorf("Field 'is_test_data' in %s has incorrect value. Expected: %v, Got: %v", fieldPath, expectedValue, boolValue)
	}
}

// VerifyIsTestDataInArray verifies is_test_data field in array of objects
func VerifyIsTestDataInArray(t *testing.T, items []interface{}, expectedValue bool, arrayName string) {
	for i, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			t.Errorf("Item %d in %s is not a map", i, arrayName)
			continue
		}
		VerifyIsTestDataField(t, itemMap, expectedValue, fmt.Sprintf("%s[%d]", arrayName, i))
	}
}

// ParseJSONResponse parses HTTP response as JSON
func ParseJSONResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Fatalf("Request failed with status %d: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	return response
}

// GetDataFromResponse extracts the 'data' field from API response
func GetDataFromResponse(t *testing.T, response map[string]interface{}) map[string]interface{} {
	data, exists := response["data"]
	if !exists {
		t.Fatalf("Response does not contain 'data' field")
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		t.Fatalf("Response 'data' field is not an object")
	}

	return dataMap
}

// Cleanup cleans up test resources (mainly closes DB connection)
// Individual tests should use CleanupTestUser/CleanupTestQuiz for specific cleanup
func Cleanup(t *testing.T) {
	// Close database connections
	if database.DB != nil {
		database.DB.Close()
	}
}

// TestMain can be used in individual test files for setup/teardown
func TestMain(m *testing.M) {
	// Setup
	code := m.Run()
	// Teardown
	os.Exit(code)
}

// Helper functions for pointer conversion
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

// Test cleanup helper functions

// CleanupTestUser deletes a specific test user and all related data
func CleanupTestUser(userID uuid.UUID) {
	if database.DB == nil {
		return
	}

	// Delete in reverse order to respect foreign key constraints
	cleanupQueries := []struct {
		query string
		desc  string
	}{
		{"DELETE FROM user_quiz_favorites WHERE user_id = $1", "favorites"},
		{"DELETE FROM quiz_sessions WHERE user_id = $1", "quiz sessions"},
		{"DELETE FROM quiz_attempts WHERE user_id = $1", "quiz attempts"},
		{"DELETE FROM user_achievements WHERE user_id = $1", "user achievements"},
		{"DELETE FROM friend_notifications WHERE user_id = $1 OR related_user_id = $1", "friend notifications"},
		{"DELETE FROM challenges WHERE challenger_id = $1 OR challenged_id = $1", "challenges"},
		{"DELETE FROM friendships WHERE user1_id = $1 OR user2_id = $1", "friendships"},
		{"DELETE FROM friend_requests WHERE requester_id = $1 OR requested_id = $1", "friend requests"},
		{"DELETE FROM discussion_replies WHERE user_id = $1", "discussion replies"},
		{"DELETE FROM discussions WHERE user_id = $1", "discussions"},
		{"DELETE FROM user_preferences WHERE user_id = $1", "user preferences"},
		{"DELETE FROM refresh_tokens WHERE user_id = $1", "refresh tokens"},
		{"DELETE FROM users WHERE id = $1", "user"},
	}

	for _, cleanup := range cleanupQueries {
		_, err := database.DB.Exec(cleanup.query, userID)
		if err != nil {
			// Log error but continue cleanup
			fmt.Printf("Warning: Failed to cleanup %s for user %s: %v\n", cleanup.desc, userID, err)
		}
	}
}

// CleanupTestQuiz deletes a specific test quiz and all related data
func CleanupTestQuiz(quizID uuid.UUID) {
	if database.DB == nil {
		return
	}

	cleanupQueries := []struct {
		query string
		desc  string
	}{
		{"DELETE FROM user_quiz_favorites WHERE quiz_id = $1", "quiz favorites"},
		{"DELETE FROM quiz_sessions WHERE quiz_id = $1", "quiz sessions"},
		{"DELETE FROM quiz_attempts WHERE quiz_id = $1", "quiz attempts"},
		{"DELETE FROM challenges WHERE quiz_id = $1", "challenges"},
		{"DELETE FROM quiz_statistics WHERE quiz_id = $1", "quiz statistics"},
		{"DELETE FROM questions WHERE quiz_id = $1", "questions"},
		{"DELETE FROM quizzes WHERE id = $1", "quiz"},
	}

	for _, cleanup := range cleanupQueries {
		_, err := database.DB.Exec(cleanup.query, quizID)
		if err != nil {
			fmt.Printf("Warning: Failed to cleanup %s for quiz %s: %v\n", cleanup.desc, quizID, err)
		}
	}
}
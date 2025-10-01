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
	"github.com/lib/pq"
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

// createTestCategoryDirect creates a test category if it doesn't exist
func createTestCategoryDirect(t *testing.T) string {
	categoryID := "test_category"

	// Check if category already exists
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM interests WHERE id = $1", categoryID).Scan(&count)
	if err == nil && count > 0 {
		return categoryID
	}

	// Create test category
	_, err = database.DB.Exec(`
		INSERT INTO interests (id, name, description, icon_name, color_hex, is_test_data)
		VALUES ($1, $2, $3, $4, $5, true)
		ON CONFLICT (id) DO UPDATE SET is_test_data = true
	`, categoryID, "Test Category", "Category for testing purposes", "test", "#FF0000")

	if err != nil {
		t.Logf("Warning: Could not create test category: %v", err)
		// Try to use an existing category
		err = database.DB.QueryRow("SELECT id FROM interests WHERE is_test_data = true LIMIT 1").Scan(&categoryID)
		if err != nil {
			// Use any existing category
			err = database.DB.QueryRow("SELECT id FROM interests LIMIT 1").Scan(&categoryID)
			if err != nil {
				t.Fatalf("No categories available and cannot create test category: %v", err)
			}
		}
	}

	return categoryID
}

// createTestDifficultyDirect creates or gets a test difficulty level
func createTestDifficultyDirect(t *testing.T) string {
	difficultyID := "test_difficulty"

	// Check if difficulty already exists
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM difficulty_levels WHERE id = $1", difficultyID).Scan(&count)
	if err == nil && count > 0 {
		return difficultyID
	}

	// Create test difficulty
	_, err = database.DB.Exec(`
		INSERT INTO difficulty_levels (id, name, description, icon_name, background_color_hex, is_test_data)
		VALUES ($1, $2, $3, $4, $5, true)
		ON CONFLICT (id) DO UPDATE SET is_test_data = true
	`, difficultyID, "Test Difficulty", "Difficulty for testing", "test", "#00FF00")

	if err != nil {
		t.Logf("Warning: Could not create test difficulty: %v", err)
		// Try to use an existing difficulty
		err = database.DB.QueryRow("SELECT id FROM difficulty_levels WHERE is_test_data = true LIMIT 1").Scan(&difficultyID)
		if err != nil {
			// Use any existing difficulty
			err = database.DB.QueryRow("SELECT id FROM difficulty_levels LIMIT 1").Scan(&difficultyID)
			if err != nil {
				t.Fatalf("No difficulty levels available and cannot create test difficulty: %v", err)
			}
		}
	}

	return difficultyID
}

// createTestQuizComprehensive creates a complete test quiz with questions
func createTestQuizComprehensive(t *testing.T, createdByUserID uuid.UUID) uuid.UUID {
	categoryID := createTestCategoryDirect(t)
	difficultyID := createTestDifficultyDirect(t)

	quizID := uuid.New()

	// Create the quiz
	_, err := database.DB.Exec(`
		INSERT INTO quizzes (
			id, title, description, category_id, difficulty, total_questions,
			time_limit_minutes, points, created_by, is_featured, is_active,
			is_public, tags, is_test_data, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, quizID, "Test Quiz", "A quiz for testing purposes", categoryID, difficultyID,
		3, 10, 100, createdByUserID, false, true, true, `{"test", "sample"}`)

	if err != nil {
		t.Fatalf("Failed to create test quiz: %v", err)
	}

	// Create test questions
	questions := []struct {
		text    string
		options []string
		answer  string
		order   int
	}{
		{"What is 2 + 2?", []string{"3", "4", "5", "6"}, "4", 1},
		{"What is the capital of France?", []string{"London", "Berlin", "Paris", "Madrid"}, "Paris", 2},
		{"Which planet is closest to the Sun?", []string{"Venus", "Mercury", "Earth", "Mars"}, "Mercury", 3},
	}

	for _, q := range questions {
		questionID := uuid.New()
		_, err = database.DB.Exec(`
			INSERT INTO questions (
				id, quiz_id, question_text, question_type, options, correct_answer,
				explanation, order_index, is_test_data, created_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, true, CURRENT_TIMESTAMP)
		`, questionID, quizID, q.text, "multipleChoice", pq.Array(q.options), q.answer,
			"Test explanation", q.order)

		if err != nil {
			t.Fatalf("Failed to create test question: %v", err)
		}
	}

	// Create quiz statistics (try both possible column names for last_updated)
	_, err = database.DB.Exec(`
		INSERT INTO quiz_statistics (
			quiz_id, total_attempts, total_completions, average_score,
			average_time_seconds, difficulty_rating, popularity_score,
			is_test_data, updated_at
		)
		VALUES ($1, 0, 0, 0.0, 0, 0.0, 0, true, CURRENT_TIMESTAMP)
	`, quizID)

	if err != nil {
		// Try with last_updated column name as fallback
		_, err = database.DB.Exec(`
			INSERT INTO quiz_statistics (
				quiz_id, total_attempts, total_completions, average_score,
				average_time_seconds, difficulty_rating, popularity_score,
				is_test_data, last_updated
			)
			VALUES ($1, 0, 0, 0.0, 0, 0.0, 0, true, CURRENT_TIMESTAMP)
		`, quizID)
	}

	if err != nil {
		t.Logf("Warning: Failed to create quiz statistics: %v", err)
	}

	return quizID
}

// getFirstAvailableQuiz gets an available quiz for testing
func getFirstAvailableQuiz(t *testing.T, tc *TestConfig, token string) uuid.UUID {
	// Try to get an existing test quiz first
	var quizID uuid.UUID
	err := database.DB.QueryRow(`
		SELECT id FROM quizzes WHERE is_test_data = true LIMIT 1
	`).Scan(&quizID)

	if err == nil {
		return quizID
	}

	// Try to get any existing quiz
	err = database.DB.QueryRow(`
		SELECT id FROM quizzes WHERE is_active = true LIMIT 1
	`).Scan(&quizID)

	if err == nil {
		return quizID
	}

	// No quizzes exist, create one
	// First get a test user to create the quiz
	userID, _ := CreateTestUser(t, tc)

	// Create comprehensive test quiz
	quizID = createTestQuizComprehensive(t, userID)

	return quizID
}

// setupFriendship creates a friendship between two users via API or direct database
func setupFriendship(t *testing.T, tc *TestConfig, user1ID, user2ID uuid.UUID, user1Token, user2Token string) {
	// Try API approach first
	if tryAPIFriendshipSetup(t, tc, user1ID, user2ID, user1Token, user2Token) {
		return
	}

	// Fall back to direct database approach
	t.Logf("API friendship setup failed, using direct database approach")
	setupFriendshipDirect(t, user1ID, user2ID)
}

// tryAPIFriendshipSetup attempts to create friendship via API, returns true if successful
func tryAPIFriendshipSetup(t *testing.T, tc *TestConfig, user1ID, user2ID uuid.UUID, user1Token, user2Token string) bool {
	// Create friend request
	friendReq := models.SendFriendRequestRequest{
		RequestedUserID: user2ID,
		Message:         stringPtr("Test friendship"),
	}

	reqBody, _ := json.Marshal(friendReq)
	w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/friends/requests", user1Token, reqBody)
	if w.Code != http.StatusCreated {
		t.Logf("Friend request creation failed with status: %d", w.Code)
		return false
	}

	// Get the friend request ID
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Logf("Failed to parse friend request response: %v", err)
		return false
	}

	// Check if response contains request (the correct field name)
	friendRequestData, exists := response["request"]
	if !exists || friendRequestData == nil {
		t.Logf("Friend request API response format not as expected: %v", response)
		return false
	}

	friendRequest := friendRequestData.(map[string]interface{})
	requestIDInterface, exists := friendRequest["id"]
	if !exists {
		t.Logf("Friend request response missing id field: %v", friendRequest)
		return false
	}
	requestID := requestIDInterface.(string)

	// Accept the friend request
	acceptUrl := fmt.Sprintf("/api/v1/friends/requests/%s/accept", requestID)
	w2 := MakeAuthenticatedRequest(t, tc, "PUT", acceptUrl, user2Token, nil)
	if w2.Code != http.StatusOK {
		t.Logf("Friend request acceptance failed with status: %d, response: %s", w2.Code, w2.Body.String())
		return false
	}

	return true
}

// setupFriendshipDirect creates friendship directly in database
func setupFriendshipDirect(t *testing.T, user1ID, user2ID uuid.UUID) {
	// Ensure consistent ordering for friendship
	if user1ID.String() > user2ID.String() {
		user1ID, user2ID = user2ID, user1ID
	}

	_, err := database.DB.Exec(`
		INSERT INTO friendships (user1_id, user2_id, is_test_data)
		VALUES ($1, $2, true)
		ON CONFLICT (user1_id, user2_id) DO NOTHING
	`, user1ID, user2ID)

	if err != nil {
		t.Fatalf("Failed to create direct friendship: %v", err)
	}
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
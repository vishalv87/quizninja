package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"quizninja-api/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler(t *testing.T) {
	tc := SetupTestServer(t)
	defer CleanupWithSupabase(t, tc)

	t.Run("Register", func(t *testing.T) {
		registerReq := models.RegisterRequest{
			SupabaseUserID: uuid.New().String(),
			Email:          fmt.Sprintf("testregister_%s@example.com", uuid.New().String()[:8]),
			Name:           "Test Register User",
		}

		reqBody, _ := json.Marshal(registerReq)
		w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/auth/register", "", reqBody)

		if w.Code == http.StatusCreated {
			// Parse response to get user ID for cleanup
			var response models.AuthResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err == nil {
				// Cleanup this specific test user
				defer CleanupTestUser(response.User.ID)
			}
		}

		if w.Code == http.StatusCreated {
			var response models.AuthResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should parse register response")

			// Verify user data has is_test_data field
			userMap := map[string]interface{}{
				"is_test_data": response.User.IsTestData,
			}
			VerifyIsTestDataField(t, userMap, true, "registered user")

			// Verify preferences if they exist
			if response.User.Preferences != nil {
				preferencesMap := map[string]interface{}{
					"is_test_data": response.User.Preferences.IsTestData,
				}
				VerifyIsTestDataField(t, preferencesMap, true, "user preferences")
			}
		}
	})

	// Create a test user for subsequent tests
	_, token, _, cleanup := CreateTestUserWithCleanup(t, tc, "Auth Handler Main User")
	defer cleanup()

	t.Run("Login", func(t *testing.T) {
		uniqueEmail := fmt.Sprintf("testlogin_%s@example.com", uuid.New().String()[:8])
		loginReq := models.LoginRequest{
			SupabaseUserID: uuid.New().String(),
			Email:          uniqueEmail,
			Name:           "Test Login User",
		}

		// First register a user to login with
		registerReq := models.RegisterRequest{
			SupabaseUserID: loginReq.SupabaseUserID,
			Email:          loginReq.Email,
			Name:           "Test Login User",
		}

		reqBody, _ := json.Marshal(registerReq)
		w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/auth/register", "", reqBody)
		assert.Equal(t, http.StatusCreated, w.Code)

		// Parse response to get user ID for cleanup
		var registerResponse models.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &registerResponse)
		if err == nil {
			// Cleanup this specific test user
			defer CleanupTestUser(registerResponse.User.ID)
		}

		// Now login
		reqBody, _ = json.Marshal(loginReq)
		w = MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/auth/login", "", reqBody)

		if w.Code == http.StatusOK {
			var response models.AuthResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should parse login response")

			// Verify user data has is_test_data field
			userMap := map[string]interface{}{
				"is_test_data": response.User.IsTestData,
			}
			VerifyIsTestDataField(t, userMap, true, "logged in user")

			// Verify preferences if they exist
			if response.User.Preferences != nil {
				preferencesMap := map[string]interface{}{
					"is_test_data": response.User.Preferences.IsTestData,
				}
				VerifyIsTestDataField(t, preferencesMap, true, "login user preferences")
			}
		}
	})

	t.Run("GetProfile", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/profile", token, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		response := ParseJSONResponse(t, w)

		// Verify is_test_data field in profile
		VerifyIsTestDataField(t, response, true, "user profile")

		// Check preferences if they exist
		preferences, prefsExist := response["preferences"]
		if prefsExist && preferences != nil {
			preferencesMap, ok := preferences.(map[string]interface{})
			if ok {
				VerifyIsTestDataField(t, preferencesMap, true, "profile preferences")
			}
		}
	})

	t.Run("UpdateProfile", func(t *testing.T) {
		updateReq := models.UpdateProfileRequest{
			Name: stringPtr("Updated Test User"),
		}

		reqBody, _ := json.Marshal(updateReq)
		w := MakeAuthenticatedRequest(t, tc, "PUT", "/api/v1/profile", token, reqBody)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// Verify is_test_data field in updated profile
			VerifyIsTestDataField(t, response, true, "updated user profile")

			// Check preferences if they exist
			preferences, prefsExist := response["preferences"]
			if prefsExist && preferences != nil {
				preferencesMap, ok := preferences.(map[string]interface{})
				if ok {
					VerifyIsTestDataField(t, preferencesMap, true, "updated profile preferences")
				}
			}
		}
	})

	t.Run("GetUserStats", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/stats", token, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		response := ParseJSONResponse(t, w)

		// Verify data contains statistics with is_test_data field
		data, dataExists := response["data"]
		assert.True(t, dataExists, "Response should contain 'data' field")

		statisticsMap, ok := data.(map[string]interface{})
		assert.True(t, ok, "Data field should be an object")

		VerifyIsTestDataField(t, statisticsMap, true, "user statistics")
	})

	t.Run("Logout", func(t *testing.T) {
		// Create a test user for logout
		_, token, _, cleanup := CreateTestUserWithCleanup(t, tc, "Auth Handler Logout User")
		defer cleanup()

		// Logout only requires authentication header, no body needed
		w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/auth/logout", token, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		response := ParseJSONResponse(t, w)
		message, exists := response["message"]
		assert.True(t, exists, "Response should contain logout message")
		assert.Contains(t, message, "Successfully logged out", "Should confirm successful logout")
	})

	t.Run("LogoutWithoutAuth", func(t *testing.T) {
		// Test that logout fails without authentication
		w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/auth/logout", "", nil)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Parse error response manually since ParseJSONResponse expects 200/201
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Should parse error response JSON")

		errorMsg, exists := response["error"]
		assert.True(t, exists, "Response should contain error message")
		assert.Contains(t, errorMsg, "Authorization header required", "Should indicate missing authorization")
	})

	t.Run("RegisterWithIdempotencyKey_FirstRequest", func(t *testing.T) {
		// Test registration with idempotency key - first request should succeed
		idempotencyKey := "test-idempotency-key-" + uuid.New().String()
		registerReq := models.RegisterRequest{
			SupabaseUserID: uuid.New().String(),
			Email:          fmt.Sprintf("testidempotency_%s@example.com", uuid.New().String()[:8]),
			Name:           "Test Idempotency User",
		}

		reqBody, _ := json.Marshal(registerReq)
		headers := map[string]string{
			"X-Idempotency-Key": idempotencyKey,
		}

		w := MakeRequestWithHeaders(t, tc, "POST", "/api/v1/auth/register", "", reqBody, headers)
		assert.Equal(t, http.StatusCreated, w.Code, "First request with idempotency key should succeed")

		// Parse response to get user ID for cleanup
		var response models.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Should parse register response")

		// Cleanup test user
		defer CleanupTestUser(response.User.ID)

		// Verify response
		assert.Equal(t, registerReq.Email, response.User.Email, "Email should match")
		assert.Equal(t, registerReq.Name, response.User.Name, "Name should match")
	})

	t.Run("RegisterWithIdempotencyKey_DuplicateRequest", func(t *testing.T) {
		// Test that duplicate requests with same idempotency key return cached response
		idempotencyKey := "test-idempotency-duplicate-" + uuid.New().String()
		registerReq := models.RegisterRequest{
			SupabaseUserID: uuid.New().String(),
			Email:          fmt.Sprintf("testidempdup_%s@example.com", uuid.New().String()[:8]),
			Name:           "Test Idempotency Duplicate User",
		}

		reqBody, _ := json.Marshal(registerReq)
		headers := map[string]string{
			"X-Idempotency-Key": idempotencyKey,
		}

		// First request - should create user
		w1 := MakeRequestWithHeaders(t, tc, "POST", "/api/v1/auth/register", "", reqBody, headers)
		assert.Equal(t, http.StatusCreated, w1.Code, "First request should succeed")

		var response1 models.AuthResponse
		err := json.Unmarshal(w1.Body.Bytes(), &response1)
		assert.NoError(t, err, "Should parse first response")

		// Cleanup test user
		defer CleanupTestUser(response1.User.ID)

		// Second request with SAME idempotency key - should return cached response
		w2 := MakeRequestWithHeaders(t, tc, "POST", "/api/v1/auth/register", "", reqBody, headers)
		assert.Equal(t, http.StatusCreated, w2.Code, "Second request should also return 201 (cached)")

		var response2 models.AuthResponse
		err = json.Unmarshal(w2.Body.Bytes(), &response2)
		assert.NoError(t, err, "Should parse second response")

		// Verify that both responses are identical (cached)
		assert.Equal(t, response1.User.ID, response2.User.ID, "User IDs should match (cached response)")
		assert.Equal(t, response1.User.Email, response2.User.Email, "Emails should match")
		assert.Equal(t, response1.User.Name, response2.User.Name, "Names should match")
	})

	t.Run("RegisterWithIdempotencyKey_DifferentKeys", func(t *testing.T) {
		// Test that different idempotency keys allow different requests
		baseEmail := fmt.Sprintf("testidemp_%s", uuid.New().String()[:8])

		registerReq1 := models.RegisterRequest{
			SupabaseUserID: uuid.New().String(),
			Email:          baseEmail + "_1@example.com",
			Name:           "Test Idempotency User 1",
		}

		registerReq2 := models.RegisterRequest{
			SupabaseUserID: uuid.New().String(),
			Email:          baseEmail + "_2@example.com",
			Name:           "Test Idempotency User 2",
		}

		// First request with key1
		reqBody1, _ := json.Marshal(registerReq1)
		headers1 := map[string]string{
			"X-Idempotency-Key": "test-key-1-" + uuid.New().String(),
		}
		w1 := MakeRequestWithHeaders(t, tc, "POST", "/api/v1/auth/register", "", reqBody1, headers1)
		assert.Equal(t, http.StatusCreated, w1.Code)

		var response1 models.AuthResponse
		err := json.Unmarshal(w1.Body.Bytes(), &response1)
		assert.NoError(t, err)
		defer CleanupTestUser(response1.User.ID)

		// Second request with different key2
		reqBody2, _ := json.Marshal(registerReq2)
		headers2 := map[string]string{
			"X-Idempotency-Key": "test-key-2-" + uuid.New().String(),
		}
		w2 := MakeRequestWithHeaders(t, tc, "POST", "/api/v1/auth/register", "", reqBody2, headers2)
		assert.Equal(t, http.StatusCreated, w2.Code)

		var response2 models.AuthResponse
		err = json.Unmarshal(w2.Body.Bytes(), &response2)
		assert.NoError(t, err)
		defer CleanupTestUser(response2.User.ID)

		// Verify that different users were created
		assert.NotEqual(t, response1.User.ID, response2.User.ID, "Different keys should create different users")
		assert.NotEqual(t, response1.User.Email, response2.User.Email, "Users should have different emails")
	})

	t.Run("RegisterWithoutIdempotencyKey", func(t *testing.T) {
		// Test that registration still works without idempotency key
		registerReq := models.RegisterRequest{
			SupabaseUserID: uuid.New().String(),
			Email:          fmt.Sprintf("testnoidempkey_%s@example.com", uuid.New().String()[:8]),
			Name:           "Test No Idempotency Key User",
		}

		reqBody, _ := json.Marshal(registerReq)
		// No idempotency key header
		w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/auth/register", "", reqBody)
		assert.Equal(t, http.StatusCreated, w.Code, "Registration should work without idempotency key")

		var response models.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Cleanup test user
		defer CleanupTestUser(response.User.ID)

		assert.Equal(t, registerReq.Email, response.User.Email)
		assert.Equal(t, registerReq.Name, response.User.Name)
	})

}

// Helper functions for creating pointers to primitive types

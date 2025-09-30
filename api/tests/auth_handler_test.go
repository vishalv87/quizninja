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
	defer Cleanup(t)

	t.Run("Register", func(t *testing.T) {
		registerReq := models.RegisterRequest{
			Email:    fmt.Sprintf("testregister_%s@example.com", uuid.New().String()[:8]),
			Password: "testpassword123",
			Name:     "Test Register User",
			Age:      intPtr(25),
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
	userID, token := CreateTestUser(t, tc)

	t.Run("Login", func(t *testing.T) {
		uniqueEmail := fmt.Sprintf("testlogin_%s@example.com", uuid.New().String()[:8])
		loginReq := models.LoginRequest{
			Email:    uniqueEmail,
			Password: "testpassword123",
		}

		// First register a user to login with
		registerReq := models.RegisterRequest{
			Email:    loginReq.Email,
			Password: loginReq.Password,
			Name:     "Test Login User",
			Age:      intPtr(25),
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
			Age:  intPtr(30),
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

	t.Run("RefreshToken", func(t *testing.T) {
		// First get a refresh token by registering a new user
		registerReq := models.RegisterRequest{
			Email:    fmt.Sprintf("testrefresh_%s@example.com", uuid.New().String()[:8]),
			Password: "testpassword123",
			Name:     "Test Refresh User",
			Age:      intPtr(25),
		}

		reqBody, _ := json.Marshal(registerReq)
		w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/auth/register", "", reqBody)
		assert.Equal(t, http.StatusCreated, w.Code)

		var registerResponse models.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &registerResponse)
		assert.NoError(t, err)

		// Cleanup this specific test user
		defer CleanupTestUser(registerResponse.User.ID)

		// Use refresh token to get new tokens
		refreshReq := models.RefreshTokenRequest{
			RefreshToken: registerResponse.RefreshToken,
		}

		reqBody, _ = json.Marshal(refreshReq)
		w = MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/auth/refresh", "", reqBody)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// Verify response contains new tokens
			_, accessTokenExists := response["access_token"]
			assert.True(t, accessTokenExists, "Response should contain new access token")

			_, refreshTokenExists := response["refresh_token"]
			assert.True(t, refreshTokenExists, "Response should contain new refresh token")
		}
	})

	t.Run("Logout", func(t *testing.T) {
		// First get a refresh token by registering a new user
		registerReq := models.RegisterRequest{
			Email:    fmt.Sprintf("testlogout_%s@example.com", uuid.New().String()[:8]),
			Password: "testpassword123",
			Name:     "Test Logout User",
			Age:      intPtr(25),
		}

		reqBody, _ := json.Marshal(registerReq)
		w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/auth/register", "", reqBody)
		assert.Equal(t, http.StatusCreated, w.Code)

		var registerResponse models.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &registerResponse)
		assert.NoError(t, err)

		// Cleanup this specific test user
		defer CleanupTestUser(registerResponse.User.ID)

		// Logout using refresh token - now requires authentication
		logoutReq := models.RefreshTokenRequest{
			RefreshToken: registerResponse.RefreshToken,
		}

		reqBody, _ = json.Marshal(logoutReq)
		w = MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/auth/logout", registerResponse.AccessToken, reqBody)
		assert.Equal(t, http.StatusOK, w.Code)

		response := ParseJSONResponse(t, w)
		message, exists := response["message"]
		assert.True(t, exists, "Response should contain logout message")
		assert.Contains(t, message, "Successfully logged out", "Should confirm successful logout")
	})

	t.Run("LogoutWithoutAuth", func(t *testing.T) {
		// Test that logout fails without authentication
		logoutReq := models.RefreshTokenRequest{
			RefreshToken: "dummy-token",
		}

		reqBody, _ := json.Marshal(logoutReq)
		w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/auth/logout", "", reqBody)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Parse error response manually since ParseJSONResponse expects 200/201
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Should parse error response JSON")

		errorMsg, exists := response["error"]
		assert.True(t, exists, "Response should contain error message")
		assert.Contains(t, errorMsg, "Authorization header required", "Should indicate missing authorization")
	})

	t.Run("LogoutWithWrongUserToken", func(t *testing.T) {
		// Create two users to test cross-user token validation
		registerReq1 := models.RegisterRequest{
			Email:    fmt.Sprintf("testlogout1_%s@example.com", uuid.New().String()[:8]),
			Password: "testpassword123",
			Name:     "Test Logout User 1",
			Age:      intPtr(25),
		}

		reqBody, _ := json.Marshal(registerReq1)
		w1 := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/auth/register", "", reqBody)
		assert.Equal(t, http.StatusCreated, w1.Code)

		var user1Response models.AuthResponse
		err := json.Unmarshal(w1.Body.Bytes(), &user1Response)
		assert.NoError(t, err)
		defer CleanupTestUser(user1Response.User.ID)

		registerReq2 := models.RegisterRequest{
			Email:    fmt.Sprintf("testlogout2_%s@example.com", uuid.New().String()[:8]),
			Password: "testpassword123",
			Name:     "Test Logout User 2",
			Age:      intPtr(25),
		}

		reqBody, _ = json.Marshal(registerReq2)
		w2 := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/auth/register", "", reqBody)
		assert.Equal(t, http.StatusCreated, w2.Code)

		var user2Response models.AuthResponse
		err = json.Unmarshal(w2.Body.Bytes(), &user2Response)
		assert.NoError(t, err)
		defer CleanupTestUser(user2Response.User.ID)

		// Try to logout user1's refresh token using user2's access token
		logoutReq := models.RefreshTokenRequest{
			RefreshToken: user1Response.RefreshToken,
		}

		reqBody, _ = json.Marshal(logoutReq)
		w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/auth/logout", user2Response.AccessToken, reqBody)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Parse error response manually since ParseJSONResponse expects 200/201
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Should parse error response JSON")

		errorMsg, exists := response["error"]
		assert.True(t, exists, "Response should contain error message")
		assert.Contains(t, errorMsg, "Token does not belong to authenticated user", "Should indicate token ownership mismatch")
	})

	_ = userID // Use userID to avoid unused variable warning
}

// Helper functions for creating pointers to primitive types

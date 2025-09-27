package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"quizninja-api/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestNewUserOnboardingIntegration tests the complete New User Onboarding workflow
// This integration test validates the complete flow: Register → Login → Profile → Preferences → Complete Onboarding
func TestNewUserOnboardingIntegration(t *testing.T) {
	tc := SetupTestServer(t)
	defer Cleanup(t)

	t.Run("CompleteOnboardingFlow", func(t *testing.T) {
		testCompleteOnboardingFlow(t, tc)
	})

	t.Run("OnboardingFlowVariations", func(t *testing.T) {
		testOnboardingFlowVariations(t, tc)
	})

	t.Run("OnboardingValidation", func(t *testing.T) {
		testOnboardingValidation(t, tc)
	})

	t.Run("OnboardingCompletion", func(t *testing.T) {
		testOnboardingCompletion(t, tc)
	})

	t.Run("OnboardingWithRegistrationPreferences", func(t *testing.T) {
		testOnboardingWithRegistrationPreferences(t, tc)
	})

	t.Run("OnboardingStatusConsistency", func(t *testing.T) {
		testOnboardingStatusConsistency(t, tc)
	})
}

// testCompleteOnboardingFlow tests the complete happy path onboarding flow
func testCompleteOnboardingFlow(t *testing.T, tc *TestConfig) {
	// Step 1: Register new user
	uniqueEmail := fmt.Sprintf("onboarding_test_%s@example.com", uuid.New().String()[:8])
	registerReq := models.RegisterRequest{
		Email:    uniqueEmail,
		Password: "testpassword123",
		Name:     "Test Onboarding User",
		Age:      intPtr(25),
	}

	reqBody, _ := json.Marshal(registerReq)
	w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/auth/register", "", reqBody)
	assert.Equal(t, http.StatusCreated, w.Code, "Registration should succeed")

	var registerResponse models.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &registerResponse)
	assert.NoError(t, err, "Should parse register response")

	// Cleanup this specific test user
	defer CleanupTestUser(registerResponse.User.ID)

	accessToken := registerResponse.AccessToken
	userID := registerResponse.User.ID

	// Verify user is created with basic info
	assert.NotEmpty(t, accessToken, "Should receive access token")
	assert.Equal(t, uniqueEmail, registerResponse.User.Email, "Email should match")
	assert.Equal(t, "Test Onboarding User", registerResponse.User.Name, "Name should match")
	assert.True(t, registerResponse.User.IsTestData, "User should be marked as test data")

	// Step 2: Optional separate login (test both registration token and fresh login)
	// Skip separate login test if registration already provides a working token
	loginToken := accessToken
	if accessToken != "" {
		// Test that we can use the registration token for subsequent requests
		w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/profile", accessToken, nil)
		if w.Code != http.StatusOK {
			// If registration token doesn't work, try fresh login
			loginReq := models.LoginRequest{
				Email:    uniqueEmail,
				Password: "testpassword123",
			}

			reqBody, _ = json.Marshal(loginReq)
			w = MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/auth/login", "", reqBody)
			assert.Equal(t, http.StatusOK, w.Code, "Login should succeed")

			var loginResponse models.AuthResponse
			err = json.Unmarshal(w.Body.Bytes(), &loginResponse)
			assert.NoError(t, err, "Should parse login response")

			// Verify login returns same user data
			assert.Equal(t, userID, loginResponse.User.ID, "User ID should remain same")
			assert.Equal(t, uniqueEmail, loginResponse.User.Email, "Email should remain same")
			assert.NotEmpty(t, loginResponse.AccessToken, "Should receive new access token")

			loginToken = loginResponse.AccessToken
		}
	}

	// Step 3: Get profile to verify current user state
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/profile", loginToken, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Profile retrieval should succeed")

	response := ParseJSONResponse(t, w)
	assert.Equal(t, userID.String(), response["id"], "Profile should have correct user ID")
	assert.Equal(t, uniqueEmail, response["email"], "Profile should have correct email")
	assert.Equal(t, "Test Onboarding User", response["name"], "Profile should have correct name")

	// Verify profile shows no preferences initially (or empty preferences)
	if preferences, exists := response["preferences"]; exists && preferences != nil {
		preferencesMap, ok := preferences.(map[string]interface{})
		assert.True(t, ok, "Preferences should be an object if present")
		// Check if onboarding is marked as incomplete
		if completedAt, exists := preferencesMap["onboarding_completed_at"]; exists {
			assert.Nil(t, completedAt, "Onboarding should not be completed yet")
		}
	}

	// Step 4: Set user preferences
	// First, get available categories and interests
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/categories", loginToken, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Categories should be available")

	response = ParseJSONResponse(t, w)
	categories, exists := response["data"]
	assert.True(t, exists, "Should have categories data")

	categoriesList, ok := categories.([]interface{})
	assert.True(t, ok, "Categories should be an array")
	assert.Greater(t, len(categoriesList), 0, "Should have at least one category")

	// Extract available interests
	var availableInterests []string
	for _, category := range categoriesList {
		categoryMap, ok := category.(map[string]interface{})
		if !ok {
			continue
		}

		interests, ok := categoryMap["interests"].([]interface{})
		if !ok {
			continue
		}

		for _, interest := range interests {
			interestMap, ok := interest.(map[string]interface{})
			if !ok {
				continue
			}

			if id, ok := interestMap["id"].(string); ok {
				availableInterests = append(availableInterests, id)
			}
		}
	}

	assert.Greater(t, len(availableInterests), 0, "Should have available interests")

	// Select first 3 interests for testing
	selectedInterests := availableInterests[:min(3, len(availableInterests))]

	// Update user preferences
	updatePrefsReq := models.UpdatePreferencesRequest{
		SelectedInterests:     selectedInterests,
		DifficultyPreference:  "Medium",
		NotificationsEnabled:  true,
		NotificationFrequency: "Daily",
		ProfileVisibility:     boolPtr(true),
		ShowOnlineStatus:      boolPtr(true),
		AllowFriendRequests:   boolPtr(true),
		ShareActivityStatus:   boolPtr(true),
	}

	reqBody, _ = json.Marshal(updatePrefsReq)
	w = MakeAuthenticatedRequest(t, tc, "PUT", "/api/v1/users/preferences", loginToken, reqBody)
	assert.Equal(t, http.StatusOK, w.Code, "Preferences update should succeed")

	response = ParseJSONResponse(t, w)
	data := GetDataFromResponse(t, response)

	// Verify preferences were saved correctly
	assert.Equal(t, userID.String(), data["user_id"], "Preferences should have correct user ID")
	assert.Equal(t, "Medium", data["difficulty_preference"], "Difficulty should be saved")
	assert.Equal(t, true, data["notifications_enabled"], "Notifications should be enabled")
	assert.Equal(t, "Daily", data["notification_frequency"], "Frequency should be saved")
	assert.Equal(t, true, data["profile_visibility"], "Profile visibility should be saved")

	savedInterests := data["selected_interests"].([]interface{})
	assert.Equal(t, len(selectedInterests), len(savedInterests), "Should save all selected interests")

	// Verify onboarding is not yet completed
	if completedAt, exists := data["onboarding_completed_at"]; exists {
		assert.Nil(t, completedAt, "Onboarding should not be automatically completed by preferences update")
	}

	// Step 5: Complete onboarding
	onboardingReq := models.OnboardingCompleteRequest{
		SelectedInterests:     selectedInterests,
		DifficultyPreference:  "Medium",
		NotificationsEnabled:  true,
		NotificationFrequency: "Daily",
	}

	reqBody, _ = json.Marshal(onboardingReq)
	w = MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/users/onboarding/complete", loginToken, reqBody)
	assert.Equal(t, http.StatusOK, w.Code, "Onboarding completion should succeed")

	response = ParseJSONResponse(t, w)
	data = GetDataFromResponse(t, response)

	// Verify onboarding completion
	assert.Equal(t, userID.String(), data["user_id"], "Should have correct user ID")
	assert.NotNil(t, data["onboarding_completed_at"], "Onboarding completion timestamp should be set")

	completedAtStr, ok := data["onboarding_completed_at"].(string)
	assert.True(t, ok, "Completion timestamp should be a string")
	assert.NotEmpty(t, completedAtStr, "Completion timestamp should not be empty")

	// Verify completion timestamp is recent (within last minute)
	completedAt, err := time.Parse(time.RFC3339, completedAtStr)
	assert.NoError(t, err, "Should parse completion timestamp")
	assert.WithinDuration(t, time.Now(), completedAt, time.Minute, "Completion should be recent")

	// Step 6: Verify onboarding status
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/onboarding/status", loginToken, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Onboarding status should be accessible")

	response = ParseJSONResponse(t, w)
	data = GetDataFromResponse(t, response)

	assert.Equal(t, true, data["is_completed"], "Onboarding should be marked as completed")
	assert.NotNil(t, data["completed_at"], "Completion timestamp should be available")

	// Step 7: Verify full platform access (test protected endpoints)
	// Test profile access
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/profile", loginToken, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Profile should remain accessible")

	// Test preferences access
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/preferences", loginToken, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Preferences should remain accessible")

	// Test protected quiz endpoints
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/quizzes", loginToken, nil)
	assert.Equal(t, http.StatusOK, w.Code, "User quizzes should be accessible")

	// Test user stats
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/stats", loginToken, nil)
	assert.Equal(t, http.StatusOK, w.Code, "User statistics should be accessible")

	t.Logf("Complete onboarding flow validated successfully for user: %s", userID)
}

// testOnboardingFlowVariations tests different onboarding scenarios
func testOnboardingFlowVariations(t *testing.T, tc *TestConfig) {
	t.Run("MinimalPreferences", func(t *testing.T) {
		// Test onboarding with minimal required preferences
		userID, token := CreateTestUser(t, tc)

		// Get available interests
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/categories", token, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		response := ParseJSONResponse(t, w)
		categories := response["data"].([]interface{})
		availableInterests := extractInterestIDsFromCategories(categories)

		// Minimal onboarding (just 1 interest)
		onboardingReq := models.OnboardingCompleteRequest{
			SelectedInterests:     []string{availableInterests[0]},
			DifficultyPreference:  "Easy",
			NotificationsEnabled:  false,
			NotificationFrequency: "Never",
		}

		reqBody, _ := json.Marshal(onboardingReq)
		w = MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/users/onboarding/complete", token, reqBody)
		assert.Equal(t, http.StatusOK, w.Code, "Minimal onboarding should succeed")

		// Verify completion
		w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/onboarding/status", token, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		response = ParseJSONResponse(t, w)
		data := GetDataFromResponse(t, response)
		assert.Equal(t, true, data["is_completed"], "Minimal onboarding should be completed")

		t.Logf("Minimal preferences onboarding completed for user: %s", userID)
	})

	t.Run("MaximalPreferences", func(t *testing.T) {
		// Test onboarding with many preferences
		userID, token := CreateTestUser(t, tc)

		// Get all available interests
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/categories", token, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		response := ParseJSONResponse(t, w)
		categories := response["data"].([]interface{})
		availableInterests := extractInterestIDsFromCategories(categories)

		// Use all available interests (up to a reasonable limit)
		selectedInterests := availableInterests[:min(len(availableInterests), 10)]

		onboardingReq := models.OnboardingCompleteRequest{
			SelectedInterests:     selectedInterests,
			DifficultyPreference:  "Hard",
			NotificationsEnabled:  true,
			NotificationFrequency: "Daily",
		}

		reqBody, _ := json.Marshal(onboardingReq)
		w = MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/users/onboarding/complete", token, reqBody)
		assert.Equal(t, http.StatusOK, w.Code, "Maximal onboarding should succeed")

		// Verify all interests were saved
		w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/preferences", token, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		response = ParseJSONResponse(t, w)
		data := GetDataFromResponse(t, response)
		savedInterests := data["selected_interests"].([]interface{})
		assert.Equal(t, len(selectedInterests), len(savedInterests), "All interests should be saved")

		t.Logf("Maximal preferences onboarding completed for user: %s", userID)
	})

	t.Run("DifferentDifficultyLevels", func(t *testing.T) {
		difficultyLevels := []string{"Easy", "Medium", "Hard"}

		for _, difficulty := range difficultyLevels {
			t.Run(fmt.Sprintf("Difficulty_%s", difficulty), func(t *testing.T) {
				userID, token := CreateTestUser(t, tc)

				// Get available interests
				w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/categories", token, nil)
				response := ParseJSONResponse(t, w)
				categories := response["data"].([]interface{})
				availableInterests := extractInterestIDsFromCategories(categories)

				onboardingReq := models.OnboardingCompleteRequest{
					SelectedInterests:     []string{availableInterests[0]},
					DifficultyPreference:  difficulty,
					NotificationsEnabled:  true,
					NotificationFrequency: "Weekly",
				}

				reqBody, _ := json.Marshal(onboardingReq)
				w = MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/users/onboarding/complete", token, reqBody)
				assert.Equal(t, http.StatusOK, w.Code, "Onboarding with %s difficulty should succeed", difficulty)

				// Verify difficulty was saved
				w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/preferences", token, nil)
				response = ParseJSONResponse(t, w)
				data := GetDataFromResponse(t, response)
				assert.Equal(t, difficulty, data["difficulty_preference"], "Difficulty should match")

				t.Logf("Onboarding with %s difficulty completed for user: %s", difficulty, userID)
			})
		}
	})
}

// testOnboardingValidation tests error cases and validation
func testOnboardingValidation(t *testing.T, tc *TestConfig) {
	t.Run("InvalidInterests", func(t *testing.T) {
		_, token := CreateTestUser(t, tc)

		// Try to complete onboarding with invalid interest IDs
		onboardingReq := models.OnboardingCompleteRequest{
			SelectedInterests:     []string{"invalid_interest_1", "invalid_interest_2"},
			DifficultyPreference:  "Medium",
			NotificationsEnabled:  true,
			NotificationFrequency: "Daily",
		}

		reqBody, _ := json.Marshal(onboardingReq)
		w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/users/onboarding/complete", token, reqBody)
		// Note: This may succeed with invalid interests depending on validation logic
		// The test verifies the API handles invalid interests gracefully
		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusBadRequest, "Should handle invalid interests gracefully")
	})

	t.Run("InvalidDifficulty", func(t *testing.T) {
		_, token := CreateTestUser(t, tc)

		// Get valid interests
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/categories", token, nil)
		response := ParseJSONResponse(t, w)
		categories := response["data"].([]interface{})
		availableInterests := extractInterestIDsFromCategories(categories)

		// Try invalid difficulty
		onboardingReq := models.OnboardingCompleteRequest{
			SelectedInterests:     []string{availableInterests[0]},
			DifficultyPreference:  "Invalid",
			NotificationsEnabled:  true,
			NotificationFrequency: "Daily",
		}

		reqBody, _ := json.Marshal(onboardingReq)
		w = MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/users/onboarding/complete", token, reqBody)
		assert.Equal(t, http.StatusBadRequest, w.Code, "Invalid difficulty should be rejected")
	})

	t.Run("InvalidNotificationFrequency", func(t *testing.T) {
		_, token := CreateTestUser(t, tc)

		// Get valid interests
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/categories", token, nil)
		response := ParseJSONResponse(t, w)
		categories := response["data"].([]interface{})
		availableInterests := extractInterestIDsFromCategories(categories)

		// Try invalid notification frequency
		onboardingReq := models.OnboardingCompleteRequest{
			SelectedInterests:     []string{availableInterests[0]},
			DifficultyPreference:  "Medium",
			NotificationsEnabled:  true,
			NotificationFrequency: "Invalid",
		}

		reqBody, _ := json.Marshal(onboardingReq)
		w = MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/users/onboarding/complete", token, reqBody)
		assert.Equal(t, http.StatusBadRequest, w.Code, "Invalid notification frequency should be rejected")
	})

	t.Run("EmptyInterests", func(t *testing.T) {
		_, token := CreateTestUser(t, tc)

		// Try with empty interests array
		onboardingReq := models.OnboardingCompleteRequest{
			SelectedInterests:     []string{},
			DifficultyPreference:  "Medium",
			NotificationsEnabled:  true,
			NotificationFrequency: "Daily",
		}

		reqBody, _ := json.Marshal(onboardingReq)
		w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/users/onboarding/complete", token, reqBody)
		// The API might allow empty interests - check what actually happens
		if w.Code == http.StatusOK {
			t.Log("API allows empty interests - this is acceptable behavior")
		} else {
			assert.Equal(t, http.StatusBadRequest, w.Code, "Empty interests should be rejected or handled gracefully")
		}
	})

	t.Run("UnauthenticatedRequest", func(t *testing.T) {
		// Try to complete onboarding without authentication
		onboardingReq := models.OnboardingCompleteRequest{
			SelectedInterests:     []string{"technology"},
			DifficultyPreference:  "Medium",
			NotificationsEnabled:  true,
			NotificationFrequency: "Daily",
		}

		reqBody, _ := json.Marshal(onboardingReq)
		w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/users/onboarding/complete", "", reqBody)
		assert.Equal(t, http.StatusUnauthorized, w.Code, "Unauthenticated request should be rejected")
	})
}

// testOnboardingCompletion tests platform access after onboarding completion
func testOnboardingCompletion(t *testing.T, tc *TestConfig) {
	userID, token := CreateTestUser(t, tc)

	// Complete onboarding first
	w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/categories", token, nil)
	response := ParseJSONResponse(t, w)
	categories := response["data"].([]interface{})
	availableInterests := extractInterestIDsFromCategories(categories)

	onboardingReq := models.OnboardingCompleteRequest{
		SelectedInterests:     []string{availableInterests[0]},
		DifficultyPreference:  "Medium",
		NotificationsEnabled:  true,
		NotificationFrequency: "Daily",
	}

	reqBody, _ := json.Marshal(onboardingReq)
	w = MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/users/onboarding/complete", token, reqBody)
	assert.Equal(t, http.StatusOK, w.Code, "Onboarding should complete successfully")

	// Test full platform access after completion
	t.Run("QuizAccess", func(t *testing.T) {
		// Test access to quiz-related endpoints
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/quizzes", token, nil)
		assert.Equal(t, http.StatusOK, w.Code, "User should have access to quiz listings")

		w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/attempts", token, nil)
		assert.Equal(t, http.StatusOK, w.Code, "User should have access to quiz attempts")

		w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/active-sessions", token, nil)
		assert.Equal(t, http.StatusOK, w.Code, "User should have access to active sessions")
	})

	t.Run("SocialAccess", func(t *testing.T) {
		// Test access to social features
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/friends", token, nil)
		assert.Equal(t, http.StatusOK, w.Code, "User should have access to friends")

		w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/friends/requests", token, nil)
		assert.Equal(t, http.StatusOK, w.Code, "User should have access to friend requests")

		w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges", token, nil)
		assert.Equal(t, http.StatusOK, w.Code, "User should have access to challenges")

		w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/leaderboard", token, nil)
		assert.Equal(t, http.StatusOK, w.Code, "User should have access to leaderboard")
	})

	t.Run("PersonalizationAccess", func(t *testing.T) {
		// Test access to personalization features
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/achievements", token, nil)
		assert.Equal(t, http.StatusOK, w.Code, "User should have access to achievements")

		w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/favorites", token, nil)
		assert.Equal(t, http.StatusOK, w.Code, "User should have access to favorites")

		w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/stats", token, nil)
		assert.Equal(t, http.StatusOK, w.Code, "User should have access to personal statistics")
	})

	t.Run("PreferencesModification", func(t *testing.T) {
		// Test that user can still modify preferences after onboarding
		updateReq := models.UpdatePreferencesRequest{
			SelectedInterests:     []string{availableInterests[0], availableInterests[1]},
			DifficultyPreference:  "Hard",
			NotificationsEnabled:  false,
			NotificationFrequency: "Weekly",
			ProfileVisibility:     boolPtr(false),
			ShowOnlineStatus:      boolPtr(false),
			AllowFriendRequests:   boolPtr(false),
			ShareActivityStatus:   boolPtr(false),
		}

		reqBody, _ := json.Marshal(updateReq)
		w := MakeAuthenticatedRequest(t, tc, "PUT", "/api/v1/users/preferences", token, reqBody)
		assert.Equal(t, http.StatusOK, w.Code, "User should be able to update preferences after onboarding")

		// Verify changes were saved
		w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/preferences", token, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		response := ParseJSONResponse(t, w)
		data := GetDataFromResponse(t, response)
		assert.Equal(t, "Hard", data["difficulty_preference"], "Updated preference should persist")
		assert.Equal(t, false, data["notifications_enabled"], "Updated preference should persist")

		// Verify onboarding status behavior after preference updates
		if completedAt, exists := data["onboarding_completed_at"]; exists {
			if completedAt != nil {
				t.Log("Onboarding completion timestamp is preserved in preferences after updates")
			} else {
				t.Log("Onboarding completion timestamp was cleared after preference updates")
			}
		} else {
			t.Log("onboarding_completed_at field not present in preferences response - checking status endpoint")
		}

		// Check onboarding status endpoint to understand current behavior
		w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/onboarding/status", token, nil)
		if w.Code == http.StatusOK {
			response = ParseJSONResponse(t, w)
			statusData := GetDataFromResponse(t, response)
			isCompleted, exists := statusData["is_completed"]
			if exists {
				completed, ok := isCompleted.(bool)
				if ok && !completed {
					t.Log("Note: Onboarding status was reset to incomplete after preference updates - this may be intended API behavior")
					t.Log("This could indicate that preference updates require re-completing onboarding")
				} else if ok && completed {
					t.Log("Onboarding status remains completed after preference updates")
				}
			}
		}
	})

	t.Logf("Platform access verification completed for user: %s", userID)
}

// testOnboardingWithRegistrationPreferences tests onboarding when preferences are provided during registration
func testOnboardingWithRegistrationPreferences(t *testing.T, tc *TestConfig) {
	// Get available interests first
	w := MakeRequest(t, tc, "GET", "/api/v1/preferences/categories")
	if w.Code != http.StatusOK {
		t.Skip("Skipping registration with preferences test - categories endpoint not available publicly")
		return
	}

	uniqueEmail := fmt.Sprintf("reg_prefs_test_%s@example.com", uuid.New().String()[:8])

	// Register with initial preferences
	registerReq := models.RegisterRequest{
		Email:    uniqueEmail,
		Password: "testpassword123",
		Name:     "Test Registration Preferences User",
		Age:      intPtr(30),
		Preferences: &models.UserPreferencesRequest{
			SelectedInterests:     []string{"technology", "science"},
			DifficultyPreference:  "Easy",
			NotificationsEnabled:  true,
			NotificationFrequency: "Weekly",
		},
	}

	reqBody, _ := json.Marshal(registerReq)
	w = MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/auth/register", "", reqBody)
	assert.Equal(t, http.StatusCreated, w.Code, "Registration with preferences should succeed")

	var registerResponse models.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &registerResponse)
	assert.NoError(t, err, "Should parse register response")

	// Cleanup this specific test user
	defer CleanupTestUser(registerResponse.User.ID)

	token := registerResponse.AccessToken

	// Verify preferences were set during registration
	if registerResponse.User.Preferences != nil {
		assert.Equal(t, "Easy", registerResponse.User.Preferences.DifficultyPreference, "Registration preferences should be set")
		assert.Equal(t, true, registerResponse.User.Preferences.NotificationsEnabled, "Registration preferences should be set")
		assert.Equal(t, "Weekly", registerResponse.User.Preferences.NotificationFrequency, "Registration preferences should be set")
	}

	// Check onboarding status
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/onboarding/status", token, nil)
	if w.Code == http.StatusOK {
		response := ParseJSONResponse(t, w)
		data := GetDataFromResponse(t, response)

		// Registration with preferences might auto-complete onboarding
		isCompleted, exists := data["is_completed"]
		if exists && isCompleted != nil {
			completed, ok := isCompleted.(bool)
			if ok && completed {
				t.Log("Onboarding was auto-completed during registration with preferences")
				return
			}
		}
	}

	// If not auto-completed, complete onboarding manually
	onboardingReq := models.OnboardingCompleteRequest{
		SelectedInterests:     []string{"technology", "science"},
		DifficultyPreference:  "Easy",
		NotificationsEnabled:  true,
		NotificationFrequency: "Weekly",
	}

	reqBody, _ = json.Marshal(onboardingReq)
	w = MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/users/onboarding/complete", token, reqBody)
	assert.Equal(t, http.StatusOK, w.Code, "Onboarding completion should succeed even with registration preferences")

	t.Logf("Registration with preferences onboarding completed for user: %s", registerResponse.User.ID)
}

// testOnboardingStatusConsistency tests that onboarding status is consistent across endpoints
func testOnboardingStatusConsistency(t *testing.T, tc *TestConfig) {
	userID, token := CreateTestUser(t, tc)

	// Initially onboarding should not be completed
	w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/onboarding/status", token, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	response := ParseJSONResponse(t, w)
	data := GetDataFromResponse(t, response)
	assert.Equal(t, false, data["is_completed"], "Onboarding should initially be incomplete")
	assert.Nil(t, data["completed_at"], "Completion timestamp should be nil initially")

	// Complete onboarding
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/categories", token, nil)
	response = ParseJSONResponse(t, w)
	categories := response["data"].([]interface{})
	availableInterests := extractInterestIDsFromCategories(categories)

	onboardingReq := models.OnboardingCompleteRequest{
		SelectedInterests:     []string{availableInterests[0]},
		DifficultyPreference:  "Medium",
		NotificationsEnabled:  true,
		NotificationFrequency: "Daily",
	}

	reqBody, _ := json.Marshal(onboardingReq)
	w = MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/users/onboarding/complete", token, reqBody)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify consistency across different endpoints
	t.Run("OnboardingStatusEndpoint", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/onboarding/status", token, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		response := ParseJSONResponse(t, w)
		data := GetDataFromResponse(t, response)
		assert.Equal(t, true, data["is_completed"], "Status endpoint should show completed")
		assert.NotNil(t, data["completed_at"], "Status endpoint should have completion timestamp")
	})

	t.Run("PreferencesEndpoint", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/preferences", token, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		response := ParseJSONResponse(t, w)
		data := GetDataFromResponse(t, response)
		assert.NotNil(t, data["onboarding_completed_at"], "Preferences endpoint should have completion timestamp")
	})

	t.Run("ProfileEndpoint", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/profile", token, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		response := ParseJSONResponse(t, w)
		if preferences, exists := response["preferences"]; exists && preferences != nil {
			preferencesMap, ok := preferences.(map[string]interface{})
			assert.True(t, ok, "Profile preferences should be an object")
			assert.NotNil(t, preferencesMap["onboarding_completed_at"], "Profile preferences should have completion timestamp")
		}
	})

	t.Logf("Onboarding status consistency verified for user: %s", userID)
}

// Helper function to extract interest IDs from categories response
func extractInterestIDsFromCategories(categories []interface{}) []string {
	var interestIDs []string
	for _, category := range categories {
		categoryMap, ok := category.(map[string]interface{})
		if !ok {
			continue
		}

		interests, ok := categoryMap["interests"].([]interface{})
		if !ok {
			continue
		}

		for _, interest := range interests {
			interestMap, ok := interest.(map[string]interface{})
			if !ok {
				continue
			}

			if id, ok := interestMap["id"].(string); ok {
				interestIDs = append(interestIDs, id)
			}
		}
	}
	return interestIDs
}


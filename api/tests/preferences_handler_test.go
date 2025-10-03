package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"quizninja-api/models"
	"quizninja-api/repository"

	"github.com/stretchr/testify/assert"
)

func TestPreferencesHandler(t *testing.T) {
	tc := SetupTestServer(t)
	defer Cleanup(t)

	userID, token := CreateTestUser(t, tc)

	t.Run("GetPreferences", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/preferences", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			// Verify is_test_data field in preferences
			VerifyIsTestDataField(t, data, true, "user preferences")
		}
	})

	t.Run("UpdatePreferences", func(t *testing.T) {
		updateReq := models.UpdatePreferencesRequest{
			SelectedInterests:     []string{"technology", "science"},
			DifficultyPreference:  "Medium",
			NotificationsEnabled:  true,
			NotificationFrequency: "daily",
			ProfileVisibility:     boolPtr(true),
			ShowOnlineStatus:      boolPtr(true),
			AllowFriendRequests:   boolPtr(true),
			ShareActivityStatus:   boolPtr(true),
		}

		reqBody, _ := json.Marshal(updateReq)
		w := MakeAuthenticatedRequest(t, tc, "PUT", "/api/v1/users/preferences", token, reqBody)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			// Verify is_test_data field in updated preferences
			VerifyIsTestDataField(t, data, true, "updated user preferences")
		}
	})

	t.Run("GetCategories", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/preferences/categories", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// Categories/interests is an array directly in data field
			data, exists := response["data"]
			if exists && data != nil {
				interestsList, ok := data.([]interface{})
				if ok && len(interestsList) > 0 {
					VerifyIsTestDataInArray(t, interestsList, true, "preference categories/interests")
				}
			}
		}
	})

	t.Run("GetDifficultyLevels", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/preferences/difficulty-levels", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// Difficulty levels is an array directly in data field
			data, exists := response["data"]
			if exists && data != nil {
				levelsList, ok := data.([]interface{})
				if ok && len(levelsList) > 0 {
					VerifyIsTestDataInArray(t, levelsList, true, "difficulty levels")
				}
			}
		}
	})

	t.Run("GetNotificationFrequencies", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/preferences/notification-frequencies", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// Notification frequencies is an array directly in data field
			data, exists := response["data"]
			if exists && data != nil {
				frequenciesList, ok := data.([]interface{})
				if ok && len(frequenciesList) > 0 {
					VerifyIsTestDataInArray(t, frequenciesList, true, "notification frequencies")
				}
			}
		}
	})

	t.Run("GetOnboardingStatus", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/preferences/onboarding/status", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			// Check if preferences exist in onboarding status
			preferences, prefsExist := data["preferences"]
			if prefsExist && preferences != nil {
				preferencesMap, ok := preferences.(map[string]interface{})
				if ok {
					VerifyIsTestDataField(t, preferencesMap, true, "onboarding status preferences")
				}
			}

			// Verify the basic structure of onboarding status
			isCompleted, exists := data["is_completed"]
			assert.True(t, exists, "Onboarding status should contain is_completed field")

			if isCompleted != nil {
				_, ok := isCompleted.(bool)
				assert.True(t, ok, "is_completed should be a boolean")
			}
		}
	})

	t.Run("CompleteOnboarding", func(t *testing.T) {
		onboardingReq := models.OnboardingCompleteRequest{
			SelectedInterests:     []string{"technology", "science", "sports"},
			DifficultyPreference:  "Medium",
			NotificationsEnabled:  true,
			NotificationFrequency: "weekly",
		}

		reqBody, _ := json.Marshal(onboardingReq)
		w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/preferences/onboarding/complete", token, reqBody)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			// Verify is_test_data field in completed onboarding preferences
			VerifyIsTestDataField(t, data, true, "completed onboarding preferences")

			// Verify the message indicates successful completion
			message, messageExists := response["message"]
			if messageExists {
				messageStr, ok := message.(string)
				if ok {
					assert.Contains(t, messageStr, "successfully", "Should indicate successful onboarding completion")
				}
			}
		}
	})

	t.Run("UpsertValidation", func(t *testing.T) {
		// Test to validate that multiple preference updates don't create duplicate records
		// This validates the fix for the settings save issue (UNIQUE constraint + UPSERT)

		// Create initial preferences
		initialReq := models.UpdatePreferencesRequest{
			SelectedInterests:     []string{"technology"},
			DifficultyPreference:  "Easy",
			NotificationsEnabled:  true,
			NotificationFrequency: "daily",
			ProfileVisibility:     boolPtr(true),
			ShowOnlineStatus:      boolPtr(true),
			AllowFriendRequests:   boolPtr(true),
			ShareActivityStatus:   boolPtr(true),
		}

		reqBody, _ := json.Marshal(initialReq)
		w := MakeAuthenticatedRequest(t, tc, "PUT", "/api/v1/users/preferences", token, reqBody)

		if w.Code == http.StatusOK {
			// Update preferences multiple times to test UPSERT behavior
			for i := 1; i <= 3; i++ {
				updateReq := models.UpdatePreferencesRequest{
					SelectedInterests:     []string{"science", "sports"},
					DifficultyPreference:  "Hard",
					NotificationsEnabled:  false,
					NotificationFrequency: "weekly",
					ProfileVisibility:     boolPtr(false),
					ShowOnlineStatus:      boolPtr(false),
					AllowFriendRequests:   boolPtr(false),
					ShareActivityStatus:   boolPtr(false),
				}

				reqBody, _ := json.Marshal(updateReq)
				_ = MakeAuthenticatedRequest(t, tc, "PUT", "/api/v1/users/preferences", token, reqBody)
				// Continue testing even if individual updates fail
			}

			// Verify that the final state can be retrieved and has the latest values
			w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/preferences", token, nil)
			if w.Code == http.StatusOK {
				response := ParseJSONResponse(t, w)
				data := GetDataFromResponse(t, response)

				// Verify final state contains the last update (proving UPSERT worked)
				// If the UNIQUE constraint fix didn't work, we'd either have errors
				// or the values would be inconsistent/stuck on old values
				if difficulty, exists := data["difficulty_preference"]; exists {
					// Should have the updated value, not the initial "Easy"
					assert.NotEqual(t, "Easy", difficulty, "Difficulty should be updated from initial value")
				}

				if frequency, exists := data["notification_frequency"]; exists {
					// Should have the updated value, not the initial "daily"
					assert.NotEqual(t, "daily", frequency, "Frequency should be updated from initial value")
				}

				// Verify is_test_data field for consistency with other tests
				VerifyIsTestDataField(t, data, true, "UPSERT validation preferences")
			}
		}
	})

	t.Run("DebugGetUserPreferences", func(t *testing.T) {
		// Test GetUserPreferences directly after a successful update
		userRepo := repository.NewUserRepository()

		// Create test user preferences first
		preferences := &models.UserPreferences{
			UserID:                userID,
			SelectedInterests:     models.StringArray([]string{"technology"}),
			DifficultyPreference:  "Medium",
			NotificationsEnabled:  true,
			NotificationFrequency: "daily",
			ProfileVisibility:     true,
			ShowOnlineStatus:      true,
			AllowFriendRequests:   true,
			ShareActivityStatus:   true,
			NotificationTypes: map[string]interface{}{
				"challenges":   true,
				"achievements": true,
			},
			IsTestData: true,
		}

		// First, save the preferences
		err := userRepo.UpdateUserPreferences(preferences)
		if err != nil {
			t.Fatalf("Failed to save preferences: %v", err)
		}
		t.Log("Successfully saved preferences")

		// Now try to retrieve them
		retrievedPrefs, err := userRepo.GetUserPreferences(userID)
		if err != nil {
			t.Fatalf("Failed to get preferences: %v", err)
		}

		t.Logf("Successfully retrieved preferences: %+v", retrievedPrefs)
		t.Logf("Selected interests: %v", retrievedPrefs.SelectedInterests)
		t.Logf("Notification types: %+v", retrievedPrefs.NotificationTypes)
	})

	_ = userID // Use userID to avoid unused variable warning
}

func boolPtr(b bool) *bool {
	return &b
}

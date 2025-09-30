package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"quizninja-api/models"

	"github.com/stretchr/testify/assert"
)

// TestCategoriesPreferencesIntegration tests the complete Categories & Preferences flow
// This simulates the Settings → Preferences → Interest Categories workflow
func TestCategoriesPreferencesIntegration(t *testing.T) {
	tc := SetupTestServer(t)
	defer Cleanup(t)

	t.Run("NewUserFlow", func(t *testing.T) {
		testCategoriesPreferencesFlowNewUser(t, tc)
	})

	t.Run("ExistingUserFlow", func(t *testing.T) {
		testCategoriesPreferencesFlowExistingUser(t, tc)
	})

	t.Run("DataConsistency", func(t *testing.T) {
		testCategoriesCrossEndpointConsistency(t, tc)
	})

	t.Run("SettingsWorkflow", func(t *testing.T) {
		testSettingsScreenWorkflow(t, tc)
	})

	t.Run("CategoriesStructure", func(t *testing.T) {
		testDatabaseDrivenCategoriesStructure(t, tc)
	})
}

// Helper functions for assertions
func assertCategoryStructure(t *testing.T, category map[string]interface{}) {
	// Verify required fields for category
	assert.Contains(t, category, "id", "Category should have 'id' field")
	assert.Contains(t, category, "name", "Category should have 'name' field")
	assert.Contains(t, category, "display_name", "Category should have 'display_name' field")
	assert.Contains(t, category, "description", "Category should have 'description' field")
	assert.Contains(t, category, "interests", "Category should have 'interests' field")

	// Verify interests array
	interests, ok := category["interests"].([]interface{})
	assert.True(t, ok, "Interests should be an array")
	assert.Greater(t, len(interests), 0, "Category should have at least one interest")

	// Verify each interest has required fields
	for _, interest := range interests {
		interestMap, ok := interest.(map[string]interface{})
		assert.True(t, ok, "Each interest should be an object")

		assert.Contains(t, interestMap, "id", "Interest should have 'id' field")
		assert.Contains(t, interestMap, "name", "Interest should have 'name' field")
		assert.Contains(t, interestMap, "display_name", "Interest should have 'display_name' field")
		assert.Contains(t, interestMap, "description", "Interest should have 'description' field")
		assert.Contains(t, interestMap, "is_test_data", "Interest should have 'is_test_data' field")
	}
}

func assertPreferencesStructure(t *testing.T, preferences map[string]interface{}) {
	// Verify required preference fields
	expectedFields := []string{
		"user_id", "selected_interests", "difficulty_preference",
		"notifications_enabled", "notification_frequency",
		"profile_visibility", "show_online_status",
		"allow_friend_requests", "share_activity_status",
		"created_at", "updated_at",
	}

	for _, field := range expectedFields {
		assert.Contains(t, preferences, field, "Preferences should have '%s' field", field)
	}

	// Verify selected_interests is an array
	selectedInterests, ok := preferences["selected_interests"].([]interface{})
	assert.True(t, ok, "Selected interests should be an array")

	// If interests are selected, verify they are strings
	for _, interest := range selectedInterests {
		_, ok := interest.(string)
		assert.True(t, ok, "Each selected interest should be a string")
	}
}

func extractInterestIDs(categories []interface{}) []string {
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

// testCategoriesPreferencesFlowNewUser tests the complete flow for a new user
func testCategoriesPreferencesFlowNewUser(t *testing.T, tc *TestConfig) {
	// Step 1: Create new user (no existing preferences)
	userID, token := CreateTestUser(t, tc)

	// Step 2: Fetch categories from database-driven endpoint
	w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/categories", token, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Categories endpoint should return 200 OK")

	response := ParseJSONResponse(t, w)
	categories, exists := response["data"]
	assert.True(t, exists, "Response should contain 'categories' field")

	categoriesList, ok := categories.([]interface{})
	assert.True(t, ok, "Categories field should be an array")
	assert.Greater(t, len(categoriesList), 0, "Should have at least one category")

	// Step 3: Verify grouped interests display structure
	expectedCategories := map[string]bool{
		"general": false, "science": false, "sports": false, "entertainment": false,
	}

	for _, category := range categoriesList {
		categoryMap, ok := category.(map[string]interface{})
		assert.True(t, ok, "Each category should be an object")
		assertCategoryStructure(t, categoryMap)

		// Mark expected categories as found
		if categoryID, ok := categoryMap["id"].(string); ok {
			if _, expected := expectedCategories[categoryID]; expected {
				expectedCategories[categoryID] = true
			}
		}
	}

	// Verify all expected categories were found
	for categoryName, found := range expectedCategories {
		assert.True(t, found, "Should find expected category: %s", categoryName)
	}

	// Step 4: Fetch preference options (difficulty, frequencies)
	w = MakeRequest(t, tc, "GET", "/api/v1/preferences/difficulty-levels")
	assert.Equal(t, http.StatusOK, w.Code, "Difficulty levels endpoint should return 200 OK")

	w = MakeRequest(t, tc, "GET", "/api/v1/preferences/notification-frequencies")
	assert.Equal(t, http.StatusOK, w.Code, "Notification frequencies endpoint should return 200 OK")

	// Step 5: Create initial preferences with selected interests from categories
	availableInterests := extractInterestIDs(categoriesList)
	assert.Greater(t, len(availableInterests), 0, "Should have available interests")

	// Select first 3 interests for testing
	selectedInterests := availableInterests[:min(3, len(availableInterests))]

	updateReq := models.UpdatePreferencesRequest{
		SelectedInterests:     selectedInterests,
		DifficultyPreference:  "Medium",
		NotificationsEnabled:  true,
		NotificationFrequency: "Daily",
		ProfileVisibility:     boolPtr(true),
		ShowOnlineStatus:      boolPtr(true),
		AllowFriendRequests:   boolPtr(true),
		ShareActivityStatus:   boolPtr(true),
	}

	reqBody, _ := json.Marshal(updateReq)
	w = MakeAuthenticatedRequest(t, tc, "PUT", "/api/v1/users/preferences", token, reqBody)
	assert.Equal(t, http.StatusOK, w.Code, "Preferences update should return 200 OK")

	// Step 6: Verify preferences saved correctly
	response = ParseJSONResponse(t, w)
	data := GetDataFromResponse(t, response)
	assertPreferencesStructure(t, data)

	// Verify specific values
	savedInterests := data["selected_interests"].([]interface{})
	assert.Equal(t, len(selectedInterests), len(savedInterests), "Should save all selected interests")

	assert.Equal(t, "Medium", data["difficulty_preference"], "Should save difficulty preference")
	assert.Equal(t, true, data["notifications_enabled"], "Should save notifications enabled")
	assert.Equal(t, "Daily", data["notification_frequency"], "Should save notification frequency")

	// Step 7: Fetch preferences back and verify consistency
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/preferences", token, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Get preferences should return 200 OK")

	response = ParseJSONResponse(t, w)
	data = GetDataFromResponse(t, response)
	assertPreferencesStructure(t, data)

	// Verify consistency
	retrievedInterests := data["selected_interests"].([]interface{})
	assert.Equal(t, len(selectedInterests), len(retrievedInterests), "Retrieved interests should match saved")

	// Step 8: Verify profile includes new preferences
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/profile", token, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Profile endpoint should return 200 OK")

	response = ParseJSONResponse(t, w)
	if preferences, exists := response["preferences"]; exists {
		preferencesMap, ok := preferences.(map[string]interface{})
		assert.True(t, ok, "Preferences in profile should be an object")
		assertPreferencesStructure(t, preferencesMap)
	}

	t.Logf("New user flow completed successfully for user: %s", userID)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// testCategoriesPreferencesFlowExistingUser tests updating preferences for an existing user
func testCategoriesPreferencesFlowExistingUser(t *testing.T, tc *TestConfig) {
	// Step 1: Create user with existing preferences
	userID, token := CreateTestUser(t, tc)

	// First create initial preferences
	initialReq := models.UpdatePreferencesRequest{
		SelectedInterests:     []string{"technology", "science"},
		DifficultyPreference:  "Easy",
		NotificationsEnabled:  false,
		NotificationFrequency: "Weekly",
		ProfileVisibility:     boolPtr(false),
		ShowOnlineStatus:      boolPtr(false),
		AllowFriendRequests:   boolPtr(false),
		ShareActivityStatus:   boolPtr(false),
	}

	reqBody, _ := json.Marshal(initialReq)
	w := MakeAuthenticatedRequest(t, tc, "PUT", "/api/v1/users/preferences", token, reqBody)
	assert.Equal(t, http.StatusOK, w.Code, "Initial preferences should save successfully")

	// Record initial timestamp
	initialResponse := ParseJSONResponse(t, w)
	initialData := GetDataFromResponse(t, initialResponse)
	initialUpdatedAt := initialData["updated_at"].(string)

	// Wait a moment to ensure timestamp difference
	time.Sleep(100 * time.Millisecond)

	// Step 2: Fetch current preferences
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/preferences", token, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Get preferences should return 200 OK")

	response := ParseJSONResponse(t, w)
	data := GetDataFromResponse(t, response)
	assertPreferencesStructure(t, data)

	// Verify initial values
	assert.Equal(t, "Easy", data["difficulty_preference"], "Should have initial difficulty")
	assert.Equal(t, false, data["notifications_enabled"], "Should have initial notifications setting")

	// Step 3: Fetch updated categories (verify consistency)
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/categories", token, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Categories endpoint should return 200 OK")

	response = ParseJSONResponse(t, w)
	categories, exists := response["data"]
	assert.True(t, exists, "Should have categories")

	categoriesList, ok := categories.([]interface{})
	assert.True(t, ok, "Categories should be an array")

	// Step 4: Update preferences with new interests from categories
	availableInterests := extractInterestIDs(categoriesList)
	newSelectedInterests := availableInterests[:min(4, len(availableInterests))]

	updateReq := models.UpdatePreferencesRequest{
		SelectedInterests:     newSelectedInterests,
		DifficultyPreference:  "Hard",        // Changed
		NotificationsEnabled:  true,          // Changed
		NotificationFrequency: "Daily",       // Changed
		ProfileVisibility:     boolPtr(true), // Changed
		ShowOnlineStatus:      boolPtr(true), // Changed
		AllowFriendRequests:   boolPtr(true), // Changed
		ShareActivityStatus:   boolPtr(true), // Changed
	}

	reqBody, _ = json.Marshal(updateReq)
	w = MakeAuthenticatedRequest(t, tc, "PUT", "/api/v1/users/preferences", token, reqBody)
	assert.Equal(t, http.StatusOK, w.Code, "Preferences update should return 200 OK")

	// Step 5: Verify UPSERT logic works correctly (UPDATE, not INSERT)
	response = ParseJSONResponse(t, w)
	data = GetDataFromResponse(t, response)
	assertPreferencesStructure(t, data)

	// Verify user_id hasn't changed (confirming UPDATE, not INSERT)
	assert.Equal(t, userID.String(), data["user_id"], "User ID should remain the same")

	// Step 6: Verify all fields updated (not just interests)
	assert.Equal(t, "Hard", data["difficulty_preference"], "Difficulty should be updated")
	assert.Equal(t, true, data["notifications_enabled"], "Notifications should be updated")
	assert.Equal(t, "Daily", data["notification_frequency"], "Frequency should be updated")
	assert.Equal(t, true, data["profile_visibility"], "Profile visibility should be updated")
	assert.Equal(t, true, data["show_online_status"], "Online status should be updated")
	assert.Equal(t, true, data["allow_friend_requests"], "Friend requests should be updated")
	assert.Equal(t, true, data["share_activity_status"], "Activity status should be updated")

	// Verify interests were updated
	updatedInterests := data["selected_interests"].([]interface{})
	assert.Equal(t, len(newSelectedInterests), len(updatedInterests), "Should have updated interests count")

	// Step 7: Verify updated_at timestamp changed
	newUpdatedAt := data["updated_at"].(string)
	assert.NotEqual(t, initialUpdatedAt, newUpdatedAt, "Updated timestamp should change")

	// Step 8: Fetch preferences again to verify persistence
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/preferences", token, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Get preferences should return 200 OK")

	response = ParseJSONResponse(t, w)
	finalData := GetDataFromResponse(t, response)

	// Verify all changes persisted
	assert.Equal(t, "Hard", finalData["difficulty_preference"], "Difficulty should persist")
	assert.Equal(t, true, finalData["notifications_enabled"], "Notifications should persist")
	assert.Equal(t, "Daily", finalData["notification_frequency"], "Frequency should persist")

	t.Logf("Existing user update flow completed successfully for user: %s", userID)
}

// testCategoriesCrossEndpointConsistency tests data consistency across different endpoints
func testCategoriesCrossEndpointConsistency(t *testing.T, tc *TestConfig) {
	// Step 1: Create user and set initial preferences
	userID, token := CreateTestUser(t, tc)

	// Step 2: Fetch categories and extract interests
	w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/categories", token, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Categories endpoint should return 200 OK")

	response := ParseJSONResponse(t, w)
	categories, exists := response["data"]
	assert.True(t, exists, "Response should contain 'categories' field")

	categoriesList, ok := categories.([]interface{})
	assert.True(t, ok, "Categories field should be an array")
	availableInterests := extractInterestIDs(categoriesList)
	assert.Greater(t, len(availableInterests), 0, "Should have available interests")

	// Step 3: Save preferences with interests from categories
	selectedInterests := availableInterests[:min(2, len(availableInterests))]
	updateReq := models.UpdatePreferencesRequest{
		SelectedInterests:     selectedInterests,
		DifficultyPreference:  "Medium",
		NotificationsEnabled:  true,
		NotificationFrequency: "Daily",
		ProfileVisibility:     boolPtr(true),
		ShowOnlineStatus:      boolPtr(true),
		AllowFriendRequests:   boolPtr(true),
		ShareActivityStatus:   boolPtr(true),
	}

	reqBody, _ := json.Marshal(updateReq)
	w = MakeAuthenticatedRequest(t, tc, "PUT", "/api/v1/users/preferences", token, reqBody)
	assert.Equal(t, http.StatusOK, w.Code, "Preferences update should return 200 OK")

	// Step 4: Verify categories endpoint still returns same data
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/categories", token, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Categories should still be accessible")

	response = ParseJSONResponse(t, w)
	newCategories, exists := response["data"]
	assert.True(t, exists, "Should still have categories")

	newCategoriesList, ok := newCategories.([]interface{})
	assert.True(t, ok, "Categories should still be an array")
	newAvailableInterests := extractInterestIDs(newCategoriesList)

	// Verify categories data consistency
	assert.Equal(t, len(availableInterests), len(newAvailableInterests), "Categories count should remain consistent")

	// Step 5: Verify preferences endpoint returns saved data
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/preferences", token, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Get preferences should return 200 OK")

	response = ParseJSONResponse(t, w)
	preferencesData := GetDataFromResponse(t, response)
	assertPreferencesStructure(t, preferencesData)

	// Verify saved interests are subset of available interests
	savedInterests := preferencesData["selected_interests"].([]interface{})
	for _, savedInterest := range savedInterests {
		interestID := savedInterest.(string)
		found := false
		for _, availableInterest := range newAvailableInterests {
			if availableInterest == interestID {
				found = true
				break
			}
		}
		assert.True(t, found, "Saved interest %s should exist in available interests", interestID)
	}

	// Step 6: Verify profile endpoint includes consistent data
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/profile", token, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Profile endpoint should return 200 OK")

	response = ParseJSONResponse(t, w)
	assert.Equal(t, userID.String(), response["id"], "Profile should have correct user ID")

	// Check if preferences are included in profile response
	if profilePreferences, exists := response["preferences"]; exists {
		profilePreferencesMap, ok := profilePreferences.(map[string]interface{})
		assert.True(t, ok, "Profile preferences should be an object")

		// Verify profile preferences match direct preferences endpoint
		assert.Equal(t, preferencesData["difficulty_preference"], profilePreferencesMap["difficulty_preference"], "Difficulty should match across endpoints")
		assert.Equal(t, preferencesData["notifications_enabled"], profilePreferencesMap["notifications_enabled"], "Notifications setting should match across endpoints")
	}

	t.Logf("Data consistency validation completed successfully for user: %s", userID)
}

// testSettingsScreenWorkflow simulates the complete settings screen workflow
func testSettingsScreenWorkflow(t *testing.T, tc *TestConfig) {
	// Step 1: Create user (simulating app startup)
	userID, token := CreateTestUser(t, tc)

	// Step 2: Simulate settings screen loading - fetch all required data
	// This simulates what the Flutter app does when opening settings

	// Fetch categories for interest selection
	w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/categories", token, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Categories should load for settings")

	categoriesResponse := ParseJSONResponse(t, w)
	categories, exists := categoriesResponse["data"]
	assert.True(t, exists, "Should have categories for settings")

	categoriesList, ok := categories.([]interface{})
	assert.True(t, ok, "Categories should be a list")
	assert.Greater(t, len(categoriesList), 0, "Should have categories to display")

	// Fetch difficulty levels for dropdown
	w = MakeRequest(t, tc, "GET", "/api/v1/preferences/difficulty-levels")
	assert.Equal(t, http.StatusOK, w.Code, "Difficulty levels should load")

	// Fetch notification frequencies for dropdown
	w = MakeRequest(t, tc, "GET", "/api/v1/preferences/notification-frequencies")
	assert.Equal(t, http.StatusOK, w.Code, "Notification frequencies should load")

	// Fetch current user preferences (will be empty for new user)
	_ = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/preferences", token, nil)
	// For new user, this might return 404 or empty preferences

	// Step 3: User selects interests and other preferences (simulating form interaction)
	availableInterests := extractInterestIDs(categoriesList)
	selectedInterests := availableInterests[:min(3, len(availableInterests))]

	// Simulate user filling the settings form
	initialPreferences := models.UpdatePreferencesRequest{
		SelectedInterests:     selectedInterests,
		DifficultyPreference:  "Easy",
		NotificationsEnabled:  true,
		NotificationFrequency: "Daily",
		ProfileVisibility:     boolPtr(true),
		ShowOnlineStatus:      boolPtr(true),
		AllowFriendRequests:   boolPtr(true),
		ShareActivityStatus:   boolPtr(true),
	}

	// Step 4: Save preferences (simulating save button click)
	reqBody, _ := json.Marshal(initialPreferences)
	w = MakeAuthenticatedRequest(t, tc, "PUT", "/api/v1/users/preferences", token, reqBody)
	assert.Equal(t, http.StatusOK, w.Code, "Settings should save successfully")

	// Step 5: Simulate navigating away and back to settings (testing persistence)
	time.Sleep(50 * time.Millisecond) // Brief pause to simulate navigation

	// Load settings again (simulating returning to settings screen)
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/preferences", token, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Settings should reload successfully")

	response := ParseJSONResponse(t, w)
	reloadedData := GetDataFromResponse(t, response)
	assertPreferencesStructure(t, reloadedData)

	// Verify all settings persisted correctly
	assert.Equal(t, "Easy", reloadedData["difficulty_preference"], "Difficulty should persist")
	assert.Equal(t, true, reloadedData["notifications_enabled"], "Notifications should persist")
	assert.Equal(t, "Daily", reloadedData["notification_frequency"], "Frequency should persist")
	assert.Equal(t, true, reloadedData["profile_visibility"], "Profile visibility should persist")

	reloadedInterests := reloadedData["selected_interests"].([]interface{})
	assert.Equal(t, len(selectedInterests), len(reloadedInterests), "Interest count should persist")

	// Step 6: Simulate user making changes to existing settings
	updatedInterests := availableInterests[:min(4, len(availableInterests))] // Add one more interest
	updateReq := models.UpdatePreferencesRequest{
		SelectedInterests:     updatedInterests,
		DifficultyPreference:  "Hard",         // Changed
		NotificationsEnabled:  false,          // Changed
		NotificationFrequency: "Weekly",       // Changed
		ProfileVisibility:     boolPtr(false), // Changed
		ShowOnlineStatus:      boolPtr(false), // Changed
		AllowFriendRequests:   boolPtr(false), // Changed
		ShareActivityStatus:   boolPtr(false), // Changed
	}

	reqBody, _ = json.Marshal(updateReq)
	w = MakeAuthenticatedRequest(t, tc, "PUT", "/api/v1/users/preferences", token, reqBody)
	assert.Equal(t, http.StatusOK, w.Code, "Settings update should save successfully")

	// Step 7: Verify updates persisted
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/preferences", token, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Updated settings should reload")

	response = ParseJSONResponse(t, w)
	finalData := GetDataFromResponse(t, response)

	// Verify all changes took effect
	assert.Equal(t, "Hard", finalData["difficulty_preference"], "Updated difficulty should persist")
	assert.Equal(t, false, finalData["notifications_enabled"], "Updated notifications should persist")
	assert.Equal(t, "Weekly", finalData["notification_frequency"], "Updated frequency should persist")
	assert.Equal(t, false, finalData["profile_visibility"], "Updated privacy should persist")

	finalInterests := finalData["selected_interests"].([]interface{})
	assert.Equal(t, len(updatedInterests), len(finalInterests), "Updated interest count should persist")

	t.Logf("Complete settings workflow validated successfully for user: %s", userID)
}

// testDatabaseDrivenCategoriesStructure validates the database-driven categories structure
func testDatabaseDrivenCategoriesStructure(t *testing.T, tc *TestConfig) {
	// Step 1: Create user for authenticated requests
	_, token := CreateTestUser(t, tc)

	// Step 2: Fetch categories and validate structure
	w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/categories", token, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Categories endpoint should return 200 OK")

	response := ParseJSONResponse(t, w)
	categories, exists := response["data"]
	assert.True(t, exists, "Response should contain 'categories' field")

	categoriesList, ok := categories.([]interface{})
	assert.True(t, ok, "Categories field should be an array")
	assert.Greater(t, len(categoriesList), 0, "Should have at least one category")

	// Step 3: Validate each category structure
	expectedCategoryFields := []string{"id", "name", "display_name", "description", "interests"}
	expectedInterestFields := []string{"id", "name", "display_name", "description", "is_test_data"}

	categoryIds := make(map[string]bool)
	interestIds := make(map[string]bool)
	totalInterests := 0

	for _, category := range categoriesList {
		categoryMap, ok := category.(map[string]interface{})
		assert.True(t, ok, "Each category should be an object")

		// Validate required category fields
		for _, field := range expectedCategoryFields {
			assert.Contains(t, categoryMap, field, "Category should have '%s' field", field)
		}

		// Validate category ID uniqueness
		categoryID := categoryMap["id"].(string)
		assert.False(t, categoryIds[categoryID], "Category ID '%s' should be unique", categoryID)
		categoryIds[categoryID] = true

		// Validate category has non-empty name and display_name
		assert.NotEmpty(t, categoryMap["name"], "Category name should not be empty")
		assert.NotEmpty(t, categoryMap["display_name"], "Category display_name should not be empty")

		// Validate interests array
		interests, ok := categoryMap["interests"].([]interface{})
		assert.True(t, ok, "Category interests should be an array")
		assert.Greater(t, len(interests), 0, "Each category should have at least one interest")

		// Validate each interest
		for _, interest := range interests {
			interestMap, ok := interest.(map[string]interface{})
			assert.True(t, ok, "Each interest should be an object")

			// Validate required interest fields
			for _, field := range expectedInterestFields {
				assert.Contains(t, interestMap, field, "Interest should have '%s' field", field)
			}

			// Validate interest ID uniqueness across all categories
			interestID := interestMap["id"].(string)
			assert.False(t, interestIds[interestID], "Interest ID '%s' should be unique across all categories", interestID)
			interestIds[interestID] = true

			// Validate interest has non-empty names
			assert.NotEmpty(t, interestMap["name"], "Interest name should not be empty")
			assert.NotEmpty(t, interestMap["display_name"], "Interest display_name should not be empty")

			// Validate is_test_data is boolean
			_, ok = interestMap["is_test_data"].(bool)
			assert.True(t, ok, "Interest is_test_data should be a boolean")

			totalInterests++
		}
	}

	// Step 4: Validate overall structure expectations
	assert.GreaterOrEqual(t, len(categoryIds), 3, "Should have at least 3 categories")
	assert.GreaterOrEqual(t, totalInterests, 10, "Should have at least 10 total interests across all categories")

	// Verify we have expected categories based on seed data
	expectedCategoryNames := []string{"general", "science", "sports", "entertainment"}
	foundExpectedCategories := 0

	for _, category := range categoriesList {
		categoryMap := category.(map[string]interface{})
		categoryID := categoryMap["id"].(string)

		for _, expectedName := range expectedCategoryNames {
			if categoryID == expectedName {
				foundExpectedCategories++
				break
			}
		}
	}

	assert.GreaterOrEqual(t, foundExpectedCategories, 3, "Should find at least 3 expected categories from seed data")

	// Step 5: Validate categories can be used for preferences
	// Extract all interest IDs
	allInterestIDs := extractInterestIDs(categoriesList)
	assert.Equal(t, totalInterests, len(allInterestIDs), "Extracted interest IDs should match total count")

	// Try to use a few interests in preferences to ensure they're valid
	testInterests := allInterestIDs[:min(2, len(allInterestIDs))]
	updateReq := models.UpdatePreferencesRequest{
		SelectedInterests:     testInterests,
		DifficultyPreference:  "Medium",
		NotificationsEnabled:  true,
		NotificationFrequency: "Daily",
		ProfileVisibility:     boolPtr(true),
		ShowOnlineStatus:      boolPtr(true),
		AllowFriendRequests:   boolPtr(true),
		ShareActivityStatus:   boolPtr(true),
	}

	reqBody, _ := json.Marshal(updateReq)
	w = MakeAuthenticatedRequest(t, tc, "PUT", "/api/v1/users/preferences", token, reqBody)
	assert.Equal(t, http.StatusOK, w.Code, "Should be able to save preferences with category interests")

	// Step 6: Verify saved preferences match available interests
	w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/preferences", token, nil)
	assert.Equal(t, http.StatusOK, w.Code, "Should be able to retrieve saved preferences")

	response = ParseJSONResponse(t, w)
	preferencesData := GetDataFromResponse(t, response)
	savedInterests := preferencesData["selected_interests"].([]interface{})

	for _, savedInterest := range savedInterests {
		interestID := savedInterest.(string)
		assert.True(t, interestIds[interestID], "Saved interest '%s' should exist in available categories", interestID)
	}

	t.Logf("Database-driven categories structure validation completed. Found %d categories with %d total interests", len(categoryIds), totalInterests)
}

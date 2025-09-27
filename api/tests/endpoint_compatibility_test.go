package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndpointCompatibility(t *testing.T) {
	tc := SetupTestServer(t)
	defer Cleanup(t)

	userID, token := CreateTestUser(t, tc)

	t.Run("CategoriesEndpointCompatibility", func(t *testing.T) {
		// Test that the categories endpoint returns the same structure as the original mock
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/categories", token, nil)
		assert.Equal(t, http.StatusOK, w.Code, "Categories endpoint should work")

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			// Verify the response structure matches original mock expectations
			categories, exists := data["categories"]
			assert.True(t, exists, "Should have 'categories' field")

			categoriesList, ok := categories.([]interface{})
			assert.True(t, ok, "Categories should be an array")
			assert.Greater(t, len(categoriesList), 0, "Should have at least one category")

			// Verify we have the expected categories from original mock
			expectedCategoryIDs := []string{"general", "science", "sports"}
			foundCategories := make(map[string]bool)

			for _, category := range categoriesList {
				if categoryMap, ok := category.(map[string]interface{}); ok {
					if id, exists := categoryMap["id"]; exists {
						if idStr, ok := id.(string); ok {
							foundCategories[idStr] = true

							// Verify each category has the expected structure
							assert.Contains(t, categoryMap, "name", "Category should have 'name'")
							assert.Contains(t, categoryMap, "display_name", "Category should have 'display_name'")
							assert.Contains(t, categoryMap, "description", "Category should have 'description'")
							assert.Contains(t, categoryMap, "interests", "Category should have 'interests'")

							// Verify interests structure
							if interests, ok := categoryMap["interests"].([]interface{}); ok {
								assert.Greater(t, len(interests), 0, "Category should have interests")

								// Check first interest structure
								if len(interests) > 0 {
									if interest, ok := interests[0].(map[string]interface{}); ok {
										VerifyIsTestDataField(t, interest, true, "category interest")
									}
								}
							}
						}
					}
				}
			}

			// Verify all expected categories were found
			for _, expectedID := range expectedCategoryIDs {
				assert.True(t, foundCategories[expectedID], "Should find category: %s", expectedID)
			}
		}
	})

	t.Run("AppSettingsEndpointCompatibility", func(t *testing.T) {
		// Test that the app settings endpoint returns the same structure as the original mock
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/config/app-settings", token, nil)
		assert.Equal(t, http.StatusOK, w.Code, "App settings endpoint should work")

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			// Verify original mock fields are present
			originalMockFields := map[string]interface{}{
				"app_name":                  "QuizNinja",
				"max_questions_per_quiz":    20.0, // Should be 20 after migration
				"quiz_categories_enabled":   true,
				"leaderboard_enabled":       true,
				"achievements_enabled":      true,
			}

			for fieldName, expectedValue := range originalMockFields {
				value, exists := data[fieldName]
				assert.True(t, exists, "Should have field from original mock: %s", fieldName)
				assert.Equal(t, expectedValue, value, "Field %s should match original mock value", fieldName)
			}

			// Verify computed fields that were in original mock
			supportedLanguages, exists := data["supported_languages"]
			assert.True(t, exists, "Should have 'supported_languages' field")
			if languages, ok := supportedLanguages.([]interface{}); ok {
				assert.Contains(t, languages, "en", "Should support English")
				assert.Contains(t, languages, "es", "Should support Spanish")
				assert.Contains(t, languages, "fr", "Should support French")
				assert.Contains(t, languages, "de", "Should support German")
			}

			difficultyLevels, exists := data["difficulty_levels"]
			assert.True(t, exists, "Should have 'difficulty_levels' field")
			if levels, ok := difficultyLevels.([]interface{}); ok {
				assert.Len(t, levels, 3, "Should have 3 difficulty levels")
				// Verify structure of first difficulty level
				if len(levels) > 0 {
					if level, ok := levels[0].(map[string]interface{}); ok {
						assert.Contains(t, level, "id", "Difficulty level should have 'id'")
						assert.Contains(t, level, "name", "Difficulty level should have 'name'")
						assert.Contains(t, level, "description", "Difficulty level should have 'description'")
					}
				}
			}

			notificationFrequencies, exists := data["notification_frequencies"]
			assert.True(t, exists, "Should have 'notification_frequencies' field")
			if frequencies, ok := notificationFrequencies.([]interface{}); ok {
				assert.Len(t, frequencies, 3, "Should have 3 notification frequencies")
			}
		}
	})

	t.Run("EndToEndUserJourney", func(t *testing.T) {
		// Test a complete user journey: get categories -> get app settings -> complete quiz -> check stats

		// Step 1: Get categories
		categoriesW := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/categories", token, nil)
		assert.Equal(t, http.StatusOK, categoriesW.Code, "Should get categories")

		// Step 2: Get app settings
		settingsW := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/config/app-settings", token, nil)
		assert.Equal(t, http.StatusOK, settingsW.Code, "Should get app settings")

		// Step 3: Get initial user stats
		initialStatsW := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/stats", token, nil)
		assert.Equal(t, http.StatusOK, initialStatsW.Code, "Should get initial user stats")

		var initialQuizCount int
		if initialStatsW.Code == http.StatusOK {
			response := ParseJSONResponse(t, initialStatsW)
			data := GetDataFromResponse(t, response)
			VerifyIsTestDataField(t, data, true, "user stats")

			if count, exists := data["total_quizzes_completed"]; exists {
				if val, ok := count.(float64); ok {
					initialQuizCount = int(val)
				}
			}
		}

		// Step 4: Get available quizzes
		quizzesW := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes?limit=1", token, nil)
		assert.Equal(t, http.StatusOK, quizzesW.Code, "Should get quizzes")

		var quizID string
		if quizzesW.Code == http.StatusOK {
			response := ParseJSONResponse(t, quizzesW)
			data := GetDataFromResponse(t, response)

			if quizzes, exists := data["quizzes"]; exists {
				if quizzesList, ok := quizzes.([]interface{}); ok && len(quizzesList) > 0 {
					if quiz, ok := quizzesList[0].(map[string]interface{}); ok {
						if id, exists := quiz["id"]; exists {
							quizID = id.(string)
						}
					}
				}
			}
		}

		if quizID == "" {
			t.Skip("No quizzes available for end-to-end test")
			return
		}

		// Step 5: Start a quiz attempt
		startW := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/quizzes/"+quizID+"/attempts", token, nil)
		assert.Equal(t, http.StatusOK, startW.Code, "Should start quiz attempt")

		var attemptID string
		if startW.Code == http.StatusOK {
			response := ParseJSONResponse(t, startW)
			if attempt, exists := response["attempt_id"]; exists {
				attemptID = attempt.(string)
			}
		}

		if attemptID == "" {
			t.Fatal("Failed to start quiz attempt")
			return
		}

		// Step 6: Submit quiz (with minimal answers)
		submitData := map[string]interface{}{
			"attemptId": attemptID,
			"answers":   []map[string]interface{}{}, // Empty answers for simplicity
			"timeSpent": 60,
		}

		submitJSON, _ := json.Marshal(submitData)
		submitW := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/quizzes/"+quizID+"/attempts/"+attemptID+"/submit", token, submitJSON)
		assert.Equal(t, http.StatusOK, submitW.Code, "Should submit quiz")

		// Step 7: Verify user stats were updated
		finalStatsW := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/stats", token, nil)
		assert.Equal(t, http.StatusOK, finalStatsW.Code, "Should get final user stats")

		if finalStatsW.Code == http.StatusOK {
			response := ParseJSONResponse(t, finalStatsW)
			data := GetDataFromResponse(t, response)
			VerifyIsTestDataField(t, data, true, "final user stats")

			if count, exists := data["total_quizzes_completed"]; exists {
				if val, ok := count.(float64); ok {
					finalQuizCount := int(val)
					assert.Equal(t, initialQuizCount+1, finalQuizCount, "Quiz count should have increased by 1")
				}
			}
		}
	})

	t.Run("VerifyTestDataConsistency", func(t *testing.T) {
		// Verify that all endpoints consistently return is_test_data = true for test data

		// Test categories endpoint
		categoriesW := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/categories", token, nil)
		if categoriesW.Code == http.StatusOK {
			response := ParseJSONResponse(t, categoriesW)
			data := GetDataFromResponse(t, response)

			if categories, exists := data["categories"]; exists {
				if categoriesList, ok := categories.([]interface{}); ok {
					for _, category := range categoriesList {
						if categoryMap, ok := category.(map[string]interface{}); ok {
							if interests, exists := categoryMap["interests"]; exists {
								if interestsList, ok := interests.([]interface{}); ok {
									for _, interest := range interestsList {
										if interestMap, ok := interest.(map[string]interface{}); ok {
											VerifyIsTestDataField(t, interestMap, true, "category interest test data consistency")
										}
									}
								}
							}
						}
					}
				}
			}
		}

		// Test user stats endpoint
		statsW := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/stats", token, nil)
		if statsW.Code == http.StatusOK {
			response := ParseJSONResponse(t, statsW)
			data := GetDataFromResponse(t, response)
			VerifyIsTestDataField(t, data, true, "user stats test data consistency")
		}
	})

	_ = userID // Use userID to avoid unused variable warning
}
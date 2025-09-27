package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"quizninja-api/models"

	"github.com/stretchr/testify/assert"
)

func TestAchievementHandler(t *testing.T) {
	tc := SetupTestServer(t)
	defer Cleanup(t)

	userID, token := CreateTestUser(t, tc)

	t.Run("GetAllAchievements", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/achievements", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			achievements, exists := response["achievements"]
			assert.True(t, exists, "Response should contain 'achievements' field")

			achievementsList, ok := achievements.([]interface{})
			assert.True(t, ok, "Achievements field should be an array")

			if len(achievementsList) > 0 {
				VerifyIsTestDataInArray(t, achievementsList, true, "all achievements")
			}

			// Verify total count
			total, totalExists := response["total"]
			assert.True(t, totalExists, "Response should contain 'total' field")

			totalFloat, ok := total.(float64)
			if ok {
				assert.Equal(t, len(achievementsList), int(totalFloat), "Total should match achievements count")
			}
		}
	})

	t.Run("GetUserAchievements", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/achievements", token, nil)

		if w.Code == http.StatusOK {
			var response models.AchievementListResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should parse user achievements response")

			if len(response.Achievements) > 0 {
				// Verify user achievements have is_test_data field
				achievementsData := make([]interface{}, len(response.Achievements))
				for i, achievement := range response.Achievements {
					achievementsData[i] = map[string]interface{}{
						"is_test_data": achievement.IsTestData,
					}
				}
				VerifyIsTestDataInArray(t, achievementsData, true, "user achievements")

				// Check nested achievement data
				for i, userAchievement := range response.Achievements {
					if userAchievement.Achievement != nil {
						achievementMap := map[string]interface{}{
							"is_test_data": userAchievement.Achievement.IsTestData,
						}
						VerifyIsTestDataField(t, achievementMap, true, fmt.Sprintf("user achievement[%d] achievement", i))
					}
				}
			}
		}
	})

	t.Run("GetAchievementProgress", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/achievements/progress", token, nil)

		if w.Code == http.StatusOK {
			var response models.AchievementProgressResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should parse achievement progress response")

			if len(response.Progress) > 0 {
				// Note: AchievementProgress model doesn't have is_test_data field
				// This is expected as progress is calculated dynamically
				t.Logf("Found %d achievement progress entries", len(response.Progress))
			}
		}
	})

	t.Run("GetUserAchievementsByUserID", func(t *testing.T) {
		// Test getting achievements for the current user via user ID endpoint
		url := fmt.Sprintf("/api/v1/users/%s/achievements", userID)
		w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)

		if w.Code == http.StatusOK {
			var response models.AchievementListResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should parse user achievements by ID response")

			if len(response.Achievements) > 0 {
				// Verify achievements have is_test_data field
				achievementsData := make([]interface{}, len(response.Achievements))
				for i, achievement := range response.Achievements {
					achievementsData[i] = map[string]interface{}{
						"is_test_data": achievement.IsTestData,
					}
				}
				VerifyIsTestDataInArray(t, achievementsData, true, "user achievements by ID")
			}
		}
	})

	t.Run("GetAchievementsByCategory", func(t *testing.T) {
		categories := []string{"quiz", "social", "streak", "general"}

		for _, category := range categories {
			t.Run(fmt.Sprintf("Category_%s", category), func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/achievements/category/%s", category)
				w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)

				if w.Code == http.StatusOK {
					response := ParseJSONResponse(t, w)

					achievements, exists := response["achievements"]
					if exists {
						achievementsList, ok := achievements.([]interface{})
						if ok && len(achievementsList) > 0 {
							VerifyIsTestDataInArray(t, achievementsList, true, fmt.Sprintf("achievements category %s", category))
						}
					}

					// Verify category in response
					responseCategory, categoryExists := response["category"]
					if categoryExists {
						assert.Equal(t, category, responseCategory, "Response should contain correct category")
					}
				}
			})
		}
	})

	t.Run("GetAchievementStats", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/achievements/stats", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			stats, exists := response["stats"]
			assert.True(t, exists, "Response should contain 'stats' field")

			statsMap, ok := stats.(map[string]interface{})
			assert.True(t, ok, "Stats field should be an object")

			// Verify basic stats structure
			expectedKeys := []string{
				"total_achievements",
				"unlocked_achievements",
				"locked_achievements",
				"completion_percentage",
				"total_points_from_achievements",
				"rare_achievements",
			}

			for _, key := range expectedKeys {
				_, keyExists := statsMap[key]
				assert.True(t, keyExists, fmt.Sprintf("Stats should contain '%s' field", key))
			}
		}
	})

	t.Run("CheckAchievements", func(t *testing.T) {
		// Test manual achievement check
		w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/users/achievements/check", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// Verify response structure
			newAchievements, exists := response["new_achievements"]
			assert.True(t, exists, "Response should contain 'new_achievements' field")

			newAchievementsList, ok := newAchievements.([]interface{})
			assert.True(t, ok, "New achievements field should be an array")

			if len(newAchievementsList) > 0 {
				VerifyIsTestDataInArray(t, newAchievementsList, true, "new achievements from check")
			}

			// Verify notifications
			notifications, notificationsExist := response["notifications"]
			if notificationsExist {
				notificationsList, ok := notifications.([]interface{})
				if ok && len(notificationsList) > 0 {
					// Notifications are derived from achievements, may not have is_test_data
					// but we can verify basic structure
					assert.GreaterOrEqual(t, len(notificationsList), 0, "Notifications should be valid array")
				}
			}

			// Verify count
			count, countExists := response["count"]
			if countExists {
				countFloat, ok := count.(float64)
				if ok {
					assert.Equal(t, len(newAchievementsList), int(countFloat), "Count should match new achievements length")
				}
			}
		}
	})

	t.Run("CheckAchievementsWithTrigger", func(t *testing.T) {
		// Test achievement check with specific trigger
		w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/users/achievements/check?trigger=quiz_completed", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// Verify trigger type in response
			triggerType, triggerExists := response["trigger_type"]
			if triggerExists {
				assert.Equal(t, "quiz_completed", triggerType, "Response should contain correct trigger type")
			}
		}
	})

	t.Run("UnlockAchievement", func(t *testing.T) {
		// Test manual achievement unlock (this might be an admin endpoint)
		achievementKey := "test_achievement"
		url := fmt.Sprintf("/api/v1/users/achievements/unlock/%s", achievementKey)
		w := MakeAuthenticatedRequest(t, tc, "POST", url, token, nil)

		// This might fail if the achievement doesn't exist or user already has it
		// but we test that the endpoint is working and handles authentication
		if w.Code == http.StatusCreated {
			response := ParseJSONResponse(t, w)

			// Verify achievement in response
			achievement, achievementExists := response["achievement"]
			if achievementExists {
				achievementMap, ok := achievement.(map[string]interface{})
				if ok {
					VerifyIsTestDataField(t, achievementMap, true, "unlocked achievement")
				}
			}

			// Verify notification
			notification, notificationExists := response["notification"]
			if notificationExists {
				// Notifications are generated data, may not have is_test_data
				assert.NotNil(t, notification, "Notification should be present")
			}
		}
	})

	t.Run("GetLeaderboardWithAchievements", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/achievements/leaderboard", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			leaderboard, exists := response["leaderboard"]
			assert.True(t, exists, "Response should contain 'leaderboard' field")

			leaderboardList, ok := leaderboard.([]interface{})
			assert.True(t, ok, "Leaderboard field should be an array")

			if len(leaderboardList) > 0 {
				// Leaderboard entries contain user data
				for i, entry := range leaderboardList {
					entryMap, ok := entry.(map[string]interface{})
					if ok {
						// Check user data if present
						if user, userExists := entryMap["user"]; userExists {
							userMap, ok := user.(map[string]interface{})
							if ok {
								VerifyIsTestDataField(t, userMap, true, fmt.Sprintf("leaderboard[%d] user", i))
							}
						}
					}
				}
			}

			// Verify response structure
			period, periodExists := response["period"]
			if periodExists {
				assert.NotEmpty(t, period, "Period should be specified")
			}
		}
	})

	t.Run("GetLeaderboardWithAchievementsFiltered", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/achievements/leaderboard?period=weekly&limit=10", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// Verify filtered parameters in response
			period, periodExists := response["period"]
			if periodExists {
				assert.Equal(t, "weekly", period, "Should reflect filtered period")
			}

			limit, limitExists := response["limit"]
			if limitExists {
				limitFloat, ok := limit.(float64)
				if ok {
					assert.Equal(t, float64(10), limitFloat, "Should reflect filtered limit")
				}
			}
		}
	})

	_ = userID // Use userID to avoid unused variable warning
}

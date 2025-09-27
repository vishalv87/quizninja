package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLeaderboardHandler(t *testing.T) {
	tc := SetupTestServer(t)
	defer Cleanup(t)

	userID, token := CreateTestUser(t, tc)

	t.Run("GetGlobalLeaderboard", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/leaderboard", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			leaderboard, exists := data["leaderboard"]
			if exists {
				leaderboardList, ok := leaderboard.([]interface{})
				if ok && len(leaderboardList) > 0 {
					// Leaderboard entries contain user data
					for i, entry := range leaderboardList {
						entryMap, ok := entry.(map[string]interface{})
						if ok {
							// Check user data if present
							user, userExists := entryMap["user"]
							if userExists && user != nil {
								userMap, ok := user.(map[string]interface{})
								if ok {
									VerifyIsTestDataField(t, userMap, true, fmt.Sprintf("global leaderboard[%d] user", i))
								}
							}
						}
					}
				}
			}
		}
	})

	t.Run("GetGlobalLeaderboardWithPeriod", func(t *testing.T) {
		periods := []string{"daily", "weekly", "monthly", "alltime"}

		for _, period := range periods {
			t.Run(fmt.Sprintf("Period_%s", period), func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/leaderboard?period=%s&limit=10", period)
				w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)

				if w.Code == http.StatusOK {
					response := ParseJSONResponse(t, w)
					data := GetDataFromResponse(t, response)

					leaderboard, exists := data["leaderboard"]
					if exists {
						leaderboardList, ok := leaderboard.([]interface{})
						if ok && len(leaderboardList) > 0 {
							// Check user data in leaderboard entries
							for i, entry := range leaderboardList {
								entryMap, ok := entry.(map[string]interface{})
								if ok {
									user, userExists := entryMap["user"]
									if userExists && user != nil {
										userMap, ok := user.(map[string]interface{})
										if ok {
											VerifyIsTestDataField(t, userMap, true, fmt.Sprintf("%s leaderboard[%d] user", period, i))
										}
									}
								}
							}
						}
					}

					// Verify period is reflected in response
					responsePeriod, periodExists := data["period"]
					if periodExists {
						assert.Equal(t, period, responsePeriod, "Response should reflect requested period")
					}
				}
			})
		}
	})

	t.Run("GetQuizLeaderboard", func(t *testing.T) {
		// First get a quiz to test leaderboard for
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes", token, nil)
		if w.Code != http.StatusOK {
			t.Skip("No quizzes available to test quiz leaderboard")
			return
		}

		response := ParseJSONResponse(t, w)
		data := GetDataFromResponse(t, response)

		quizzes, exists := data["quizzes"]
		if !exists {
			t.Skip("No quizzes field in response")
			return
		}

		quizzesList, ok := quizzes.([]interface{})
		if !ok || len(quizzesList) == 0 {
			t.Skip("No quizzes available")
			return
		}

		firstQuiz, ok := quizzesList[0].(map[string]interface{})
		if !ok {
			t.Skip("Invalid quiz data")
			return
		}

		quizID, exists := firstQuiz["id"]
		if !exists {
			t.Skip("Quiz missing ID")
			return
		}

		url := fmt.Sprintf("/api/v1/leaderboard/quiz/%s", quizID)
		w = MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			leaderboard, exists := data["leaderboard"]
			if exists {
				leaderboardList, ok := leaderboard.([]interface{})
				if ok && len(leaderboardList) > 0 {
					// Quiz leaderboard entries contain user and attempt data
					for i, entry := range leaderboardList {
						entryMap, ok := entry.(map[string]interface{})
						if ok {
							// Check user data
							user, userExists := entryMap["user"]
							if userExists && user != nil {
								userMap, ok := user.(map[string]interface{})
								if ok {
									VerifyIsTestDataField(t, userMap, true, fmt.Sprintf("quiz leaderboard[%d] user", i))
								}
							}

							// Check attempt data if present
							attempt, attemptExists := entryMap["attempt"]
							if attemptExists && attempt != nil {
								attemptMap, ok := attempt.(map[string]interface{})
								if ok {
									VerifyIsTestDataField(t, attemptMap, true, fmt.Sprintf("quiz leaderboard[%d] attempt", i))
								}
							}
						}
					}
				}
			}
		}
	})

	t.Run("GetUserLeaderboardPosition", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/leaderboard/position", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			// Verify user position data
			position, exists := data["position"]
			if exists {
				_, ok := position.(float64) // JSON numbers are float64
				assert.True(t, ok, "position should be a number")
			}

			// Check user data if present
			user, userExists := data["user"]
			if userExists && user != nil {
				userMap, ok := user.(map[string]interface{})
				if ok {
					VerifyIsTestDataField(t, userMap, true, "leaderboard position user")
				}
			}
		}
	})

	t.Run("GetUserLeaderboardPositionWithPeriod", func(t *testing.T) {
		periods := []string{"daily", "weekly", "monthly"}

		for _, period := range periods {
			t.Run(fmt.Sprintf("Position_%s", period), func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/leaderboard/position?period=%s", period)
				w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)

				if w.Code == http.StatusOK {
					response := ParseJSONResponse(t, w)
					data := GetDataFromResponse(t, response)

					// Verify period is reflected in response
					responsePeriod, periodExists := data["period"]
					if periodExists {
						assert.Equal(t, period, responsePeriod, "Response should reflect requested period")
					}

					// Check user data
					user, userExists := data["user"]
					if userExists && user != nil {
						userMap, ok := user.(map[string]interface{})
						if ok {
							VerifyIsTestDataField(t, userMap, true, fmt.Sprintf("%s position user", period))
						}
					}
				}
			})
		}
	})

	t.Run("GetFriendsLeaderboard", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/leaderboard/friends", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			leaderboard, exists := data["leaderboard"]
			if exists {
				leaderboardList, ok := leaderboard.([]interface{})
				if ok && len(leaderboardList) > 0 {
					// Friends leaderboard entries contain user data
					for i, entry := range leaderboardList {
						entryMap, ok := entry.(map[string]interface{})
						if ok {
							user, userExists := entryMap["user"]
							if userExists && user != nil {
								userMap, ok := user.(map[string]interface{})
								if ok {
									VerifyIsTestDataField(t, userMap, true, fmt.Sprintf("friends leaderboard[%d] user", i))
								}
							}
						}
					}
				}
			}
		}
	})

	t.Run("GetCategoryLeaderboard", func(t *testing.T) {
		categories := []string{"science", "technology", "sports", "history"}

		for _, category := range categories {
			t.Run(fmt.Sprintf("Category_%s", category), func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/leaderboard/category/%s", category)
				w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)

				if w.Code == http.StatusOK {
					response := ParseJSONResponse(t, w)
					data := GetDataFromResponse(t, response)

					leaderboard, exists := data["leaderboard"]
					if exists {
						leaderboardList, ok := leaderboard.([]interface{})
						if ok && len(leaderboardList) > 0 {
							// Category leaderboard entries contain user data
							for i, entry := range leaderboardList {
								entryMap, ok := entry.(map[string]interface{})
								if ok {
									user, userExists := entryMap["user"]
									if userExists && user != nil {
										userMap, ok := user.(map[string]interface{})
										if ok {
											VerifyIsTestDataField(t, userMap, true, fmt.Sprintf("category %s leaderboard[%d] user", category, i))
										}
									}
								}
							}
						}
					}

					// Verify category is reflected in response
					responseCategory, categoryExists := data["category"]
					if categoryExists {
						assert.Equal(t, category, responseCategory, "Response should reflect requested category")
					}
				}
			})
		}
	})

	t.Run("GetLeaderboardStats", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/leaderboard/stats", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			// Verify basic stats structure
			totalUsers, exists := data["total_users"]
			if exists {
				_, ok := totalUsers.(float64)
				assert.True(t, ok, "total_users should be a number")
			}

			// Check top performer if present
			topPerformer, exists := data["top_performer"]
			if exists && topPerformer != nil {
				performerMap, ok := topPerformer.(map[string]interface{})
				if ok {
					user, userExists := performerMap["user"]
					if userExists && user != nil {
						userMap, ok := user.(map[string]interface{})
						if ok {
							VerifyIsTestDataField(t, userMap, true, "leaderboard stats top performer")
						}
					}
				}
			}
		}
	})

	_ = userID // Use userID to avoid unused variable warning
}

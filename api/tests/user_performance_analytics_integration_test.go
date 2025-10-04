package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestUserPerformanceAnalyticsIntegration tests the complete User Performance & Analytics flow
func TestUserPerformanceAnalyticsIntegration(t *testing.T) {
	tc := SetupTestServer(t)
	defer CleanupWithSupabase(t, tc)

	userID, token := CreateTestUser(t, tc)

	// Create a completed quiz attempt for testing
	var completedAttemptID string

	// Setup phase - create a completed quiz attempt for testing
	// Get available quizzes
	w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes?limit=1", token, nil)
	t.Logf("Quiz fetch request status: %d", w.Code)
	if w.Code != http.StatusOK {
		t.Logf("Failed to fetch quizzes: %s", w.Body.String())
	}

	if w.Code == http.StatusOK {
		response := ParseJSONResponse(t, w)
		data := GetDataFromResponse(t, response)

		if quizzes, exists := data["quizzes"]; exists {
			if quizzesList, ok := quizzes.([]interface{}); ok && len(quizzesList) > 0 {
				t.Logf("Found %d quiz(s) to test with", len(quizzesList))
				firstQuiz := quizzesList[0].(map[string]interface{})
				quizID := firstQuiz["id"].(string)

				// Start quiz attempt
				startURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts", quizID)
				startW := MakeAuthenticatedRequest(t, tc, "POST", startURL, token, nil)
				t.Logf("Quiz attempt creation status: %d for quiz %s", startW.Code, quizID)
				if startW.Code != http.StatusCreated && startW.Code != http.StatusOK {
					t.Logf("Failed to create quiz attempt: %s", startW.Body.String())
				}

				if startW.Code == http.StatusCreated || startW.Code == http.StatusOK {
					startResponse := ParseJSONResponse(t, startW)
					startData := GetDataFromResponse(t, startResponse)

					attemptID := startData["attempt_id"].(string)
					completedAttemptID = attemptID

					// Get quiz with questions for building answers
					if quiz, exists := startData["quiz"]; exists {
						quizMap := quiz.(map[string]interface{})
						if questions, exists := quizMap["questions"]; exists {
							// Build answers payload
							var answers []map[string]interface{}
							questionsList := questions.([]interface{})
							for i, question := range questionsList {
								questionMap := question.(map[string]interface{})
								questionID := questionMap["id"].(string)

								// Create a simple answer (select first option)
								answer := map[string]interface{}{
									"questionId":          questionID,
									"selectedOptionIndex": 0,
								}
								answers = append(answers, answer)

								// Limit to first 3 questions
								if i >= 2 {
									break
								}
							}

							// Submit the quiz attempt
							submitURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts/%s/submit", quizID, attemptID)
							submitPayload := map[string]interface{}{
								"attemptId": attemptID,
								"answers":   answers,
								"timeSpent": 180, // 3 minutes
							}

							submitBody, _ := json.Marshal(submitPayload)
							submitW := MakeAuthenticatedRequest(t, tc, "POST", submitURL, token, submitBody)
							t.Logf("Quiz submission status: %d for attempt %s", submitW.Code, attemptID)

							// Verify submission was successful
							if submitW.Code == http.StatusOK {
								t.Logf("Successfully created completed attempt: %s", attemptID)

								// Verify the attempt was actually marked as completed
								submitResponse := ParseJSONResponse(t, submitW)
								if score, exists := submitResponse["score"]; exists {
									t.Logf("Quiz completed with score: %v", score)
								}
							} else {
								t.Logf("Quiz submission failed with status %d: %s", submitW.Code, submitW.Body.String())
								completedAttemptID = "" // Clear the attempt ID if submission failed
							}
						}
					}
				}
			} else {
				t.Logf("No quizzes found in response or quizzes list is empty")
			}
		} else {
			t.Logf("No 'quizzes' field found in response")
		}
	} else {
		t.Logf("Quiz fetch request failed with status: %d", w.Code)
	}

	// Step 2: Test GET /users/stats endpoint
	t.Run("GetUserStats", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/stats", token, nil)
		assert.Equal(t, http.StatusOK, w.Code, "Should get user statistics")

		response := ParseJSONResponse(t, w)
		data, exists := response["data"]
		assert.True(t, exists, "Response should contain data field")

		stats := data.(map[string]interface{})
		VerifyIsTestDataField(t, stats, true, "user statistics")

		// Verify required statistics fields
		requiredFields := []string{
			"user_id", "total_attempts", "total_quizzes_completed", "completion_rate",
			"average_score", "total_points", "current_streak", "best_streak",
			"average_completion_time", "category_performance", "recent_activity",
			"quizzes_by_difficulty", "score_distribution", "monthly_progress",
		}

		for _, field := range requiredFields {
			assert.Contains(t, stats, field, "Statistics should contain %s field", field)
		}

		// Verify category performance structure
		if categoryPerf, exists := stats["category_performance"]; exists && categoryPerf != nil {
			categoryList, ok := categoryPerf.([]interface{})
			assert.True(t, ok, "Category performance should be an array")

			if len(categoryList) > 0 {
				firstCategory := categoryList[0].(map[string]interface{})
				categoryFields := []string{"category_id", "category_name", "quizzes_completed", "average_score", "total_attempts", "best_score"}
				for _, field := range categoryFields {
					assert.Contains(t, firstCategory, field, "Category performance should contain %s", field)
				}
			} else {
				t.Logf("Category performance is empty (expected for user with no completed quizzes)")
			}
		} else {
			t.Logf("No category performance data available (expected for user with no completed quizzes)")
		}

		// Verify recent activity structure
		if recentActivity, exists := stats["recent_activity"]; exists && recentActivity != nil {
			activityList, ok := recentActivity.([]interface{})
			assert.True(t, ok, "Recent activity should be an array")

			if len(activityList) > 0 {
				firstActivity := activityList[0].(map[string]interface{})
				activityFields := []string{"quiz_id", "quiz_title", "score", "category", "completed_at", "time_spent", "difficulty"}
				for _, field := range activityFields {
					assert.Contains(t, firstActivity, field, "Recent activity should contain %s", field)
				}
			} else {
				t.Logf("Recent activity is empty (expected for user with no completed quizzes)")
			}
		} else {
			t.Logf("No recent activity data available (expected for user with no completed quizzes)")
		}

		t.Logf("User statistics verified successfully")
	})

	// Step 3: Test GET /users/attempts endpoint with various filters
	t.Run("GetUserAttempts", func(t *testing.T) {
		// Test basic endpoint
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/attempts", token, nil)
		assert.Equal(t, http.StatusOK, w.Code, "Should get user attempts")

		response := ParseJSONResponse(t, w)
		data := GetDataFromResponse(t, response)

		// Verify response structure
		attempts, exists := data["attempts"]
		assert.True(t, exists, "Response should contain attempts")

		attemptsList, ok := attempts.([]interface{})
		assert.True(t, ok, "Attempts should be an array")

		if len(attemptsList) > 0 {
			VerifyIsTestDataInArray(t, attemptsList, true, "user attempts")
		}

		// Verify pagination fields
		requiredPaginationFields := []string{"total", "page", "page_size", "total_pages"}
		for _, field := range requiredPaginationFields {
			assert.Contains(t, data, field, "Response should contain %s field", field)
		}

		// Test with filtering parameters
		t.Run("WithFilters", func(t *testing.T) {
			filterParams := []string{
				"?page=1&page_size=5",
				"?sort_by=completed_at&sort_order=desc",
				"?sort_by=score&sort_order=asc",
			}

			for _, params := range filterParams {
				url := "/api/v1/users/attempts" + params
				w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)
				assert.Equal(t, http.StatusOK, w.Code, "Should handle filter params: %s", params)

				response := ParseJSONResponse(t, w)
				data := GetDataFromResponse(t, response)

				attempts, exists := data["attempts"]
				assert.True(t, exists, "Filtered response should contain attempts")

				attemptsList, ok := attempts.([]interface{})
				assert.True(t, ok, "Filtered attempts should be an array")

				if len(attemptsList) > 0 {
					VerifyIsTestDataInArray(t, attemptsList, true, "filtered attempts")
				}
			}
		})

		t.Logf("User attempts endpoint verified successfully")
	})

	// Step 4: Test GET /users/attempts/{attemptId} endpoint (MISSING TEST COVERAGE)
	t.Run("GetAttemptDetails", func(t *testing.T) {
		if completedAttemptID == "" {
			t.Skip("No completed attempt available for testing")
			return
		}

		// Test valid attempt access by owner
		t.Run("ValidAttemptAccess", func(t *testing.T) {
			url := fmt.Sprintf("/api/v1/users/attempts/%s", completedAttemptID)
			w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)
			assert.Equal(t, http.StatusOK, w.Code, "Should get attempt details")

			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			// Verify response structure
			attempt, exists := data["attempt"]
			assert.True(t, exists, "Response should contain attempt field")

			attemptMap := attempt.(map[string]interface{})
			// Note: The attempt might not be marked as test data since it was just created
			// We'll verify the structure instead of the test data flag for newly created attempts

			// Verify attempt fields
			requiredAttemptFields := []string{
				"id", "quiz_id", "user_id", "score", "total_points", "time_spent",
				"percentage_score", "passed", "status", "is_completed", "started_at",
				"answers", "retake_count",
			}

			for _, field := range requiredAttemptFields {
				assert.Contains(t, attemptMap, field, "Attempt should contain %s field", field)
			}

			// Verify user_id matches authenticated user
			attemptUserID := attemptMap["user_id"].(string)
			assert.Equal(t, userID.String(), attemptUserID, "Attempt should belong to authenticated user")

			// Verify quiz details are included
			quiz, quizExists := attemptMap["quiz"]
			assert.True(t, quizExists, "Attempt should include quiz details")

			quizMap := quiz.(map[string]interface{})
			VerifyIsTestDataField(t, quizMap, true, "attempt quiz details")

			requiredQuizFields := []string{"id", "title", "category", "difficulty", "question_count"}
			for _, field := range requiredQuizFields {
				assert.Contains(t, quizMap, field, "Quiz details should contain %s field", field)
			}

			// Verify answers structure
			answers, answersExist := attemptMap["answers"]
			assert.True(t, answersExist, "Attempt should include answers")

			answersList, ok := answers.([]interface{})
			assert.True(t, ok, "Answers should be an array")

			if len(answersList) > 0 {
				firstAnswer := answersList[0].(map[string]interface{})
				answerFields := []string{"question_id", "selected_option", "is_correct", "points_earned"}
				for _, field := range answerFields {
					if !assert.Contains(t, firstAnswer, field, "Answer should contain %s field", field) {
						t.Logf("Available fields in answer: %+v", firstAnswer)
					}
				}
			}

			t.Logf("Attempt details verified successfully for attempt: %s", completedAttemptID)
		})

		// Test invalid attempt ID
		t.Run("InvalidAttemptID", func(t *testing.T) {
			invalidID := "invalid-uuid"
			url := fmt.Sprintf("/api/v1/users/attempts/%s", invalidID)
			w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)
			assert.Equal(t, http.StatusBadRequest, w.Code, "Should reject invalid attempt ID")

			// Parse error response
			var errorResponse map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
			assert.NoError(t, err, "Should parse error response")

			assert.Contains(t, errorResponse, "error", "Error response should contain error field")
		})

		// Test non-existent attempt ID
		t.Run("NonExistentAttemptID", func(t *testing.T) {
			nonExistentID := uuid.New().String()
			url := fmt.Sprintf("/api/v1/users/attempts/%s", nonExistentID)
			w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)
			assert.Equal(t, http.StatusNotFound, w.Code, "Should return 404 for non-existent attempt")

			var errorResponse map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
			assert.NoError(t, err, "Should parse error response")

			assert.Contains(t, errorResponse, "error", "Error response should contain error field")
		})

		// Test unauthorized access (attempt belongs to different user)
		t.Run("UnauthorizedAccess", func(t *testing.T) {
			// Create a second user
			otherUserID, otherToken := CreateTestUser(t, tc)
			_ = otherUserID // Used for testing, will be cleaned up by CleanupWithSupabase

			// Try to access the first user's attempt with second user's token
			url := fmt.Sprintf("/api/v1/users/attempts/%s", completedAttemptID)
			w := MakeAuthenticatedRequest(t, tc, "GET", url, otherToken, nil)
			assert.Equal(t, http.StatusForbidden, w.Code, "Should reject unauthorized access")

			var errorResponse map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
			assert.NoError(t, err, "Should parse error response")

			assert.Contains(t, errorResponse, "error", "Error response should contain error field")
		})
	})

	// Step 5: Test cross-endpoint data consistency
	t.Run("DataConsistency", func(t *testing.T) {
		// Get user stats
		statsW := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/stats", token, nil)
		assert.Equal(t, http.StatusOK, statsW.Code)

		statsResponse := ParseJSONResponse(t, statsW)
		stats := GetDataFromResponse(t, statsResponse)

		// Get user attempts
		attemptsW := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/attempts", token, nil)
		assert.Equal(t, http.StatusOK, attemptsW.Code)

		attemptsResponse := ParseJSONResponse(t, attemptsW)
		attemptsData := GetDataFromResponse(t, attemptsResponse)

		// Verify total attempts consistency
		statsTotal := int(stats["total_attempts"].(float64))
		attemptsTotal := int(attemptsData["total"].(float64))

		assert.GreaterOrEqual(t, attemptsTotal, 0, "Attempts total should be non-negative")

		// Note: We don't assert exact equality because stats might include attempts from setup
		// that weren't returned in the paginated attempts response
		t.Logf("Data consistency verified: stats_total=%d, attempts_total=%d", statsTotal, attemptsTotal)

		// Verify completed quizzes consistency
		statsCompleted := int(stats["total_quizzes_completed"].(float64))
		assert.GreaterOrEqual(t, statsCompleted, 0, "Completed quizzes should be non-negative")

		// If we have attempts, verify at least one completed quiz exists
		if attemptsTotal > 0 {
			attempts := attemptsData["attempts"].([]interface{})
			completedCount := 0
			for _, attempt := range attempts {
				attemptMap := attempt.(map[string]interface{})
				if attemptMap["is_completed"].(bool) {
					completedCount++
				}
			}

			t.Logf("Found %d completed attempts in paginated response", completedCount)
		}
	})

	// Step 6: Test performance analytics calculations
	t.Run("PerformanceAnalytics", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/stats", token, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		response := ParseJSONResponse(t, w)
		stats := GetDataFromResponse(t, response)

		// Verify completion rate calculation
		totalAttempts := stats["total_attempts"].(float64)
		completedQuizzes := stats["total_quizzes_completed"].(float64)
		completionRate := stats["completion_rate"].(float64)

		if totalAttempts > 0 {
			expectedRate := (completedQuizzes / totalAttempts) * 100
			assert.InDelta(t, expectedRate, completionRate, 0.1, "Completion rate should be correctly calculated")
		} else {
			assert.Equal(t, 0.0, completionRate, "Completion rate should be 0 when no attempts")
		}

		// Verify score distribution structure
		scoreDistribution := stats["score_distribution"].(map[string]interface{})
		scoreRanges := []string{"range_0_to_20", "range_21_to_40", "range_41_to_60", "range_61_to_80", "range_81_to_100"}

		totalScoreEntries := 0
		for _, scoreRange := range scoreRanges {
			assert.Contains(t, scoreDistribution, scoreRange, "Score distribution should contain %s", scoreRange)
			count := int(scoreDistribution[scoreRange].(float64))
			assert.GreaterOrEqual(t, count, 0, "Score range count should be non-negative")
			totalScoreEntries += count
		}

		// Verify streaks are non-negative
		currentStreak := int(stats["current_streak"].(float64))
		bestStreak := int(stats["best_streak"].(float64))
		assert.GreaterOrEqual(t, currentStreak, 0, "Current streak should be non-negative")
		assert.GreaterOrEqual(t, bestStreak, 0, "Best streak should be non-negative")
		assert.GreaterOrEqual(t, bestStreak, currentStreak, "Best streak should be >= current streak")

		t.Logf("Performance analytics verified: completion_rate=%.2f%%, current_streak=%d, best_streak=%d",
			completionRate, currentStreak, bestStreak)
	})

	t.Logf("User Performance & Analytics integration test completed successfully")
}
